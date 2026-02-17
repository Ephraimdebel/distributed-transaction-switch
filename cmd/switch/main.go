package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Ephraimdebel/transaction-switch/internal/transaction/application/command"
	"github.com/Ephraimdebel/transaction-switch/internal/transaction/application/service"
	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/event"
	"github.com/Ephraimdebel/transaction-switch/internal/transaction/infrastructure/messaging"
	"github.com/Ephraimdebel/transaction-switch/internal/transaction/infrastructure/persistence"

	"net"

	txgrpc "github.com/Ephraimdebel/transaction-switch/internal/transaction/api/grpc"
	transactionpb "github.com/Ephraimdebel/transaction-switch/internal/transaction/api/proto"
	"google.golang.org/grpc"
)

func main() {
	repo := persistence.NewMemoryTransactionRepository()
	idempotencyStore := messaging.NewProcessedEventStore()
	dlq := messaging.NewDeadLetterQueue()

	publisher := messaging.NewEventPublisher(
		4,   // worker count
		100, // queue size
		dlq,
	)

	publisher.Register(
		"TransactionReserved",
		messaging.IdempotentHandler(
			idempotencyStore,
			func(e event.DomainEvent) error {
				fmt.Println("Handling TransactionReserved:", e.EventID())

				// simulate processing
				time.Sleep(200 * time.Millisecond)

				return nil
			},
		),
	)

	txService := service.NewTransactionService(
		repo,
		publisher,
	)

	ctx := context.Background()

	cmd := command.ReserveTransaction{
		TransactionID: "tx-123",
		Amount:        1000,
	}

	if err := txService.Reserve(ctx, cmd); err != nil {
		panic(err)
	}

	fmt.Println("DLQ events:")
	for _, e := range dlq.List() {
		fmt.Printf(
			"Event=%s Reason=%s Retries=%d\n",
			e.Event.EventName(),
			e.Reason,
			e.RetryCount,
		)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	txGRPCServer := txgrpc.NewTransactionGRPCServer(txService)

	transactionpb.RegisterTransactionServiceServer(
		grpcServer,
		txGRPCServer,
	)

	fmt.Println("gRPC server running on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

}
