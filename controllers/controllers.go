package controllers

import (
	"net/http"
	"sportzone/initializers"
	"sportzone/middlewares"
	"sportzone/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Password hashing
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// password authorization
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func LoginHandler(c *gin.Context) {
	var user models.User
	_, err := c.Cookie("token")
	if err == nil {
		c.Redirect(http.StatusFound, "/home")
		return
	}
	var body struct {
		Email    string
		Password string
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, "Error to bind JSON format")
		return
	}
	initializers.DB.Where("email=?", body.Email).Find(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, "Invalid Email address")
		return
	}
	ok := CheckPasswordHash(body.Password, user.Password)
	if !ok {
		c.JSON(http.StatusUnauthorized, "Invalid Password")
		return
	}
	//calling function to get the token string
	tknstr, err := middlewares.GenerateToken(user.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Not able to generate token"})
		return
	}
	//token string set to the cookie
	c.SetCookie("token", tknstr, int(time.Now().Add(30*time.Minute).Unix()), "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/home")
}

func HomeHandler(c *gin.Context) {
	user, ok := c.Get("user")
	if ok {
		c.JSON(http.StatusOK, user)
	}
}

func LogoutHandler(c *gin.Context) {
	_, err := c.Cookie("token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"logout": "Success"})
}
