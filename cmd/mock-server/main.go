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

	// Real prover ID (20-byte address)
	prover1Fulfiller := mustDecodeHex("50685c8c3924ae1af4dd7c1e1e7e9243b5c06cba")

	// Dummy second prover (won’t be used)
	prover2Fulfiller := mustDecodeHex("2222222222222222222222222222222222222222")

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
