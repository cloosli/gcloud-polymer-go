package posts

import (
	"encoding/json"
	"log"
	"net/http"

	"appengine"
	"appengine/datastore"
)

type Post struct {
	UID      string
	Text     string
	Username string
	Avatar   string
	Favorite bool
}

var posts = []Post{
	{
		"1",
		"Go is awesome",
		"Gopher",
		"http://www.unixstickers.com/image/cache/data/stickers/golang/golang.sh-600x600.png",
		true,
	},
	{
		"2",
		"And polymer is not bad",
		"Gopher",
		"http://magdkudama.com/assets/201407/gopher-flying.jpg",
		true,
	},
}

func init() {
	http.HandleFunc("/posts", listPosts)
}

func listPosts(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	posts := []Post{}
	_, err := datastore.NewQuery("Post").GetAll(c, &posts)
	if err != nil {
		c.Errorf("fetching posts: %v", err)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(posts)
	if err != nil {
		log.Printf("encoding: %v", err)
		c.Errorf("encoding: %v", err)
	}
}
