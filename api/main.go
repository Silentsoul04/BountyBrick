package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/splinter0/api/debricked"
	"github.com/splinter0/api/master"
	"github.com/splinter0/api/miner"
	"github.com/splinter0/api/security"
	"github.com/splinter0/api/views"
)

const (
	CRT string = "server.crt"
	KEY string = "server.key"
)

func main() {
	fmt.Println("Starting BountyBrick Service...")
	go master.Start()

	miner.AddSecret("DEBRICKED_USERNAME", os.Getenv("DEBRICKED_USER"))
	miner.AddSecret("DEBRICKED_PASSWORD", os.Getenv("DEBRICKED_PASS"))

	debricked.ScaryDeamon()
	/*vulns := debricked.GetLatest("5307")
	for i := range vulns {
		fmt.Println(vulns[i].CVE, vulns[i].CVSS)
	}*/

	r := gin.Default()

	// Routes
	r.POST("/api/login", views.Login)
	// Middleware
	r.Use(security.AuthMiddleware())
	// Authorization required
	r.GET("/api/", views.Index)
	r.GET("/api/programs", views.Programs)
	r.GET("/api/programs/:id", views.GetProgram)
	r.POST("/api/programs/:action", views.ProgAction)
	r.GET("/api/repos", views.Repositories)
	r.GET("/api/repos/:id", views.GetRepository)
	r.POST("/api/repos/:action", views.RepoAction)
	r.GET("/api/actions", views.Actions)
	// Start HTTPS
	err := http.ListenAndServeTLS(":443", CRT, KEY, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
