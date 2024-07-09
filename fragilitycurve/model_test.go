package fragilitycurve

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

func TestModelUnMarshal(t *testing.T) {
	file, err := os.Open("/workspaces/fragilitycurveplugin/configs/fragilitycurve.json")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	body, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fcm := Model{}
	errjson := json.Unmarshal(body, &fcm)
	if errjson != nil {
		fmt.Println(errjson)
		t.Fail()
	}
}
func TestModelCSV(t *testing.T) {
	file, err := os.Open("/workspaces/fragilitycurveplugin/configs/levee1_fragilitycurve.csv")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	body, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fcl := InitFragilityCurveLocation(body)
	fmt.Print(fcl)
}
func TestModelMarshal(t *testing.T) {
	filepath := "/workspaces/fragilitycurveplugin/configs/fc_test.json"
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	defer file.Close()
	systemResponseCurve := createSampleData()
	fcl1 := FragilityCurveLocation{
		Name:           "levee1",
		FragilityCurve: systemResponseCurve,
	}
	locations := make([]FragilityCurveLocation, 1)
	locations[0] = fcl1
	fcm := Model{Locations: locations}
	bytes, errjson := json.Marshal(fcm)
	if errjson != nil {
		fmt.Println(errjson)
		t.Fail()
	}
	_, err = file.Write(bytes)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
func createSampleData() paireddata.UncertaintyPairedData {
	xs := []float64{0.0, .1, .2, .3, .4, .5, .6, .7, .8, .9, 1.0}
	ydists := make([]statistics.ContinuousDistribution, 11)
	ydists[0] = statistics.TriangularDistribution{Min: 98, MostLikely: 99, Max: 100}
	ydists[1] = statistics.TriangularDistribution{Min: 99, MostLikely: 100, Max: 101}
	ydists[2] = statistics.TriangularDistribution{Min: 100, MostLikely: 101, Max: 102}
	ydists[3] = statistics.TriangularDistribution{Min: 101, MostLikely: 102, Max: 103}
	ydists[4] = statistics.TriangularDistribution{Min: 102, MostLikely: 103, Max: 104}
	ydists[5] = statistics.TriangularDistribution{Min: 103, MostLikely: 104, Max: 105}
	ydists[6] = statistics.TriangularDistribution{Min: 104, MostLikely: 105, Max: 106}
	ydists[7] = statistics.TriangularDistribution{Min: 105, MostLikely: 106, Max: 107}
	ydists[8] = statistics.TriangularDistribution{Min: 106, MostLikely: 107, Max: 108}
	ydists[9] = statistics.TriangularDistribution{Min: 107, MostLikely: 108, Max: 109}
	ydists[10] = statistics.TriangularDistribution{Min: 108, MostLikely: 109, Max: 110}

	return paireddata.UncertaintyPairedData{Xvals: xs, Yvals: ydists}
}
func TestSampleFragilityCurve(t *testing.T) {
	file, err := os.Open("/workspaces/fragilitycurveplugin/configs/fragilitycurve.json")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	body, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fcm := Model{}
	errjson := json.Unmarshal(body, &fcm)
	if errjson != nil {
		fmt.Println(errjson)
		t.Fail()
	}
	modelResult, err := fcm.Compute(1234, 1234)
	bytes, err := json.Marshal(modelResult)
	fmt.Println(string(bytes))
}
