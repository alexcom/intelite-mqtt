package main

import (
	"github.com/alexcom/intelite-mqtt/mqtt"
	framework "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type config struct {
	Topic     string `yaml:"topic"`
	BrokerURL string `yaml:"broker_url"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	ClientID  string `yaml:"client_id"`
}

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("[ERROR]", "expected config file path as the only argument")
	} else if len(os.Args) > 2 {
		log.Fatalln("[ERROR]", "expected single argument")
	}
	var file *os.File
	var err error
	if file, err = os.Open(os.Args[1]); os.IsNotExist(err) {
		log.Fatalf("%s file '%s' does not exist", "[ERROR]", os.Args[1])
	}
	cfg := config{}
	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		log.Fatalln("[ERROR]", "error reading config YAML file")
	}

	opts := framework.NewClientOptions().
		AddBroker(cfg.BrokerURL).
		SetClientID(cfg.ClientID)
	client := framework.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalln("[ERROR]", "failed to connect to MQTT broker:", token.Error())
	}
	log.Println("[INFO]", "connect successful")

	topic := cfg.Topic
	if token := client.Subscribe(topic, 0, mqtt.Handler{}.Handle); token.Wait() && token.Error() != nil {
		log.Fatalln("[ERROR]", "failed to subscribe to topic", token.Error())
	}
	log.Println("[INFO]", "subscribe successful")

	<-waitForShutdown()

	if token := client.Unsubscribe(topic); token.Wait() || token.Error() != nil {
		log.Fatalln("[ERROR]", "failed to unsubscribe from topic", token.Error())
	}
	log.Println("[INFO]", "unsubscribe successful")
	client.Disconnect(250)
}

func waitForShutdown() <-chan os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGKILL)
	return ch
}
