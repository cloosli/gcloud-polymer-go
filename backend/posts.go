package posts

import (
	"fmt"
	"log"

	"appengine"
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
	api, err := endpoints.RegisterService(PostAPI{}, "posts", "v1", "posts API", true)
	if err != nil {
		log.Fatal(err)
	}

	info := api.MethodByName("SetFavorite").Info()
	info.Name = "setFavorite"

	info = api.MethodByName("List").Info()
	info.Name = "getPosts"

	info = api.MethodByName("Add").Info()
	info.Name = "addPost"

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

type UpdateRequest struct {
	UID      *datastore.Key `json:"uid"`
	Favorite bool           `json:"favorite"`
}

func (PostAPI) SetFavorite(c endpoints.Context, r *UpdateRequest) error {
	return datastore.RunInTransaction(c, func(c appengine.Context) error {
		var post Post
		err := datastore.Get(c, r.UID, &post)
		if err != nil {
			return fmt.Errorf("get post: %v", err)
		}

		post.Favorite = r.Favorite
		_, err = datastore.Put(c, r.UID, &post)
		if err != nil {
			return fmt.Errorf("update post: %v", err)
		}
		return nil
	}, nil)
}
