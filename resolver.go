package go_graphql_demo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/markbates/going/randx"
	"github.com/ridhamtarpara/go-graphql-demo/api"
	"github.com/ridhamtarpara/go-graphql-demo/api/dal"
	"github.com/ridhamtarpara/go-graphql-demo/api/dataloaders"
	"github.com/ridhamtarpara/go-graphql-demo/api/errors"
	"strconv"
	"time"
)

var videoPublishedChannel map[string]chan api.Video

func init() {
	videoPublishedChannel = map[string]chan api.Video{}
}

type contextKey string

var (
	UserIDCtxKey = contextKey("userID")
)

type Resolver struct {
	db *sql.DB
}

func NewRootResolvers(db *sql.DB) Config {
	c := Config{
		Resolvers: &Resolver{
			db: db,
		},
	}

	// Complexity
	countComplexity := func(childComplexity int, limit *int, offset *int) int {
		return *limit * childComplexity
	}
	c.Complexity.Query.Videos = countComplexity
	c.Complexity.Video.Related = countComplexity

	// Schema Directive
	c.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		ctxUserID := ctx.Value(UserIDCtxKey)
		if ctxUserID != nil {
			return next(ctx)
		} else {
			return nil, errors.UnauthorisedError
		}
	}
	return c
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Review() ReviewResolver {
	return &reviewResolver{r}
}
func (r *Resolver) Video() VideoResolver {
	return &videoResolver{r}
}
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateVideo(ctx context.Context, input NewVideo) (api.Video, error) {
	userID := UserFromContext(ctx)
	newVideo := api.Video{
		URL:         input.URL,
		Description: input.Description,
		Name:        input.Name,
		CreatedAt:   time.Now().UTC(),
		UserID:      userID,
	}

	rows, err := dal.LogAndQuery(r.db, "INSERT INTO videos (name, description, url, user_id, created_at) VALUES($1, $2, $3, $4, $5) RETURNING id",
		input.Name, input.Description, input.URL, userID, newVideo.CreatedAt)
	if err != nil || !rows.Next() {
		return api.Video{}, err
	}
	defer rows.Close()
	
	if err := rows.Scan(&newVideo.ID); err != nil {
		errors.DebugPrintf(err)
		if errors.IsForeignKeyError(err) {
			return api.Video{}, errors.UserNotExist
		}
		return api.Video{}, errors.InternalServerError
	}

	for _, observer := range videoPublishedChannel {
		observer <- newVideo
	}

	return newVideo, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Videos(ctx context.Context, limit *int, offset *int) ([]api.Video, error) {
	var video api.Video
	var videos []api.Video

	rows, err := dal.LogAndQuery(r.db, "SELECT id, name, description, url, created_at, user_id FROM videos ORDER BY created_at desc limit $1 offset $2", limit, offset)
	defer rows.Close();
	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}
	for rows.Next() {
		if err := rows.Scan(&video.ID, &video.Name, &video.Description, &video.URL, &video.CreatedAt, &video.UserID); err != nil {
			errors.DebugPrintf(err)
			return nil, errors.InternalServerError
		}
		videos = append(videos, video)
	}

	return videos, nil
}

type reviewResolver struct{ *Resolver }

func (r *reviewResolver) User(ctx context.Context, obj *api.Review) (api.User, error) {
	panic("not implemented")
}

type videoResolver struct{ *Resolver }

func (r *videoResolver) User(ctx context.Context, obj *api.Video) (api.User, error) {
	// Raw

	//rows, _ := dal.LogAndQuery(r.db, "SELECT id, name, email FROM users where id = $1", obj.UserID)
	//defer rows.Close()
	//
	//if !rows.Next() {
	//	return api.User{}, nil
	//}
	//var user api.User
	//if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
	//	errors.DebugPrintf(err)
	//	return api.User{}, errors.InternalServerError
	//}
	//
	//return user, nil

	// DataLoader
	user, err := ctx.Value(dataloaders.CtxKey).(*dataloaders.UserLoader).Load(obj.UserID)
	return *user, err
	//return api.User{}, nil
}
func (r *videoResolver) Screenshots(ctx context.Context, obj *api.Video) ([]*api.Screenshot, error) {
	panic("not implemented")
}

func (r *videoResolver) Related(ctx context.Context, obj *api.Video, limit *int, offset *int) ([]api.Video, error) {
	fmt.Println("Related Called")
	return []api.Video{}, nil
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) VideoPublished(ctx context.Context) (<-chan api.Video, error) {
	id := randx.String(8)

	videoEvent := make(chan api.Video, 1)
	go func() {
		<-ctx.Done()
	}()
	videoPublishedChannel[id] = videoEvent
	return videoEvent, nil
}

func UserFromContext(ctx context.Context) (int) {
	userIDStr, _ := ctx.Value(UserIDCtxKey).(string)
	userID, _ := strconv.Atoi(userIDStr)
	return userID
}
