package lib

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetTransactions(t *testing.T) {
	now := time.Now()
	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)
	client.EXPECT().request(&clientOption{
		endpoint: publicApiEndpoint,
		method:   http.MethodGet,
		path:     fmt.Sprintf(formatTransactionsAll, entity.PairBtcJpy),
	}).Return(
		[]byte(`{"success":1,"data":{"transactions":[{"transaction_id":10654992,"side":"sell","price":"962005","amount":"0.0100","executed_at":1526390858690},{"transaction_id":10654991,"side":"sell","price":"962005","amount":"0.0090","executed_at":1526390856716},{"transaction_id":10654990,"side":"sell","price":"962010","amount":"0.0060","executed_at":1526390856658},{"transaction_id":10654989,"side":"sell","price":"962020","amount":"0.0050","executed_at":1526390856596}]}}`),
		nil,
	)
	client.EXPECT().request(&clientOption{
		endpoint: publicApiEndpoint,
		method:   http.MethodGet,
		path:     fmt.Sprintf(formatTransactions, entity.PairBtcJpy, now.Format("20060102")),
	}).Return(
		[]byte(`{"success":1,"data":{"transactions":[{"transaction_id":10507587,"side":"sell","price":"956290","amount":"0.0500","executed_at":1526256003017},{"transaction_id":10507588,"side":"sell","price":"956290","amount":"0.0453","executed_at":1526256003458}]}}`),
		nil,
	)

	api := &APIImpl{
		client: client,
	}

	transactionsAll, err := api.GetTransactions(entity.PairBtcJpy, nil)
	assert.Nil(t, err)
	assert.NotNil(t, transactionsAll[0])
	assert.NotEqual(t, 0, transactionsAll[0].TransactionId)
	assert.NotEqual(t, 0, transactionsAll[0].Side)
	assert.NotEqual(t, 0, transactionsAll[0].Price)
	assert.NotEqual(t, 0, transactionsAll[0].Amount)
	assert.NotEqual(t, 0, transactionsAll[0].ExecutedAt)

	transactions, err := api.GetTransactions(entity.PairBtcJpy, &now)
	assert.Nil(t, err)
	assert.NotNil(t, transactions[0])
	assert.NotEqual(t, 0, transactions[0].TransactionId)
	assert.NotEqual(t, 0, transactions[0].Side)
	assert.NotEqual(t, 0, transactions[0].Price)
	assert.NotEqual(t, 0, transactions[0].Amount)
	assert.NotEqual(t, 0, transactions[0].ExecutedAt)
}
