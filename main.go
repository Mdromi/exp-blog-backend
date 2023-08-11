package main

import "github.com/Mdromi/exp-blog-backend/api"

func main() {

	api.Run()

}

// git log --author="mdromi" --oneline --shortstat | grep "files\? changed" | awk '{files+=$1; inserted+=$4; deleted+=$6} END {print "Files changed:", files, "Inserted lines:", inserted, "Deleted lines:", deleted}'
// fmt.Printf("Post: %+v\n", post)
