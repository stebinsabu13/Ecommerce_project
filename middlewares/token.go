package middlewares

import (
	"net/http"
	"sportzone/initializers"
	"sportzone/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/securecookie"
)

var jwtkey = []byte(securecookie.GenerateRandomKey(64))

type Claims struct {
	UserName interface{}
	jwt.RegisteredClaims
}

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		tknstr, err := c.Cookie("usertoken")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Needs to login"})
			return
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tknstr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtkey, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Claims not parsed"})
			return
		}
		if token.Valid {
			//checking the expiry of the token
			if time.Now().Unix() > claims.ExpiresAt.Unix() {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				return
			}
			var user models.User
			//Finding the user
			initializers.DB.Where("email=?", claims.UserName).Find(&user)

			if user.ID == 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized User"})
				return
			}
			c.Set("user_id", user.ID)
		}
	}
}

func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tknstr, err := c.Cookie("admintoken")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Needs to login"})
			return
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tknstr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtkey, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Claims not parsed"})
			return
		}
		if token.Valid {
			//checking the expiry of the token
			if time.Now().Unix() > claims.ExpiresAt.Unix() {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				return
			}
			var user models.User
			//Finding the user
			initializers.DB.Where("email=?", claims.UserName).Find(&user)

			if user.ID == 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized User"})
				return
			}
			c.Set("user_id", user.ID)
		}
	}
}
func GenerateToken(userName string) (tokenString string, err error) {
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := Claims{UserName: userName, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expirationTime)}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtkey)
	return
}
