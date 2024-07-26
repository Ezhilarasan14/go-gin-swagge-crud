package main

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/swaggo/files" // swagger embed files
    "github.com/swaggo/gin-swagger" // gin-swagger middleware

    _ "gin_test/docs" // replace with your actual module path
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// User represents a user model
type User struct {
    ID        uint      `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
}

// PongResponse represents the ping response
type PongResponse struct {
    Message string `json:"message"`
}

var db *gorm.DB
var err error

// @title Gin-GORM Example API
// @version 1.0
// @description This is a sample server for a Gin-GORM application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.basic BasicAuth

func main() {
    // Initialize the database
    db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect to the database")
    }

    // Auto-migrate the User model
    db.AutoMigrate(&User{})

    // Initialize the Gin router
    r := gin.Default()

    // Define routes
    r.POST("/users", createUser)
    r.GET("/users", getUsers)
    r.GET("/users/:id", getUser)
    r.PUT("/users/:id", updateUser)
    r.DELETE("/users/:id", deleteUser)

    // Swagger endpoint
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Start the server
    r.Run(":8080")
}

// createUser handles the creation of a new user
// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body User true "User"
// @Success 200 {object} User
// @Failure 400 {object} map[string]interface{} "error message"
// @Router /users [post]
func createUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
        return
    }
    db.Create(&user)
    c.JSON(http.StatusOK, user)
}

// getUsers retrieves all users
// @Summary Get users
// @Description Get all users
// @Tags users
// @Produce  json
// @Success 200 {array} User
// @Router /users [get]
func getUsers(c *gin.Context) {
    var users []User
    db.Find(&users)
    c.JSON(http.StatusOK, users)
}

// getUser retrieves a user by ID
// @Summary Get user by ID
// @Description Get a user by ID
// @Tags users
// @Produce  json
// @Param id path uint true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} map[string]interface{} "error message"
// @Router /users/{id} [get]
func getUser(c *gin.Context) {
    var user User
    if err := db.First(&user, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, map[string]interface{}{"error": "User not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}

// updateUser updates a user's details by ID
// @Summary Update user
// @Description Update a user's details by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path uint true "User ID"
// @Param user body User true "User"
// @Success 200 {object} User
// @Failure 400 {object} map[string]interface{} "error message"
// @Failure 404 {object} map[string]interface{} "error message"
// @Router /users/{id} [put]
func updateUser(c *gin.Context) {
    var user User
    if err := db.First(&user, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, map[string]interface{}{"error": "User not found"})
        return
    }

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
        return
    }

    db.Save(&user)
    c.JSON(http.StatusOK, user)
}

// deleteUser deletes a user by ID
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Param id path uint true "User ID"
// @Success 204
// @Failure 404 {object} map[string]interface{} "error message"
// @Router /users/{id} [delete]
func deleteUser(c *gin.Context) {
    var user User
    if err := db.First(&user, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, map[string]interface{}{"error": "User not found"})
        return
    }
    db.Delete(&user)
    c.Status(http.StatusNoContent)
}
