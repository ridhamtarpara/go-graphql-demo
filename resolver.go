package go_graphql_demo

import (
	"context"
	"github.com/ridhamtarpara/go-graphql-demo/api/dal"
	"github.com/ridhamtarpara/go-graphql-demo/api/errors"
	"time"

	"github.com/ridhamtarpara/go-graphql-demo/api"
)

type Resolver struct{}

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

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateVideo(ctx context.Context, input NewVideo) (api.Video, error) {
	currentDateTimeUTC := time.Now().Format("2006-01-02 15:04:05")
	row := dal.DBConn.QueryRow("INSERT INTO videos (name, description, url, user_id, created_at) VALUES($1, $2, $3, $4, $5) RETURNING id",
		input.Name, input.Description, input.URL, input.UserID, currentDateTimeUTC)
	newVideo := api.Video{
		URL:         input.URL,
		Description: input.Description,
		Name:        input.Name,
		CreatedAt:   currentDateTimeUTC,
	}
	if err := row.Scan(&newVideo.ID); err != nil {
		errors.DebugPrintf(err)
		if errors.IsForeignKeyError(err) {
			return api.Video{}, errors.UserNotExist
		}
		return api.Video{}, errors.InternalServerError
	}
	return newVideo, nil}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Videos(ctx context.Context, limit *int, offset *int) ([]api.Video, error) {
	var video api.Video
	var videos []api.Video
	rows, err := dal.DBConn.Query("SELECT id, name, description, url, created_at, user_id FROM videos ORDER BY created_at desc limit $1 offset $2", limit, offset)
	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}
	defer func(){
		if err := rows.Close();err != nil{
			errors.DebugPrintf(err)
		}
	}()
	for rows.Next() {
		if err := rows.Scan(&video.ID, &video.Name, &video.Description, &video.URL,&video.CreatedAt, &video.UserID); err != nil {
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
	var user api.User
	row := dal.DBConn.QueryRow("SELECT id, name, email FROM persons where id = $1", obj.UserID)
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		errors.DebugPrintf(err)
		return api.User{}, errors.InternalServerError
	}

	return user, nil
}
func (r *videoResolver) Screenshots(ctx context.Context, obj *api.Video) ([]*api.Screenshot, error) {
	panic("not implemented")
}
