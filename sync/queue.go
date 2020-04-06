package sync

import (
	"github.com/google/uuid"
	"github.com/nsqio/go-nsq"
)

// Queue provides a task queue
type Queue struct {
	Name string

	source   *nsq.Consumer
	handlers []nsq.Handler

	inTopic    string
	outTopic   string
	inChannel  chan []byte
	outChannel chan []byte
	producer   *nsq.Producer
}

// Handle adds a nmodifier handler
func (q *Queue) Handle(h nsq.Handler) *Queue {
	q.handlers = append(q.handlers, h)
	return q
}

func (q *Queue) dispatch(m *nsq.Message) error {
	for _, h := range q.handlers {
		err := h.HandleMessage(m)
		if err != nil {
			return err
		}
	}
	q.outChannel <- m.Body
	return nil
}

func (q *Queue) publish(ch chan []byte, topic string) error {
	b := <-ch
	return q.producer.Publish(topic, b)

}

// Push pushes a new message to the inTopic. Useful for piping in data while running
func (q *Queue) Push(b []byte) error {
	return q.producer.Publish(q.inTopic, b)
}

// PushFunc runs a handler function to
func (q *Queue) PushFunc(h nsq.Handler) error {
	var body []byte
	m := nsq.NewMessage(nsq.MessageID(uuid.New()), body)
	go func(m *nsq.Message, h nsq.Handler, ch chan []byte) {
		for {
			if err := h.HandleMessage(m); err != nil {
				return
			}
			ch <- m.Body
		}
	}(m, h, q.inChannel)
	return nil
}

// Run registers all handlers with source
func (q *Queue) Run() error {
	if q.source == nil {
		c := nsq.NewConfig()
		consumer, err := nsq.NewConsumer(q.inTopic, q.Name, c)
		if err != nil {
			return err
		}
		q.source = consumer
	}

	q.source.AddHandler(nsq.HandlerFunc(q.dispatch))
	go q.publish(q.inChannel, q.inTopic)
	go q.publish(q.outChannel, q.outTopic)
	return nil
}

// Pipe connects two queues
func (q Queue) Pipe(q2 Queue) (Queue, error) {
	c := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(q.outTopic, q2.Name, c)
	if err != nil {
		return q2, err
	}
	q2.source = consumer
	return q2, nil
}

// New returns a new Queue
func New(name string) (*Queue, error) {
	c := nsq.NewConfig()
	p, err := nsq.NewProducer("https://127.0.0.1:4150", c) // TODO: Use config for nsqd address
	if err != nil {
		return nil, err
	}
	err = p.Ping()
	if err != nil {
		return nil, err
	}
	outchan := make(chan []byte)
	inchan := make(chan []byte)
	q := Queue{
		Name:       name,
		producer:   p,
		inTopic:    name + "_in",
		outTopic:   name + "_out",
		outChannel: outchan,
		inChannel:  inchan,
	}
	return &q, nil
}
