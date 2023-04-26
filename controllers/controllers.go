package controllers

import (
	"fmt"
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

func AdminLoginHandler(c *gin.Context) {
	var user models.User
	_, err := c.Cookie("admintoken")
	if err == nil {
		c.Redirect(http.StatusFound, "/admin_panel")
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
	if user.Role == "ADMIN" {
		//calling function to get the token string
		tknstr, err := middlewares.GenerateToken(user.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Not able to generate token"})
			return
		}
		//token string set to the cookie
		c.SetCookie("admintoken", tknstr, int(time.Now().Add(30*time.Minute).Unix()), "/", "localhost", false, true)
		c.Redirect(http.StatusFound, "/admin_panel")
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unauthorized User"})
		return
	}
}
func HomeAdmin(c *gin.Context) {
	var user models.User
	user_id, ok := c.Get("user_id")
	fmt.Println(user_id)
	if ok {
		initializers.DB.Where("id=?", user_id).Find(&user)
		c.JSON(http.StatusOK, user)
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unauthorized access"})
}
func AdminLogoutHandler(c *gin.Context) {
	_, err := c.Cookie("admintoken")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	c.SetCookie("admintoken", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"logout": "Success"})
}
func UserLoginHandler(c *gin.Context) {
	var user models.User
	_, err := c.Cookie("usertoken")
	if err == nil {
		c.Redirect(http.StatusFound, "/user_home")
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
	if user.Role == "USER" {
		//calling function to get the token string
		tknstr, err := middlewares.GenerateToken(user.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Not able to generate token"})
			return
		}
		//token string set to the cookie
		c.SetCookie("usertoken", tknstr, int(time.Now().Add(30*time.Minute).Unix()), "/", "localhost", false, true)
		c.Redirect(http.StatusFound, "/user_home")
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unauthorized User"})
	}
}

func HomeUser(c *gin.Context) {
	var user models.User
	user_id, ok := c.Get("user_id")
	if ok {
		initializers.DB.Where("id=?", user_id).Find(&user)
		c.JSON(http.StatusOK, user)
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unauthorized access"})
}

func UserLogoutHandler(c *gin.Context) {
	_, err := c.Cookie("usertoken")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	c.SetCookie("usertoken", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"logout": "Success"})
}
func SigninHandler(c *gin.Context) {
	var body struct {
		FirstName string
		LastName  string
		Email     string
		MobileNum string
		Password  string
	}
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Not able to bind JSON body"})
		return
	}
	hashpassword, _ := HashPassword(body.Password)
	user := models.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		MobileNum: body.MobileNum,
		Password:  hashpassword,
	}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.MobileNum == "" || user.Password == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Fill the required fields"})
			return
		} else {
			initializers.DB.Where("email=?", user.Email).Find(&user)
			if user.ID != 0 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not registered"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Register": "success"})
}
