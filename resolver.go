//go:generate gorunpkg github.com/99designs/gqlgen

package go_graphql_jobs

import (
	context "context"

	dal "github.com/ridhamtarpara/go-graphql-jobs/dal"
	graph "github.com/ridhamtarpara/go-graphql-jobs/graph"
	models "github.com/ridhamtarpara/go-graphql-jobs/models"
)

type Resolver struct{}

func (r *Resolver) Application() graph.ApplicationResolver {
	return &applicationResolver{r}
}
func (r *Resolver) Job() graph.JobResolver {
	return &jobResolver{r}
}
func (r *Resolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}

type applicationResolver struct{ *Resolver }

func (r *applicationResolver) Job(ctx context.Context, obj *dal.Application) (dal.Job, error) {
	panic("not implemented")
}

type jobResolver struct{ *Resolver }

func (r *jobResolver) CreatedBy(ctx context.Context, obj *dal.Job) (dal.User, error) {
	panic("not implemented")
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateJob(ctx context.Context, input models.NewJob) (dal.Job, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteJob(ctx context.Context, id string) (string, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateApplication(ctx context.Context, input models.NewApplication) (dal.Application, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Jobs(ctx context.Context) ([]dal.Job, error) {
	panic("not implemented")
}
func (r *queryResolver) Applications(ctx context.Context, jobID string) ([]dal.Application, error) {
	panic("not implemented")
}
