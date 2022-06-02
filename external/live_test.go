package external

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type testScoreUpdateReceiver struct {
	receivedUpdates []ScoreUpdate
}

func (t *testScoreUpdateReceiver) Receive(update ScoreUpdate) {
	if t.receivedUpdates == nil {
		t.receivedUpdates = make([]ScoreUpdate, 0)
	}
	t.receivedUpdates = append(t.receivedUpdates, update)
}

type testWinningTeamUpdateReceiver struct {
	receivedUpdates []WinningTeamUpdate
}

func (t *testWinningTeamUpdateReceiver) Receive(update WinningTeamUpdate) {
	if t.receivedUpdates == nil {
		t.receivedUpdates = make([]WinningTeamUpdate, 0)
	}
	t.receivedUpdates = append(t.receivedUpdates, update)
}

func TestRegisterScoreUpdateReceivers(t *testing.T) {

	tearDown := func() {
		registeredScoreUpdateReceivers = make([]ScoreUpdateReceiver, 0)
	}

	t.Run("when single receiver registers receiver", func(t *testing.T) {
		receiver := &testScoreUpdateReceiver{}

		RegisterScoreUpdateReceivers(receiver)

		assert.Equal(t, 1, len(registeredScoreUpdateReceivers))
		assert.Equal(t, receiver, registeredScoreUpdateReceivers[0])

		tearDown()
	})

	t.Run("when single receiver registers receiver", func(t *testing.T) {
		receiver1 := &testScoreUpdateReceiver{}
		receiver2 := &testScoreUpdateReceiver{}

		RegisterScoreUpdateReceivers(receiver1, receiver2)

		assert.Equal(t, 2, len(registeredScoreUpdateReceivers))

		tearDown()
	})
}

func TestWinningTeamUpdateReceivers(t *testing.T) {

	tearDown := func() {
		registeredWinningTeamUpdateReceivers = make([]WinningTeamUpdateReceiver, 0)
	}

	t.Run("when single receiver registers receiver", func(t *testing.T) {
		receiver := &testWinningTeamUpdateReceiver{}

		RegisterWinningTeamUpdateReceivers(receiver)

		assert.Equal(t, 1, len(registeredWinningTeamUpdateReceivers))
		assert.Equal(t, receiver, registeredWinningTeamUpdateReceivers[0])

		tearDown()
	})

	t.Run("when single receiver registers receiver", func(t *testing.T) {
		receiver1 := &testWinningTeamUpdateReceiver{}
		receiver2 := &testWinningTeamUpdateReceiver{}

		RegisterWinningTeamUpdateReceivers(receiver1, receiver2)

		assert.Equal(t, 2, len(registeredWinningTeamUpdateReceivers))

		tearDown()
	})
}

func TestRandomLiveScorePublisher(t *testing.T) {

	var scoreUpdateReceiver *testScoreUpdateReceiver
	var winningTeamUpdateReceiver *testWinningTeamUpdateReceiver
	var decisionProvider *decisionProviderMock

	setup := func() *randomLiveScorePublisher {
		scoreUpdateReceiver = &testScoreUpdateReceiver{}
		winningTeamUpdateReceiver = &testWinningTeamUpdateReceiver{}
		decisionProvider = new(decisionProviderMock)
		RegisterScoreUpdateReceivers(scoreUpdateReceiver)
		RegisterWinningTeamUpdateReceivers(winningTeamUpdateReceiver)
		return &randomLiveScorePublisher{
			fixtures:         make([]*fixture, 0),
			fixtureScores:    make(map[string][]int),
			decisionProvider: decisionProvider,
			teamScoreLimit:   2,
		}
	}

	tearDown := func() {
		registeredScoreUpdateReceivers = make([]ScoreUpdateReceiver, 0)
	}

	defaultFixture := &fixture{
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
	}

	t.Run("doGenerateRandomScoreAndPublish", func(t *testing.T) {

		t.Run("when should update fixture and should update team publishes score update", func(t *testing.T) {
			publisher := setup()

			testFixture := defaultFixture
			publisher.fixtures = []*fixture{testFixture}

			decisionProvider.On("TrueFalse", mock.Anything).Times(2).Return(true)
			decisionProvider.On("TrueFalse", mock.Anything).Return(false)

			publisher.doGenerateRandomScoreAndPublish()

			assert.Equal(t, 1, len(scoreUpdateReceiver.receivedUpdates))

			tearDown()
		})

		t.Run("when publishes score has matching fixture id", func(t *testing.T) {
			publisher := setup()

			testFixture := defaultFixture
			publisher.fixtures = []*fixture{testFixture}

			decisionProvider.On("TrueFalse", mock.Anything).Return(true)

			publisher.doGenerateRandomScoreAndPublish()

			resultScoreUpdate := scoreUpdateReceiver.receivedUpdates[0]

			assert.Equal(t, testFixture.Id, resultScoreUpdate.FixtureId())

			tearDown()
		})

		t.Run("when publishes score has matching team id", func(t *testing.T) {
			publisher := setup()

			testFixture := defaultFixture
			publisher.fixtures = []*fixture{testFixture}

			decisionProvider.On("TrueFalse", mock.Anything).Return(true)

			publisher.doGenerateRandomScoreAndPublish()

			resultScoreUpdateTeam1 := scoreUpdateReceiver.receivedUpdates[0]
			resultScoreUpdateTeam2 := scoreUpdateReceiver.receivedUpdates[1]

			assert.Equal(t, testFixture.Teams[0].Id, resultScoreUpdateTeam1.TeamId())
			assert.Equal(t, testFixture.Teams[1].Id, resultScoreUpdateTeam2.TeamId())

			tearDown()
		})

		t.Run("when publishes score score is incremented by 1", func(t *testing.T) {
			publisher := setup()

			testFixture := defaultFixture
			publisher.fixtures = []*fixture{testFixture}
			publisher.teamScoreLimit = 3

			decisionProvider.On("TrueFalse", mock.Anything).Return(true)

			publisher.doGenerateRandomScoreAndPublish()
			publisher.doGenerateRandomScoreAndPublish()

			resultScoreUpdate := scoreUpdateReceiver.receivedUpdates[2]

			assert.Equal(t, 2, resultScoreUpdate.Score())

			tearDown()
		})

		t.Run("when decision provider returns false for should update fixture does not publish score update", func(t *testing.T) {
			publisher := setup()

			testFixture := defaultFixture
			publisher.fixtures = []*fixture{testFixture}
			publisher.teamScoreLimit = 3

			decisionProvider.On("TrueFalse", mock.Anything).Return(false)

			publisher.doGenerateRandomScoreAndPublish()

			assert.Equal(t, 0, len(scoreUpdateReceiver.receivedUpdates))

			tearDown()
		})

		t.Run("when decision provider returns false for should update team does not publish score update", func(t *testing.T) {
			publisher := setup()

			testFixture := defaultFixture
			publisher.fixtures = []*fixture{testFixture}
			publisher.teamScoreLimit = 3

			decisionProvider.On("TrueFalse", mock.Anything).Times(1).Return(true)
			decisionProvider.On("TrueFalse", mock.Anything).Return(false)

			publisher.doGenerateRandomScoreAndPublish()

			assert.Equal(t, 0, len(scoreUpdateReceiver.receivedUpdates))

			tearDown()
		})

		t.Run("when fixture scheduled to start in future should not publish update", func(t *testing.T) {
			publisher := setup()

			testFixture := defaultFixture
			testFixture.ScheduledStartTime = time.Now().UTC().Add(1 * time.Hour).Unix()
			publisher.fixtures = []*fixture{testFixture}
			publisher.teamScoreLimit = 3

			decisionProvider.On("TrueFalse", mock.Anything).Return(true)

			publisher.doGenerateRandomScoreAndPublish()

			assert.Equal(t, 0, len(scoreUpdateReceiver.receivedUpdates))

			tearDown()
		})

		t.Run("when fixture scheduled start time is in past should publish update", func(t *testing.T) {
			publisher := setup()

			testFixture := defaultFixture
			testFixture.ScheduledStartTime = time.Now().UTC().Add(-1 * time.Hour).Unix()
			publisher.fixtures = []*fixture{testFixture}
			publisher.teamScoreLimit = 3

			decisionProvider.On("TrueFalse", mock.Anything).Return(true)

			publisher.doGenerateRandomScoreAndPublish()

			assert.Equal(t, 2, len(scoreUpdateReceiver.receivedUpdates))

			tearDown()
		})

		t.Run("when team has won fixture should not publish update", func(t *testing.T) {
			publisher := setup()

			testFixture := defaultFixture
			testFixture.ScheduledStartTime = time.Now().UTC().Add(-1 * time.Hour).Unix()
			publisher.fixtures = []*fixture{testFixture}
			publisher.teamScoreLimit = 1

			decisionProvider.On("TrueFalse", mock.Anything).Return(true)

			publisher.doGenerateRandomScoreAndPublish()
			publisher.doGenerateRandomScoreAndPublish()

			assert.Equal(t, 1, len(scoreUpdateReceiver.receivedUpdates))

			tearDown()
		})

		t.Run("when team reaches score limit publishes winning team update with correct values", func(t *testing.T) {
			publisher := setup()

			expectedFixtureId := "expected-fixture-id"
			expectedTeamId := "expected-team-id"

			testFixture := defaultFixture
			testFixture.Id = expectedFixtureId
			testFixture.Teams[0].Id = expectedTeamId
			publisher.fixtures = []*fixture{testFixture}
			publisher.teamScoreLimit = 1

			decisionProvider.On("TrueFalse", mock.Anything).Return(true)

			publisher.doGenerateRandomScoreAndPublish()

			assert.Equal(t, 1, len(winningTeamUpdateReceiver.receivedUpdates))
			assert.Equal(t, expectedFixtureId, winningTeamUpdateReceiver.receivedUpdates[0].FixtureId())
			assert.Equal(t, expectedTeamId, winningTeamUpdateReceiver.receivedUpdates[0].TeamId())

			tearDown()
		})

		t.Run("when team has already reached score limit does not publish winning team update", func(t *testing.T) {
			publisher := setup()

			expectedFixtureId := "expected-fixture-id"
			expectedTeamId := "expected-team-id"

			testFixture := defaultFixture
			testFixture.Id = expectedFixtureId
			testFixture.Teams[0].Id = expectedTeamId
			publisher.fixtures = []*fixture{testFixture}
			publisher.teamScoreLimit = 1

			decisionProvider.On("TrueFalse", mock.Anything).Return(true)

			publisher.doGenerateRandomScoreAndPublish()
			publisher.doGenerateRandomScoreAndPublish()

			assert.Equal(t, 1, len(winningTeamUpdateReceiver.receivedUpdates))

			tearDown()
		})
	})
}
