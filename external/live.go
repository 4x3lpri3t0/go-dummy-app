package external

import (
	"sync"
	"time"
)

var receiversLock sync.RWMutex
var registeredScoreUpdateReceivers []ScoreUpdateReceiver
var registeredWinningTeamUpdateReceivers []WinningTeamUpdateReceiver

func RegisterScoreUpdateReceivers(receivers ...ScoreUpdateReceiver) {
	receiversLock.Lock()
	defer receiversLock.Unlock()

	for _, receiver := range receivers {
		registeredScoreUpdateReceivers = append(registeredScoreUpdateReceivers, receiver)
	}
}

func RegisterWinningTeamUpdateReceivers(receivers ...WinningTeamUpdateReceiver) {
	receiversLock.Lock()
	defer receiversLock.Unlock()

	for _, receiver := range receivers {
		registeredWinningTeamUpdateReceivers = append(registeredWinningTeamUpdateReceivers, receiver)
	}
}

type WinningTeamUpdateReceiver interface {
	Receive(update WinningTeamUpdate)
}

type WinningTeamUpdate interface {
	FixtureId() string
	TeamId() string
}

type winningTeamUpdate struct {
	fixtureId string
	teamId    string
}

func (u *winningTeamUpdate) FixtureId() string {
	return u.fixtureId
}

func (u *winningTeamUpdate) TeamId() string {
	return u.teamId
}

type ScoreUpdateReceiver interface {
	Receive(update ScoreUpdate)
}

type ScoreUpdate interface {
	FixtureId() string
	TeamId() string
	Score() int
}

type scoreUpdate struct {
	fixtureId string
	teamId    string
	score     int
}

func (u *scoreUpdate) FixtureId() string {
	return u.fixtureId
}

func (u *scoreUpdate) TeamId() string {
	return u.teamId
}

func (u *scoreUpdate) Score() int {
	return u.score
}

type randomLiveScorePublisher struct {
	sync.Mutex
	fixtures         []*fixture
	fixtureScores    map[string][]int
	teamScoreLimit   int
	tickerDuration   time.Duration
	decisionProvider decisionProvider
}

func (p *randomLiveScorePublisher) StartRandomPublish() {
	ticker := time.NewTicker(p.tickerDuration)
	go func() {
		for {
			<-ticker.C
			p.doGenerateRandomScoreAndPublish()
		}
	}()
}

func (p *randomLiveScorePublisher) doGenerateRandomScoreAndPublish() {
	scoreUpdates := p.generateRandomScoreUpdates()
	for _, scoreUpdate := range scoreUpdates {
		p.publish(scoreUpdate)
	}
	p.calculateAndPublishWinningTeamUpdates(scoreUpdates)
}

func (p *randomLiveScorePublisher) calculateAndPublishWinningTeamUpdates(scoreUpdates []ScoreUpdate) {
	for _, scoreUpdate := range scoreUpdates {
		if scoreUpdate.Score() == p.teamScoreLimit {
			p.publishWinningTeamUpdate(&winningTeamUpdate{
				fixtureId: scoreUpdate.FixtureId(),
				teamId:    scoreUpdate.TeamId(),
			})
		}
	}
}

func (p *randomLiveScorePublisher) generateRandomScoreUpdates() []ScoreUpdate {
	p.Lock()
	defer p.Unlock()

	scoreUpdates := make([]ScoreUpdate, 0)

	for _, fixture := range p.fixtures {
		fixtureScores, found := p.fixtureScores[fixture.Id]
		if !found {
			fixtureScores = p.initialiseFixtureScores(*fixture)
			p.fixtureScores[fixture.Id] = fixtureScores
		}

		if !p.isFixtureLive(*fixture) || p.hasTeamWonFixture(fixtureScores) {
			continue
		}

		shouldUpdateFixture := p.decisionProvider.TrueFalse(3)
		if shouldUpdateFixture {
			for i, team := range fixture.Teams {
				shouldUpdateTeam := p.decisionProvider.TrueFalse(5)
				if shouldUpdateTeam {
					fixtureScores[i] = fixtureScores[i] + 1
					p.fixtureScores[fixture.Id] = fixtureScores
					scoreUpdates = append(scoreUpdates, &scoreUpdate{
						fixtureId: fixture.Id,
						teamId:    team.Id,
						score:     fixtureScores[i],
					})

					hasTeamWon := fixtureScores[i] >= p.teamScoreLimit
					if hasTeamWon {
						break
					}
				}
			}
		}
	}

	return scoreUpdates
}

func (p *randomLiveScorePublisher) isFixtureLive(fixture fixture) bool {
	return fixture.ScheduledStartTime < time.Now().UTC().Unix()
}

func (p *randomLiveScorePublisher) hasTeamWonFixture(fixtureScores []int) bool {
	for _, score := range fixtureScores {
		if score >= p.teamScoreLimit {
			return true
		}
	}
	return false
}

func (p *randomLiveScorePublisher) initialiseFixtureScores(fixture fixture) []int {
	scores := make([]int, len(fixture.Teams))
	for i, _ := range fixture.Teams {
		scores[i] = 0
	}
	return scores
}

func (p *randomLiveScorePublisher) publish(scoreUpdate ScoreUpdate) {
	receiversLock.RLock()
	defer receiversLock.RUnlock()

	for _, receiver := range registeredScoreUpdateReceivers {
		receiver.Receive(scoreUpdate)
	}
}

func (p *randomLiveScorePublisher) publishWinningTeamUpdate(update WinningTeamUpdate) {
	receiversLock.RLock()
	defer receiversLock.RUnlock()

	for _, receiver := range registeredWinningTeamUpdateReceivers {
		receiver.Receive(update)
	}
}

func newRandomLiveScorePublisher(
	fixtures []*fixture,
	publishTickDuration time.Duration,
	decisionProvider decisionProvider,
	teamScoreLimit int) *randomLiveScorePublisher {

	return &randomLiveScorePublisher{
		fixtures:         fixtures,
		fixtureScores:    make(map[string][]int),
		tickerDuration:   publishTickDuration,
		decisionProvider: decisionProvider,
		teamScoreLimit:   teamScoreLimit,
	}
}
