package main

import (
	"sportzone/controllers"
	"sportzone/initializers"
	"sportzone/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.Loadenvariable()
	initializers.Initdb()
	initializers.Syncdatabase()
}
func main() {
	r := gin.Default()
	r.POST("/login", controllers.LoginHandler)
	r.GET("/home", middlewares.VerifyToken(), controllers.HomeHandler)
	r.GET("/logout", controllers.LogoutHandler)
	r.Run(":8080")
}
