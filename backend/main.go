package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "user-management/docs"
)

var db *gorm.DB
var sqlDB *sql.DB

type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Cause   error  `swaggertype:"string"`
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

// @title			User Management API
// @version		1.0
// @description	This is a simple User Management API.
// @BasePath		/api
func main() {
	initDB()
	defer closeDB()

	router := gin.Default()

	router.GET("/api/users", getUsers)
	router.GET("/api/users/:id", getUser)
	router.POST("/api/users", saveUser)
	router.PUT("/api/users/:id", updateUser)
	router.DELETE("/api/users/:id", deleteUser)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":8080")
}

// @Summary		Get all users
// @Description	Get a list of all users
// @Produce		json
// @Success		200	{array}	User
// @Router			/api/users [get]
func getUsers(c *gin.Context) {
	var users []User
	db.Find(&users)
	c.JSON(http.StatusOK, users)
}

// @Summary		Get a user by ID
// @Description	Get a user by ID
// @Produce		json
// @Param			id	path		int	true	"User ID"
// @Success		200	{object}	User
// @Failure		404	{object}	ErrorResponse
// @Router			/api/users/{id} [get]
func getUser(c *gin.Context) {
	var user User
	id := c.Param("id")
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary		Create a new user
// @Description	Create a new user
// @Accept			json
// @Produce		json
// @Param			user	body		User	true	"User object"
// @Success		201		{object}	User
// @Failure		400		{object}	ErrorResponse
// @Router			/api/users [post]
func saveUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	db.Create(&user)
	c.JSON(http.StatusCreated, user)
}

// @Summary		Update a user by ID
// @Description	Update a user by ID
// @Accept			json
// @Produce		json
// @Param			id		path		int		true	"User ID"
// @Param			user	body		User	true	"User object"
// @Success		200		{object}	User
// @Failure		400		{object}	ErrorResponse
// @Router			/api/users/{id} [put]
func updateUser(c *gin.Context) {
	var user User
	id := c.Param("id")

	// Check if the user with the given ID exists
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Bind the JSON request payload to the existing user
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Update the user in the database
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary		Delete a user by ID
// @Description	Delete a user by ID
// @Produce		json
// @Param			id	path	int	true	"User ID"
// @Success		204	"No Content"
// @Failure		404	{object}	ErrorResponse
// @Router			/api/users/{id} [delete]
func deleteUser(c *gin.Context) {
	var user User
	id := c.Param("id")

	// Check if the user with the given ID exists
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := db.Delete(&User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
