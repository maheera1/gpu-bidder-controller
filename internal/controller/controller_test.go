package controller

import (
	"context"
	"testing"
)

type fakeSystemd struct {
	active map[string]bool
	starts []string
	stops  []string
}

func newFakeSystemd() *fakeSystemd {
	return &fakeSystemd{active: map[string]bool{}}
}

func (f *fakeSystemd) Start(unit string) error {
	f.starts = append(f.starts, unit)
	f.active[unit] = true
	return nil
}
func (f *fakeSystemd) Stop(unit string) error {
	f.stops = append(f.stops, unit)
	f.active[unit] = false
	return nil
}
func (f *fakeSystemd) IsActive(unit string) bool { return f.active[unit] }

type fakeNotifier struct {
	events []struct {
		event  string
		prover string
	}
}

func (n *fakeNotifier) Send(event string, activeProver string) error {
	n.events = append(n.events, struct {
		event  string
		prover string
	}{event, activeProver})
	return nil
}

func TestController_AssignedStopsBothAndNotifies(t *testing.T) {
	sys := newFakeSystemd()
	sys.active["b1"] = true
	sys.active["b2"] = true

	n := &fakeNotifier{}
	c := New("b1", "b2", n)
	c.Systemd = sys
	c.Notifier = n

	c.OnAssigned(context.Background(), Prover1, "order-1")

	if len(sys.stops) != 2 {
		t.Fatalf("expected 2 stops, got %d: %#v", len(sys.stops), sys.stops)
	}
	if sys.stops[0] != "b1" || sys.stops[1] != "b2" {
		t.Fatalf("expected stops [b1 b2], got %#v", sys.stops)
	}
	if len(n.events) != 1 || n.events[0].event != "bidders_stopped" || n.events[0].prover != "prover1" {
		t.Fatalf("unexpected notifier events: %#v", n.events)
	}
}

func TestController_CompletedStartsBothWhenNoOrders(t *testing.T) {
	sys := newFakeSystemd()
	sys.active["b1"] = false
	sys.active["b2"] = false

	n := &fakeNotifier{}
	c := New("b1", "b2", n)
	c.Systemd = sys
	c.Notifier = n

	c.OnAssigned(context.Background(), Prover1, "order-1")
	c.OnCompleted(context.Background(), Prover1, "order-1")

	if len(sys.starts) != 2 {
		t.Fatalf("expected 2 starts, got %d: %#v", len(sys.starts), sys.starts)
	}
	if len(n.events) < 2 || n.events[1].event != "bidders_started" || n.events[1].prover != "none" {
		t.Fatalf("unexpected notifier events: %#v", n.events)
	}
}
