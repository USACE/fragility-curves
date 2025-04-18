package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/usace/cc-go-sdk"
	tiledb "github.com/usace/cc-go-sdk/tiledb-store"
	"github.com/usace/fragility-curves/fragilitycurve"
	"github.com/usace/fragility-curves/utils"
)

func main() {
	fmt.Println("Fragility Curves!")
	//register tiledb
	cc.DataStoreTypeRegistry.Register("TILEDB", tiledb.TileDbEventStore{})
	pm, err := cc.InitPluginManager()
	if err != nil {
		log.Fatalf("Unable to initialize the plugin manager: %s\n", err)
	}
	payload := pm.Payload
	for _, a := range payload.Actions {
		switch a.Type {
		case "single-sample":
			err = computeAction(a)
		case "all-samples":
			err = computeAllAction(a)
		}
	}

	if err != nil {
		pm.Logger.Error(err.Error())
	} else {
		pm.Logger.Info("payload compute complete")
	}

}

func computeAction(a cc.Action) error {

	if len(a.Outputs) != 1 {
		err := errors.New("more than one output was defined")
		return err
	}

	var fcm fragilitycurve.Model
	modelReader, err := a.GetReader(cc.DataSourceOpInput{DataSourceName: "fragilitycurve", PathKey: "default"})
	if err != nil {
		return err
	}
	defer modelReader.Close()
	err = json.NewDecoder(modelReader).Decode(&fcm)
	if err != nil {
		return err
	}
	var seedSet utils.SeedSet
	var ec utils.EventConfiguration
	eventConfigurationReader, err := a.GetReader(cc.DataSourceOpInput{DataSourceName: "seeds", PathKey: "default"})
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
		DataSourceOpInput: cc.DataSourceOpInput{DataSourceName: a.Outputs[0].Name, PathKey: "default"},
	}
	_, err = a.Put(input)
	if err != nil {
		return err
	}

	return nil
}
func computeAllAction(a cc.Action) error {
	//
	readSeedsFromTiledb := a.Attributes.GetBooleanOrDefault("seeds_format", false)
	writeSamplesToTiledb := a.Attributes.GetBooleanOrDefault("elevations_format", false)
	if len(a.Outputs) != 1 {
		err := errors.New("more than one output was defined")
		return err
	}

	var fcm fragilitycurve.Model
	modelReader, err := a.GetReader(cc.DataSourceOpInput{DataSourceName: "fragilitycurve", PathKey: "default"})
	if err != nil {
		return err
	}
	defer modelReader.Close()
	err = json.NewDecoder(modelReader).Decode(&fcm)
	if err != nil {
		return err
	}
	//seeds
	seeds := make([]utils.SeedSet, 0)
	if readSeedsFromTiledb {
		seeds, err = utils.ReadSeedsFromTiledb(a.IOManager, "store", "seeds", "fragilitycurveplugin") //improve this to not be hard coded.
		if err != nil {
			return err
		}
	} else {
		//json
		eventConfigurationReader, err := a.GetReader(cc.DataSourceOpInput{DataSourceName: "seeds", PathKey: "default"})
		if err != nil {
			return err
		}
		var ecs []utils.EventConfiguration
		defer eventConfigurationReader.Close()
		err = json.NewDecoder(eventConfigurationReader).Decode(&ecs)
		if err != nil {
			return err
		}
		for _, ec := range ecs {
			seeds = append(seeds, ec.Seeds["fragilitycurveplugin"])
		}
	}

	modelResult, err := fcm.ComputeAll(seeds)
	if err != nil {
		return err
	}
	if writeSamplesToTiledb {
		err = fragilitycurve.WriteFailureElevationsToTiledb(a.IOManager, "store", "failure_elevations", modelResult)
		if err != nil {
			return err
		}
	} else {
		data, err := json.Marshal(modelResult)
		if err != nil {
			return err
		}
		//fmt.Println(string(data))
		input := cc.PutOpInput{
			SrcReader:         bytes.NewReader(data),
			DataSourceOpInput: cc.DataSourceOpInput{DataSourceName: a.Outputs[0].Name, PathKey: "default"},
		}
		_, err = a.Put(input)
		if err != nil {
			return err
		}
	}

	return nil
}
