package security

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/splinter0/api/database"
	"github.com/splinter0/api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println(err)
	}
	return string(bytes)
}

func VerifyPassword(hashed string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return false
	}
	return true
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

type claim struct {
	Username string
	Role     string
	jwt.StandardClaims
}

func GenerateToken(username, role string) string {
	c := &claim{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(SECRET_KEY))
	return token
}

// Returns *claim and bool if expired or not
func ValidateToken(t string) *claim {
	token, err := jwt.ParseWithClaims(
		t,
		&claim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		return nil
	}
	c, ok := token.Claims.(*claim)
	if !ok || c.ExpiresAt < time.Now().Local().Unix() {
		return nil
	}
	check := database.GetUserToken(c.Username)
	if check != t {
		return nil
	}
	return c
}

func NotAuth(c *gin.Context) {
	c.JSON(403, gin.H{
		"message": "Unauthorized.",
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			NotAuth(c)
			c.Abort()
			return
		}
		claims := ValidateToken(token)
		if claims == nil {
			NotAuth(c)
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
	}
}

func CreateAdmin(username, password string) {
	admin := models.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: HashPassword(password),
		Role:     "root",
	}
	if database.AddUser(admin) {
		fmt.Println("Created admin user! " + admin.Username + ":" + password)
	} else {
		fmt.Println("Could not create admin user!")
	}
}
