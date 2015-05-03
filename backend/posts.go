package posts

import (
	"encoding/json"

	"net/http"
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
	enc := json.NewEncoder(w)
	enc.Encode(posts)
}
