# Parking Lot Management System

A scalable, clean-architecture based Parking Lot System with priority allocation, re-entry constraints, and agent-service communication.

---

# Overview

This project implements a Parking Lot Allocation System that assigns optimal parking slots to vehicles based on multiple real-world constraints such as:

- Capacity balancing across levels
- Vehicle-type segregation
- Re-entry restrictions
- Priority-based allocation (VIP/Emergency)

The system is designed using Clean Architecture and simulates a real-world distributed setup using:

- Service (Daemon) → Core allocation engine
- Agent (CLI Tool) → Sends parking requests

---

# Architecture

```text
+------------------+        HTTP/JSON        +----------------------+
|   Agent (CLI)     |  --------------------> |   Parking Service     |
| (Simulates user)  |                        |   (Core Logic)        |
+------------------+                        +----------------------+
                                                    |
                                                    v
                                          In-Memory Data Store
```

---

# Features

- Optimal slot allocation using Min Heap
- Re-entry restriction (configurable time window)
- Vehicle-type based slot segregation
- Priority handling for VIP/Emergency vehicles
- Concurrent request handling (Goroutines + Mutex)
- Clean Architecture (Domain, Usecase, Repository)
- Unit test support with mocks

---

# Vehicle Types

| Type   | Example    |
| ------ | ---------- |
| Small  | Motorcycle |
| Medium | Car        |
| Large  | Truck      |

---

# Constraints Implemented

1️⃣ Capacity Constraints
- Each level has 10–100 slots
- Balanced allocation across levels
- Prevents overloading a single level

2️⃣ Re-Entry Restriction
- A vehicle cannot re-enter the same level within a fixed duration
- Allowed - (After time expiry ,Or in a different level)

3️⃣ Vehicle Type Segregation
- Each level contains - (Small slots, Medium slots, Large slots)

4️⃣ Priority Parking
- VIP/Emergency vehicles get higher priority
- Regular vehicles are queued

---

# Allocation Strategy

- Receive request
- Check vehicle history
- Validate re-entry rule
- Filter eligible levels
- Select best level using Min Heap
- Allocate slot
- Update state
- Return response

---

# Data Structures Used

| Use Case          | Data Structure |
| ----------------- | -------------- |
| Slot Management   | Stack / Queue  |
| Priority Handling | Min Heap       |
| Vehicle History   | Hash Map       |
| Levels            | Slice          |

---

# Project Structure

```text
parking-lot/
│
├── cmd/
│   ├── service/        # Service binary (HTTP server)
│   └── agent/          # CLI tool
│
├── internal/
│   ├── config/             # Load env
│   ├── domain/             # Models & interfaces
│   ├── usecase/            # Business logic (Allocator,Dispatcher)
│   ├── infrastructure/     # In-memory implementation
│   └── delivery/            # HTTP handlers
│
├── tests/
│   └── mocks/          # Mock implementations
│
└── README.md
```

---

# Getting Started

Prerequisites
- Go 1.20+
- Git
---

## Running the Project

```bash
# Clone Repository
git clone https://github.com/Varunjp/Parking-lot.git
cd Parking-lot

# Set up environment variables
cp .env.example .env

# Install dependencies
go mod tidy

```

Create .env file
```env
PARKING_LEVELS=3
SMALL_SLOTS_PER_LEVEL=10
MEDIUM_SLOTS_PER_LEVEL=7
LARGE_SLOTS_PER_LEVEL=4
REENTRY_SECONDS=3600
HTTP_PORT=8080
```

# Run the Service

```bash
go run cmd/service/main.go
```
Server runs at:
```bash
http://localhost:8080
```

# Run the Agent

```text
go run cmd/agent/main.go
```

## CLI Usage (Agent Commands):

The Agent CLI accepts commands in a structured format to simulate parking operations.

## Park a Vehicle
```bash
PARK <vehicleID> <vehicleType> <customerType>
```
### Parameters
| Field        | Description               |
| ------------ | ------------------------- |
| vehicleID    | Unique vehicle identifier |
| vehicleType  | Type of vehicle           |
| customerType | Priority category         |

### Vehicle Types 

- SMALL → Motorcycle
- MEDIUM → Car
- LARGE → Truck

### Customer Types

- EMERGENCY → Highest priority
- VIP → High priority
- REGULAR → Normal priority

### Example

```bash
PARK KL01 MEDIUM VIP
```

## Exit a Vehicle

```bash
EXIT <vehicleID>
```

### Example

```bash
EXIT KL01
```

```bash
vehicle KL01 exited successfully
```

---

## Concurrency
- Handles multiple simultaneous requests
- Uses: (Goroutines, Mutex locks for safe updates)
---

## Running Tests

```text
go test ./tests
```

### Test Coverage Includes:
- Allocation correctness
- Priority handling
- Re-entry validation
- Edge cases
---
