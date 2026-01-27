# GPU Bidder Controller

A **Go-based controller** that manages two bidder services sharing the same GPU cluster.

The controller listens to **ProofRequest** updates over **gRPC**, stops both bidders when a prover is assigned work, restarts them when work completes, and sends **HTTP notifications** on bidder state changes.

---

## Prerequisites

* **Go 1.22+**
* **systemd** (Linux)
* **protoc + Go plugins** (only if regenerating protobufs)
* Access to a **Succinct ProverNetwork gRPC endpoint**

  * Or use the **mock server** for local testing

---

## Clone the Repository

```bash
git clone https://github.com/maheera1/gpu-bidder-controller.git
cd gpu-bidder-controller
```

---

## Environment Configuration

### Required Variables

```bash
export NETWORK_GRPC_ADDR="localhost:50051"     # gRPC endpoint

export PROVER1_ID="0x50685c8c3924ae1af4dd7c1e1e7e9243b5c06cba"
export PROVER2_ID="0x2222222222222222222222222222222222222222"

export BIDDER1_UNIT="bidder1.service"
export BIDDER2_UNIT="bidder2.service"

export NOTIFY_URL="http://127.0.0.1:8081/signal"

# Dry-run mode (no real systemd calls)
export DRY_RUN="1"   # set to 0 for real services
```

---

## Local Testing (Recommended First)

### Terminal 1 — Receiver

```bash
go run ./cmd/receiver
```

### Terminal 2 — Mock gRPC Server

```bash
go run ./cmd/mock-server
```

### Terminal 3 — Controller

```bash
go run ./cmd/controller
```

### Expected Output

* Controller logs showing **bidders stopping/starting**
* Receiver logs showing:

  * `bidders_stopped`
  * `bidders_started`
* Mock server emitting test **ProofRequests**

---

## Run with Real Services (Production-style)

1. Replace `NETWORK_GRPC_ADDR` with the real **Succinct gRPC endpoint**
2. Set actual **systemd bidder unit names**
3. Disable dry-run mode

```bash
export DRY_RUN="0"
go run ./cmd/controller
```

---

## Tests

Run all unit tests:

```bash
go test ./...
```

### Test Coverage

Tests validate:

* Correct bidder **stop/start behavior**
* Correct handling of **order assignment and completion**
* No real **systemd** or **network calls** during testing

---

## Notes

* `gen/` contains **auto-generated protobuf code** — do **not** edit manually
* `cmd/mock-server` is **only for local testing**
* Controller logic is **fully testable** and **production-ready**

---

