package service

import "encoding/json"

type EventSender interface {
	SendEvent(topic string, message []byte) error
}
type EventService struct {
	sender EventSender
}

func NewEventSender(sender EventSender) *EventService {
	return &EventService{
		sender: sender,
	}
}

func (s *EventService) SendEvent(topic string, message any) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return s.sender.SendEvent(topic, msg)
}
