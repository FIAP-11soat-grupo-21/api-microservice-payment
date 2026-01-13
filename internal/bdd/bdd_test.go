package bdd

import (
	"os"
	"payment_microservice/internal/bdd/suites"
	"testing"

	"github.com/cucumber/godog"
)

func TestMain(_ *testing.M) {
	status := godog.TestSuite{
		ScenarioInitializer: suites.InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"features"},
		},
	}.Run()
	os.Exit(status)
}
