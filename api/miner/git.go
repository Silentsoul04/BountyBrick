package miner

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/splinter0/api/database"
	"github.com/splinter0/api/models"
)

func Fork(repo models.Repo) {
	body := []byte(`{"organization": "` + os.Getenv("GITHUB_ORG") + `"}`)
	req, _ := http.NewRequest(
		"POST",
		os.Getenv("GITHUB_API")+"repos/"+repo.Short+"/forks",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_OAUTH"))
	client := &http.Client{}
	client.Do(req)
	database.SetForked(repo.ID, true)
}

func Delete(repo models.Repo) {
	req, _ := http.NewRequest(
		"DELETE",
		os.Getenv("GITHUB_API")+"repos/"+os.Getenv("GITHUB_ORG")+"/"+repo.Name,
		nil,
	)
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_OAUTH"))
	client := &http.Client{}
	client.Do(req)
	database.SetForked(repo.ID, false)
}

// Fetch forks from Github and compare them with our DB
// cuz github is never wrong and we probably are :)
func GetForks() []models.Repo {
	var forks []models.Repo
	req, _ := http.NewRequest(
		"GET",
		os.Getenv("GITHUB_API")+"orgs/"+os.Getenv("GITHUB_ORG")+"/repos",
		nil,
	)
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_OAUTH"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return forks
	}

	var result []struct {
		Name string `json:"name"`
		Fork bool   `json:"bool"`
		Link string `json:""`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalln(err)
		return forks
	}
	for r := range result {
		if result[r].Fork {
			if repo := database.GetRepoByName(result[r].Name); repo.Name == result[r].Name {
				forks = append(forks, repo)
			}
		}
	}
	return forks
}

// Remove all repositories from github profile
func WipeOut() {
	repos := GetForks()
	for r := range repos {
		Delete(repos[r])
		database.SetForked(repos[r].ID, false)
	}
}
