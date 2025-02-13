package internal

import (
	"fmt"
	"log"
	"text/template"
	"time"
)

var Templates *template.Template

func ParseTemplates() {
	var err error
	Templates, err = template.ParseGlob("./web/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	Templates, err = Templates.ParseGlob("./web/components/*.html")
	if err != nil {
		log.Fatal(err)
	}
}

func TimeFormatter(timestamp string) string {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return "Invalid time format"
	}

	// Calculate the difference
	now := time.Now().UTC()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		seconds := int(diff.Seconds())
		if seconds < 1 {
			return "just now"
		}
		return fmt.Sprintf("%d seconds ago", seconds)
	case diff < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	case diff < 7*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	case diff < 30*24*time.Hour:
		return fmt.Sprintf("%d weeks ago", int(diff.Hours()/168))
	case diff < 12*30*24*time.Hour:
		return fmt.Sprintf("%d months ago", int(diff.Hours()/720))
	default:
		return fmt.Sprintf("%d years ago", int(diff.Hours()/8760))
	}
}
