//go:generate gorunpkg github.com/99designs/gqlgen

package api

import (
	"context"
	"github.com/ridhamtarpara/go-graphql-jobs/models"
	"github.com/ridhamtarpara/go-graphql-jobs/graph"
	"fmt"
	"time"
	"github.com/ridhamtarpara/go-graphql-jobs/dal"
	"github.com/ridhamtarpara/go-graphql-jobs/dal/firebase"
)

type Resolver struct{
	db *firebase.DBConnection
}

func (r *Resolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Application() graph.ApplicationResolver {
	return &applicationResolver{r}
}
func (r *Resolver) Job() graph.JobResolver {
	return &jobResolver{r}
}

// Mutations
type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateJob(ctx context.Context, input models.NewJob) (dal.Job, error) {
	jobRepository, err := dal.NewJobFirebaseRepository(r.db.Conn, r.db.Context)
	if err != nil {
		fmt.Printf("firebase error: ", err)
		return dal.Job{}, err
	}

	var jobID string
	if jobID, err = jobRepository.GetID(); err != nil {
		fmt.Printf("jobRepository GetID error: ", err)
		return dal.Job{}, err
	}
	// Create job object from request
	job := dal.Job{
		ID:          jobID,
		Name:        input.Name,
		Description: input.Description,
		Location: input.Location,
		CreatedBy: input.CreatedByID,
		CreatedAt:   time.Now().UTC(),
	}

	// Set the values in the DB
	if err = jobRepository.Insert(job); err != nil {
		fmt.Printf("firebase error: ", err)
		return dal.Job{}, err
	}

	return job, nil
}
func (r *mutationResolver) DeleteJob(ctx context.Context, id string) (string, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateApplication(ctx context.Context, input models.NewApplication) (dal.Application, error) {
	panic("not implemented")
}

// Queries
type queryResolver struct{ *Resolver }

func (r *queryResolver) Jobs(ctx context.Context) ([]dal.Job, error) {
	var allJobs []dal.Job

	jobRepository, err := dal.NewJobFirebaseRepository(r.db.Conn, r.db.Context)

	if allJobs, err = jobRepository.GetAll(); err != nil {
		fmt.Printf("firebase error", err)
	}

	return allJobs, nil
}
func (r *queryResolver) Applications(ctx context.Context, jobID string) ([]dal.Application, error) {
	panic("not implemented")
}

// Resolvers
type jobResolver struct{ *Resolver }

func (r *jobResolver) CreatedBy(ctx context.Context, obj *dal.Job) (dal.User, error) {
	userRepository, err := dal.NewUserFirebaseRepository(r.db.App, r.db.Context)
	if err != nil {
		fmt.Printf("Error Fetching NewTeamFirebaseRepository", err)
		return dal.User{}, err
	}
	user, err := userRepository.GetByID(obj.CreatedBy)
	if err != nil {
		fmt.Printf("Error Fetching user", err)
		return dal.User{}, err
	}
	return user, err
}

type applicationResolver struct{ *Resolver }

func (r *applicationResolver) Job(ctx context.Context, obj *dal.Application) (dal.Job, error) {
	panic("not implemented")
}

func NewRootResolvers() graph.Config {
	resolver := Resolver{}
	resolver.db = firebase.Connect()

	c := graph.Config{
		Resolvers: &resolver,
	}
	return c
}