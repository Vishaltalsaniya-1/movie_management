package producer

type ProducerInterface interface {
	Initialize() error
	Publish(message []byte) error
}
   