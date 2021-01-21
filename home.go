package main

import (
	"fmt"
	"technique/controllers"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	router = gin.Default()

	racine := router.Group("/")
	{
		racine.POST("/add/users", controllers.AddUsers)
		racine.POST("/login", controllers.CheckLogin)
		racine.DELETE("/delete/user/:id", controllers.DeleteUser)
		racine.GET("/users/list", controllers.GetUsers)
		racine.GET("/user/:id", controllers.GetUserById)
		racine.PUT("/user/:id", controllers.EditUserByID)
	}
	//démarrer le serveur
	router.Run()
	fmt.Println("Hello azeaez!")
}
