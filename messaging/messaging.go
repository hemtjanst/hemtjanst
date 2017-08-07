package messaging

// Publisher publishes messages on a transport
type Publisher interface {
	Publish(destination string, message []byte, qos int, persist bool)
}

// Subscriber receives messages from a transport
type Subscriber interface {
	Subscribe(source string, qos int, callback func(Message))
	Unsubscribe(sources ...string)
}

// PublishSubscriber can both publish and receives messages from a transport
type PublishSubscriber interface {
	Publisher
	Subscriber
}

type Message interface {
	Topic() string
	Payload() []byte
}
