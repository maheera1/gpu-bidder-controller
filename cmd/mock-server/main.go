package main

import (
	"encoding/hex"
	"log"
	"net"
	"time"

	netpb "github.com/maheera1/gpu-bidder-controller/gen/network"
	typepb "github.com/maheera1/gpu-bidder-controller/gen/types"
	"google.golang.org/grpc"
)

type MockServer struct {
	netpb.UnimplementedProverNetworkServer
}

func mustDecodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func (s *MockServer) SubscribeProofRequests(req *typepb.GetFilteredProofRequestsRequest, stream netpb.ProverNetwork_SubscribeProofRequestsServer) error {
	ctx := stream.Context()
	log.Println("mock: client subscribed; sending test ProofRequest updates")

	// IMPORTANT: Your mapper compares hex(pr.Fulfiller) to PROVER*_ID.
	// So we send fulfiller BYTES that hex-encode cleanly:
	// 0x1111 => bytes {0x11, 0x11} => "1111"
	prover1Fulfiller := mustDecodeHex("1111")
	prover2Fulfiller := mustDecodeHex("2222")

	// Order 1 assigned to prover1
	pr1 := &typepb.ProofRequest{
		RequestId:         []byte("order-001"),
		Fulfiller:         prover1Fulfiller,
		FulfillmentStatus: typepb.FulfillmentStatus_ASSIGNED,
	}
	log.Println("mock: sending order-001 ASSIGNED to prover1")
	if err := stream.Send(pr1); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)

	// Order 1 fulfilled
	pr2 := &typepb.ProofRequest{
		RequestId:         []byte("order-001"),
		Fulfiller:         prover1Fulfiller,
		FulfillmentStatus: typepb.FulfillmentStatus_FULFILLED,
	}
	log.Println("mock: sending order-001 FULFILLED")
	if err := stream.Send(pr2); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)

	// Order 2 assigned to prover2
	pr3 := &typepb.ProofRequest{
		RequestId:         []byte("order-002"),
		Fulfiller:         prover2Fulfiller,
		FulfillmentStatus: typepb.FulfillmentStatus_ASSIGNED,
	}
	log.Println("mock: sending order-002 ASSIGNED to prover2")
	if err := stream.Send(pr3); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)

	// Order 2 unfulfillable (terminal)
	pr4 := &typepb.ProofRequest{
		RequestId:         []byte("order-002"),
		Fulfiller:         prover2Fulfiller,
		FulfillmentStatus: typepb.FulfillmentStatus_UNFULFILLABLE,
	}
	log.Println("mock: sending order-002 UNFULFILLABLE")
	if err := stream.Send(pr4); err != nil {
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	srv := grpc.NewServer()
	netpb.RegisterProverNetworkServer(srv, &MockServer{})

	log.Println("mock grpc server listening on :50051")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
