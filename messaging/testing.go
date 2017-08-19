package messaging

import (
	mq "github.com/eclipse/paho.mqtt.golang"
)

// TestingMessenger is a no-op messenger useful for when running tests
type TestingMessenger struct {
	Client   mq.Client
	Action   string
	Topic    []string
	Message  []byte
	Qos      int
	Persist  bool
	Callback func(Message)
}

func (tm *TestingMessenger) Publish(topic string, message []byte, qos int, persist bool) {
	tm.Action = "publish"
	tm.Topic = []string{topic}
	tm.Message = message
	tm.Qos = qos
	tm.Persist = persist
}
func (tm *TestingMessenger) Subscribe(topic string, qos int, callback func(Message)) {
	tm.Action = "subscribe"
	tm.Topic = []string{topic}
	tm.Qos = qos
	tm.Callback = callback
}
func (tm *TestingMessenger) Unsubscribe(topics ...string) {
	tm.Action = "unsubscribe"
	tm.Topic = topics
}
