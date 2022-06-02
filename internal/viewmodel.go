package internal

import (
	"fmt"
	"log"
	"sort"

	"github.com/Zedronar/go-dummy-app.git/external"
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

type ViewModel []fixture

func (viewModel *ViewModel) PublishViewModel() {
	err := external.PublishViewModel(viewModel)

	if err != nil {
		log.Print(fmt.Sprintf("@updateAndPublishViewModel -> error publishing viewmodel: %s", err.Error()))
		return
	}
}

func (viewModel *ViewModel) Sort() {
	// Sort viewmodel's fixtures by fixture id
	sort.Slice(*viewModel, func(i, j int) bool {
		return (*viewModel)[i].Id < (*viewModel)[j].Id
	})

	// Sort each fixture's teams by team id
	for fixture := range *viewModel {
		(*viewModel)[fixture].sortTeams()
	}
}

// Sort teams by team id
func (fixture *fixture) sortTeams() {
	sort.Slice(fixture.Teams, func(i, j int) bool {
		return fixture.Teams[i].Id < fixture.Teams[j].Id
	})
}
