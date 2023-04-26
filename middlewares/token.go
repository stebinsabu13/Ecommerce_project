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

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tknstr, err := c.Cookie("token")
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
		}
		if token.Valid {
			if time.Now().Unix() > claims.ExpiresAt.Unix() {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			}
			var user models.User

			initializers.DB.Where("email=?", claims.UserName).Find(&user)

			if user.ID == 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized User"})
			}

			c.Set("user", user)
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

// func GenerateToken(c *gin.Context) {
// 	user_email, ok := c.Get("user_email")
// 	if !ok {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized User"})
// 	}
// 	expirationTime := time.Now().Add(30 * time.Minute)
// 	claims := Claims{UserName: user_email, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expirationTime)}}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tknstr, err := token.SignedString(jwtkey)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Not able to generate token"})
// 		return
// 	}
// 	c.SetCookie("token", tknstr, int(time.Now().Add(30*time.Minute).Unix()), "/", "localhost", false, true)
// 	c.Redirect(http.StatusFound, "/home")
// }
