package transationAnalizer

import "bitgo.com/pkg/models"

type AncestorFinder struct {
	blockTxns map[string]*models.Transaction
}

func (af *AncestorFinder) FindAncestors(targetNode *models.Transaction) []*models.Transaction {

	visited := make(map[*models.Transaction]bool)
	ancestors := make(map[*models.Transaction]bool)

	af.helper(targetNode, visited, ancestors)

	delete(ancestors, targetNode)

	result := []*models.Transaction{}
	for ancestor := range ancestors {
		result = append(result, ancestor)
	}
	return result
}

func (af *AncestorFinder) FindAncestorsForAllTxn(blockTxns map[string]*models.Transaction) map[string][]*models.Transaction {

	af.blockTxns = blockTxns

	results := make(map[string][]*models.Transaction)
	for _, txn := range blockTxns {
		results[txn.Id] = af.FindAncestors(txn)
	}
	return results

}

func (af *AncestorFinder) helper(node *models.Transaction, visited map[*models.Transaction]bool,
	ancestors map[*models.Transaction]bool) {

	visited[node] = true

	ancestors[node] = true

	for _, parent := range node.Parents {
		_, ok := af.blockTxns[parent.Id]
		if !ok {
			continue
		}
		if !visited[parent] {
			af.helper(af.blockTxns[parent.Id], visited, ancestors)
		}
	}
}
