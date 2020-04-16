package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strings"

	"github.com/charlysan/goneo4jgql/internal/app/graph/generated"
	"github.com/charlysan/goneo4jgql/internal/app/graph/model"
	"github.com/charlysan/goneo4jgql/internal/app/models"
	validator "github.com/go-playground/validator/v10"
)

func (r *movieResolver) Directors(ctx context.Context, obj *models.Movie) ([]*models.Person, error) {
	ds, err := r.Service.FindDirectorsByMovieUUID(ctx, obj.UUID)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func (r *movieResolver) Writers(ctx context.Context, obj *models.Movie) ([]*models.Person, error) {
	ws, err := r.Service.FindWritersByMovieUUID(ctx, obj.UUID)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (r *movieResolver) Cast(ctx context.Context, obj *models.Movie) ([]*models.Person, error) {
	c, err := r.Service.FindCastByMovieUUID(ctx, obj.UUID)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *personResolver) Participated(ctx context.Context, obj *models.Person) ([]*model.Participation, error) {
	p, err := r.Service.FindMovieParticipationsByPersonUUID(ctx, obj.UUID)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *queryResolver) Movie(ctx context.Context, uuid string) (*models.Movie, error) {
	mv, err := r.Service.FindMovieByUUID(ctx, uuid)

	if err != nil {
		return nil, err
	}

	return mv, nil
}

func (r *queryResolver) Movies(ctx context.Context, title *string, actor *string) ([]*models.Movie, error) {
	validator := validator.New()

	// validate input
	if title != nil {
		err := validator.Var(strings.Trim(*title, " "), "alphanum,min=3")
		if err != nil {
			return nil, err
		}
	}

	if actor != nil {
		err := validator.Var(strings.Trim(*actor, " "), "alphanum,min=3")
		if err != nil {
			return nil, err
		}
	}

	movies, err := r.Service.FindMovies(ctx, title, actor)

	if err != nil {
		return nil, err
	}

	return movies, nil
}

// Movie returns generated.MovieResolver implementation.
func (r *Resolver) Movie() generated.MovieResolver { return &movieResolver{r} }

// Person returns generated.PersonResolver implementation.
func (r *Resolver) Person() generated.PersonResolver { return &personResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type movieResolver struct{ *Resolver }
type personResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
