package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	typepb "github.com/maheera1/gpu-bidder-controller/gen/types"
	"github.com/maheera1/gpu-bidder-controller/internal/controller"
	"github.com/maheera1/gpu-bidder-controller/internal/network"
	"github.com/maheera1/gpu-bidder-controller/internal/notifier"
)

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func main() {
	// Config
	notifyURL := getenv("NOTIFY_URL", "http://127.0.0.1:8081/signal")
	b1 := getenv("BIDDER1_UNIT", "bidder1.service")
	b2 := getenv("BIDDER2_UNIT", "bidder2.service")

	grpcAddr := getenv("NETWORK_GRPC_ADDR", "")
	prover1ID := getenv("PROVER1_ID", "")
	prover2ID := getenv("PROVER2_ID", "")

	// Core controller
	c := controller.New(b1, b2, notifier.Client{URL: notifyURL})

	// Mapper: turns ProofRequest updates into controller actions
	mapper := network.NewMapper(network.MapperConfig{
		Prover1ID: prover1ID,
		Prover2ID: prover2ID,
	}, c)

	// Shutdown handling
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("controller started (GRPC MODE)")
	log.Printf("NETWORK_GRPC_ADDR=%s", grpcAddr)
	log.Printf("DRY_RUN=%s", getenv("DRY_RUN", "0"))

	// Stream + process updates
	err := network.StartProofRequestStream(ctx, network.StreamConfig{Addr: grpcAddr}, func(pr *typepb.ProofRequest) {
		mapper.HandleUpdate(ctx, pr)
	})
	if err != nil {
		log.Fatalf("stream error: %v", err)
	}
}
