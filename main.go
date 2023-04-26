package main

import (
	"sportzone/controllers"
	"sportzone/initializers"

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
	r.GET("/home", controllers.HomeHandler)
	r.Run(":8080")
}
