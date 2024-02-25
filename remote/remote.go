package remote

import (
	"log"

	"cloud.google.com/go/firestore"
	"github.com/DerMaddis/z2m/config"
	"github.com/DerMaddis/z2m/device"
	"github.com/DerMaddis/z2m/firestoreState"
	"github.com/DerMaddis/z2m/util"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Remote struct {
	DeviceIdx int
	Device    *device.Device
	Mode      int
	Devices   []*device.Device
}

func New(deviceNames []string, mqttClient *mqtt.Client, collection *firestore.CollectionRef) Remote {
	devices := []*device.Device{}

	for _, name := range deviceNames {
		name := name
		state, err := firestoreState.GetDeviceState(name, collection)
		if err != nil {
			log.Println(err)
			continue
		}
        newDevice := device.New(name, state, mqttClient, *collection)
		devices = append(devices, &newDevice)
	}

	return Remote{
		Devices:   devices,
		DeviceIdx: 0,
		Device:    devices[0],
		Mode:      config.BrightnessMode,
	}
}

func (r *Remote) StateUpdate(deviceName string, state firestoreState.FirestoreState) {
	log.Println(deviceName, state)
	searchFunc := func(device *device.Device) bool {
		return device.Name == deviceName
	}

	device, err := util.Find(r.Devices, searchFunc)
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
func (r *Remote) SwitchPress(action string) {
	log.Println("SwitchPress", action)
	switch action {
	case "toggle":
		r.toggle()
	case "brightness_up_click":
		r.nextDevice()
	case "arrow_left_click":
		r.valueDown()
	case "arrow_right_click":
		r.valueUp()
	case "brightness_down_click":
		r.nextMode()
	}
}

