package helpers

import (
	"errors"
	"fmt"
	"time"
)

// LogInfo function is used to print formatted information logs to the console.
func LogInfo(message string) (i int, err error) {
	return fmt.Println(formatLog(message))
}

// LogError function is used to print formatted error logs to the console.
// Second return value returns the error occured while using fmt.Println().
func LogError(message string) (n int, err error) {
	errMsg := errors.New(formatLog(message))
	return fmt.Println(errMsg)
}

// Format log string.
func formatLog(str string) string {
	currentTime := time.Now()
	dateText := fmt.Sprintf(
		"%d/%d/%d %d:%d:%d ",
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		currentTime.Hour(),
		currentTime.Hour(),
		currentTime.Second(),
	)
	return dateText + str
}
