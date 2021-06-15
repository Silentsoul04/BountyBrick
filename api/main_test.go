package main

import (
	"os"
	"testing"

	"github.com/splinter0/api/security"
)

/* ENVIRONMENT VARIABLES TESTS */

func TestEnvironment(t *testing.T) {
	variables := []string{
		"MONGO_URI",
		"SECRET_KEY",
		"MAGIC_LINK",
		"ROOT_LINK",
		"GITHUB_API",
		"GITHUB_OAUTH",
		"GITHUB_ORG",
		"DEBRICKED_API",
		"DEBRICKED_USER",
		"DEBRICKED_PASS",
	}
	for e := range variables {
		if _, present := os.LookupEnv(variables[e]); !present {
			t.Errorf("Expected environment variable %s to be set", variables[e])
		}
	}
}

/* SECURITY RELATED TESTS */

func TestPasswordHashing(t *testing.T) {
	const password string = "yayDebricked123!"
	var hashed string = security.HashPassword(password)
	if !security.VerifyPassword(hashed, password) {
		t.Errorf("Expected password %s with hash %s to pass the verification", password, hashed)
	}
}

/*func TestTokenGeneration(t *testing.T) {
	const username string = "splinter"
	const role string = "root"
	var token string = security.GenerateToken(username, role)
	claim := security.ValidateToken(token)
	if claim == nil {
		t.Errorf("Token validation with value %s resulted in nil claim", token)
	} else {
		assert.Equal(t, username, claim.Username)
		assert.Equal(t, role, claim.Role)
	}
}*/
