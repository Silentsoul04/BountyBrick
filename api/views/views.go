package views

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/splinter0/api/database"
	"github.com/splinter0/api/miner"
	"github.com/splinter0/api/models"
	"github.com/splinter0/api/security"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

func bad(c *gin.Context) {
	c.JSON(400, gin.H{
		"message": "Bad request",
	})
}

type login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role"`
}

func Login(c *gin.Context) {
	var l login
	err := c.ShouldBindJSON(&l)
	if err != nil {
		bad(c)
		return
	}
	validationErr := validate.Struct(l)
	if validationErr != nil {
		bad(c)
		return
	}
	user := database.FindUser(l.Username)
	if user.Username == l.Username && security.VerifyPassword(user.Password, l.Password) {
		token := security.GenerateToken(user.Username, user.Role)
		database.AddUserToken(user.Username, token)
		c.JSON(200, gin.H{
			"message":  "success",
			"token":    token,
			"username": user.Username,
		})
	} else {
		security.NotAuth(c)
	}
}

func Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "success",
		"role":    c.GetString("role"),
		"user":    c.GetString("username"),
	})
}

// Returns all the programs
func Programs(c *gin.Context) {
	c.JSON(200, gin.H{
		"message":  "success",
		"programs": database.GetAllPrograms(),
	})
}

// Returns program based on id
func GetProgram(c *gin.Context) {
	p := c.Param("id")
	id, err := primitive.ObjectIDFromHex(p)
	if err != nil {
		bad(c)
		return
	}
	program := database.GetProgramById(id)
	if program.Link != "" {
		c.JSON(200, gin.H{
			"message": "success",
			"program": program,
		})
	} else {
		c.JSON(400, gin.H{
			"message": "No program with id: " + p + " found!",
		})
	}
}

var filters = map[string]string{
	"forked":     "Show only the ones that were forked",
	"scanned":    "Show only the ones that have a scan completed or ongoing",
	"bookmarked": "Show only bookmarked ones",
	"completed":  "Show only the ones with a complete scan",
	"critical":   "Show only the ones which have a critical vulnerability",
}

func Repositories(c *gin.Context) {
	//filter := c.Query("filter")
	c.JSON(200, gin.H{
		"message":  "success",
		"programs": database.GetAllRepos(),
	})
}

// Helper function to validate repo
func valRepo(p string) (bool, models.Repo) {
	var repo models.Repo
	id, err := primitive.ObjectIDFromHex(p)
	if err != nil {
		return false, repo
	}
	repo = database.GetRepo(id)
	if repo.Link == "" {
		return false, repo
	}
	return true, repo
}

func GetRepository(c *gin.Context) {
	p := c.Param("id")
	if ok, repo := valRepo(p); ok {
		c.JSON(200, gin.H{
			"message":    "success",
			"repository": repo,
		})
	} else {
		c.JSON(400, gin.H{
			"message": "No repository with id: " + p + " found!",
		})
	}
}

var actions = map[string]string{
	"fork":     "Fork the repository",
	"remove":   "Remove repository from github page",
	"scan":     "Run a Debricked scan on repository",
	"bookmark": "Bookmark repository to personal profile",
}

func RepoAction(c *gin.Context) {
	p := c.Param("id")
	if ok, repo := valRepo(p); ok {
		action := c.Query("action")
		switch action {
		case "fork":
			miner.Fork(repo)
		case "remove":
			miner.Delete(repo)
		default:
			c.JSON(400, gin.H{
				"message": "The action: " + action + " isn't valid!",
				"actions": actions,
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "Successfully executed action: " + action + " on " + repo.Name,
		})
	} else {
		c.JSON(400, gin.H{
			"message": "No repository with id: " + p + " found!",
		})
	}
}
