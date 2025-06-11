package event

type Publisher interface {
	Publish(queue string, body []byte) error
}
