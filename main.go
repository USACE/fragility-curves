package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/usace/cc-go-sdk"
	"github.com/usace/cc-go-sdk/plugin"
	"github.com/usace/fragility-curves/fragilitycurve"
)

func main() {
	fmt.Println("Fragility Curves!")
	pm, err := cc.InitPluginManager()
	if err != nil {
		log.Fatalf("Unable to initialize the plugin manager: %s\n", err)
	}
	payload := pm.Payload
	err = computePayload(payload, pm)
	if err != nil {
		pm.Logger.Error(err.Error())
	}

}

func computePayload(payload cc.Payload, pm *cc.PluginManager) error {

	if len(payload.Outputs) != 1 {
		err := errors.New("more than one output was defined")
		return err
	}

	var fcm fragilitycurve.Model
	modelReader, err := pm.GetReader(cc.DataSourceOpInput{DataSourceName: "fragilitycurve", PathKey: "default"})
	if err != nil {
		return err
	}
	defer modelReader.Close()
	err = json.NewDecoder(modelReader).Decode(&fcm)
	if err != nil {
		return err
	}
	var seedSet plugin.SeedSet
	var ec plugin.EventConfiguration
	eventConfigurationReader, err := pm.GetReader(cc.DataSourceOpInput{DataSourceName: "seeds", PathKey: "default"})
	if err != nil {
		return err
	}
	defer eventConfigurationReader.Close()
	err = json.NewDecoder(eventConfigurationReader).Decode(&ec)
	if err != nil {
		return err
	}

	seedSetName := "fragilitycurveplugin"
	seedSet, seedsFound := ec.Seeds[seedSetName]
	if !seedsFound {
		return fmt.Errorf("no seeds found by name of %v", seedSetName)
	}
	modelResult, err := fcm.Compute(seedSet.EventSeed, seedSet.RealizationSeed)
	if err != nil {
		return err
	}
	data, err := json.Marshal(modelResult)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	input := cc.PutOpInput{
		SrcReader:         bytes.NewReader(data),
		DataSourceOpInput: cc.DataSourceOpInput{DataSourceName: payload.Outputs[0].Name, PathKey: "default"},
	}
	_, err = pm.Put(input)
	if err != nil {
		return err
	}
	pm.Logger.Info("payload compute complete")
	return nil
}
