package stateManager

import (
	"context"
	"log"

	"dermaddis.de/z2m/device"
	"dermaddis.de/z2m/firestoreState"
	"dermaddis.de/z2m/remote"
	"dermaddis.de/z2m/util"

	"cloud.google.com/go/firestore"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type StateManager struct {
	Remote  *remote.Remote
}

func New(deviceNames []string, mqttClient *mqtt.Client, collection *firestore.CollectionRef) *StateManager {
	bgCtx := context.Background()
	devices := []*device.Device{}

	for _, name := range deviceNames {
		name := name
		var state firestoreState.FirestoreState
		snapshot, err := collection.Doc(name).Get(bgCtx)
		if err != nil {
			log.Fatalln(err)
			continue
		}

		err = snapshot.DataTo(&state)
		if err != nil {
			log.Fatalln(err)
			continue
		}
		devices = append(devices, device.New(name, state, mqttClient, *collection))
	}

	return &StateManager{
		remote.New(devices),
	}
}

func (m *StateManager) StateUpdate(deviceName string, state firestoreState.FirestoreState) {
	log.Println(deviceName, state)
	searchFunc := func(device *device.Device) bool {
		return device.Name == deviceName
	}

	device, err := util.Find(m.Remote.Devices, searchFunc)
	if err != nil {
		return
	}
	device.State = state
	device.SendMqtt()
}

/*
toggle (center)
brightness_up_click (top)
brightness_down_click (bottom)
arrow_left_click (left)
arrow_right_click (right)
*/
func (m StateManager) SwitchPress(action string) {
	log.Println("SwitchPress", action)
	switch action {
	case "toggle":
		m.toggle()
	case "brightness_up_click":
		m.nextDevice()
    case "arrow_left_click":
        m.valueDown()
    case "arrow_right_click":
        m.valueUp()
    case "brightness_down_click":
        m.nextMode()
    }
}

func (m StateManager) toggle() {
    m.Remote.Toggle()
}

func (m StateManager) nextDevice() {
    m.Remote.NextDevice()
}

func (m StateManager) nextMode() {
    m.Remote.NextMode()
}

func (m StateManager) valueDown() {
    m.Remote.ValueDown()
}

func (m StateManager) valueUp() {
    m.Remote.ValueUp()
}
