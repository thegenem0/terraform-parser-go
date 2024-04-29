package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/thegenem0/terraspect_server/pkg/appctx"
	"github.com/thegenem0/terraspect_server/pkg/changes"
	"github.com/thegenem0/terraspect_server/pkg/db"
	"github.com/thegenem0/terraspect_server/pkg/depsgraph"
)

func main() {

	database, err := db.NewDBService()
	if err != nil {
		panic(err)
	}

	changeService := changes.NewChangeService()

	appCtx, err := appctx.Init(database, changeService)
	if err != nil {
		panic(err)
	}

	appCtx.Service().Database.AutoMigrate()

	routes := gin.Default()

	routes.GET("/graph", func(ctx *gin.Context) {
		graph, err := GetHandleGraphRoute(database, changeService)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"graph": graph,
		})
	})

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

		graph, err := PostHandleGraphRoute(payload, changeService)
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

func GetHandleGraphRoute(database *db.DBService, changeSerice *changes.ChangeService) (*depsgraph.DepsGraph, error) {
	var plan db.Plan
	database.Connection.First(&plan)

	graph := depsgraph.NewDepsGraph()

	var storedPlan *tfjson.Plan

	err := json.Unmarshal(plan.TerraformPlan, &storedPlan)
	if err != nil {
		return nil, err
	}

	BuildTree(graph, storedPlan.PlannedValues.RootModule, changeSerice, storedPlan.ResourceChanges)

	return graph, nil
}

func PostHandleGraphRoute(plan *tfjson.Plan, changeService *changes.ChangeService) (*depsgraph.DepsGraph, error) {
	graph := depsgraph.NewDepsGraph()

	BuildTree(graph, plan.PlannedValues.RootModule, changeService, plan.ResourceChanges)

	return graph, nil
}
