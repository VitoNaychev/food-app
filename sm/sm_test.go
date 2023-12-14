package sm_test

import (
	"errors"
	"testing"

	"github.com/VitoNaychev/food-app/sm"
	"github.com/VitoNaychev/food-app/testutil"
)

const (
	Locked sm.State = iota
	Unlocked
)

const (
	Push sm.Event = iota
	Coin
)

func TestSimpleSMTransisions(t *testing.T) {
	deltas := []sm.Delta{
		{Locked, Push, Locked, nil, nil},
		{Locked, Coin, Unlocked, nil, nil},
		{Unlocked, Push, Locked, nil, nil},
	}

	testsm := sm.New(Locked, deltas, nil)

	t.Run("transisions from Locked to Unlocked on Coin", func(t *testing.T) {
		err := testsm.Exec(Coin)

		testutil.AssertNil(t, err)
		testutil.AssertEqual(t, testsm.Current, Unlocked)

	})

	t.Run("returns ErrInvalidEvent on state Unlocked and event Coin", func(t *testing.T) {
		err := testsm.Exec(Coin)

		testutil.AssertEqual(t, err, sm.ErrInvalidEvent)
	})

	t.Run("transisions from Unlocked to Locked on Push", func(t *testing.T) {
		err := testsm.Exec(Push)

		testutil.AssertNil(t, err)
		testutil.AssertEqual(t, testsm.Current, Locked)
	})

	t.Run("remains on Locked on Push", func(t *testing.T) {
		err := testsm.Exec(Push)

		testutil.AssertNil(t, err)
		testutil.AssertEqual(t, testsm.Current, Locked)
	})
}

var DummyPredicateError = errors.New("dummy predicate error")

type StubContext struct {
	retValue    bool
	shouldError bool
}

func StubPredicate(delta sm.Delta, ictx sm.Context) (bool, error) {
	ctx := ictx.(*StubContext)

	if ctx.shouldError {
		return false, DummyPredicateError
	}

	return ctx.retValue, nil
}

func NotedStubPredicate(delta sm.Delta, ictx sm.Context) (bool, error) {
	ctx := ictx.(*StubContext)

	return !ctx.retValue, nil
}

func TestSMPredicates(t *testing.T) {
	context := StubContext{shouldError: false}

	deltas := []sm.Delta{
		{Locked, Coin, Unlocked, sm.Predicate(StubPredicate), nil},
		{Locked, Coin, Locked, sm.Predicate(NotedStubPredicate), nil},
		{Unlocked, Push, Locked, nil, nil},
	}

	testsm := sm.New(Locked, deltas, &context)

	t.Run("remains in Locked when DummyPredicate fails", func(t *testing.T) {
		context.retValue = false

		err := testsm.Exec(Coin)

		testutil.AssertNil(t, err)
		testutil.AssertEqual(t, testsm.Current, Locked)
	})

	t.Run("transisitions to Unlocked when DummyPredicate succeedes", func(t *testing.T) {
		context.retValue = true

		err := testsm.Exec(Coin)

		testutil.AssertNil(t, err)
		testutil.AssertEqual(t, testsm.Current, Unlocked)
	})

	t.Run("returns predicate's error", func(t *testing.T) {
		context.shouldError = true

		testsm.Reset(Locked)
		err := testsm.Exec(Coin)

		testutil.AssertEqual(t, err, DummyPredicateError)
	})
}

var DummyCallbackError = errors.New("dummy callback error")

type SpyContext struct {
	wasCallbackCalled bool
	shouldError       bool
}

func SpyCallback(delta sm.Delta, ictx sm.Context) error {
	ctx := ictx.(*SpyContext)

	ctx.wasCallbackCalled = true

	if ctx.shouldError {
		return DummyCallbackError
	}

	return nil
}

func TestSMCallbacks(t *testing.T) {
	context := SpyContext{wasCallbackCalled: false, shouldError: false}

	deltas := []sm.Delta{
		{Locked, Coin, Unlocked, nil, sm.Callback(SpyCallback)},
		{Unlocked, Push, Locked, nil, nil},
	}

	testsm := sm.New(Locked, deltas, &context)

	t.Run("calls delta's callback", func(t *testing.T) {
		err := testsm.Exec(Coin)

		testutil.AssertNil(t, err)
		testutil.AssertEqual(t, context.wasCallbackCalled, true)
	})

	t.Run("returns callback's error", func(t *testing.T) {
		context.shouldError = true

		testsm.Reset(Locked)
		err := testsm.Exec(Coin)

		testutil.AssertEqual(t, err, DummyCallbackError)
		testutil.AssertEqual(t, context.wasCallbackCalled, true)
	})
}
