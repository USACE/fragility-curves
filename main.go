package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/usace/fragility-curves/fragilitycurve"
	"github.com/usace/wat-go"
	"github.com/usace/wat-go/plugin"
)

func main() {
	fmt.Println("fragility curves!")
	pm, err := wat.InitPluginManager()
	if err != nil {
		pm.LogMessage(wat.Message{
			Message: err.Error(),
		})
	}
	payload, err := pm.GetPayload()
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return
	}
	err = computePayload(payload, pm)
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return
	}
}
func computePayload(payload wat.Payload, pm wat.PluginManager) error {

	if len(payload.Outputs) != 1 {
		err := errors.New("more than one output was defined")
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	if len(payload.Inputs) != 2 {
		err := errors.New("more than two inputs were defined")
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	var modelResourceInfo wat.DataSource
	var eventConfigurationResourceInfo wat.DataSource
	foundModel := false
	foundEventConfig := false
	//seedSetName := ""
	tmpName, ok := payload.Attributes["ModelName"]
	modelName := ""
	if ok {
		modelName = tmpName.(string)
		//seedSetName = modelName
	}
	for _, input := range payload.Inputs {
		if strings.Contains(input.Name, modelName+".json") {
			modelResourceInfo = input
			foundModel = true
		}
		if strings.Contains(input.Name, "eventconfiguration.json") { //not sure this is how we will make this work.
			eventConfigurationResourceInfo = input
			foundEventConfig = true
		}
	}
	if !foundModel {
		err := fmt.Errorf("could not find %s.json", modelName)
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	if !foundEventConfig {
		err := fmt.Errorf("could not find eventconfiguration.json")
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}

	modelBytes, err := pm.GetObject(modelResourceInfo)
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	var fcm fragilitycurve.Model
	err = json.Unmarshal(modelBytes, &fcm)
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	eventConfiguration, err := pm.GetObject(eventConfigurationResourceInfo)
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	//then we need to get the specific set of seeds.
	var seedSet plugin.SeedSet
	var ec plugin.EventConfiguration
	err = json.Unmarshal(eventConfiguration, &ec)
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	seedSetName := "fragilitycurveplugin" //not sure this is right
	seedSet, seedsFound := ec.Seeds[seedSetName]
	if !seedsFound {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      fmt.Errorf("no seeds found by name of %v", seedSetName).Error(),
		})
		return err
	}
	modelResult, err := fcm.Compute(seedSet)
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	bytes, err := json.Marshal(modelResult)
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	fmt.Println(string(bytes))
	err = pm.PutObject(payload.Outputs[0], bytes)
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	pm.ReportProgress(wat.StatusReport{
		Status:   wat.SUCCEEDED,
		Progress: 100,
	})
	return nil
}
