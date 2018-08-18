package lib

import (
	"fmt"
	"net/http"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetDepth(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)
	client.EXPECT().request(&clientOption{
		endpoint: publicApiEndpoint,
		method:   http.MethodGet,
		path:     fmt.Sprintf(formatDepth, entity.PairBtcJpy),
	}).Return(
		[]byte(`{"success":1,"data":{"asks":[["964745","0.0004"],["964753","0.0010"]],"bids":[["964254","0.0060"],["964249","0.0200"]],"timestamp":1526387708186}}`),
		nil,
	)
	api := &APIImpl{
		client: client,
	}
	depth, err := api.GetDepth(entity.PairBtcJpy)

	assert.Nil(t, err)
	assert.NotNil(t, depth)
	assert.Equal(t, []float64{964745.0, 0.0004}, depth.Asks[0])
	assert.Equal(t, []float64{964254.0, 0.0060}, depth.Bids[0])
}
