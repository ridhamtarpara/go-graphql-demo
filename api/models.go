package api

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"io"
	"strconv"
)

type Review struct {
	ID          int64  `json:"id"`
	VideoID   int64  `json:"videoId"`
	Description string `json:"description"`
	Rating      int    `json:"rating"`
	CreatedAt   string `json:"createdAt"`
}

type Screenshot struct {
	ID        int64  `json:"id"`
	VideoID int64  `json:"videoId"`
	URL       string `json:"url"`
}

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Video struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	UserID      int64          `json:"-"`
	URL         string        `json:"url"`
	CreatedAt   string        `json:"createdAt"`
}

// Lets redefine the base ID type to use an id from an external library
func MarshalID(id int64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(fmt.Sprintf("%d", id)))
	})
}

// And the same for the unmarshaler
func UnmarshalID(v interface{}) (int64, error) {
	id, ok := v.(string)
	if !ok {
		return 0, fmt.Errorf("ids must be strings")
	}
	i, e := strconv.Atoi(id)
	return int64(i), e
}
