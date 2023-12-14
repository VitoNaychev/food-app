package sm

import "errors"

var (
	ErrSpurious     = errors.New("state machine is in a spurious state")
	ErrInvalidEvent = errors.New("state doesn't support this event")
)

type Context interface{}

type State int
type Event int

type Predicate func(Delta, Context) (bool, error)
type Callback func(Delta, Context) error

type Delta struct {
	Current   State
	Event     Event
	Next      State
	Predicate Predicate
	Callback  Callback
}

type SM struct {
	Current State
	Deltas  []Delta
	Context Context
}

func New(initial State, deltas []Delta, context Context) SM {
	sm := SM{
		Current: initial,
		Deltas:  deltas,
		Context: context,
	}

	return sm
}

func (s *SM) Exec(event Event) error {
	for _, delta := range s.Deltas {
		if delta.Event == event && delta.Current == s.Current {
			if delta.Predicate != nil {
				if ok, err := delta.Predicate(delta, s.Context); err != nil {
					return err
				} else if ok == false {
					continue
				}
			}
			if delta.Callback != nil {
				if err := delta.Callback(delta, s.Context); err != nil {
					return err
				}
			}
			s.Current = delta.Next
			return nil
		}
	}

	return ErrInvalidEvent
}

func (s *SM) Reset(current State) {
	s.Current = current
}
