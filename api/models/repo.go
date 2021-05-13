package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repo struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `json:"name"`  // Name of the repo
	Link        string             `json:"link"`  // Github link
	Short       string             `json:"short"` // Short Github link (last part)
	Debricked   string             `json:"brick"` // Debricked ID
	ScanStatus  string             `json:"scan_status"`
	ScanResult  string             `json:"scan_result"`
	Commit      string             `json:"commit"` // Commit corresponding to the scan
	Vulns       []Vuln             `json:"vulns"`
	Program     primitive.ObjectID `json:"program"`      // Owned by what program?
	ProgramName string             `json:"program_name"` // Name of program
	Forked      bool               `json:"forked"`       // If repo has been forked already
	GitForks    int                `json:"git_forks"`
	GitStars    int                `json:"git_stars"`
	Size        int                `json:"size"`
	Created     time.Time          `json:"created"` // When repo was added to scope
	Updated     time.Time          `json:"updated"` // Last time the repo was updated
}

type Vuln struct {
	CVE  string  `json:"cve"`
	CVSS float32 `json:"cvss3"`
}
