package remote

import (
	"log"

	"cloud.google.com/go/firestore"
	"github.com/DerMaddis/z2m/config"
	"github.com/DerMaddis/z2m/util"
)

func (r Remote) toggle() {
	currentState := r.Device.State.State
	if currentState == "ON" {
		r.Device.Update([]firestore.Update{
			{Path: "state", Value: "OFF"},
		})
	} else if currentState == "OFF" {
		r.Device.Update([]firestore.Update{
			{Path: "state", Value: "ON"},
		})
	}
}

func (r *Remote) nextDevice() {
	r.DeviceIdx = (r.DeviceIdx + 1) % len(r.Devices)
	r.Device = r.Devices[r.DeviceIdx]
	r.Device.SelectedAnimation()
}

func (r *Remote) nextMode() {
    log.Println("current mode", r.Mode)
	if r.Mode == config.BrightnessMode {
		r.Mode = config.ColorMode
		r.Device.ColorModeAnimation()
	} else if r.Mode == config.ColorMode {
		r.Mode = config.BrightnessMode
		r.Device.BrightnessModeAnimation()
	}
    log.Println("mode now", r.Mode)
}

func (r *Remote) valueDown() {
	switch r.Mode {
	case config.BrightnessMode:
		r.brightnessChange(false)
	case config.ColorMode:
		r.colorChange(false)
	}
}

func (r *Remote) valueUp() {
	switch r.Mode {
	case config.BrightnessMode:
		r.brightnessChange(true)
	case config.ColorMode:
		r.colorChange(true)
	}
}

func (r *Remote) brightnessChange(up bool) {
	var sign float32
	if up {
		sign = 1.0
	} else {
		sign = -1.0
	}

	currentBrightness := r.Device.State.Brightness

	newBrightness := util.Clamp(currentBrightness+(config.BrightnessStepSize*sign), 0.0, 255.0)
	r.Device.Update([]firestore.Update{
		{Path: "brightness", Value: newBrightness},
	})
}

func (r *Remote) colorChange(up bool) {
	var sign int
	if up {
		sign = 1
	} else {
		sign = -1
	}

	currentColor := r.Device.State.Color

	var newColor string
	if i, err := util.IndexOf(config.Colors[:], currentColor); err == nil {
		newIndex := (i + sign) % len(config.Colors)
		if newIndex == -1 {
			newIndex = len(config.Colors) - 1
		}
		newColor = config.Colors[newIndex]
	} else {
		newColor = config.Colors[0]
	}

	log.Println(newColor)

	r.Device.Update([]firestore.Update{
		{Path: "color", Value: newColor},
	})
}
