package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/DerMaddis/z2m/firestoreState"
	"github.com/DerMaddis/z2m/stateManager"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/eclipse/paho.mqtt.golang"
)

type MqttSwitchMessage struct {
	Action      string
	Battery     int
	linkquality int
}

func modFloor(a, n int) int {
	return ((a % n) + n) % n
}

func main() {
	bgCtx := context.Background()

	credFile := "./cred.json"
	options := option.WithCredentialsFile(credFile)

	app, err := firebase.NewApp(bgCtx, &firebase.Config{}, options)
	if err != nil {
		panic(err)
	}
	log.Println("firebase opened")

	fClient, err := app.Firestore(bgCtx)
	if err != nil {
		panic(err)
	}

	mqttClient, err := NewMqttClient()
	if err != nil {
		panic(err)
	}
	defer mqttClient.Disconnect(250)
	log.Println("mqtt opened")

	deviceNames := []string{"strip01", "light01"}
	sManager := stateManager.New(deviceNames, &mqttClient, fClient.Collection("state"))

	mqttClient.Subscribe("z2m/switch01", 0, func(c mqtt.Client, m mqtt.Message) {
		var message MqttSwitchMessage
		err := json.Unmarshal(m.Payload(), &message)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(message)

		sManager.SwitchPress(message.Action)
	})

	allDocsIter := fClient.Collection("state").Snapshots(bgCtx)
	for {
		querySnap, err := allDocsIter.Next()
		if err != nil {
			if err == iterator.Done {
				continue
			}
			log.Fatalln(err)
		}
		for _, change := range querySnap.Changes {
			state, err := handleDocChange(change)
			if err != nil {
				continue
			}
			sManager.StateUpdate(change.Doc.Ref.ID, state)
		}
	}
}

func handleDocChange(change firestore.DocumentChange) (firestoreState.FirestoreState, error) {
	var state firestoreState.FirestoreState

	if change.Kind != firestore.DocumentModified {
		return state, errors.New("not of kind DocumentModified")
	}

	err := change.Doc.DataTo(&state)
	return state, err
}
