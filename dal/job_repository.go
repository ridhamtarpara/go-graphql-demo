package dal

import (
	"context"
	"fmt"
	"firebase.google.com/go/db"
	"github.com/ridhamtarpara/go-graphql-jobs/dal/firebase"
)

type JobFirebaseRepository struct {
	Conn    *db.Client
	Context context.Context
}

var positionPath = firebase.PositionsPath

func NewJobFirebaseRepository(db *db.Client, ctx context.Context) (*JobFirebaseRepository, error) {
	return &JobFirebaseRepository{
		Conn:    db,
		Context: ctx,
	}, nil
}

func (jobRepo JobFirebaseRepository) GetAll() ([]Job, error) {
	var results map[string]Job
	err := jobRepo.Conn.NewRef(positionPath).Get(jobRepo.Context, &results)
	if err != nil {
		fmt.Printf("firebase error", err)
		return nil, err
	}

	var jobs []Job
	for _, r := range results {
		jobs = append(jobs, r)
	}
	return jobs, nil
}

func (jobRepo JobFirebaseRepository) Insert(newJob Job) error {
	err := jobRepo.Conn.NewRef(positionPath).Child(newJob.ID).Set(jobRepo.Context, newJob)
	return err
}

func (jobRepo JobFirebaseRepository) GetID() (string, error) {
	ID, err := GetID(jobRepo.Conn, jobRepo.Context, positionPath)
	return ID, err
}