package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var sqlDB *sql.DB

type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	initDB()
	defer closeDB()

	router := gin.Default()

	router.GET("/api/users", getUsers)
	router.GET("/api/users/:id", getUser)
	router.POST("/api/users", saveUser)
	router.PUT("/api/users/:id", updateUser)
	router.DELETE("/api/users/:id", deleteUser)

	router.Run(":8080")
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("./db.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err = db.DB()
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&User{})
}

func closeDB() {
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func getUsers(c *gin.Context) {
	var users []User
	db.Find(&users)
	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	var user User
	id := c.Param("id")
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func saveUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	db.Create(&user)
	c.JSON(http.StatusCreated, user)
}

func updateUser(c *gin.Context) {
	var user User
	id := c.Param("id")

	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	db.Save(&user)
	c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
