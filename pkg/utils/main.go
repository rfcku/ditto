package utils

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DateToString(date primitive.DateTime) string {
	d := date.Time()
	currentTime := time.Now()
	diff := currentTime.Sub(d)
	days := int(diff.Hours() / 24)
	hours := int(diff.Hours())
	minutes := int(diff.Minutes())
	seconds := int(diff.Seconds())

	if days > 0 && days < 7 { // 7 days in a week
		return fmt.Sprintf("%d days ago", days)
	} else if days > 7 && days < 30 { // weeks in a month
		return fmt.Sprintf("%d weeks ago", days/7)
	} else if days > 30 && days < 365 { // days in a year
		return fmt.Sprintf("%d months ago", days/30)
	} else if days > 365 {
		return fmt.Sprintf("%d years ago", days/365)
	} else if hours > 0 {
		return fmt.Sprintf("%d hours ago", hours)
	} else if minutes > 0 {
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if seconds > 0 {
		return fmt.Sprintf("%d seconds ago", seconds)
	}
	return date.Time().String()
}