package internal

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Zedronar/go-dummy-app.git/external"
)

var liveDataServer *LiveDataServer

func InitLiveServer() {
	// Query for initial fixtures
	viewModel := getStaticFixtures()

	// Sort viewmodel's fixtures and teams, so we can access them using binary search from now on.
	// TODO: In production, if fixtures are added dynamically, we should sort on every addition (see ADR.md).
	viewModel.Sort()

	// Initialize live data server
	winningTeamUpdateReceiver := &winningTeamUpdateReceiver{}
	scoreUpdateReceiver := &scoreUpdateReceiver{}
	liveDataServer = newLiveDataServer(viewModel, *winningTeamUpdateReceiver, *scoreUpdateReceiver)

	// Register live score receivers
	external.RegisterWinningTeamUpdateReceivers(&liveDataServer.winningTeamUpdateReceiver)
	external.RegisterScoreUpdateReceivers(&liveDataServer.scoreUpdateReceiver)

	// Initial viewModel publish
	liveDataServer.viewModel.PublishViewModel()

	// Handle live data requests
	http.HandleFunc("/livedata", liveDataServer.handleLiveDataRequest)
}

// 1- Use hash/digest?
// --> Getting a hash string from the client
// --> Compare it to the current hash (app level / redis)
// ^ read heavy scenario

// 2- Diffs/deltas

// * Cache
// -> MISS ?
// --> Retrieve from DB
// --> Update cache
// -> HIT
// --> return

// Massive data? -> LRU on cache (redis)

func getStaticFixtures() *ViewModel {

	// TODO: Implement retry mechanism
	// * Exponential backoff
	// * Max 3 retries
	// try / catch
	resp, err := http.Get("http://localhost:8080/fixtures")
	callFailureMax := 3
	failureCount := 0
	exponentialBackoff := 1

	for err != nil && failureCount < callFailureMax {
		resp, err = http.Get("http://localhost:8080/fixtures")

		// sleep
		// TODO: Fix exponential
		time.Sleep(time.Second * 1)

		// inc
		failureCount++
		exponentialBackoff *= 2

		// On each failure (might be transient)
		log.Fatal(err)
	}

	if err != nil {
		// Passed max attempts
		log.Fatal(err)
		panic("Fatal error - called too many times...")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var viewmodel ViewModel
	json.Unmarshal(body, &viewmodel)
	return &viewmodel
}

type scoreUpdateReceiver struct{}

func (t *scoreUpdateReceiver) Receive(update external.ScoreUpdate) {
	// Update viewmodel with new team score and publish it
	liveDataServer.updateScoreAndPublish(update.FixtureId(), update.TeamId(), update.Score())
}

type winningTeamUpdateReceiver struct{}

func (t *winningTeamUpdateReceiver) Receive(update external.WinningTeamUpdate) {
	// Update viewmodel with new winning team and publish it
	liveDataServer.updateWinnerAndPublish(update.FixtureId(), update.TeamId())
}
