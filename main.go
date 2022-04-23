package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/henrygeorgist/fragilitycurveplugin/model"
	"github.com/usace/wat-api/utils"
)

func main() {
	fmt.Println("fragility curves!")
	payload := "/data/fragilitycurveplugin/watModelPayload.yml"
	fmt.Println("initializing filestore")
	loader, err := utils.InitLoader("")
	if err != nil {
		log.Fatal(err)
		return
	}
	fs, err := loader.InitStore()
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = loader.InitRedis()
	if err != nil {
		log.Fatal(err)
		return
	}
	queue, err := loader.InitQueue()
	if err != nil {
		log.Fatal(err)
		return
	}
	payloadInstructions, err := utils.LoadModelPayloadFromS3(payload, fs)
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
	err = utils.LoadJsonPluginModelFromS3(payloadInstructions.ModelConfigurationPaths[0], fs, &fcm)
	if err != nil {
		fmt.Println("error:", err)
		return
	} else {
		fmt.Println("computing model")
		//fmt.Println(hsm)
		err = fcm.Compute(&payloadInstructions, fs)
		if err != nil {
			log.Fatal(err)
		}
	}
	//}
	message := "Fragility Curve Complete"
	fmt.Println("sending message: " + message)
	queueURL := fmt.Sprintf("%v/queue/messages", queue.Endpoint)
	fmt.Println("sending message to:", queueURL)
	_, err = queue.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(1),
		MessageBody:  aws.String(message),
		QueueUrl:     &queueURL,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(message)
}
