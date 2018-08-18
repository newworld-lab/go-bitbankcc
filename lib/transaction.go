package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	time "time"

	entity "github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/pkg/errors"
)

const (
	formatTransactionsAll = "/%s/transactions"
	formatTransactions    = "/%s/transactions/%s"
)

type transactionsResponse struct {
	baseResponse
	Data struct {
		Transactions transactions `json:"transactions"`
	} `json:"data"`
}

type transaction struct {
	entity.Transaction
	ExecutedAt int64 `json:"executed_at"`
}

type transactions []transaction

func (ts transactions) convert() entity.Transactions {
	transactions := make(entity.Transactions, 0)
	for _, t := range ts {
		transactions = append(transactions, entity.Transaction{
			TransactionId: t.TransactionId,
			Side:          t.Side,
			Price:         t.Price,
			Amount:        t.Amount,
			ExecutedAt:    time.Unix(t.ExecutedAt/1000, t.ExecutedAt%1000*1000000),
		})
	}
	return transactions
}

func (api *APIImpl) GetTransactions(pair entity.TypePair, t *time.Time) (entity.Transactions, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	var path string
	if t == nil {
		path = fmt.Sprintf(formatTransactionsAll, pair)
	} else {
		path = fmt.Sprintf(formatTransactions, pair, t.Format("20060102"))
	}
	bytes, err := api.client.request(&clientOption{
		endpoint: publicApiEndpoint,
		method:   http.MethodGet,
		path:     path,
	})

	if err != nil {
		return nil, err
	}

	res := new(transactionsResponse)
	json.Unmarshal(bytes, res)

	err = res.parseError()
	if err != nil {
		return nil, err
	}

	return res.Data.Transactions.convert(), nil
}
