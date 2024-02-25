package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/DerMaddis/z2m/firestoreState"
	"github.com/DerMaddis/z2m/remote"
	"google.golang.org/api/iterator"

	"github.com/eclipse/paho.mqtt.golang"
)

type MqttSwitchMessage struct {
	Action      string
	Battery     int
	linkquality int
}

func main() {
	bgCtx := context.Background()

    fClient, err := NewFirestoreClient()
    if err != nil {
        panic(err)
    }
    log.Println("firestore opened")

	mqttClient, err := NewMqttClient()
	if err != nil {
		panic(err)
	}
	defer mqttClient.Disconnect(250)
	log.Println("mqtt opened")

	deviceNames := []string{"strip01", "light01"}
	remote := remote.New(deviceNames, &mqttClient, fClient.Collection("state"))

	mqttClient.Subscribe("z2m/switch01", 0, func(c mqtt.Client, m mqtt.Message) {
		var message MqttSwitchMessage
		err := json.Unmarshal(m.Payload(), &message)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(message)

		remote.SwitchPress(message.Action)
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
			remote.StateUpdate(change.Doc.Ref.ID, state)
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
