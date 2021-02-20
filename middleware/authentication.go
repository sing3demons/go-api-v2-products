package middleware

import (
	"app/config"
	"app/models"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type formLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

//Login - sign in
func Login(c *gin.Context) {
	var form formLogin
	var user models.User
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unable to bind data"})
		return
	}

	copier.Copy(&user, &form)
	db := config.GetDB()
	if err := db.First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	jwt, _ := jwtSign(user)
	serializedUser := jwt
	// copier.Copy(&form, &user)
	c.JSON(http.StatusOK, gin.H{"token": serializedUser})
}

func jwtSign(user models.User) (string, error) {
	atClaims := jwt.MapClaims{}

	atClaims["id"] = user.ID
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Local().Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	return at.SignedString([]byte(os.Getenv("SECRET_KEY")))

}

//JwtVerify - call this methos to add interceptor
func JwtVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Split(c.Request.Header["Authorization"][0], " ")[1] == "" {
			fmt.Print("err")
		}
		tokenString := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok && !token.Valid {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		}

		var user models.User
		id := fmt.Sprintf("%v", claims["id"])
		db := config.GetDB()
		if err := db.First(&user, id).Error; err != nil {
			fmt.Println(err.Error())
		}

		var serializedUser userResponse
		copier.Copy(&serializedUser, &user)
		c.Set("jwt_id", serializedUser)

		c.Next()
	}
}
