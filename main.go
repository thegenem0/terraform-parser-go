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
	"github.com/thegenem0/terraspect_server/pkg/reflector"
)

func main() {

	database, err := db.NewDBService()
	if err != nil {
		panic(err)
	}

	changeService := changes.NewChangeService()
	reflectorService := reflector.NewReflectorService(changeService)

	appCtx, err := appctx.Init(database, changeService, reflectorService)
	if err != nil {
		panic(err)
	}

	appCtx.Service().Database.AutoMigrate()

	routes := gin.Default()

	routes.GET("/graph", func(ctx *gin.Context) {
		graph, err := GetHandleGraphRoute(appCtx)
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

		graph, err := PostHandleGraphRoute(appCtx, payload)
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

func GetHandleGraphRoute(ctx appctx.AppContext) (*depsgraph.DepsGraph, error) {
	var plan db.Plan
	ctx.Service().Database.Connection().First(&plan)

	graph := depsgraph.NewDepsGraph()

	var storedPlan *tfjson.Plan

	err := json.Unmarshal(plan.TerraformPlan, &storedPlan)
	if err != nil {
		return nil, err
	}

	BuildTree(ctx, graph, storedPlan.PlannedValues.RootModule, storedPlan.ResourceChanges)

	return graph, nil
}

func PostHandleGraphRoute(ctx appctx.AppContext, plan *tfjson.Plan) (*depsgraph.DepsGraph, error) {
	graph := depsgraph.NewDepsGraph()

	BuildTree(ctx, graph, plan.PlannedValues.RootModule, plan.ResourceChanges)

	return graph, nil
}
