package models

type Transaction struct {
	Id      string         `json:"txid"`
	Parents []*Transaction `json:"vin"`
}

func NewTransaction(id string, parents []*Transaction) *Transaction {
	return &Transaction{
		Id:      id,
		Parents: parents,
	}
}

type TransactionAncestors struct {
	Id        string
	Ancestors []Transaction
}

func NewTransactionAncestors(id string, ancestors []Transaction) *TransactionAncestors {
	return &TransactionAncestors{
		Id:        id,
		Ancestors: ancestors,
	}
}
