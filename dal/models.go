package dal

import (
	"time"
	"github.com/99designs/gqlgen/graphql"
	"io"
	"strconv"
	"errors"
)

type Application struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	CoverLetter *string `json:"coverLetter"`
	CvURL       string  `json:"cvURL"`
	LinkedInURL *string `json:"linkedInURL"`
	JobId       string  `json:"job"`
	CreatedAt   time.Time  `json:"createdAt"`
}
type Job struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	CreatedBy   string  `json:"createdBy"`
	CreatedAt   time.Time  `json:"createdAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func MarshalTimestamp(t time.Time) graphql.Marshaler {
	timestamp := t.Unix()
	if timestamp < 0 {
		timestamp = 0
	}
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.FormatInt(timestamp, 10))
	})
}

func UnmarshalTimestamp(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(int); ok {
		return time.Unix(int64(tmpStr), 0), nil
	}
	return time.Time{}, errors.New("time should be a unix timestamp")
}
