package handlers

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/payments"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	router *gin.Engine
	ln     net.Listener
	srv    *http.Server

	// dependencies
	config pkg.Config
	maker  pkg.JWTMaker
	repo   *mysql.MySQLRepo

	// services
	payments services.PaymentService
	worker services.WorkerService
	cache services.CacheService
	report services.ReportService

	// otel components
	meter metric.Meter
	logger *slog.Logger
	tracer trace.Tracer
}

func NewServer(config pkg.Config, maker pkg.JWTMaker, repo *mysql.MySQLRepo, payment *payments.PaymentService, worker services.WorkerService, cache services.CacheService, report services.ReportService) *Server {
	r := gin.Default()

	meter := otel.Meter("")
	logger := otelslog.NewLogger("")
	tracer := otel.Tracer("kokomed-fin")


	s := &Server{
		router:   r,
		ln:       nil,

		config:   config,
		maker:    maker,
		repo:     repo,

		payments: payment,
		worker: worker,
		cache: cache,
		report: report,

		meter: meter,
		logger: logger,
		tracer: tracer,
	}

	s.setUpRoutes()

	return s
}

func (s *Server) setUpRoutes() {
	s.router.Use(otelgin.Middleware(s.config.SERVER_NAME)) // otelgin middleware
	s.router.Use(CORSmiddleware())
	
	v1 := s.router.Group("/api/v1")
	v1Auth := s.router.Group("/api/v1")

	// protected routes
	authRoute := v1Auth.Use(authMiddleware(s.maker))

	// cached routes
	cachedRoutes := authRoute.Use(redisCacheMiddleware(s.cache))

	// health check
	s.router.GET("/health-check", s.healthCheckHandler)

	// users routes
	v1.POST("/login", s.loginUser)
	v1.GET("/refreshToken", s.refreshToken)
	authRoute.GET("/logout", s.logoutUser)
	authRoute.POST("/user", s.createUser)
	cachedRoutes.GET("/user", s.listUsers)
	authRoute.GET("/user/:id", s.getUser)
	v1.POST("/forgot-password", s.forgotPassword)
	v1.PATCH("/user/reset-password/:token", s.updateUserCredentials)
	authRoute.PATCH("/user/:id", s.updateUser)

	// clients routes
	authRoute.POST("/client", s.createClient)
	cachedRoutes.GET("/client", s.listClients)
	authRoute.GET("/client/:id", s.getClient)
	authRoute.PATCH("/client/:id", s.updateClient)

	// product routes
	authRoute.POST("/product", s.createProduct)
	cachedRoutes.GET("/product", s.listProducts)
	authRoute.GET("/product/:id", s.getProduct)

	// non-posted routes
	cachedRoutes.GET("/non-posted/all", s.listAllNonPostedPayments)
	authRoute.POST("/non-posted/clients", s.listClientsNonPosted)

	// branches routes
	cachedRoutes.GET("/branch", s.listBranches)
	authRoute.GET("/branch/:id", s.getBranch)
	authRoute.POST("/branch", s.createBranch)
	authRoute.PATCH("/branch/:id", s.updateBranch)

	// loans routes
	authRoute.POST("/loan", s.createLoan)
	authRoute.PATCH("/loan/:id/disburse", s.disburseLoan)
	authRoute.GET("/loan/:id/installments", s.getLoanInstallments)
	cachedRoutes.GET("/loan", s.listLoansByCategory)
	cachedRoutes.GET("/loan/expected-payments", s.listExpectedPayments)
	cachedRoutes.GET("/loan/unpaid-installments", s.listUnpaidInstallmentsData)

	// payments routes
	v1.POST("/payment/callback", s.paymentCallback)
	v1.POST("/payment/validation", s.validationCallback)
	authRoute.PATCH("/payment/:id/assign", s.paymentByAdmin)


	// helper routes
	authRoute.GET("/helper/dashboard", s.getDashboardData)
	authRoute.GET("/helper/formData", s.getLoanFormData)
	authRoute.GET("/helper/loanEvents", s.getLoanEvents)
	authRoute.GET("/mpesa/token", s.getMPESAAccesToken)
	authRoute.POST("/helper/client-payments", s.getClientNonPosted)

	// reports routes
	authRoute.POST("/report", s.generateReport)

	s.srv = &http.Server{
		Addr:         s.config.HTTP_PORT,
		Handler:      s.router.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (s *Server) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) Start() error {
	var err error
	if s.ln, err = net.Listen("tcp", s.config.HTTP_PORT); err != nil {
		return err
	}

	go func(s *Server) {
		err := s.srv.Serve(s.ln)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}(s)

	return nil
}

func (s *Server) Stop() error {
	log.Println("Shutting down http server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.srv.Shutdown(ctx)
}

func (s *Server) GetPort() int {
	if s.ln == nil {
		return 0
	}

	return s.ln.Addr().(*net.TCPAddr).Port
}

func errorResponse(err error) gin.H {
	return gin.H{
		"status_code": pkg.ErrorCode(err),
		"message":     pkg.ErrorMessage(err),
	}
}

func constructCacheKey(path string, queryParams map[string][]string) string {
	path = strings.TrimPrefix(path, "/api/v1/")

	var queryParts []string
	for key, values := range queryParams {
		for _, value := range values {
			queryParts = append(queryParts, fmt.Sprintf("%s=%s", key, value))
		}
	}
	sort.Strings(queryParts) // Sort to ensure cache key consistency

	return fmt.Sprintf("%s:%s", path, strings.Join(queryParts, ":"))
}
