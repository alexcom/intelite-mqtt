package mqtt

import (
	"encoding/json"
	"github.com/alexcom/intelite-mqtt/core"
	"github.com/alexcom/intelite-mqtt/intelite"
	"github.com/alexcom/intelite-mqtt/tasmota"
	framework "github.com/eclipse/paho.mqtt.golang"
	"log"
)

type Handler struct{}

func (Handler) Handle(client framework.Client, msg framework.Message) {
	log.Println("[INFO]", "INCOMING", string(msg.Payload()))
	var update core.LampUpdate
	err := json.Unmarshal(msg.Payload(), &update)
	if err != nil {
		log.Println("[ERROR]", "parsing incoming message:", err.Error())
		return
	}
	s := intelite.NewSignal()
	switch update.Mode {
	case core.LampModeOff:
		s.Mode(intelite.ModeOff)
	case core.LampModeOn:
		s.Mode(intelite.ModeNormal)
		s.Brightness(update.Brightness)
		s.Color(update.Color)
	case core.LampModeNight:
		s.Mode(intelite.ModeNight)
		s.NightBrightness(update.Brightness)
	case core.LampModeMax:
		s.Mode(intelite.ModeNormal)
		s.MaxMode(true)
	}
	resp := tasmota.Payload(s.GenPattern())
	log.Println("[INFO]", "OUTGOING", resp)
	client.Publish(update.RespondTo, 0, false, resp)

	msg.Ack()
}
