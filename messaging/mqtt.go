package messaging

import (
	"fmt"
	mq "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

type mqttMessenger struct {
	client mq.Client
}

type Handler struct {
	Ann           chan Message
	Leave         chan Message
	AnnounceTopic string
	LeaveTopic    string
	DiscoverTopic string
	DiscoverDelay time.Duration
	DiscoverStart chan bool
}

// RetryWithBackoff will retry the operation for the amount of attempts. The
// backoff time gets multiplied by the attempt to create an exponential backoff.
//
// Returns an error if the operation still failed after the specified amount of
// attempts have been executed.
func RetryWithBackoff(attempts int, backoff time.Duration, callback func() error) (err error) {
	for i := 1; ; i++ {
		err = callback()
		if err == nil {
			log.Print("Operation succeed at attempt: ", i)
			return
		}

		if i >= attempts {
			break
		}
		backoff = time.Duration(i) * backoff
		log.Printf("Operation failed with error: %s. Going to reattempt in %d seconds", err, backoff/time.Second)
		time.Sleep(backoff)
	}
	return fmt.Errorf("Operation failed after %d attempts, last error: %s", attempts, err)
}

// onConnect gets executed when we've established a connection with the MQTT
// broker, regardless of if this was our first attempt or after a reconnect.
func (h *Handler) OnConnect(c mq.Client) {
	log.Print("Connected to MQTT broker")

	if h.Ann != nil && h.AnnounceTopic != "" {
		log.Print("Attempting to subscribe to announce topic")
		err := RetryWithBackoff(5, 2*time.Second, func() error {
			token := c.Subscribe(h.AnnounceTopic, 1, func(client mq.Client, msg mq.Message) {
				h.Ann <- msg
			})
			token.Wait()
			return token.Error()
		})
		if err != nil {
			log.Fatal("Could not subscribe to announce topic")
		}
		log.Print("Subscribed to announce topic")
	}

	if h.Leave != nil && h.LeaveTopic != "" {
		log.Print("Attempting to subscribe to leave topic")
		err := RetryWithBackoff(5, 2*time.Second, func() error {
			token := c.Subscribe(h.LeaveTopic, 1, func(client mq.Client, msg mq.Message) {
				h.Leave <- msg
			})
			token.Wait()
			return token.Error()
		})
		if err != nil {
			log.Fatal("Could not subscribe to leave topic")
		}
		log.Print("Subscribed to leave topic")
	}

	if h.DiscoverDelay > 0 {
		<-time.After(h.DiscoverDelay)
		h.DiscoverDelay = 0
	}
	if h.DiscoverStart != nil {
		h.DiscoverStart <- true
		h.DiscoverStart = nil
	}

	if h.DiscoverTopic != "" {
		log.Print("Attempting to publish to discover topic")
		err := RetryWithBackoff(5, 2*time.Second, func() error {
			token := c.Publish(h.DiscoverTopic, 1, true, "1")
			token.Wait()
			return token.Error()
		})
		if err != nil {
			log.Fatal("Could not publish to discover topic")
		}
		log.Print("Initiated discovery")
	}
}

// onConnectionLost gets triggered whenver we unexpectedly lose connection with
// the MQTT broker.
func (h *Handler) OnConnectionLost(c mq.Client, e error) {
	log.Print("Unexpectedly lost connection to MQTT broker, attempting to reconnect")
}

// NewMQTTMessenger returns a PublishSubscriber.
//
// It expects to be given something that looks like an MQTT Client and
// a channel on which it will publish any messages from topics to which
// we have subscribed.
//
// It allows for publishing messages to a topic on an MQTT broker, to
// subscribe to messages published to topics and to unsubscribe from topic.
func NewMQTTMessenger(client mq.Client) PublishSubscriber {
	return &mqttMessenger{
		client: client,
	}
}

// Publish publishes a msg on the specified topic. qos represents the MQTT QoS
// level and retain informs the broker that it needs to persist this message so
// that when a new client subscribes to the topic we published on they will
// automatically get that message.
func (m *mqttMessenger) Publish(topic string, msg []byte, qos int, retain bool) {
	m.client.Publish(topic, byte(qos), retain, msg)
}

// Subscribe subscribes to the specified topic with a certain qos. The topic
// and message are then passed into this messenger's recv channel and can be
// read from by any interested consumer.
func (m *mqttMessenger) Subscribe(topic string, qos int, callback func(Message)) {
	m.client.Subscribe(topic, byte(qos), func(c mq.Client, msg mq.Message) {
		callback(msg)
	})
}

// Unsubscribe unsubscribes from one or multiple topics.
func (m *mqttMessenger) Unsubscribe(topics ...string) {
	m.client.Unsubscribe(topics...)
}
