package main

import (
	"sportzone/pkg/controllers"
	"sportzone/pkg/initializers"
	"sportzone/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.Loadenvariable()
	initializers.Initdb()
	initializers.Syncdatabase()
}
func main() {
	r := gin.Default()
	r.POST("/admin_login", controllers.AdminLoginHandler)
	r.GET("/admin_panel", middlewares.AuthAdmin(), controllers.HomeAdmin)
	r.GET("/admin_logout", controllers.AdminLogoutHandler)
	r.POST("/user_login", controllers.UserLoginHandler)
	r.GET("/user_home", middlewares.AuthUser(), controllers.HomeUser)
	r.GET("/user_logout", controllers.UserLogoutHandler)
	r.POST("/user_registration", controllers.SigninHandler)
	r.Run(":8080")
}
