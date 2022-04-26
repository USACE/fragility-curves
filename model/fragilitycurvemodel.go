package model

import (
	"fmt"
	"math/rand"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
	"github.com/USACE/filestore"
	"github.com/usace/wat-api/wat"
)

type FragilityCurveModel struct {
	Name      string                   `json:"name"`
	Locations []FragilityCurveLocation `json:"locations"`
}
type FragilityCurveLocation struct {
	Name           string                           `json:"location"`
	FragilityCurve paireddata.UncertaintyPairedData `json:"stage-probability"`
}

func (fcm FragilityCurveModel) Compute(modelpayload *wat.ModelPayload, fs filestore.FileStore) (float64, error) {
	realizationSeed := modelpayload.Realization.Seed
	eventSeed := modelpayload.Event.Seed
	realizationRandom := rand.New(rand.NewSource(realizationSeed))
	eventRandom := rand.New(rand.NewSource(eventSeed))
	failure_elevation := 0.0
	for _, fcl := range fcm.Locations {
		//sample fragility curve for a location with knowledge uncertianty
		pd := fcl.FragilityCurve.SampleValueSampler(realizationRandom.Float64())
		//sample sampledfragility curve at a location with natural variability
		failure_elevation = pd.SampleValue(eventRandom.Float64())
		fmt.Println(failure_elevation)
	}
	return failure_elevation, nil
}
