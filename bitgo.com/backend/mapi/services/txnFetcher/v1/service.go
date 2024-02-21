package v1

import (
	"bitgo.com/pkg/models"
	"context"
	"fmt"
	"sort"
)

type (
	TransactionAnalyzer interface {
		FetchBlockHash(height int) (string, error)
		FetchTransactionsForBlock(hash string) (map[string]*models.Transaction, error)
	}
	AncestorFinder interface {
		FindAncestorsForAllTxn(txns map[string]*models.Transaction) map[string][]*models.Transaction
	}
	Server struct {
		transactionAnalyzer TransactionAnalyzer
		ancestorFinder      AncestorFinder
	}
)

func NewServer(transactionAnalyzer TransactionAnalyzer, ancestorFinder AncestorFinder) *Server {
	return &Server{
		transactionAnalyzer: transactionAnalyzer,
		ancestorFinder:      ancestorFinder,
	}
}

type (
	FindTopTransactionsRequest struct {
		BlockNumber int
	}
	Transaction struct {
		Id           string
		AncestorSize int
	}
	FindTopTransactionsResponse struct {
		Transactions []Transaction
	}
)

func (s *Server) FindTopTransactions(ctx context.Context, req *FindTopTransactionsRequest) (*FindTopTransactionsResponse, error) {

	blockHash, err := s.transactionAnalyzer.FetchBlockHash(req.BlockNumber)
	if err != nil {
		fmt.Println("Error fetching block hash:", err)
		return nil, err
	}

	blockTransactions, err := s.transactionAnalyzer.FetchTransactionsForBlock(blockHash)
	if err != nil {
		fmt.Println("Error fetching transactions:", err)
		return nil, err
	}

	trasactionAnsisterList := s.ancestorFinder.FindAncestorsForAllTxn(blockTransactions)
	if err != nil {
		fmt.Println("Error analyzing transactions:", err)
		return nil, err
	}
	var txns []Transaction
	for id, list := range trasactionAnsisterList {
		txns = append(txns, Transaction{Id: id, AncestorSize: len(list)})
	}
	sort.Slice(txns, func(i, j int) bool {
		return txns[i].AncestorSize > txns[j].AncestorSize
	})

	return &FindTopTransactionsResponse{
		Transactions: txns[:10],
	}, nil
}
