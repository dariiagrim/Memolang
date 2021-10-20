package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"memolang/configuration"
	"memolang/store"
)

type Server struct {
	router *gin.Engine
	config *configuration.Configuration
	store  *store.Store
}

func NewServer(config *configuration.Configuration) *Server {
	return &Server{
		config: config,
	}
}

func (s *Server) Start() error {
	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	if err := s.router.Run(s.config.ServerConfig.Port); err != nil {
		return s.router.Run(":8080")
	}
	return nil
}

func (s *Server) configureRouter() {
	s.router = gin.Default()
	s.router.MaxMultipartMemory = 100 << 20

	s.router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}))
	s.router.POST("/register", s.handleRegistration())
}

func (s *Server) configureStore() error {
	st := store.NewStore(&s.config.StoreConfig)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st

	return nil
}
