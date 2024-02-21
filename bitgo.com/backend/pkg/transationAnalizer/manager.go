package transationAnalizer

import (
	"bitgo.com/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

type Analyzer struct {
	Transactions []models.Transaction
}

var (
	dbOnce       sync.Once
	thisAnalyzer *Analyzer
)

func NewAnalyzer() *Analyzer {

	dbOnce.Do(func() {
		thisAnalyzer = &Analyzer{}
	})
	return thisAnalyzer
}

var (
	ErrTestErr = errors.New("test error")
)

const (
	baseURL      = "https://blockstream.info/api"
	blockHeight  = 680000
	maxAncestors = 10

	cachePath = "./cache/"
)

func (ta *Analyzer) FetchBlockHash(height int) (string, error) {
	url := fmt.Sprintf("%s/block-height/%d", baseURL, height)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	blockHash := string(body)

	return blockHash, nil
}

func (ta *Analyzer) FetchTransactionsForBlock(hash string) (map[string]*models.Transaction, error) {

	cachedFilePath := hash + ".json"
	if _, err := os.Stat(cachedFilePath); err == nil {
		return ta.readCachedTransactions(cachedFilePath)
	}

	var allTransactions []models.Transaction

	startIndex := 0
	for {
		transactions, err := ta.fetchTransactionsPage(hash, startIndex)
		if err != nil {
			return nil, err
		}

		allTransactions = append(allTransactions, transactions...)

		if len(transactions) < 25 {
			break
		}

		startIndex += 25
	}

	if err := ta.writeTransactionsToFile(cachedFilePath, allTransactions); err != nil {
		return nil, err
	}

	return ta.readCachedTransactions(cachedFilePath)
}

func (ta *Analyzer) readCachedTransactions(filePath string) (map[string]*models.Transaction, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction
	if err := json.Unmarshal(data, &transactions); err != nil {
		return nil, err
	}
	transactionsMap := make(map[string]*models.Transaction)
	for i, _ := range transactions {
		transactionsMap[transactions[i].Id] = &transactions[i]
	}

	return transactionsMap, nil
}

func (ta *Analyzer) writeTransactionsToFile(filePath string, transactions []models.Transaction) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(transactions)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (ta *Analyzer) fetchTransactionsPage(hash string, startIndex int) ([]models.Transaction, error) {
	url := fmt.Sprintf("%s/block/%s/txs/%d", baseURL, hash, startIndex)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//TODO: handle 201 and other error code

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, nil
	}
	var transactions []models.Transaction
	if err := json.Unmarshal(body, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (ta *Analyzer) GetAncestorsTxns(txnId string) ([]models.Transaction, error) {
	return []models.Transaction{}, nil
}
