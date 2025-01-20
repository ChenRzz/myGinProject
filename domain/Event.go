package domain

type Event struct {
	Name string
	Body []byte
}

type EventProducer interface {
	Publish(event Event) error
}
