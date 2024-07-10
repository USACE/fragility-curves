package fragilitycurve

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

type Model struct {
	Name      string                   `json:"name"`
	Locations []FragilityCurveLocation `json:"locations"`
}
type FragilityCurveLocation struct {
	Name           string                           `json:"location"`
	NLD_System_ID  string                           `json:"nld_system_id`
	NLD_Segment_ID string                           `json:"nld_segment_id`
	NIDID          string                           `json:"nidid`
	FailureMode    string                           `json:"failure_mode`
	Source         string                           `json:"source`
	FragilityCurve paireddata.UncertaintyPairedData `json:"probability-stage"`
}
type FragilityCurveLocationResult struct {
	Name             string  `json:"location"`
	FailureElevation float64 `json:"failure_elevation"`
}
type ModelResult struct {
	Results []FragilityCurveLocationResult `json:"results"`
}

func InitFragilityCurveLocation(locationbytes []byte) FragilityCurveLocation {
	file := string(locationbytes)
	lines := strings.Split(file, "\r\n")
	name := parseLine(lines[0])
	NLD_System_Id := parseLine(lines[1])
	NLD_Segment_ID := parseLine(lines[2])
	NIDID := parseLine(lines[3])
	failure_mode := parseLine(lines[4])
	source := parseLine(lines[5])
	lines = lines[7:] //skip stage,probability
	xvals := make([]float64, 0)
	yvals := make([]statistics.ContinuousDistribution, 0)

	for _, line := range lines {
		row := strings.Split(line, ",")
		if len(row) >= 2 {
			xstring := row[0]
			ystring := row[1]
			xval, err := strconv.ParseFloat(xstring, 64)
			if err != nil {
				return FragilityCurveLocation{}
			}
			yval, err := strconv.ParseFloat(ystring, 64)
			if err != nil {
				return FragilityCurveLocation{}
			}
			xvals = append(xvals, xval)
			yvals = append(yvals, statistics.DeterministicDistribution{Value: yval})
		}

	}
	return FragilityCurveLocation{
		Name:           name,
		NLD_System_ID:  NLD_System_Id,
		NLD_Segment_ID: NLD_Segment_ID,
		NIDID:          NIDID,
		FailureMode:    failure_mode,
		Source:         source,
		FragilityCurve: paireddata.UncertaintyPairedData{Xvals: xvals, Yvals: yvals},
	}
}
func parseLine(line string) string {
	parts := strings.Split(line, ",")
	if len(parts) < 1 {
		return ""
	} else {
		if len(parts) == 1 {
			return ""
		} else {
			if len(parts) == 2 {
				return parts[1]
			} else { //what if the name has commas in it?
				fullname := ""
				parts = parts[1:]
				for _, p := range parts {
					fullname += p
				}
				return fullname
			}

		}
	}
}
func (fcm Model) Compute(eventSeed int64, realizationSeed int64) (ModelResult, error) {
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
