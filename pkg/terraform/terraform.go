package terraform

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

type TFService struct {
	Core *tfexec.Terraform
}

func NewTFService(tmpDir string) (*TFService, error) {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.0.0")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		return nil, err
	}

	tfCore, err := tfexec.NewTerraform(tmpDir, execPath)
	if err != nil {
		return nil, err
	}

	err = tfCore.Init(context.Background())
	if err != nil {
		return nil, err
	}

	return &TFService{
		Core: tfCore,
	}, nil
}

func (tfs *TFService) ParsePlan(plan []byte) (*tfjson.Plan, error) {
	var planData *tfjson.Plan

	err := json.Unmarshal(plan, &planData)
	if err != nil {
		return nil, err
	}

	return planData, nil
}
