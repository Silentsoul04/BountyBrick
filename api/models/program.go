package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Program struct {
	ID      primitive.ObjectID   `bson:"_id"`
	Name    string               `json:"name"` // Name of the program
	Link    string               `json:"link"` // Link to the original program (HackerOne, BugCrowd)
	Repos   []primitive.ObjectID `json:"repos"`
	Created time.Time            `json:"created"` // When program was added to scope
	Updated time.Time            `json:"updated"` // Last time the program was updated
}

func (p *Program) Equals(other *Program) bool {
	for i := range p.Repos {
		found := false
		for j := range other.Repos {
			if p.Repos[i] == other.Repos[j] {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
