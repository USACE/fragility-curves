package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/usace/fragility-curves/fragilitycurve"
	"github.com/usace/wat-go"
)

var pluginName string = "fragilitycurveplugin"

func main() {
	fmt.Println("fragility curves!")
	ws, err := wat.NewS3WatStore()
	if err != nil {
		wat.Log(wat.Message{
			Status:    wat.FAILED,
			Progress:  0,
			Level:     wat.ERROR,
			Message:   err.Error(),
			Sender:    pluginName,
			PayloadId: "unknown",
		})
	}
	payload, err := ws.GetPayload()
	if err != nil {
		wat.Log(wat.Message{
			Status:    wat.FAILED,
			Progress:  0,
			Level:     wat.ERROR,
			Message:   err.Error(),
			Sender:    pluginName,
			PayloadId: "unknown",
		})
		return
	}
	err = computePayload(payload, ws)
	if err != nil {
		wat.Log(wat.Message{
			Status:    wat.FAILED,
			Progress:  0,
			Level:     wat.ERROR,
			Message:   err.Error(),
			Sender:    pluginName,
			PayloadId: "unknown",
		})
		return
	}
}
func computePayload(payload wat.Payload, ws wat.WatStore) error {

	if len(payload.Outputs) != 1 {
		err := errors.New("more than one output was defined")
		wat.Log(wat.Message{
			Status:    wat.FAILED,
			Progress:  0,
			Level:     wat.ERROR,
			Message:   err.Error(),
			Sender:    pluginName,
			PayloadId: "unknown",
		})
		return err
	}
	if len(payload.Inputs) != 2 {
		err := errors.New("more than two inputs were defined")
		wat.Log(wat.Message{
			Status:    wat.FAILED,
			Progress:  0,
			Level:     wat.ERROR,
			Message:   err.Error(),
			Sender:    pluginName,
			PayloadId: "unknown",
		})
		return err
	}
	var modelResourceInfo wat.DataSource
	//var eventConfigurationResourceInfo wat.DataSource
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
			//eventConfigurationResourceInfo = input
			foundEventConfig = true
		}
	}
	if !foundModel {
		err := fmt.Errorf("could not find %s.json", modelName)
		wat.Log(wat.Message{
			Status:    wat.FAILED,
			Progress:  0,
			Level:     wat.ERROR,
			Message:   err.Error(),
			Sender:    pluginName,
			PayloadId: "unknown",
		})
		return err
	}
	if !foundEventConfig {
		err := fmt.Errorf("could not find eventconfiguration.json")
		wat.Log(wat.Message{
			Status:    wat.FAILED,
			Progress:  0,
			Level:     wat.ERROR,
			Message:   err.Error(),
			Sender:    pluginName,
			PayloadId: "unknown",
		})
		return err
	}

	modelBytes, err := ws.GetObject(modelResourceInfo.Name) //GetObject(modelResourceInfo)
	if err != nil {
		wat.Log(wat.Message{
			Status:    wat.FAILED,
			Progress:  0,
			Level:     wat.ERROR,
			Message:   err.Error(),
			Sender:    pluginName,
			PayloadId: "unknown",
		})
		return err
	}
	var fcm fragilitycurve.Model
	err = json.Unmarshal(modelBytes, &fcm)
	if err != nil {
		wat.Log(wat.Message{
			Status:    wat.FAILED,
			Progress:  0,
			Level:     wat.ERROR,
			Message:   err.Error(),
			Sender:    pluginName,
			PayloadId: "unknown",
		})
		return err
	}
	/*eventConfiguration, err := ws.GetObject(eventConfigurationResourceInfo.Name)
	if err != nil {
		wat.Log(wat.Message{
			Status:    wat.FAILED,
			Progress:  0,
			Level:     wat.ERROR,
			Message:   err.Error(),
			Sender:    pluginName,
			PayloadId: "unknown",
		})
		return err
	}
	//then we need to get the specific set of seeds.
	//currently event configuration is not defined in the go-wat package. seedset is not a member of anything specific to cc.
	var seedSet plugin.SeedSet
	var ec plugin.EventConfiguration
	err = json.Unmarshal(eventConfiguration, ec) //this is not right.
	seedSet, err = ec.SeedSet(seedSetName)
	if err != nil {
		wat.Log(wat.Message{
			Status:    wat.FAILED,
			Progress:  0,
			Level:     wat.ERROR,
			Message:   err.Error(),
			Sender:    pluginName,
			PayloadId: "unknown",
		})
		return err
	}
	modelResult, err := fcm.Compute(seedSet)
	if err != nil {
		wat.Log(wat.Message{
			Status:    wat.FAILED,
			Progress:  0,
			Level:     wat.ERROR,
			Message:   err.Error(),
			Sender:    pluginName,
			PayloadId: "unknown",
		})
		return err
	}
	bytes, err := json.Marshal(modelResult)
	if err != nil {
		wat.Log(wat.Message{
			Status:    wat.FAILED,
			Progress:  0,
			Level:     wat.ERROR,
			Message:   err.Error(),
			Sender:    pluginName,
			PayloadId: "unknown",
		})
		return err
	}
	//plugin.UpLoadFile(payload.Outputs[0].ResourceInfo, bytes)
	ioutil.WriteFile(fmt.Sprintf("%v/%v", ws.RootPath(), "temp.out"), bytes, fs.ModeAppend) //dont have access to local root path
	err = ws.PushObject("temp.out")
	if err != nil {
		wat.Log(wat.Message{
			Status:    wat.FAILED,
			Progress:  0,
			Level:     wat.ERROR,
			Message:   err.Error(),
			Sender:    pluginName,
			PayloadId: "unknown",
		})
		return err
	}
	plugin.Log(plugin.Message{
		Status:    plugin.SUCCEEDED,
		Progress:  100,
		Level:     plugin.INFO,
		Message:   "fragility curves complete",
		Sender:    "fragilitycurve",
		PayloadId: "unknown",
	})
	*/
	return nil
}
