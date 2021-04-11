package debricked

import (
	"github.com/splinter0/api/miner"
	"github.com/splinter0/api/models"
)

func RunScan(repo models.Repo) {
	miner.Fork(repo)
	miner.AddFile(
		ACTION,
		"Debricked Vulnerability Scan",
		repo.Name,
		"action.yml",
	)
}
