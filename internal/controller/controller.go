package controller

import (
	"context"
	"log"
	"sync"

	"github.com/maheera1/gpu-bidder-controller/internal/systemd"
)

type Prover string

const (
	None    Prover = "none"
	Prover1 Prover = "prover1"
	Prover2 Prover = "prover2"
)

type Controller struct {
	mu sync.Mutex

	Bidder1Unit string
	Bidder2Unit string

	Systemd  Systemd
	Notifier Notifier

	Active Prover
	// Track assigned-but-not-completed orders per prover
	Assigned map[Prover]map[string]bool
}

func New(b1, b2 string, n Notifier) *Controller {
	return &Controller{
		Bidder1Unit: b1,
		Bidder2Unit: b2,
		Systemd:     systemd.NewManager(),
		Notifier:    n,
		Active:      None,
		Assigned: map[Prover]map[string]bool{
			Prover1: make(map[string]bool),
			Prover2: make(map[string]bool),
		},
	}
}

func (c *Controller) stopBothBidders() {
	if c.Systemd.IsActive(c.Bidder1Unit) {
		_ = c.Systemd.Stop(c.Bidder1Unit)
	}
	if c.Systemd.IsActive(c.Bidder2Unit) {
		_ = c.Systemd.Stop(c.Bidder2Unit)
	}
}

func (c *Controller) startBothBidders() {
	if !c.Systemd.IsActive(c.Bidder1Unit) {
		_ = c.Systemd.Start(c.Bidder1Unit)
	}
	if !c.Systemd.IsActive(c.Bidder2Unit) {
		_ = c.Systemd.Start(c.Bidder2Unit)
	}
}

func (c *Controller) OnAssigned(ctx context.Context, prover Prover, orderID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Assigned[prover][orderID] = true
	c.Active = prover

	log.Printf("ASSIGNED: %s order=%s -> STOP BOTH bidders", prover, orderID)
	c.stopBothBidders()
	_ = c.Notifier.Send("bidders_stopped", string(prover))
}

func (c *Controller) OnCompleted(ctx context.Context, prover Prover, orderID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.Assigned[prover], orderID)
	log.Printf("COMPLETED: %s order=%s", prover, orderID)

	// Start both bidders only when *no* assigned orders remain
	if len(c.Assigned[Prover1]) == 0 && len(c.Assigned[Prover2]) == 0 {
		c.Active = None
		log.Printf("NO assigned orders -> START BOTH bidders")
		c.startBothBidders()
		_ = c.Notifier.Send("bidders_started", string(None))
	}
}
