// main.go
package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	initDB()
	defer db.Close()

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
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Emai string `json:"email"`
}

func getUsers(c *gin.Context) {
	var users []User

	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Name, &user.Emai)
		users = append(users, user)
	}

	c.JSON(200, users)
}

func getUser(c *gin.Context) {
	var user User
	id := c.Param("id")

	row := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)
	err := row.Scan(&user.ID, &user.Name, &user.Emai)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, user)
}

func saveUser(c *gin.Context) {
	var user User

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	stmt, _ := db.Prepare("INSERT INTO users (name, email) VALUES (?, ?)")
	result, err := stmt.Exec(user.Name, user.Emai)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	id, _ := result.LastInsertId()

	user.ID = int(id)

	c.JSON(200, user)
}

func updateUser(c *gin.Context) {
	var user User
	id := c.Param("id")

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	stmt, _ := db.Prepare("UPDATE users SET name = ?, email = ? WHERE id = ?")
	_, err = stmt.Exec(user.Name, user.Emai, id)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, user)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")

	stmt, _ := db.Prepare("DELETE FROM users WHERE id = ?")
	_, err := stmt.Exec(id)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted"})
}
