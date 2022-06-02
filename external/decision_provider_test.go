package external

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomDecisionProvider(t *testing.T) {

	setup := func() *randomDecisionProvider {
		return &randomDecisionProvider{}
	}

	t.Run("TrueFalse", func(t *testing.T) {

		t.Run("when chance 1 always returns true", func(t *testing.T) {
			decisionProvider := setup()

			hasPublishedFalse := false
			for i := 0; i < 100; i++ {
				decisionIsTrue := decisionProvider.TrueFalse(1)
				if !decisionIsTrue {
					hasPublishedFalse = true
				}
			}

			assert.False(t, hasPublishedFalse)
		})

		t.Run("when chance 0 always returns false", func(t *testing.T) {
			decisionProvider := setup()

			hasPublishedTrue := false
			for i := 0; i < 100; i++ {
				decisionIsTrue := decisionProvider.TrueFalse(0)
				if decisionIsTrue {
					hasPublishedTrue = true
				}
			}

			assert.False(t, hasPublishedTrue)
		})

		t.Run("when chance 3 returns approx 1/3 true", func(t *testing.T) {
			decisionProvider := setup()

			trueCounter := 0
			for i := 0; i < 10000; i++ {
				decisionIsTrue := decisionProvider.TrueFalse(3)
				if decisionIsTrue {
					trueCounter += 1
				}
			}

			assert.True(t, trueCounter > 3000)
			assert.True(t, trueCounter < 4000)
		})
	})
}
