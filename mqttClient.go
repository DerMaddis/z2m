package main

import (
	"github.com/eclipse/paho.mqtt.golang"
)

func NewMqttClient() (mqtt.Client, error) {
    options := mqtt.NewClientOptions().AddBroker("mqtt://192.168.178.28:1883/")

    client := mqtt.NewClient(options)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }

    return client, nil
}
