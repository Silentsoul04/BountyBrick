package miner

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

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
	database.SetCommit(repo.ID, "")
}

// Fetch forks from Github and compare them with our DB
// cuz github is never wrong and we probably are :)
func GetForks(strict bool) []models.Repo {
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
		if !strict || result[r].Fork {
			if repo := database.GetRepoByName(result[r].Name); repo.Name == result[r].Name {
				forks = append(forks, repo)
			}
		}
	}
	return forks
}

// Remove all repositories from github profile
func WipeOut() {
	repos := GetForks(false)
	for r := range repos {
		Delete(repos[r])
		database.SetForked(repos[r].ID, false)
	}
}

// Public key for secrets
func GetPublicKey() (key, id string) {
	req, _ := http.NewRequest(
		"GET",
		os.Getenv("GITHUB_API")+"orgs/"+os.Getenv("GITHUB_ORG")+"/actions/secrets/public-key",
		nil,
	)
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_OAUTH"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return
	}
	var result struct {
		Key string `json:"key"`
		ID  string `json:"key_id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalln(err)
		return
	}
	key = result.Key
	id = result.ID
	return
}

// Add secret to organization (used for Debricked login)
func AddSecret(name, value string) {
	key, id := GetPublicKey()
	// TODO: Implement this in Golang
	cmd := exec.Command("python3", "security/secret.py", key, value)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln(err, string(out))
	}
	secret := string(out[:len(out)-1])
	body := []byte(
		`{
			"encrypted_value": "` + secret + `", 
			"key_id": "` + id + `",
			"visibility": "all"
		}`,
	)
	req, _ := http.NewRequest(
		"PUT",
		os.Getenv("GITHUB_API")+"orgs/"+os.Getenv("GITHUB_ORG")+"/actions/secrets/"+name,
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_OAUTH"))
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
}

func AddFile(content []byte, commit, repo, name string) {
	encoded := b64.StdEncoding.EncodeToString(content)
	body := []byte(
		`{
			"message": "` + commit + `",
			"content": "` + encoded + `"
		}`,
	)
	req, _ := http.NewRequest(
		"PUT",
		os.Getenv("GITHUB_API")+"repos/"+os.Getenv("GITHUB_ORG")+"/"+repo+"/contents/"+name,
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_OAUTH"))
	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetRun(repo, commit string) (status, title, url, summary string) {
	req, _ := http.NewRequest(
		"GET",
		os.Getenv("GITHUB_API")+"repos/"+os.Getenv("GITHUB_ORG")+"/"+repo+"/commits/"+commit+"/check-runs",
		nil,
	)
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_OAUTH"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return
	}
	var result struct {
		Runs []struct {
			Status string `json:"status"`
			Url    string `json:"details_url"`
			Output struct {
				Title string `json:"title"`
				Summ  string `json:"summary"`
			} `json:"output"`
		} `json:"check_runs"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalln(err)
		return
	}
	status = result.Runs[0].Status
	title = result.Runs[0].Output.Title
	url = result.Runs[0].Url
	summary = result.Runs[0].Output.Summ
	return
}

func GetLastCommit(repo string) (sha string) {
	stamp := time.Now().Local().Add(time.Minute * -1).Format("2006-01-02T15:04:05.000+0000")
	req, _ := http.NewRequest(
		"GET",
		os.Getenv("GITHUB_API")+"repos/"+os.Getenv("GITHUB_ORG")+"/"+repo+"/commits?since="+stamp,
		nil,
	)
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_OAUTH"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return
	}
	var result []struct {
		Sha string `json:"sha"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalln(err)
		return
	}
	if len(result) > 0 {
		sha = result[0].Sha
	}
	return
}

// This turns out to be useless, but lets keep it for the future
func GetWorkflows(repo string) {
	req, _ := http.NewRequest(
		"GET",
		os.Getenv("GITHUB_API")+"repos/"+os.Getenv("GITHUB_ORG")+"/"+repo+"/actions/runs",
		nil,
	)
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_OAUTH"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(string(b))
}
