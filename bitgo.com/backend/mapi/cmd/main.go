package main

import (
	v1 "bitgo.com/mapi/services/txnFetcher/v1"
	"bitgo.com/pkg/transationAnalizer"
	"context"
	"fmt"
)

func main() {
	txnAnalyzer := transationAnalizer.NewAnalyzer()
	ancestorFinder := transationAnalizer.AncestorFinder{}

	s := v1.NewServer(txnAnalyzer, &ancestorFinder)
	resp, _ := s.FindTopTransactions(context.Background(), &v1.FindTopTransactionsRequest{
		BlockNumber: 680000,
	})
	for _, transaction := range resp.Transactions {
		fmt.Println(transaction.AncestorSize, transaction.Id)
	}

}
