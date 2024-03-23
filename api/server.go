package api

import (
	"log"

	db "github.com/dompham21/simplebank/db/sqlc"
	"github.com/dompham21/simplebank/token"
	"github.com/dompham21/simplebank/util"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.JWTMaker
	config     util.Config
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) *Server {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker:", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: *tokenMaker,
	}
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)

	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
