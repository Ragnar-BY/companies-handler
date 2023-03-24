package broker

// Broker is mock implementation of broker
type Broker struct{}

func NewBroker() *Broker {
	return &Broker{}
}
func (b *Broker) SendEvent(topic string, message []byte) error {
	return nil
}
