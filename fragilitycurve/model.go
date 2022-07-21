package fragilitycurve

import (
	"math/rand"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
	"github.com/usace/wat-go-sdk/plugin"
)

type Model struct {
	Name      string                   `json:"name"`
	Locations []FragilityCurveLocation `json:"locations"`
}
type FragilityCurveLocation struct {
	Name           string                           `json:"location"`
	FragilityCurve paireddata.UncertaintyPairedData `json:"stage-probability"`
}
type FragilityCurveLocationResult struct {
	Name             string  `json:"location"`
	FailureElevation float64 `json:"failure_elevation"`
}
type ModelResult struct {
	Results []FragilityCurveLocationResult `json:"results"`
}

func (fcm Model) Compute(seedSet plugin.SeedSet) (ModelResult, error) {
	realizationSeed := seedSet.RealizationSeed
	eventSeed := seedSet.EventSeed
	realizationRandom := rand.New(rand.NewSource(realizationSeed))
	eventRandom := rand.New(rand.NewSource(eventSeed))
	results := ModelResult{
		Results: make([]FragilityCurveLocationResult, len(fcm.Locations)),
	}
	for idx, fcl := range fcm.Locations {
		//sample fragility curve for a location with knowledge uncertianty
		pd := fcl.FragilityCurve.SampleValueSampler(realizationRandom.Float64())
		//sample sampledfragility curve at a location with natural variability
		locationResult := FragilityCurveLocationResult{
			Name:             fcl.Name,
			FailureElevation: pd.SampleValue(eventRandom.Float64()),
		}
		results.Results[idx] = locationResult
	}
	return results, nil
}
