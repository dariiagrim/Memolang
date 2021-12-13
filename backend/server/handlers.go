package server

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v4"
	"io"
	"memolang/migrations/models"
	"net/http"
	"os"
	"strconv"
	"strings"
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
		var req handleRegistrationRequest
		if err := c.BindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}
		uuid, _ := c.Get("uuid")
		user := models.User{
			Uuid:     uuid.(string),
			Username: req.Username,
			Email:    req.Email,
		}
		if err := s.store.Db.Create(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
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
				c.AbortWithStatusJSON(http.StatusBadRequest, err)
				return
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
				c.AbortWithStatusJSON(http.StatusBadRequest, err)
			}
			for en, fr := range topicWords {
				var word models.Word
				if err := s.store.Db.
					Where("english_spelling = ? and french_spelling = ?", en, fr).
					Find(&word).Error; err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, err)
				}
				wordId := int64(0)
				if word.Id == 0 {
					word = models.Word{
						EnglishSpelling: en,
						FrenchSpelling:  fr,
					}
					if err := s.store.Db.Create(&word).Error; err != nil {
						c.AbortWithStatusJSON(http.StatusBadRequest, err)
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
					c.AbortWithStatusJSON(http.StatusBadRequest, err)
				}
			}
			file.Close()
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	}
}

func (s *Server) handleCheckUsernameUniqueness() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		var users []models.User
		if err := s.store.Db.Where("username = ?", username).Find(&users).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}
		unique := true
		if len(users) > 0 {
			unique = false
		}
		c.JSON(http.StatusOK, gin.H{
			"is_unique": unique,
		})
	}
}

func (s *Server) handleGetCollectionById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var collection models.Topic
		if err := s.store.Db.Where("id = ?", id).First(&collection).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		var words []models.Word
		if err := s.store.Db.Model(models.TopicsWord{}).
			Select("w.*").
			Joins("left join words w on word_id = w.id").
			Where("topic_id = ?", collection.Id).
			Find(&words).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		respWords := make([]responseWord, 0)

		for _, w := range words {
			respWords = append(respWords, responseWord{
				English: w.EnglishSpelling,
				French:  w.FrenchSpelling,
			})
		}

		resp := handleGetCollectionByIdResponse{
			Id:    collection.Id,
			Name:  collection.Name,
			Words: respWords,
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    resp.Id,
			"name":  resp.Name,
			"words": resp.Words,
		})
	}
}

func (s *Server) handleAddCollectionToMyCollections() gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}
		uuid, _ := c.Get("uuid")
		record := models.TopicsUser{
			TopicId: int64(id),
			UserId:  uuid.(string),
		}
		if err := s.store.Db.Create(&record).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func (s *Server) handleGetMyCollections() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid, _ := c.Get("uuid")
		collections := make([]collectionShortInfo, 0)
		var collectionsDb []models.Topic
		if err := s.store.Db.Model(models.TopicsUser{}).
			Select("t.*").
			Joins("left join topics t on topic_id = t.id").
			Where("user_id = ?", uuid).
			Find(&collectionsDb).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}
		for _, collection := range collectionsDb {
			collections = append(collections, collectionShortInfo{
				Id:   collection.Id,
				Name: collection.Name,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"data": collections,
		})
	}
}

func (s *Server) handleGetOtherCollections() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid, _ := c.Get("uuid")
		collections := make([]collectionShortInfo, 0)
		var collectionsDb []models.Topic
		if err := s.store.Db.Model(models.TopicsUser{}).
			Select("t.*").
			Joins("left join topics t on topic_id = t.id").
			Where("user_id != ?", uuid).
			Find(&collectionsDb).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}
		for _, collection := range collectionsDb {
			collections = append(collections, collectionShortInfo{
				Id:   collection.Id,
				Name: collection.Name,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"data": collections,
		})
	}
}

func (s *Server) handleSetUserAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid, _ := c.Get("uuid")
		form, _ := c.MultipartForm()
		reqFiles := form.File["file"]
		if len(reqFiles) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, errors.New("no file for avatar provided"))
		}
		file, err := reqFiles[0].Open()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		fileName := fmt.Sprintf("%v.png", uuid)
		fileOs, err := os.Create(fileName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		io.Copy(fileOs, file)

		var user models.User
		if err := s.store.Db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		user.Avatar = null.StringFrom(fileName)

		if err := s.store.Db.Save(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func (s *Server) handleGetUserAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid, _ := c.Get("uuid")

		var user models.User
		if err := s.store.Db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		file, err := os.ReadFile(user.Avatar.ValueOrZero())

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		c.JSON(http.StatusOK, gin.H{
			"file": file,
		})
	}
}

func (s *Server) handleGetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid, _ := c.Get("uuid")

		var user models.User
		if err := s.store.Db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		c.JSON(http.StatusOK, gin.H{
			"first_name":    user.FirstName,
			"lastName":      user.LastName,
			"date_of_birth": user.DateOfBirth,
			"email":         user.Email,
			"username":      user.Username,
			"points":        user.Points,
		})
	}
}

func (s *Server) handleChangeUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid, _ := c.Get("uuid")

		var user models.User
		if err := s.store.Db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		var req handleChangeUserInfoRequest
		if err := c.BindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		user.FirstName = req.FirstName
		user.LastName = req.LastName
		user.DateOfBirth = req.DateOfBirth

		if err := s.store.Db.Save(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func (s *Server) handleGetAllUserWords() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid, _ := c.Get("uuid")

		var words []models.Word
		if err := s.store.Db.Model(models.TopicsWord{}).
			Select("w.*").
			Joins("left join words w on word_id = w.id").
			Where("topic_id in (select topic_id from topics_users where user_id =?)", uuid).
			Find(&words).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		respWords := make([]responseWord, 0)

		for _, w := range words {
			respWords = append(respWords, responseWord{
				English: w.EnglishSpelling,
				French:  w.FrenchSpelling,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"words": respWords,
		})
	}
}

func (s *Server) handleAddUserPoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid, _ := c.Get("uuid")

		var user models.User
		if err := s.store.Db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		user.Points += 1

		if err := s.store.Db.Save(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func (s *Server) handleGetUserPoints() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid, _ := c.Get("uuid")

		var user models.User
		if err := s.store.Db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		c.JSON(http.StatusOK, gin.H{
			"points": user.Points,
		})
	}
}

func (s *Server) handleAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqToken := c.Request.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]
		ctx := context.Background()
		client, err := s.firebaseApp.Auth(ctx)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}

		token, err := client.VerifyIDToken(ctx, reqToken)
		if token == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errors.New("invalid token"))
		}
		c.Set("uuid", token.UID)
	}
}
