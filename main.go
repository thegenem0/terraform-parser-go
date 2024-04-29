package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/thegenem0/terraspect_server/pkg/appctx"
	"github.com/thegenem0/terraspect_server/pkg/builder"
	"github.com/thegenem0/terraspect_server/pkg/changes"
	"github.com/thegenem0/terraspect_server/pkg/db"
	"github.com/thegenem0/terraspect_server/pkg/reflector"
)

func main() {

	database, err := db.NewDBService()
	if err != nil {
		panic(err)
	}

	reflectorService := reflector.NewReflectorService()
	changeService := changes.NewChangeService(reflectorService)

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

	routes.GET("/changes", func(ctx *gin.Context) {
		changes, err := GetHandleChangesRoute(appCtx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"changes": changes,
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

		tree, err := PostHandleGraphRoute(appCtx, payload)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"tree": tree.Nodes,
		})
	})

	routes.Run()
}

func GetFullData(state *tfjson.Plan) *tfjson.Plan {
	return state
}

func GetHandleChangesRoute(ctx appctx.AppContext) ([]changes.Change, error) {
	var plan db.Plan
	ctx.Service().Database.Connection().First(&plan)

	var storedPlan *tfjson.Plan

	err := json.Unmarshal(plan.TerraformPlan, &storedPlan)
	if err != nil {
		return []changes.Change{}, err
	}

	ctx.Service().ChangeService.BuildChanges(storedPlan.ResourceChanges)

	return ctx.Service().ChangeService.GetChanges(), nil
}

func GetHandleGraphRoute(ctx appctx.AppContext) (builder.TreeData, error) {
	treeBuilder := builder.NewTreeBuilder(ctx.Service().ReflectorService)

	var plan db.Plan
	ctx.Service().Database.Connection().First(&plan)

	var storedPlan *tfjson.Plan

	err := json.Unmarshal(plan.TerraformPlan, &storedPlan)
	if err != nil {
		return builder.TreeData{}, err
	}

	treeBuilder.BuildTree(storedPlan.PlannedValues.RootModule)

	return treeBuilder.GetTree(), nil
}

func PostHandleGraphRoute(ctx appctx.AppContext, plan *tfjson.Plan) (builder.TreeData, error) {
	treeBuilder := builder.NewTreeBuilder(ctx.Service().ReflectorService)

	treeBuilder.BuildTree(plan.PlannedValues.RootModule)

	return treeBuilder.GetTree(), nil
}
