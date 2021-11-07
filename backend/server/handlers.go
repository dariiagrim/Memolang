package server

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"memolang/migrations/models"
	"net/http"
	"time"
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

func (s *Server) internalHandleAddTopic() gin.HandlerFunc {
	return func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["files"]

		for _, f := range files {
			file, err := f.Open()
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
			}
			reader := csv.NewReader(file)
			reader.FieldsPerRecord = -1
			reader.Comma = ';'
			csvLines, err := reader.ReadAll()
			topic := models.Topic{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				General:   true,
			}
			topicWords := make(map[string]string, 0)
			for index, line := range csvLines {
				if index == 0 {
					topic.Name = line[0]
				} else {
					topicWords[line[1]] = line[0]
				}
			}
			if err := s.store.Db.Create(&topic).Error; err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
			}
			for en, fr := range topicWords {
				var word models.Word
				if err := s.store.Db.
					Where("english_spelling = ? and french_spelling = ?", en, fr).
					Find(&word).Error; err != nil {
					c.AbortWithError(http.StatusBadRequest, err)
				}
				wordId := int64(0)
				if word.Id == 0{
					word = models.Word{
						EnglishSpelling: en,
						FrenchSpelling:  fr,
					}
					if err := s.store.Db.Create(&word).Error; err != nil {
						c.AbortWithError(http.StatusBadRequest, err)
					}
					wordId = word.Id
				}
				topicsWord := models.TopicsWord{
					TopicId:   topic.Id,
					WordId:    wordId,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				if err := s.store.Db.Create(&topicsWord).Error; err != nil {
					c.AbortWithError(http.StatusBadRequest, err)
				}
			}
			file.Close()
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	}
}
