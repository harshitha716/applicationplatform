package sparkpost

import (
	sp "github.com/SparkPost/gosparkpost"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Send(transmission *sp.Transmission) (id string, resp *sp.Response, err error) {
	args := m.Called(transmission)
	if args.Get(1) == nil {
		return args.String(0), nil, args.Error(2)
	}
	return args.String(0), args.Get(1).(*sp.Response), args.Error(2)
}
