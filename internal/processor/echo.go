package processor

//
//// echo ...
//type echo struct {
//	counter int64
//}

// NewEcho ...
//func NewEcho() Processor {
//	return &echo{}
//}

// Handle ...
//func (p *echo) Handle(m *pb.Command, client Client) error {
//	if m.GetEcho() == nil {
//		return nil
//	}
//
//	log := logrus.WithField("func", "processor.echo.Handler")
//	p.counter++
//
//	log.WithField("counter", p.counter).Debug("show counter")
//
//	// Send back client
//	client.Send(common.Message{
//		Body: "hello",
//	})
//
//	return nil
//}
