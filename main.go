package main

import (
	"fmt"

	"github.com/henrygeorgist/fragilitycurveplugin/model"
)

func main() {
	fmt.Println("fragility curves!")
	payload := "/data/fragilitycurveplugin/watModelPayload.yml"
	fmt.Println("initializing filestore")
	fs, err := model.InitStore()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("initializing Redis")
	_, err = model.InitRedis()
	if err != nil {
		fmt.Println(err)
		return
	}
	payloadInstructions, err := model.LoadPayloadFromS3(payload, fs)
	if err != nil {
		fmt.Println("not successful", err)
		return
	}
	// verify this is the right plugin
	if payloadInstructions.TargetPlugin != "fragilitycurveplugin" {
		fmt.Println("error", "expecting", "fragilitycurveplugin", "got", payloadInstructions.TargetPlugin)
		return
	}
	//load the model data into memory.
	fcm := model.FragilityCurveModel{}
	err = model.NewPluginModelFromS3(payloadInstructions.ModelConfigurationPaths[0], fs, &fcm)
	if err != nil {
		fmt.Println("error:", err)
		return
	} else {
		fmt.Println("computing model")
		//fmt.Println(hsm)
		fcm.Compute(&payloadInstructions, fs)

	}
	//}
	fmt.Println("Made it to the end.....")
}
