package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
)

var viewModelLock sync.RWMutex

type LiveDataServer struct {
	viewModel                 *ViewModel
	winningTeamUpdateReceiver winningTeamUpdateReceiver
	scoreUpdateReceiver       scoreUpdateReceiver
}

func (server *LiveDataServer) handleLiveDataRequest(w http.ResponseWriter, _ *http.Request) {
	viewModelLock.RLock()
	defer viewModelLock.RUnlock()

	jsonBytes, err := json.Marshal(server.viewModel)
	if err != nil {
		log.Print(fmt.Sprintf("@handleLiveDataRequest -> error marshalling fixtures: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Print(fmt.Sprintf("@handleLiveDataRequest -> error writing bytes: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func newLiveDataServer(
	viewModel *ViewModel,
	winningTeamUpdateReceiver winningTeamUpdateReceiver,
	scoreUpdateReceiver scoreUpdateReceiver) *LiveDataServer {
	return &LiveDataServer{
		viewModel:                 viewModel,
		winningTeamUpdateReceiver: winningTeamUpdateReceiver,
		scoreUpdateReceiver:       scoreUpdateReceiver,
	}
}

func (server *LiveDataServer) updateScoreAndPublish(fixtureId string, teamId string, newScore int) {
	viewModelLock.Lock()
	defer viewModelLock.Unlock()

	// Find fixture team
	fixtureTeam := server.findFixtureTeam(fixtureId, teamId)
	if fixtureTeam == nil {
		log.Println(fmt.Sprintf("@updateScoreAndPublish -> fixtureId '%s' or teamId '%s' not found!", fixtureId, teamId))
		return
	}
	
	// Update fixture team score
	fixtureTeam.Score = newScore

	// Publish
	server.viewModel.PublishViewModel()

	return
}

func (server *LiveDataServer) updateWinnerAndPublish(fixtureId string, teamId string) {
	viewModelLock.Lock()
	defer viewModelLock.Unlock()

	// Find fixture
	fixture := server.findFixture(fixtureId)
	if fixture == nil {
		log.Println(fmt.Sprintf("@updateWinnerAndPublish -> fixtureId '%s' not found!", fixtureId))
		return
	}

	// Update fixture winner
	fixture.WinningTeamId = teamId

	// Publish
	server.viewModel.PublishViewModel()

	return
}

// Find fixture by id using binary search
func (server *LiveDataServer) findFixture(fixtureId string) *fixture {
	fixtureIndex := sort.Search(len(*server.viewModel), func(i int) bool {
		return (*server.viewModel)[i].Id >= fixtureId
	})

	if fixtureIndex >= len(*server.viewModel) || (*server.viewModel)[fixtureIndex].Id != fixtureId {
		return nil
	}

	return &(*server.viewModel)[fixtureIndex]
}

// Find fixtureTeam in fixture by id using binary search
func (server *LiveDataServer) findFixtureTeam(fixtureId string, teamId string) *fixtureTeam {
	fixture := server.findFixture(fixtureId)
	if fixture == nil {
		return nil
	}

	fixtureTeamIndex := sort.Search(len(fixture.Teams), func(i int) bool {
		return fixture.Teams[i].Id >= teamId
	})

	if fixtureTeamIndex >= len(fixture.Teams) || fixture.Teams[fixtureTeamIndex].Id != teamId {
		return nil
	}

	return &fixture.Teams[fixtureTeamIndex]
}