package master

import (
	"os"

	"github.com/splinter0/api/database"
	"github.com/splinter0/api/miner"
)

func worker(jobChan <-chan string) {
	for job := range jobChan {
		program := miner.Extractor(job)
		fetched := database.GetProgramByName(program.Name)
		if program.Link == fetched.Link {
			if !program.Equals(&fetched) {
				database.UpdateProgramRepos(fetched.ID, program.Repos)
			}
		} else {
			database.AddProgram(program)
		}
	}
}

func Start() {
	links, err := miner.Crawler(os.Getenv("MAGIC_LINK"))
	if err != nil {
		panic("Could not start the service! Error: " + err.Error())
	}
	channel := make(chan string, 300)
	for i := 0; i < 4; i++ {
		go worker(channel)
	}
	for l := range links {
		channel <- links[l]
	}
}
