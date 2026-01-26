package network

import (
	"context"
	"fmt"
	"log"
	"time"

	netpb "github.com/maheera1/gpu-bidder-controller/gen/network"
	typepb "github.com/maheera1/gpu-bidder-controller/gen/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// StreamConfig configures how we subscribe to proof requests.
type StreamConfig struct {
	Addr string

	// Optional: add filters later (program/requester/status etc.) if needed.
}

// StartProofRequestStream connects and streams ProofRequest updates.
// For each update, it calls onUpdate(pr).
func StartProofRequestStream(
	ctx context.Context,
	cfg StreamConfig,
	onUpdate func(pr *typepb.ProofRequest),
) error {
	if cfg.Addr == "" {
		return fmt.Errorf("NETWORK_GRPC_ADDR is empty")
	}

	// In prod you may need TLS; for now insecure is simplest for dev/testing.
	conn, err := grpc.DialContext(
		ctx,
		cfg.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("dial grpc %s: %w", cfg.Addr, err)
	}
	defer conn.Close()

	client := netpb.NewProverNetworkClient(conn)

	// Basic subscription request (no filters).
	req := &typepb.GetFilteredProofRequestsRequest{}

	stream, err := client.SubscribeProofRequests(ctx, req)
	if err != nil {
		return fmt.Errorf("SubscribeProofRequests: %w", err)
	}

	log.Printf("subscribed to proof requests at %s", cfg.Addr)

	for {
		pr, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("stream recv: %w", err)
		}
		if pr == nil {
			continue
		}
		onUpdate(pr)

		// Small sleep to avoid hot-looping if stream bursts extremely fast.
		// Safe to remove later if not needed.
		time.Sleep(5 * time.Millisecond)
	}
}
