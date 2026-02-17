package grpc


import (
	"context"

	transactionpb "github.com/Ephraimdebel/transaction-switch/internal/transaction/api/proto"
	"github.com/Ephraimdebel/transaction-switch/internal/transaction/application/command"
	"github.com/Ephraimdebel/transaction-switch/internal/transaction/application/service"
)

type TransactionGRPCServer struct {
	transactionpb.UnimplementedTransactionServiceServer
	app *service.TransactionService
}

func NewTransactionGRPCServer(app *service.TransactionService) *TransactionGRPCServer {
	return &TransactionGRPCServer{app: app}
}

func (s *TransactionGRPCServer) ReserveTransaction(
	ctx context.Context,
	req *transactionpb.ReserveTransactionRequest,
) (*transactionpb.ReserveTransactionResponse, error) {

	cmd := command.ReserveTransaction{
		TransactionID: req.TransactionId,
		Amount:        req.Amount,
	}

	if err := s.app.Reserve(ctx, cmd); err != nil {
		return nil, err
	}

	return &transactionpb.ReserveTransactionResponse{
		Status: "RESERVED",
	}, nil
}
