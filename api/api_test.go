package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/splinter0/api/models"
)

var token string
var router *gin.Engine = routerSetup()
var repo string

func TestLogin(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer([]byte(
		`{
			"username":"root",
			"password":"`+ROOT+`"
		}`,
	)))
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Login response returned status code %d", w.Code)
	} else {
		var result struct {
			Message  string `json:"message"`
			Token    string `json:"token"`
			Username string `json:"username"`
		}
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Unexpected json response from login testing")
		} else {
			token = result.Token
		}
	}
}

func TestPrograms(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/programs", nil)
	req.Header.Set("token", token)
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Returned status code %d", w.Code)
	} else {
		var result struct {
			Message  string           `json:"message"`
			Programs []models.Program `json:"programs"`
		}
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil {
			t.Error("Unexpected json response")
		}
		if len(result.Programs) < 1 {
			t.Error("Unexpected number of programs returned")
		} else {
			prog := result.Programs[rand.Intn(len(result.Programs))].ID.Hex()
			req, _ = http.NewRequest("GET", "/api/programs/"+prog, nil)
			req.Header.Set("token", token)
			router.ServeHTTP(w, req)
			if w.Code != 200 {
				t.Errorf("Returned status code %d when querying single program", w.Code)
			}
		}
	}
}

func TestRepos(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/repos", nil)
	req.Header.Set("token", token)
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Returned status code %d", w.Code)
	} else {
		var result struct {
			Message      string        `json:"message"`
			Repositories []models.Repo `json:"repositories"`
		}
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil {
			t.Error("Unexpected json response")
		}
		if len(result.Repositories) < 1 {
			t.Error("Unexpected number of repositories returned")
		} else {
			repo = result.Repositories[rand.Intn(len(result.Repositories))].ID.Hex()
			req, _ = http.NewRequest("GET", "/api/repos/"+repo, nil)
			req.Header.Set("token", token)
			router.ServeHTTP(w, req)
			if w.Code != 200 {
				t.Fatalf("Returned status code %d when querying single repo", w.Code)
			}
		}
	}
}

func TestFork(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/repos/fork", bytes.NewBuffer([]byte(
		`{
			"repos":[
				"`+repo+`"
			]
		}`,
	)))
	req.Header.Set("token", token)
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Returned status code %d", w.Code)
	} else {
		var result struct {
			Message string `json:"message"`
		}
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil || result.Message != "success" {
			t.Fatalf("Unexpected json response")
		}
	}
}

func TestRemove(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/repos/remove", bytes.NewBuffer([]byte(
		`{
			"repos":[
				"`+repo+`"
			]
		}`,
	)))
	req.Header.Set("token", token)
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Returned status code %d", w.Code)
	} else {
		var result struct {
			Message string `json:"message"`
		}
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil || result.Message != "success" {
			t.Fatalf("Unexpected json response")
		}
	}
}

func TestScan(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/repos/scan", bytes.NewBuffer([]byte(
		`{
			"repos":[
				"`+repo+`"
			]
		}`,
	)))
	req.Header.Set("token", token)
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Returned status code %d", w.Code)
	} else {
		var result struct {
			Message string `json:"message"`
		}
		err := json.NewDecoder(w.Body).Decode(&result)
		if err != nil || result.Message != "success" {
			t.Error("Unexpected json response")
		}
	}
}
