package main

import (
	"net/http"

	"github.com/Mdromi/exp-blog-backend/api"
)

func main() {
	// Serve static files from the 'static' directory
	staticDir := "/static/"
	http.Handle(staticDir, http.StripPrefix(staticDir, http.FileServer(http.Dir("static"))))

	api.Run()

	// Log to check if static file serving is working
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
}

// git log --author="mdromi" --oneline --shortstat | grep "files\? changed" | awk '{files+=$1; inserted+=$4; deleted+=$6} END {print "Files changed:", files, "Inserted lines:", inserted, "Deleted lines:", deleted}'
// fmt.Printf("Post: %+v\n", post)
