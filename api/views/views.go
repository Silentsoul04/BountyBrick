package views

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/splinter0/api/database"
	"github.com/splinter0/api/debricked"
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

// Helper function to validate program
func valProg(p string) (bool, models.Program) {
	var prog models.Program
	id, err := primitive.ObjectIDFromHex(p)
	if err != nil {
		return false, prog
	}
	prog = database.GetProgramById(id)
	if prog.Link == "" {
		return false, prog
	}
	return true, prog
}

// Returns program based on id
func GetProgram(c *gin.Context) {
	p := c.Param("id")
	if ok, prog := valProg(p); ok {
		c.JSON(200, gin.H{
			"message": "success",
			"program": prog,
		})
	} else {
		c.JSON(400, gin.H{
			"message": "No program with id: " + p + " found!",
		})
	}
}

// Execute action on program (all repos contained)
func ProgAction(c *gin.Context) {
	action := c.Param("action")
	if _, ok := actions[action]; !ok {
		c.JSON(400, gin.H{
			"message": "The action: " + action + " isn't valid!",
			"actions": displayActions(),
		})
		return
	}

	var body struct {
		Programs []string `json:"programs"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		bad(c)
		return
	}

	response := make(map[string]string, len(body.Programs))
	// Iter over all programs provided
	for i := range body.Programs {
		if ok, program := valProg(body.Programs[i]); ok {
			// Iter over all repositories in program
			for r := range program.Repos {
				if ok, repo := valRepo(program.Repos[r].Hex()); ok {
					go actions[action].Func(repo)
				}
			}
			response[program.ID.Hex()] = "Successfully started action: " + action
		} else {
			response[body.Programs[i]] = "No such program found!"
		}
	}
	c.JSON(200, gin.H{
		"message":  "success",
		"programs": response,
	})
}

// Still not implemented
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
		"message":      "success",
		"repositories": database.GetAllRepos(),
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

type Action struct {
	Desc string
	Func func(models.Repo)
}

var actions = map[string]Action{
	"fork": {
		"Fork the repository",
		miner.Fork,
	},
	"remove": {
		"Remove repository from github page",
		miner.Delete,
	},
	"scan": {
		"Run a Debricked scan on repository",
		debricked.RunScan,
	},
	"bookmark": {
		"Bookmark repository to personal profile",
		func(models.Repo) {}, //TODO has to be executed in views context
	},
}

func displayActions() map[string]string {
	a := make(map[string]string, len(actions))
	for k, v := range actions {
		a[k] = v.Desc
	}
	return a
}

func RepoAction(c *gin.Context) {
	action := c.Param("action")
	if _, ok := actions[action]; !ok {
		c.JSON(400, gin.H{
			"message": "The action: " + action + " isn't valid!",
			"actions": displayActions(),
		})
		return
	}

	var body struct {
		Repos []string `json:"repos"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		bad(c)
		return
	}

	response := make(map[string]string, len(body.Repos))
	for i := range body.Repos {
		if ok, repo := valRepo(body.Repos[i]); ok {
			go actions[action].Func(repo)
			response[repo.Name] = "Successfully started action: " + action
		} else {
			response[body.Repos[i]] = "No repository with id: " + body.Repos[i] + " found!"
		}
	}
	c.JSON(200, gin.H{
		"message": "success",
		"repos":   response,
	})
}

func Actions(c *gin.Context) {
	c.JSON(200, gin.H{
		"message":  "success",
		"actions":  displayActions(),
		"programs": "Every action can also be executed on programs, affecting all contained repos",
	})
}
