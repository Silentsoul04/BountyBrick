package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/splinter0/api/security"
	"github.com/splinter0/api/views"
)

const (
	CRT string = "server.crt"
	KEY string = "server.key"
)

func main() {
	fmt.Println("Starting BountyBrick Service...")
	//go master.Start()
	r := gin.Default()

	// Routes
	r.POST("/api/login", views.Login)
	// Middleware
	r.Use(security.AuthMiddleware())
	// Authorization required
	r.GET("/api/", views.Index)
	r.GET("/api/programs", views.Programs)
	r.GET("/api/programs/:id", views.GetProgram)
	r.GET("/api/repos", views.Repositories)
	r.GET("/api/repos/:id", views.GetRepository)
	r.POST("/api/repos/:id", views.RepoAction)
	// Start HTTPS
	err := http.ListenAndServeTLS(":443", CRT, KEY, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
