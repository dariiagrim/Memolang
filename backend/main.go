package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"memolang/migrations"
	"memolang/migrations/models"
)

func main() {
	db, err := gorm.Open(postgres.Open("host=localhost user=dariia password=DDG256 dbname=memolang port=5432"), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
		return
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetMigrations())

	if err = m.Migrate(); err != nil {
		fmt.Println(err)
		return
	}

	router := gin.Default()

	router.POST("/register", func(c *gin.Context) {

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			fmt.Println(err)
		}
		if err := db.Create(&user); err != nil {
			fmt.Println(err)
		}
	})
	router.Run(":8080")
}