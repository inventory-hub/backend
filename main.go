package main

import (
	"Smart-Machine/inventory-hub-2/controller"
	"Smart-Machine/inventory-hub-2/database"
	"Smart-Machine/inventory-hub-2/model"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load environment file
	loadEnv()
	// load database configuration and connection
	loadDatabase()
	// start the server
	serveApplication()
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed at loading .env file")
	}
	log.Println(".env file loaded successfully")
}

func loadDatabase() {
	database.InitDB()
	database.DB.AutoMigrate(&model.Role{})
	database.DB.AutoMigrate(&model.User{})
	// seedData()
}

// load seed data into the database
func seedData() {
	var roles = []model.Role{
		{Name: "admin", Description: "Administrator role"},
		{Name: "manager", Description: "Manager role"},
		{Name: "customer", Description: "Customer role"},
	}
	var user = []model.User{{
		Username: os.Getenv("ADMIN_USERNAME"),
		Email:    os.Getenv("ADMIN_EMAIL"),
		Password: os.Getenv("ADMIN_PASSWORD"),
		RoleID:   1,
	}}
	database.DB.Save(&roles)
	database.DB.Save(&user)
}

// TODO:
// * Autogenerated API docs with OpenAI, i.e. Swagger
// * Email Service linking to "/invite" endpoint

func serveApplication() {
	router := gin.Default()

	router.POST("/api/auth/login", controller.Login)
	router.POST("/api/auth/refresh", controller.Refresh)
	// router.POST("/api/auth/register", controller.Register)
	// router.POST("/api/auth/invite", middleware.ValidateAuthorization(middleware.AuthorizedRoles), controller.Invite)

	// router.GET("/api/users", middleware.ValidateAuthorization(middleware.AllRoles), controller.GetListOfUsers)
	// router.GET("/api/users/:id", middleware.ValidateAuthorization(middleware.AllRoles), controller.GetUserById)
	// router.PUT("/api/users/:id", middleware.ValidateAuthorization(middleware.AuthorizedRoles), controller.UpdateUser)
	// router.DELETE("/api/users/:id", middleware.ValidateAuthorization(middleware.AuthorizedRoles), controller.DeleteUser)

	PORT, _ := strconv.Atoi(os.Getenv("PORT"))

	router.Run(fmt.Sprintf(":%d", PORT))
	fmt.Printf("Server running on port %d\n", PORT)
}
