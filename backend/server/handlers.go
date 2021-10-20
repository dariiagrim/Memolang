package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"memolang/migrations/models"
	"net/http"
)

type Error struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func newError(message string, code int) *Error {
	return &Error{
		Error: message,
		Code:  code,
	}
}


func (s *Server) handleRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			fmt.Println(err)
		}
		if err := s.store.Db.Create(&user); err != nil {
			fmt.Println(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	}
}