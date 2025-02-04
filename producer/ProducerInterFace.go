package producer

type ProducerInterface interface {
	Initialize() error
	Publish(message []byte, taskname string) error
}
   