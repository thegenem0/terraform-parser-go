package main

import (
	"encoding/json"
	"fmt"
	"io"
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

	routes := gin.Default()

	routes.POST("/graph", func(ctx *gin.Context) {

		requestBody := ctx.Request.Body
		requestBytes, _ := io.ReadAll(requestBody)

		var payload *tfjson.Plan

		err := json.Unmarshal(requestBytes, &payload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Failed to unmarshal request body: %s", err),
			})
			return
		}

		graph, err := HandleGraphRoute(payload)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"graph": graph,
		})
	})

	routes.Run()
}

func GetFullData(state *tfjson.Plan) *tfjson.Plan {
	return state
}

func HandleGraphRoute(plan *tfjson.Plan) (*depsgraph.DepsGraph, error) {
	graph := depsgraph.NewDepsGraph()

	BuildTree(graph, plan.PlannedValues.RootModule)

	return graph, nil
}
