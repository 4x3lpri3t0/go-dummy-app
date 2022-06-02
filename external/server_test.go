package external

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStaticDataServer(t *testing.T) {

	setup := func() *staticDataServer {
		return &staticDataServer{
			fixtures: make([]*fixture, 0),
		}
	}

	t.Run("HandleFixturesRequest", func(t *testing.T) {

		t.Run("returns expected json body", func(t *testing.T) {
			server := setup()
			recorder := httptest.NewRecorder()

			server.fixtures = []*fixture{
				{
					Id:    "id",
					Title: "title",
					Teams: []fixtureTeam{
						{
							Id: "team-id-1",
						},
						{
							Id: "team-id-2",
						},
					},
				},
			}

			expectedJsonBytes, _ := json.Marshal(server.fixtures)

			server.HandleFixturesRequest(recorder, nil)

			assert.Equal(t, expectedJsonBytes, recorder.Body.Bytes())
		})

		t.Run("when fixture valid returns ok response", func(t *testing.T) {
			server := setup()
			recorder := httptest.NewRecorder()

			server.HandleFixturesRequest(recorder, nil)

			assert.Equal(t, http.StatusOK, recorder.Code)
		})
	})
}
