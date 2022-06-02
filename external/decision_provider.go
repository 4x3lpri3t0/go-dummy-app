package external

import "math/rand"

type decisionProvider interface {
	TrueFalse(chance int) bool
}

type randomDecisionProvider struct{}

func (p *randomDecisionProvider) TrueFalse(chance int) bool {
	if chance == 0 {
		return false
	}
	randomInt := rand.Int()
	return randomInt%chance == 0
}

func newDecisionProvider() decisionProvider {
	return &randomDecisionProvider{}
}
