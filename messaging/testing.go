package messaging

import (
	"time"

	mq "github.com/eclipse/paho.mqtt.golang"
)

// TestingMessenger is a no-op messenger useful for when running tests
type TestingMessenger struct {
	client   mq.Client
	Action   string
	Topic    []string
	Message  []byte
	Qos      int
	Persist  bool
	Callback func(Message)
}

func NewTestingMessenger(client mq.Client) PublishSubscriber {
	return &TestingMessenger{client: client}
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

// TestingMQTTToken can be used in place of an mq.Token. It is meant to be
// used in tests
type TestingMQTTToken struct {
	mq.Token
}

func (t *TestingMQTTToken) Wait() bool                     { return true }
func (t *TestingMQTTToken) WaitTimeout(time.Duration) bool { return true }
func (t *TestingMQTTToken) Error() error                   { return nil }

type TestingMQTTClient struct {
	ConnectionState bool
}

// TestingMQTTClient can be used in place of an mq.Client. It is meant to be
// used in tests
func (m *TestingMQTTClient) IsConnected() bool      { return m.ConnectionState }
func (m *TestingMQTTClient) IsConnectionOpen() bool { return m.ConnectionState }
func (m *TestingMQTTClient) Connect() mq.Token {
	m.ConnectionState = true
	return &TestingMQTTToken{}
}
func (m *TestingMQTTClient) Disconnect(uint) { m.ConnectionState = false }
func (m *TestingMQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) mq.Token {
	return &TestingMQTTToken{}
}
func (m *TestingMQTTClient) Subscribe(topic string, qos byte, callback mq.MessageHandler) mq.Token {
	return &TestingMQTTToken{}
}
func (m *TestingMQTTClient) SubscribeMultiple(filters map[string]byte, callback mq.MessageHandler) mq.Token {
	return &TestingMQTTToken{}
}
func (m *TestingMQTTClient) Unsubscribe(topics ...string) mq.Token {
	return &TestingMQTTToken{}
}
func (m *TestingMQTTClient) AddRoute(topic string, callback mq.MessageHandler) {}
func (m *TestingMQTTClient) OptionsReader() mq.ClientOptionsReader {
	return mq.ClientOptionsReader{}
}
