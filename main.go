package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/henrygeorgist/fragilitycurveplugin/model"
	"github.com/usace/wat-api/utils"
)

func main() {
	fmt.Println("fragility curves!")
	var payload string
	flag.StringVar(&payload, "payload", "", "please specify an input file using `-payload=pathtopayload.yml`")
	flag.Parse()

	if payload == "" {
		fmt.Println("given a blank path...")
		fmt.Println("please specify an input file using `-payload=pathtopayload.yml`")
		return
	}
	//payload := "/data/fragilitycurveplugin/watModelPayload.yml"
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
	val := 0.0
	err = utils.LoadJsonPluginModelFromS3(payloadInstructions.ModelConfigurationPaths[0], fs, &fcm)
	if err != nil {
		fmt.Println("error:", err)
		return
	} else {
		fmt.Println("computing model")
		//fmt.Println(hsm)
		val, err = fcm.Compute(&payloadInstructions, fs)
		if err != nil {
			log.Fatal(err)
		}
	}
	//}
	message := fmt.Sprintf("Fragility Curve Complete %v", val)
	fmt.Println("sending message: " + message)
	queueURL := fmt.Sprintf("%v/queue/events", queue.Endpoint)
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
	return
}
