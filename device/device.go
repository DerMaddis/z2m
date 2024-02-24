package device

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/DerMaddis/z2m/config"
	"github.com/DerMaddis/z2m/firestoreState"
	"github.com/DerMaddis/z2m/util"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"cloud.google.com/go/firestore"
)

type Device struct {
	Name    string
	State   firestoreState.FirestoreState
	Publish func(data string)
	Update  func(updates []firestore.Update)
}

func New(name string, state firestoreState.FirestoreState, mqttClient *mqtt.Client, collection firestore.CollectionRef) Device {
	bgCtx := context.Background()
	return Device{
		Name:  name,
		State: state,
		Publish: func(data string) {
			(*mqttClient).Publish(fmt.Sprintf("z2m/%s/set", name), 0, false, data)
		},
		Update: func(updates []firestore.Update) {
			_, err := collection.Doc(name).Update(bgCtx, updates)
			if err != nil {
				log.Fatalln(err)
			}
		},
	}
}

func (d Device) SendMqtt() {
	bytes, err := json.Marshal(d.State)
	if err != nil {
		log.Fatalln(err)
	}
	d.Publish(string(bytes))
}

func (d Device) SelectedAnimation() {
	currentState := d.State.State
	var blinkState string

	if currentState == "ON" {
		blinkState = "OFF"
	} else {
		blinkState = "ON"
	}

	var blinkDuration int
	// the light is a bit slower when turning on, so we up the blink duration
	if blinkState == "ON" {
		blinkDuration = 500
	} else {
		blinkDuration = 250
	}

	d.Publish(fmt.Sprintf(`{"state": "%s", "transition": 0.25}`, blinkState))
	time.Sleep(time.Duration(blinkDuration) * time.Millisecond)
	d.Publish(fmt.Sprintf(`{"state": "%s", "transition": 0.25}`, currentState))
}

func (d Device) BrightnessModeAnimation() {
	currentBrightness := d.State.Brightness
	var blinkBrightness float32

	if currentBrightness >= 128 {
		blinkBrightness = currentBrightness - config.BrightnessStepSize * 2
	} else {
		blinkBrightness = currentBrightness + config.BrightnessStepSize * 2
	}

	d.Publish(fmt.Sprintf(`{"brightness": "%.2f", "transition": 0.25}`, blinkBrightness))
	time.Sleep(time.Duration(350) * time.Millisecond)
	d.Publish(fmt.Sprintf(`{"brightness": "%.2f", "transition": 0.25}`, currentBrightness))
}

func (d Device) ColorModeAnimation() {
	currentColor := d.State.Color

	var blinkColor string
	if i, err := util.IndexOf(config.Colors[:], currentColor); err == nil {
		blinkColor = config.Colors[(i+1)%len(config.Colors)]
	} else {
		blinkColor = config.Colors[0]
	}

	d.Publish(fmt.Sprintf(`{"color": "%s", "transition": 0.25}`, blinkColor))
	time.Sleep(time.Duration(350) * time.Millisecond)
	d.Publish(fmt.Sprintf(`{"color": "%s", "transition": 0.25}`, currentColor))
}
