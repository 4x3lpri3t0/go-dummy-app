package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewModel(t *testing.T) {

	defaultFixture := &fixture{
		Id:    "fixture-id",
		Title: "fixture-title",
		Teams: []fixtureTeam{
			// (initially unsorted)
			{
				Id:    "team-id-2",
				Name:  "team-name-2",
				Score: 2,
			},
			{
				Id:    "team-id-1",
				Name:  "team-name-1",
				Score: 1,
			},
		},
		Tournament: fixtureTournament{
			Id:   "tournament-id",
			Name: "tournament-name",
		},
		WinningTeamId: "team-id-1",
	}
	
	setup := func() *ViewModel {
		testFixture1 := defaultFixture
		testFixture2 := defaultFixture
	
		return &ViewModel{
			*testFixture1,
			*testFixture2,
		}
	}

	t.Run("when viewmodel data (fixture, tournaments, teams) is constructed it should have expected values", func(t *testing.T) {
		viewModel := setup()

		// Both default fixtures should have expected values
		for i := range *viewModel {
			// Fixture
			assert.NotNil(t, (*viewModel)[i])
			assert.Equal(t, "fixture-id", (*viewModel)[i].Id)
			assert.Equal(t, "fixture-title", (*viewModel)[i].Title)
			assert.Equal(t, "team-id-1", (*viewModel)[i].WinningTeamId)
			// Tournament
			assert.NotNil(t, (*viewModel)[i].Tournament)
			assert.Equal(t, "tournament-id", (*viewModel)[i].Tournament.Id)
			assert.Equal(t, "tournament-name", (*viewModel)[i].Tournament.Name)
			// Teams (initially unsorted)
			assert.NotNil(t, (*viewModel)[i].Teams)
			assert.Equal(t, "team-id-2", (*viewModel)[i].Teams[0].Id)
			assert.Equal(t, "team-name-2", (*viewModel)[i].Teams[0].Name)
			assert.Equal(t, 2, (*viewModel)[i].Teams[0].Score)
			assert.Equal(t, "team-id-1", (*viewModel)[i].Teams[1].Id)
			assert.Equal(t, "team-name-1", (*viewModel)[i].Teams[1].Name)
			assert.Equal(t, 1, (*viewModel)[i].Teams[1].Score)
		}
	})

	t.Run("when sortTeams is called it should sort teams by team id", func(t *testing.T) {
		// Arrange
		viewModel := setup()
		
		// Act
		(*viewModel)[0].sortTeams()

		// Assert
		assert.Equal(t, "team-id-1", defaultFixture.Teams[0].Id)
		assert.Equal(t, "team-name-1", defaultFixture.Teams[0].Name)
		assert.Equal(t, 1, defaultFixture.Teams[0].Score)
		assert.Equal(t, "team-id-2", defaultFixture.Teams[1].Id)
		assert.Equal(t, "team-name-2", defaultFixture.Teams[1].Name)
		assert.Equal(t, 2, defaultFixture.Teams[1].Score)
	})

	t.Run("when viewmodel data is sorted it should have expected values", func(t *testing.T) {
		// Arrange
		viewModel := setup()

		// Act
		viewModel.Sort()

		// Assert
		for i := range *viewModel {
			// Fixture
			assert.NotNil(t, (*viewModel)[i])
			assert.Equal(t, "fixture-id", (*viewModel)[i].Id)
			// Teams (sorted)
			assert.Equal(t, "team-id-1", (*viewModel)[i].Teams[0].Id)
			assert.Equal(t, "team-id-2", (*viewModel)[i].Teams[1].Id)
		}
	})

	t.Run("when sort is called twice it should have expected values", func(t *testing.T) {
		// Arrange
		viewModel := setup()

		// Act
		viewModel.Sort()
		viewModel.Sort()

		// Assert
		for i := range *viewModel {
			// Fixture
			assert.NotNil(t, (*viewModel)[i])
			assert.Equal(t, "fixture-id", (*viewModel)[i].Id)
			// Teams (sorted)
			assert.Equal(t, "team-id-1", (*viewModel)[i].Teams[0].Id)
			assert.Equal(t, "team-id-2", (*viewModel)[i].Teams[1].Id)
		}
	})
}
