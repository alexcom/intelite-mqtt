package main

import (
	"github.com/alexcom/intelite-mqtt/mqtt"
	framework "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	opts := framework.NewClientOptions().AddBroker("tcp://10.0.1.41:1883").SetClientID("intelite-mqtt")
	client := framework.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalln("[ERROR]", "failed to connect to MQTT broker:", token.Error())
	}
	log.Println("[INFO] connect successful")

	const topic = "intelite/bedroom"
	if token := client.Subscribe(topic, 0, mqtt.Handler{}.Handle); token.Wait() && token.Error() != nil {
		log.Fatalln("[ERROR]", "failed to subscribe to topic", token.Error())
	}
	log.Println("[INFO] subscribe successful")

	<-waitForShutdown()

	if token := client.Unsubscribe(topic); token.Wait() || token.Error() != nil {
		log.Fatalln("[ERROR]", "failed to unsubscribe from topic", token.Error())
	}
	log.Println("[INFO] unsubscribe successful")
	client.Disconnect(250)
}

func waitForShutdown() <-chan os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGKILL)
	return ch
}
