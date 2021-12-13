package server

import (
	firebase "firebase.google.com/go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"memolang/configuration"
	"memolang/store"
)

type Server struct {
	router      *gin.Engine
	config      *configuration.Configuration
	store       *store.Store
	firebaseApp *firebase.App
}

func NewServer(config *configuration.Configuration, app *firebase.App) *Server {
	return &Server{
		config:      config,
		firebaseApp: app,
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
	s.router.Use(s.handleAuth())
	s.router.POST("/register", s.handleRegistration())
	s.router.POST("internal/add/topics", s.internalHandleAddTopic())
	s.router.GET("/username/unique/:username", s.handleCheckUsernameUniqueness())
	s.router.GET("/collection/:id", s.handleGetCollectionById())
	s.router.POST("/collection/add/:id", s.handleAddCollectionToMyCollections())
	s.router.GET("/collections/user", s.handleGetMyCollections())
	s.router.GET("/collections/other", s.handleGetOtherCollections())
	s.router.POST("/user/avatar", s.handleSetUserAvatar())
	s.router.GET("/user/avatar", s.handleGetUserAvatar())
	s.router.GET("/user/info", s.handleGetUserInfo())
	s.router.POST("/user/info", s.handleChangeUserInfo())
	s.router.GET("/user/all", s.handleGetAllUserWords())
	s.router.POST("/user/point", s.handleAddUserPoint())
	s.router.GET("/user/points", s.handleGetUserPoints())
}

func (s *Server) configureStore() error {
	st := store.NewStore(&s.config.StoreConfig)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st

	return nil
}
