package db

import "gorm.io/gorm"

type State struct {
	gorm.Model
	ID             uint
	TerraformState []byte
}

type Plan struct {
	gorm.Model
	ID            uint
	TerraformPlan []byte
}
