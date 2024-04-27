package main

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/thegenem0/terraspect_server/pkg/db"
)

func TestMain(t *testing.T) {
	database, err := db.NewDBService()
	if err != nil {
		panic(err)
	}
	defer database.Close()
	database.AutoMigrate()

	stateFile, _ := os.Open("terraform/output.json")

	stateBytes, _ := io.ReadAll(stateFile)

	database.Connection.Create(&db.Plan{TerraformPlan: stateBytes})

	time.Sleep(5 * time.Second)

	if database.Connection.Error != nil {
		t.Errorf("Failed to create state in database")
	}
}
