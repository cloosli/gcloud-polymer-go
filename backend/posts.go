package posts

import (
	"fmt"

	"appengine/datastore"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
)

type PostAPI struct{}

type Post struct {
	UID      *datastore.Key `json:"uid" datastore:"-"`
	Text     string         `json:"text"`
	Username string         `json:"username"`
	Avatar   string         `json:"avatar"`
	Favorite bool           `json:"favorite"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

func init() {
	endpoints.RegisterService(PostAPI{}, "posts", "v1", "posts API", true)
	endpoints.HandleHTTP()
}

func (PostAPI) List(c endpoints.Context) (*Posts, error) {
	posts := []Post{}
	keys, err := datastore.NewQuery("Post").GetAll(c, &posts)
	if err != nil {
		return nil, err
	}

	for i, k := range keys {
		posts[i].UID = k
	}

	return &Posts{posts}, nil
}

type AddRequest struct {
	Text     string
	Username string
	Avatar   string
}

func (PostAPI) Add(c endpoints.Context, r *AddRequest) (*Post, error) {
	p := Post{
		Text:     r.Text,
		Username: r.Username,
		Avatar:   r.Avatar,
	}

	k := datastore.NewIncompleteKey(c, "Post", nil)
	k, err := datastore.Put(c, k, &p)
	if err != nil {
		return nil, fmt.Errorf("put post: %v", err)
	}
	p.UID = k
	return &p, nil
}
