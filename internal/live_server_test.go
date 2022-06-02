package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiveDataServer(t *testing.T) {
	
	defaultViewModel := &ViewModel{
		fixture{
			Id:    "fixture-id-1",
			Title: "fixture-title",
			Teams: []fixtureTeam{
				{
					Id:    "team-id-1",
					Name:  "team-name-1",
					Score: 0,
				},
				{
					Id:    "team-id-2",
					Name:  "team-name-2",
					Score: 0,
				},
			},
			Tournament: fixtureTournament{
				Id:   "tournament-id",
				Name: "tournament-name",
			},
			WinningTeamId: "",
		},
	}

	setup := func() *LiveDataServer {
		winningTeamUpdateReceiver := &winningTeamUpdateReceiver{}
		scoreUpdateReceiver := &scoreUpdateReceiver{}
		
		return newLiveDataServer(defaultViewModel, *winningTeamUpdateReceiver, *scoreUpdateReceiver)
	}

	t.Run("when newLiveDataServer is called it should return a valid LiveDataServer", func(t *testing.T) {
		// Arrange
		viewModel := &ViewModel{}
		winningTeamUpdateReceiver := &winningTeamUpdateReceiver{}
		scoreUpdateReceiver := &scoreUpdateReceiver{}
		
		// Act
		server := newLiveDataServer(viewModel, *winningTeamUpdateReceiver, *scoreUpdateReceiver)

		// Assert
		assert.NotNil(t, server)
		assert.NotNil(t, server.winningTeamUpdateReceiver)
		assert.Equal(t, winningTeamUpdateReceiver, &server.winningTeamUpdateReceiver)
		assert.NotNil(t, server.scoreUpdateReceiver)
		assert.Equal(t, scoreUpdateReceiver, &server.scoreUpdateReceiver)
	})

	t.Run("when updateScoreAndPublish is called it should update the score", func(t *testing.T) {
		// Arrange
		server := setup()
		
		// Check precondition - score should be 0
		assert.Equal(t, 0, (*server.viewModel)[0].Teams[0].Score)

		// Act
		fixtureId := "fixture-id-1"
		teamID := "team-id-1"
		server.updateScoreAndPublish(fixtureId, teamID, 1)
		
		// Assert
		assert.Equal(t, 1, (*server.viewModel)[0].Teams[0].Score)
	})

	t.Run("when updateWinnerAndPublish is called it should update the winner", func(t *testing.T) {
		// Arrange
		server := setup()
		
		// Check precondition - winner should be empty string
		assert.Equal(t, "", (*server.viewModel)[0].WinningTeamId)

		// Act
		fixtureId := "fixture-id-1"
		teamID := "team-id-2"
		server.updateWinnerAndPublish(fixtureId, teamID)
		
		// Assert
		assert.Equal(t, "team-id-2", (*server.viewModel)[0].WinningTeamId)
	})

	t.Run("when findFixture is called it should return the fixture", func(t *testing.T) {
		// Arrange
		server := setup()
		
		// Act
		fixtureId := "fixture-id-1"
		fixture := server.findFixture(fixtureId)
		
		// Assert
		assert.Equal(t, fixtureId, fixture.Id)
	})

	t.Run("when findFixture is called with invalid fixtureId it should return nil", func(t *testing.T) {
		// Arrange
		server := setup()
		
		// Act
		fixtureId := "invalid-fixture-id"
		fixture := server.findFixture(fixtureId)
		
		// Assert
		assert.Nil(t, fixture)
	})

	t.Run("when findFixtureTeam is called it should return the fixture team", func(t *testing.T) {
		// Arrange
		server := setup()
		
		// Act
		fixtureId := "fixture-id-1"
		teamID := "team-id-1"
		fixtureTeam := server.findFixtureTeam(fixtureId, teamID)
		
		// Assert
		assert.Equal(t, teamID, fixtureTeam.Id)
	})

	t.Run("when findFixtureTeam is called with invalid fixtureId it should return nil", func(t *testing.T) {
		// Arrange
		server := setup()
		
		// Act
		fixtureId := "invalid-fixture-id"
		teamID := "team-id-1"
		fixtureTeam := server.findFixtureTeam(fixtureId, teamID)
		
		// Assert
		assert.Nil(t, fixtureTeam)
	})

	t.Run("when findFixtureTeam is called with invalid teamId it should return nil", func(t *testing.T) {
		// Arrange
		server := setup()
		
		// Act
		fixtureId := "fixture-id-1"
		teamID := "invalid-team-id"
		fixtureTeam := server.findFixtureTeam(fixtureId, teamID)
		
		// Assert
		assert.Nil(t, fixtureTeam)
	})
}