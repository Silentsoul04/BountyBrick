package miner

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/splinter0/api/database"
	"github.com/splinter0/api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// This function takes a page of firebounty.com and
// extracts the link to each program crawling all results
func Crawler(url string) ([]string, error) {
	var links []string
	resp, err := soup.Get(url)
	if err != nil || len(resp) == 0 {
		return links, errors.New("Could not extract page from: " + url)
	}
	dom := soup.HTMLParse(resp)
	nums := dom.Find("section", "id", "intro").Find("div", "class", "row").FindAll("a")
	pages := []soup.Root{dom}
	for n := range nums {
		if strings.Contains(nums[n].Text(), "Next") || strings.Contains(nums[n].Text(), "Prev") {
			continue
		}
		p, _ := soup.Get(os.Getenv("ROOT_LINK") + "/" + nums[n].Attrs()["href"])
		pages = append(pages, soup.HTMLParse(p))
	}
	for p := range pages {
		rows := pages[p].FindAll("div", "class", "Rtable-row")
		for r := range rows {
			if l, ok := rows[r].Attrs()["data-url"]; ok {
				links = append(links, l)
			}
		}
	}
	return links, nil
}

// This function takes a link to a program and
// returns the object containing information about it
// /https:\/\/github.com\/[^\/]+\/[^\/]+/g
func Extractor(p string) models.Program {
	var program models.Program
	program.ID = primitive.NewObjectID()

	// Build DOM object
	resp, _ := soup.Get(os.Getenv("ROOT_LINK") + p)
	dom := soup.HTMLParse(resp)

	// Extract basic info from dom
	name := dom.Find("div", "class", "rowText").Find("h2").Find("a").Text()
	link := dom.Find("div", "class", "propertyContent").Find("a", "class", "buttonColor").Attrs()["href"]
	scope := dom.Find("table", "id", "inscope").FindAll("a")
	scope = append(scope, dom.Find("table", "id", "inscope").FindAll("p")...)

	// Iter information from program scope
	repos := []primitive.ObjectID{}
	for s := range scope {
		val, ok := scope[s].Attrs()["href"] // Check if it's a link
		if !ok {
			val = scope[s].Text() // Check if it's just normal text
		}
		// For now we only want specific github repos (will have to fix so it understands weird links too)
		if strings.Contains(val, "https://github.com/") {
			exists, id := database.RepoExistsLink(val)
			if !exists {
				u := strings.Split(val, "/")
				short := u[len(u)-2] + "/" + u[len(u)-1]

				// Collect some info about repo
				req, _ := http.NewRequest(
					"GET",
					os.Getenv("GITHUB_API")+"repos/"+short,
					nil,
				)
				req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_OAUTH"))
				client := &http.Client{}
				resp, _ := client.Do(req)
				var result struct {
					Stars int `json:"stargazers_count"`
					Forks int `json:"forks_count"`
					Size  int `json:"size"`
				}
				json.NewDecoder(resp.Body).Decode(&result)

				id := primitive.NewObjectID()
				repo := models.Repo{
					ID:          id,
					Name:        u[len(u)-1],
					Link:        val,
					Program:     program.ID,
					ProgramName: name,
					Forked:      false,
					GitForks:    result.Forks,
					GitStars:    result.Stars,
					Size:        result.Size,
					Short:       short,
					Created:     time.Now(),
				}
				database.AddRepo(repo)
			}
			repos = append(repos, id)
		}
	}
	program.Created = time.Now()
	program.Repos = repos
	program.Name = name
	program.Link = link
	return program
}
