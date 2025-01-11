package typing

type PubSub interface {
	Connect() error
	Publish(topic string, payload []byte) error
	Subscribe(topic string, handler func(payload []byte)) error
}
