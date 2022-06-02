package data

import (
	"log"
	"net/http"
	"time"
)

type fixtureTournament struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type fixtureTeam struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	// Live data
	Score int `json:"score"`
}

type fixture struct {
	Id                 string            `json:"id"`
	Title              string            `json:"title"`
	Tournament         fixtureTournament `json:"tournament"`
	Teams              []fixtureTeam     `json:"teams"`
	ScheduledStartTime int64             `json:"scheduledStartTimeUnixSeconds"`

	// Live data
	WinningTeamId string `json:"winningTeamId"`
}

type DataProvider interface {
	Retrieve() string
}

type HttpDataProvider struct {
}

func NewHttpProvider() DataProvider {
	return &HttpDataProvider{}
}

// TODO: json instead of string
func (x *HttpDataProvider) Retrieve() string {
	// TODO: CAll the DB
	// mongoClient.

	resp, err := getStaticFixtures()

	if err != nil {
		log.Fatal(err)
	}

	return resp
}

// TODO:
// func getStaticFixtures() *ViewModel {
func getStaticFixtures() (string, error) {

	// Retry mechanism
	resp, err := http.Get("http://localhost:8080/fixtures")
	callFailureMax := 3
	failureCount := 0
	exponentialBackoff := 1

	for err != nil && failureCount < callFailureMax {
		resp, err = http.Get("http://localhost:8080/fixtures")

		// TODO: Fix exponential
		time.Sleep(time.Second * 1)
		failureCount++
		exponentialBackoff *= 2

		// On each failure (might be transient)
		log.Fatal(err)
	}

	if err != nil {
		// Passed max attempts
		log.Fatal(err)
		// panic("Fatal error - called too many times...")
		return "", err
	}

	defer resp.Body.Close()

	// TODO: temp
	return "response", nil

	// body, err := ioutil.ReadAll(resp.Body)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var viewmodel ViewModel
	// json.Unmarshal(body, &viewmodel)
	// return &viewmodel
}
