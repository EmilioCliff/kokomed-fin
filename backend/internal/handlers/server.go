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
}

func NewServer(config pkg.Config, maker pkg.JWTMaker, repo *mysql.MySQLRepo, payment *payments.PaymentService) *Server {
	r := gin.Default()

	s := &Server{
		router:   r,
		config:   config,
		maker:    maker,
		repo:     repo,
		payments: payment,
		ln:       nil,
	}

	s.setUpRoutes()

	return s
}

func (s *Server) setUpRoutes() {
	s.router.GET("/health-check", s.healthCheckHandler)

	// users routes
	s.router.POST("/login", s.loginUser)
	s.router.POST("/refresh-token/:email", s.refreshToken)
	s.router.POST("/user", s.createUser)
	s.router.GET("/user", s.listUsers)
	s.router.GET("/user/:id", s.getUser)
	s.router.PATCH("/user/reset-password", s.updateUserCredentials)
	s.router.PATCH("/user/:id", s.updateUser)

	// clients routes
	s.router.POST("/client", s.createClient)
	s.router.GET("/client", s.listClients)
	s.router.GET("/client/branch/:id", s.listClientsByBranch)
	s.router.GET("/client/status", s.listClientsByActive) // use query params
	s.router.GET("/client/:id", s.getClient)
	s.router.PATCH("/client/:id", s.updateClient)

	// product routes
	s.router.POST("/product", s.createProduct)
	s.router.GET("/product", s.listProducts)
	s.router.GET("/product/branch/:id", s.listProductsByBranch)
	s.router.GET("/product/:id", s.getProduct)

	// non-posted routes
	s.router.GET("/non-posted/all", s.listAllNonPostedPayments)
	s.router.GET("/non-posted/unassigned", s.listUnassignedNonPostedPayments)
	s.router.GET("/non-posted/by-id/:id", s.getNonPostedPayment)
	s.router.GET("/non-posted/by-type/:type", s.listNonPostedByTransactionSource)

	// branches routes
	s.router.GET("/branch", s.listBranches)
	s.router.GET("/branch/:id", s.getBranch)
	s.router.POST("/branch", s.createBranch)
	s.router.PATCH("/branch/:id", s.updateBranch)

	// loans routes
	s.router.POST("/loan", s.createLoan)
	s.router.PATCH("/loan/:id/disburse", s.disburseLoan)
	s.router.PATCH("/loan/:id/assign", s.transferLoanOfficer)
	s.router.GET("/loan/:id", s.getLoan)
	s.router.GET("/loan", s.listLoansByCategory)

	// payments routes
	s.router.POST("/payment/callback", s.paymentCallback)
	s.router.PATCH("/payment/:id/assign", s.paymentByAdmin)
	// payment of from credit to repay some loan(overpayment to pay loan)

	s.srv = &http.Server{
		Addr:    s.config.HTTP_PORT,
		Handler: s.router.Handler(),
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
