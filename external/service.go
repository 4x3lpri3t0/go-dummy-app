package external

import (
	"log"
	"net/http"
	"time"
)

func init() {
	seedTime := time.Now().UTC()
	fixtures := buildFixtures(seedTime, fixtureConfigs)
	server := newStaticDataServer(fixtures)
	decisionProvider := newDecisionProvider()
	publisher := newRandomLiveScorePublisher(fixtures, time.Second*1, decisionProvider, 10)
	publisher.StartRandomPublish()

	// 1- TODO: Having all routes declared on the same place
	// 2- TODO: not use localhost
	// 3- TODO: https (TLS)
	http.HandleFunc("/fixtures", server.HandleFixturesRequest)

	log.Println("test server listening at: http://localhost:8080")

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
}

type fixtureConfiguration struct {
	offset time.Duration
	values []string
}

var fixtureConfigs = []*fixtureConfiguration{
	{
		offset: -15 * time.Minute,
		values: []string{"F1", "Title1", "TO1", "Tournament1", "TE1", "Team1", "TE2", "Team2"},
	},
	{
		offset: -5 * time.Minute,
		values: []string{"F2", "Title1", "TO1", "Tournament1", "TE3", "Team3", "TE4", "Team4"},
	},
	{
		offset: -1 * time.Minute,
		values: []string{"F3", "Title2", "TO2", "Tournament2", "TE5", "Team5", "TE6", "Team6", "TE7", "Team7", "TE8", "Team8"},
	},
	{
		offset: 85 * time.Minute,
		values: []string{"F4", "Title1", "TO1", "Tournament1", "TE2", "Team2", "TE3", "Team3"},
	},
}

func buildFixtures(seedTime time.Time, fixtureConfigs []*fixtureConfiguration) []*fixture {

	fixtures := make([]*fixture, 0)

	for _, config := range fixtureConfigs {
		fixtures = append(fixtures, buildFixture(seedTime, config))
	}

	return fixtures
}

func buildFixture(seedTime time.Time, fixtureConfig *fixtureConfiguration) *fixture {

	teams := make([]fixtureTeam, 0)
	vals := fixtureConfig.values

	for i := 4; i < len(vals); i += 2 {
		teams = append(teams, fixtureTeam{
			Id:   vals[i],
			Name: vals[i+1],
		})
	}

	return &fixture{
		Id:    vals[0],
		Title: vals[1],
		Tournament: fixtureTournament{
			Id:   vals[2],
			Name: vals[3],
		},
		Teams:              teams,
		ScheduledStartTime: seedTime.Add(fixtureConfig.offset).Unix(),
	}
}
