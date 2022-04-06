package mock

import (
	"github.com/golang/mock/gomock"
)

type HasherMockParams struct {
	HashTimes    int
	HashAnyTimes bool
	HashReq      gomock.Matcher
	HashRes      string
	HashResErr   error
}

func (p HasherMockParams) NewHasher(ctrl *gomock.Controller) *MockHasher {
	mock := NewMockHasher(ctrl)

	callTimes(mock.EXPECT().Hash(p.HashReq), p.HashTimes, p.HashAnyTimes).Return(p.HashRes, p.HashResErr)

	return mock
}

func callTimes(c *gomock.Call, times int, anyTimes bool) *gomock.Call {
	if anyTimes {
		return c.AnyTimes()
	}

	return c.Times(times)
}
