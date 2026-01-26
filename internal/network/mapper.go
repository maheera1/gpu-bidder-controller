package network

import (
	"context"
	"encoding/hex"
	"log"
	"strings"

	typepb "github.com/maheera1/gpu-bidder-controller/gen/types"
	"github.com/maheera1/gpu-bidder-controller/internal/controller"
)

type MapperConfig struct {
	Prover1ID string // fulfiller address hex (with/without 0x)
	Prover2ID string // fulfiller address hex (with/without 0x)
}

type Mapper struct {
	Cfg MapperConfig
	C   *controller.Controller

	lastFulfillment map[string]typepb.FulfillmentStatus
	assigned        map[string]bool
}

func NewMapper(cfg MapperConfig, c *controller.Controller) *Mapper {
	return &Mapper{
		Cfg:             cfg,
		C:               c,
		lastFulfillment: make(map[string]typepb.FulfillmentStatus),
		assigned:        make(map[string]bool),
	}
}

func (m *Mapper) HandleUpdate(ctx context.Context, pr *typepb.ProofRequest) {
	reqID := m.getRequestID(pr)
	if reqID == "" {
		return
	}

	fulfStatus := pr.FulfillmentStatus

	// Ignore repeats
	if prev, ok := m.lastFulfillment[reqID]; ok && prev == fulfStatus {
		return
	}
	m.lastFulfillment[reqID] = fulfStatus

	prover := m.getAssignedProver(pr)
	if prover == controller.None {
		return
	}

	// Assignment
	if !m.assigned[reqID] && pr.FulfillmentStatus == typepb.FulfillmentStatus_ASSIGNED {
		m.assigned[reqID] = true
		log.Printf("MAPPER: assigned req=%s -> %s", reqID, prover)
		m.C.OnAssigned(ctx, prover, reqID)
		return
	}

	// Completion (terminal)
	if m.assigned[reqID] && isTerminalFulfillment(fulfStatus) {
		delete(m.assigned, reqID)
		log.Printf("MAPPER: completed req=%s -> %s", reqID, prover)
		m.C.OnCompleted(ctx, prover, reqID)
		return
	}
}

func (m *Mapper) getRequestID(pr *typepb.ProofRequest) string {
	if pr == nil || len(pr.RequestId) == 0 {
		return ""
	}
	return "0x" + hex.EncodeToString(pr.RequestId)
}

func (m *Mapper) getAssignedProver(pr *typepb.ProofRequest) controller.Prover {
	if pr == nil || len(pr.Fulfiller) == 0 {
		return controller.None
	}
	f := normalizeHex(hex.EncodeToString(pr.Fulfiller))
	p1 := normalizeHex(m.Cfg.Prover1ID)
	p2 := normalizeHex(m.Cfg.Prover2ID)

	switch {
	case p1 != "" && f == p1:
		return controller.Prover1
	case p2 != "" && f == p2:
		return controller.Prover2
	default:
		return controller.None
	}
}

func isTerminalFulfillment(s typepb.FulfillmentStatus) bool {
	return s == typepb.FulfillmentStatus_FULFILLED ||
		s == typepb.FulfillmentStatus_UNFULFILLABLE
}

func normalizeHex(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	return strings.ToLower(s)
}
