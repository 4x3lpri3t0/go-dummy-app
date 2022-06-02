package external

import "github.com/stretchr/testify/mock"

type decisionProviderMock struct {
	mock.Mock
}

func (m *decisionProviderMock) TrueFalse(chance int) bool {
	args := m.Called(chance)
	return args.Bool(0)
}
