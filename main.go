package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/thegenem0/terraspect_server/pkg/db"
	"github.com/thegenem0/terraspect_server/pkg/depsgraph"
)

func main() {
	database, err := db.NewDBService()
	if err != nil {
		panic(err)
	}
	defer database.Close()
	database.AutoMigrate()

	graph := depsgraph.NewDepsGraph()
	//
	// stateFile, _ := os.Open("terraform/sandbox.json")
	//
	// stateBytes, _ := io.ReadAll(stateFile)
	//
	// database.Connection.Create(&db.State{TerraformState: stateBytes})

	var plan db.Plan

	database.Connection.Last(&plan)

	var stateData *tfjson.Plan

	err = json.Unmarshal(plan.TerraformPlan, &stateData)
	if err != nil {
		panic(err)
	}

	BuildTree(graph, stateData.PlannedValues.RootModule)

	r := gin.Default()

	r.GET("/graph", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"graph": graph,
		})
	})

	r.GET("/full", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"full": stateData,
		})
	})

	r.Run()
}

func GetFullData(state *tfjson.Plan) *tfjson.Plan {
	return state
}
