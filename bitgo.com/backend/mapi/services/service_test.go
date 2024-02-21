package services

import (
	"bitgo.com/pkg/models"
	"bitgo.com/pkg/transationAnalizer"
	"fmt"
	"sort"
	"testing"
)

func TestService(t *testing.T) {

	ancestorFinder := transationAnalizer.AncestorFinder{}

	blockTxns := make(map[string]*models.Transaction)
	txn1 := models.NewTransaction("1", []*models.Transaction{})
	txn2 := models.NewTransaction("2", []*models.Transaction{txn1})
	txn3 := models.NewTransaction("3", []*models.Transaction{txn1})
	txn4 := models.NewTransaction("4", []*models.Transaction{txn2})
	txn5 := models.NewTransaction("5", []*models.Transaction{txn3})
	txn6 := models.NewTransaction("6", []*models.Transaction{txn4, txn5})
	blockTxns[txn1.Id] = txn1
	blockTxns[txn2.Id] = txn2
	blockTxns[txn3.Id] = txn3
	blockTxns[txn4.Id] = txn4
	blockTxns[txn5.Id] = txn5
	blockTxns[txn6.Id] = txn6

	result := ancestorFinder.FindAncestorsForAllTxn(blockTxns)

	type Transaction struct {
		Id           string
		AncestorSize int
	}

	var txns []Transaction
	for id, list := range result {
		txns = append(txns, Transaction{Id: id, AncestorSize: len(list)})
	}
	sort.Slice(txns, func(i, j int) bool {
		return txns[i].AncestorSize > txns[j].AncestorSize
	})
	fmt.Print(txns)
}
