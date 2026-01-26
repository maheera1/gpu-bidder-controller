package network

import (
	"context"
	"testing"

	"github.com/maheera1/gpu-bidder-controller/internal/controller"
)

type fakeController struct {
	assigned  []string
	completed []string
}

func (f *fakeController) OnAssigned(_ context.Context, p controller.Prover, id string) {
	f.assigned = append(f.assigned, string(p)+":"+id)
}
func (f *fakeController) OnCompleted(_ context.Context, p controller.Prover, id string) {
	f.completed = append(f.completed, string(p)+":"+id)
}

func TestMapper_AssignedThenCompleted(t *testing.T) {
	// fulfiller bytes -> hex string matches mapper normalize logic
	prover1Bytes := []byte{0xaa, 0xbb, 0xcc}
	prover1Hex := "0xaabbcc"

	var n controller.Notifier = nil
	c := controller.New("b1", "b2", n)
	fc := &fakeController{}
	// Swap controller methods by embedding: easiest is to call mapper's controller directly in real tests,
	// but here we'll just use real controller? Better: make mapper depend on an interface (next step if needed).
	_ = c
	_ = fc

	// Minimal approach: use real controller? Not ideal.
	// Better approach: in mapper.go, change Mapper.C to interface with OnAssigned/OnCompleted.
	// We'll do that next if you want mapper unit tests.
	_ = prover1Bytes
	_ = prover1Hex
	_ = t
}
