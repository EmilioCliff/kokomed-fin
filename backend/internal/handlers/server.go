package handlers

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/payments"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	ln     net.Listener
	srv    *http.Server

	config pkg.Config
	maker  pkg.JWTMaker
	repo   *mysql.MySQLRepo

	payments services.PaymentService
	worker services.WorkerService
}

func NewServer(config pkg.Config, maker pkg.JWTMaker, repo *mysql.MySQLRepo, payment *payments.PaymentService, worker services.WorkerService) *Server {
	r := gin.Default()

	s := &Server{
		router:   r,
		config:   config,
		maker:    maker,
		repo:     repo,
		worker: worker,
		payments: payment,
		ln:       nil,
	}

	s.setUpRoutes()

	return s
}

func (s *Server) setUpRoutes() {
	s.router.Use(CORSmiddleware())
	
	v1 := s.router.Group("/api/v1")
	v1Auth := s.router.Group("/api/v1")

	// protected routes
	authRoute := v1Auth.Use(authMiddleware(s.maker))

	// health check
	s.router.GET("/health-check", s.healthCheckHandler)

	// users routes
	v1.POST("/login", s.loginUser)
	v1.GET("/refreshToken", s.refreshToken)
	authRoute.GET("/logout", s.logoutUser)
	authRoute.POST("/user", s.createUser)
	authRoute.GET("/user", s.listUsers)
	authRoute.GET("/user/:id", s.getUser)
	v1.POST("/forgot-password", s.forgotPassword)
	v1.PATCH("/user/reset-password/:token", s.updateUserCredentials)
	authRoute.PATCH("/user/:id", s.updateUser)

	// clients routes
	authRoute.POST("/client", s.createClient)
	authRoute.GET("/client", s.listClients)
	authRoute.GET("/client/branch/:id", s.listClientsByBranch)
	authRoute.GET("/client/status", s.listClientsByActive) // use query params
	authRoute.GET("/client/:id", s.getClient)
	authRoute.PATCH("/client/:id", s.updateClient)

	// product routes
	authRoute.POST("/product", s.createProduct)
	authRoute.GET("/product", s.listProducts)
	authRoute.GET("/product/branch/:id", s.listProductsByBranch)
	authRoute.GET("/product/:id", s.getProduct)

	// non-posted routes
	authRoute.GET("/non-posted/all", s.listAllNonPostedPayments)
	authRoute.GET("/non-posted/unassigned", s.listUnassignedNonPostedPayments)
	authRoute.GET("/non-posted/by-id/:id", s.getNonPostedPayment)
	authRoute.GET("/non-posted/by-type/:type", s.listNonPostedByTransactionSource)

	// branches routes
	authRoute.GET("/branch", s.listBranches)
	authRoute.GET("/branch/:id", s.getBranch)
	authRoute.POST("/branch", s.createBranch)
	authRoute.PATCH("/branch/:id", s.updateBranch)

	// loans routes
	authRoute.POST("/loan", s.createLoan)
	authRoute.PATCH("/loan/:id/disburse", s.disburseLoan)
	authRoute.PATCH("/loan/:id/assign", s.transferLoanOfficer)
	authRoute.GET("/loan/:id", s.getLoan)
	authRoute.GET("/loan", s.listLoansByCategory)

	// payments routes
	v1.POST("/payment/callback", s.paymentCallback)
	authRoute.PATCH("/payment/:id/assign", s.paymentByAdmin)

	// payment of from credit to repay some loan(overpayment to pay loan)


	// helper routes
	authRoute.GET("/helper/dashboard", s.getDashboardData)
	authRoute.GET("/helper/formData", s.getLoanFormData)

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
