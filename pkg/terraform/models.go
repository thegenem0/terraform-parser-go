package terraform

import (
	"encoding/json"
)

type TFState struct {
	Version          int8            `json:"version"`
	TerraformVersion string          `json:"terraform_version"`
	FormatVersion    string          `json:"format_version"`
	Serial           int64           `json:"serial"`
	Lineage          string          `json:"lineage"`
	Outputs          json.RawMessage `json:"outputs"`
	Resources        []Resource      `json:"resources"`
}

type Resource struct {
	Module    string     `json:"module"`
	Mode      string     `json:"mode"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Provider  string     `json:"provider"`
	Instances []Instance `json:"instances"`
}

type Instance struct {
	Attributes interface{} `json:"attributes"`
}
