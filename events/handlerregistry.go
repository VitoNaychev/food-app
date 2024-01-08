package events

import (
	"errors"
	"reflect"
	"strconv"
)

type TopicRegistry map[string]bool

func (t TopicRegistry) RegisterTopic(topic string) {
	t[topic] = true
}

func (t TopicRegistry) GetTopics() []string {
	topics := []string{}

	for key := range t {
		topics = append(topics, key)
	}

	return topics
}

type RegistryKey string

func GetRegistryKey(topic string, eventID EventID) RegistryKey {
	return RegistryKey(topic + "-" + strconv.Itoa(int(eventID)))
}

type RegistryEntry struct {
	eventHandler InterfaceEventHandler
	eventType    reflect.Type
}

type ConsumerGroup struct {
	eventHandlerRegistry map[RegistryKey]RegistryEntry
}

func (r *ConsumerGroup) RegisterEventHandler(topic string, eventID EventID, eventHandler InterfaceEventHandler, eventType reflect.Type) {
	entry := RegistryEntry{
		eventHandler: eventHandler,
		eventType:    eventType,
	}

	r.eventHandlerRegistry[GetRegistryKey(topic, eventID)] = entry
}

type InterfaceEventHandler func(event InterfaceEvent) error

func EventHandlerWrapper[T any](eventHandler func(event Event[T]) error) InterfaceEventHandler {
	return InterfaceEventHandler(func(ievent InterfaceEvent) error {
		event := Event[T]{
			EventID:     ievent.EventID,
			AggregateID: ievent.AggregateID,
			Timestamp:   ievent.Timestamp,
		}

		if payload, ok := ievent.Payload.(T); ok {
			event.Payload = payload

			return eventHandler(event)
		} else {
			return errors.New("incorrect type")
		}
	})
}
