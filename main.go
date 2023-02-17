package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/usace/cc-go-sdk"
	"github.com/usace/cc-go-sdk/plugin"
	"github.com/usace/fragility-curves/fragilitycurve"
)

func main() {
	fmt.Println("event generator!")
	pm, err := cc.InitPluginManager()
	if err != nil {
		log.Fatalf("Unable to initialize the plugin manager: %s\n", err)
	}
	payload := pm.GetPayload()
	err = computePayload(payload, pm)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
	}
}

func computePayload(payload cc.Payload, pm *cc.PluginManager) error {

	if len(payload.Outputs) != 1 {
		err := errors.New("more than one output was defined")
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	if len(payload.Inputs) != 2 {
		err := errors.New("more than two inputs were defined")
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	modelResourceInfo, err := pm.GetInputDataSource("fragilitycurve")
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	eventConfigurationResourceInfo, err := pm.GetInputDataSource("seeds")
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}

	var fcm fragilitycurve.Model
	modelReader, err := pm.FileReader(modelResourceInfo, 0)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	defer modelReader.Close()
	err = json.NewDecoder(modelReader).Decode(&fcm)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}

	var seedSet plugin.SeedSet
	var ec plugin.EventConfiguration
	eventConfigurationReader, err := pm.FileReader(eventConfigurationResourceInfo, 0)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	defer eventConfigurationReader.Close()
	err = json.NewDecoder(eventConfigurationReader).Decode(&ec)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}

	seedSetName := "fragilitycurveplugin" //not sure this is right
	seedSet, seedsFound := ec.Seeds[seedSetName]
	if !seedsFound {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      fmt.Errorf("no seeds found by name of %v", seedSetName).Error(),
		})
		return err
	}
	modelResult, err := fcm.Compute(seedSet.EventSeed, seedSet.RealizationSeed)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	bytes, err := json.Marshal(modelResult)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	fmt.Println(string(bytes))
	err = pm.PutFile(bytes, payload.Outputs[0], 0)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	pm.ReportProgress(cc.StatusReport{
		Status:   cc.SUCCEEDED,
		Progress: 100,
	})
	return nil
}
