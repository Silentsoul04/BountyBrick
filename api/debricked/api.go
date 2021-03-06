package debricked

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/splinter0/api/database"
	"github.com/splinter0/api/miner"
	"github.com/splinter0/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func RunScan(repo models.Repo) {
	miner.Fork(repo)
	miner.AddFile(
		ACTION,
		"Debricked Vulnerability Scan",
		repo.Name,
		"action.yml",
	)
	/* FIX DB NOT UPDATING PROPERLY */
	found := false
	var id string
	for !found {
		time.Sleep(10 * time.Second)
		active := GetRepositories()
		for a := range active {
			if active[a].Name == repo.Name {
				id = strconv.Itoa(active[a].ID)
				database.SetBrick(repo.ID, id)
				found = true
				break
			}
		}
	}
	commit := miner.GetLastCommit(repo.Name)
	finished := false
	for !finished {
		time.Sleep(5 * time.Second)
		status, _, _, _ := miner.GetRun(repo.Name, commit)
		/*database.SetCommit(repo.ID, commit)
		database.SetStatus(repo.ID, status)
		/*database.EditRepo(repo.ID, bson.M{
			"commit":      commit,
			"scan_status": status,
			"scan_result": title,
		})*/
		if status != "completed" {
			continue
		}
		// Since there is no open API endpoint on Debricked to
		// check the scan status I will use the github api
		database.SetBrick(repo.ID, id)
		vulns := GetLatest(id)
		database.SetCommit(repo.ID, commit)
		database.EditRepo(repo.ID, bson.M{"commit": commit})
		database.EditRepo(repo.ID, bson.M{"vulns": vulns})
		finished = true
	}
}

type Call struct {
	Type string
	Url  string
	Body []byte
}

// Fancy worker which self updates token and accepts jobs from channels
func DebrickedAPI(callChan <-chan Call, respChan chan *io.ReadCloser) {
	var token struct {
		Value   string
		Expires int64
	}
	client := &http.Client{}
	for call := range callChan {
		// Renew expired token
		if token.Expires < time.Now().Local().Unix() {
			token.Value = Login()
			token.Expires = time.Now().Local().Add(time.Hour * time.Duration(1)).Unix()
		}
		var body *bytes.Buffer = bytes.NewBuffer(call.Body)
		req, _ := http.NewRequest(
			call.Type,
			call.Url,
			body,
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token.Value)
		resp, err := client.Do(req)
		if err != nil {
			respChan <- nil
		} else {
			respChan <- &resp.Body
		}
	}
}

// I know it's ugly
var calls chan Call
var responses chan *io.ReadCloser

// Whoo spooky scary lazy starter
func ScaryDeamon() {
	calls = make(chan Call)
	responses = make(chan *io.ReadCloser)
	go DebrickedAPI(calls, responses)
}

func Login() (token string) {
	form := url.Values{}
	form.Add("_username", os.Getenv("DEBRICKED_USER"))
	form.Add("_password", os.Getenv("DEBRICKED_PASS"))
	resp, err := http.PostForm("https://app.debricked.com/api/login_check", form)
	if err != nil {
		log.Fatalln(err)
		return
	}
	var result struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalln(err)
		return
	}
	token = result.Token
	return
}

type drepo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetRepositories() (result []drepo) {
	calls <- Call{
		"GET",
		os.Getenv("DEBRICKED_API") + "zapier/repositories",
		[]byte{},
	}
	resp := <-responses
	if resp == nil {
		return
	}
	b, err := ioutil.ReadAll(*resp)
	if err != nil {
		log.Fatalln(err)
		return
	}
	content := strings.ReplaceAll(
		string(b),
		os.Getenv("GITHUB_ORG")+"\\/",
		"",
	)
	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func GetLatest(id string) (vulns []models.Vuln) {
	calls <- Call{
		"POST",
		os.Getenv("DEBRICKED_API") + "zapier/newcve/poll",
		[]byte(
			`{
				"repo": [
					` + id + `
				]
			}`,
		),
	}
	resp := <-responses
	if resp == nil {
		return
	}
	b, err := ioutil.ReadAll(*resp)
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = json.Unmarshal(b, &vulns)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func GetStatus(id string) {

}
