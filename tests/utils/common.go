package utils

import (
	"github.com/cucumber/godog"
	"log"
	"os"
	"testing"
)

var CasperT *testing.T

func TestFeatures(t *testing.T, featureName string, scenarioInitializer func(*godog.ScenarioContext)) {

	dir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Working dir: %s ", dir)

	suite := godog.TestSuite{
		ScenarioInitializer: scenarioInitializer,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/" + featureName},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	CasperT = t

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
