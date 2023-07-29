package postformator

import (
	"fmt"
	"regexp"
	"strings"
)

func ConvertTags(tags string) []string {
	if tags != "" {
		tagsArray := strings.Split(tags, ",")
		for i, tag := range tagsArray {
			tagsArray[i] = strings.TrimSpace(tag)
		}
		return tagsArray
	}
	return nil
}

func CreatePostPermalinks(title string) string {
	// Convert the title to lowercase
	permalinks := strings.ToLower(title)

	// Replace spaces with dashes
	permalinks = strings.ReplaceAll(permalinks, " ", "-")

	// Remove all characters except letters, digits, underscores, and dashes
	reg := regexp.MustCompile(`[^\w-]+`)
	permalinks = reg.ReplaceAllString(permalinks, "")

	return permalinks
}

func CalculateReadingTime(content string) string {
	// Calculate the number of words in the content
	words := strings.Fields(content)
	wordCount := len(words)

	// Assuming an average reading speed of 225 words per minute
	// You can adjust this value based on your preference or requirement
	wordsPerMinute := 225

	// Calculate the estimated reading time in minutes
	minutes := wordCount / wordsPerMinute

	// Format the reading time as a string (e.g., "5 min read")
	readingTime := fmt.Sprintf("%d min read", minutes)

	return readingTime
}
