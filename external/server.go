package external

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type fixtureTournament struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type fixtureTeam struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type fixture struct {
	Id                 string            `json:"id"`
	Title              string            `json:"title"`
	Tournament         fixtureTournament `json:"tournament"`
	Teams              []fixtureTeam     `json:"teams"`
	ScheduledStartTime int64             `json:"scheduledStartTimeUnixSeconds"`
}

type staticDataServer struct {
	fixtures []*fixture
}

func (s *staticDataServer) HandleFixturesRequest(w http.ResponseWriter, _ *http.Request) {
	jsonBytes, err := json.Marshal(s.fixtures)
	if err != nil {
		log.Print(fmt.Sprintf("@HandleFixturesRequest -> error marshalling fixtures: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Print(fmt.Sprintf("@HandleFixturesRequest -> error writing bytes: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func newStaticDataServer(fixtures []*fixture) *staticDataServer {
	return &staticDataServer{
		fixtures: fixtures,
	}
}
