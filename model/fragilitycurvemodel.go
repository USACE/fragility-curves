package model

import (
	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
)
type FragilityCurveLocation struct{
	Name string `json:"location"`
	FragilityCurve paireddata.UncertaintyPairedData `json:"stage-probability"`
}