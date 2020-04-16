package service

import (
	"context"

	"github.com/charlysan/goneo4jgql/graph/model"
	"github.com/charlysan/goneo4jgql/internal/app/models"
	"github.com/charlysan/goneo4jgql/internal/app/repository"
)

// Service exposes application bussiness logic
type Service struct {
	repository repository.Repository
}

// NewService creates a new service
func NewService(r repository.Repository) Service {
	return Service{
		repository: r,
	}
}

// FindMovieByUUID finds a movie by its uuid
func (s *Service) FindMovieByUUID(ctx context.Context, uuid string) (*models.Movie, error) {
	return s.repository.FindMovieByUUID(ctx, uuid)
}

// FindMovies finds movies by title and actor
func (s *Service) FindMovies(ctx context.Context, title *string, actor *string) ([]*models.Movie, error) {
	return s.repository.FindMovies(ctx, title, actor)
}

// FindDirectorsByMovieUUID finds directors for a movie by movie uuid
func (s *Service) FindDirectorsByMovieUUID(ctx context.Context, uuid string) ([]*models.Person, error) {
	return s.repository.FindPersonByMovieUUID(ctx, "DIRECTED", uuid)
}

// FindWritersByMovieUUID finds writers for a movie by movie uuid
func (s *Service) FindWritersByMovieUUID(ctx context.Context, uuid string) ([]*models.Person, error) {
	return s.repository.FindPersonByMovieUUID(ctx, "WROTE", uuid)
}

// FindCastByMovieUUID finds movie cast by movie uuid
func (s *Service) FindCastByMovieUUID(ctx context.Context, uuid string) ([]*models.Person, error) {
	return s.repository.FindPersonByMovieUUID(ctx, "ACTED_IN", uuid)
}

// FindMovieParticipationsByPersonUUID finds people that participated in a movie
func (s *Service) FindMovieParticipationsByPersonUUID(ctx context.Context, uuid string) ([]*model.Participation, error) {
	return s.repository.FindMovieParticipationsByPersonUUID(ctx, uuid)
}
