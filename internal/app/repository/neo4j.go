package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/charlysan/goneo4jgql/internal/app/graph/model"
	"github.com/charlysan/goneo4jgql/internal/app/models"
	"github.com/charlysan/goneo4jgql/pkg/logger"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/spf13/viper"
)

// NewNeo4jConnection creates a new neo4j connection
func NewNeo4jConnection() (neo4j.Driver, error) {
	target := fmt.Sprintf("%s://%s:%d", viper.GetString("NEO4J_PROTO"), viper.GetString("NEO4J_HOST"), viper.GetInt("NEO4J_PORT"))

	driver, err := neo4j.NewDriver(
		target,
		neo4j.BasicAuth(viper.GetString("NEO4J_USER"), viper.GetString("NEO4J_PASS"), ""),
		func(c *neo4j.Config) {
			c.Encrypted = false
		})
	if err != nil {
		logger.Error("Cannot connect to Neo4j Server", err)
		return nil, err
	}

	logger.Info("Connected to Neo4j Server", logger.LogFields{"neo4j_server_uri": target})

	return driver, nil
}

// Neo4jRepository is a Neo4j DB repository
type Neo4jRepository struct {
	Connection neo4j.Driver
}

// FindMovieByUUID finds a movie by its uuid
func (r *Neo4jRepository) FindMovieByUUID(ctx context.Context, uuid string) (*models.Movie, error) {
	query := `
		match (m:Movie) where m.uuid = $uuid return m.uuid, m.title, m.released, m.tagline
	`
	session, err := r.Connection.Session(neo4j.AccessModeWrite)

	if err != nil {
		return nil, err
	}

	defer session.Close()

	args := map[string]interface{}{
		"uuid": uuid,
	}

	result, err := session.Run(query, args)

	if err != nil {
		logger.Error("Cannot find movie by uuid", logger.LogFields{"uuid": uuid}, err)
	}

	logger.Debug("CYPHER_QUERY", logger.LogFields{"query": query, "args": args})

	movie := models.Movie{}

	for result.Next() {
		ParseCypherQueryResult(result.Record(), "m", &movie)
	}

	return &movie, err
}

// FindMovies finds movies by title and actor
func (r *Neo4jRepository) FindMovies(ctx context.Context, title *string, actor *string) ([]*models.Movie, error) {
	movieTitle := ""
	actorName := ""

	query := `
		match (m:Movie) return m.uuid, m.title, m.released, m.tagline
	`

	if title != nil {
		query = `
			match (m:Movie) where lower(m.title) contains $movieTitle return m.uuid, m.title, m.released, m.tagline
		`
		movieTitle = *title
	}

	if actor != nil {
		query = `
			match (m:Movie)-[r:ACTED_IN]-(p:Person) where lower(p.name) contains $actor return m.uuid, m.title, m.released, m.tagline
		`
		actorName = *actor
	}

	if title != nil && actor != nil {
		query = `
			match (m:Movie)-[r:ACTED_IN]-(p:Person) where lower(m.title) contains $movieTitle and lower(p.name) contains $actor return m.uuid, m.title, m.released, m.tagline
		`
		movieTitle = *title
		actorName = *actor
	}

	session, err := r.Connection.Session(neo4j.AccessModeWrite)

	if err != nil {
		return nil, err
	}

	defer session.Close()

	args := map[string]interface{}{
		"movieTitle": strings.ToLower(movieTitle),
		"actor":      strings.ToLower(actorName),
	}

	result, err := session.Run(query, args)
	if err != nil {
		logger.Error("Cannot find movies", err)
	}

	logger.Debug("CYPHER_QUERY", logger.LogFields{"query": query, "args": args})

	var movies []*models.Movie

	for result.Next() {
		movie := models.Movie{}
		ParseCypherQueryResult(result.Record(), "m", &movie)

		movies = append(movies, &movie)
	}

	return movies, err
}

// FindMovieParticipationsByPersonUUID finds people that participated in a movie
func (r *Neo4jRepository) FindMovieParticipationsByPersonUUID(ctx context.Context, uuid string) ([]*model.Participation, error) {
	query := `
		match (m:Movie)-[relatedTo]-(p:Person) where p.uuid = $uuid return m.uuid, m.title, m.released, m.tagline, type(relatedTo) as role
	`
	session, err := r.Connection.Session(neo4j.AccessModeWrite)

	if err != nil {
		return nil, err
	}

	defer session.Close()

	args := map[string]interface{}{
		"uuid": uuid,
	}

	result, err := session.Run(query, args)
	if err != nil {
		logger.Error("Cannot find movies", err)
	}

	logger.Debug("CYPHER_QUERY", logger.LogFields{"query": query, "args": args})

	var participations []*model.Participation

	for result.Next() {
		movie := models.Movie{}
		ParseCypherQueryResult(result.Record(), "m", &movie)
		participation := model.Participation{
			Movie: &movie,
		}
		// Append Role
		if role, ok := result.Record().Get("role"); ok {
			participation.Role = role.(string)
		}

		participations = append(participations, &participation)
	}

	return participations, err
}

// FindPersonByMovieUUID finds people (actors, directors, writers) by movie uuid
func (r *Neo4jRepository) FindPersonByMovieUUID(ctx context.Context, role string, uuid string) ([]*models.Person, error) {
	query := `
		match (p:Person)-[:%s]->(m:Movie)  where m.uuid = $uuid return p.uuid, p.name, p.born
	`
	query = fmt.Sprintf(query, role)

	session, err := r.Connection.Session(neo4j.AccessModeWrite)

	if err != nil {
		return nil, err
	}

	defer session.Close()

	args := map[string]interface{}{
		"uuid": uuid,
		"role": role,
	}

	result, err := session.Run(query, args)
	if err != nil {
		logger.Error("Cannot find any person with that role", err, logger.LogFields{"role": role})
	}

	logger.Debug("CYPHER_QUERY", logger.LogFields{"query": query, "args": args})

	var people []*models.Person

	for result.Next() {
		person := models.Person{}
		ParseCypherQueryResult(result.Record(), "p", &person)
		// Append Role
		person.Role = StringPtr(role)

		people = append(people, &person)
	}

	return people, nil
}
