package utils

import (
	"fmt"
	"time"
)

func LogUserAction(action string, userID int) {
	logEntry := fmt.Sprintf("[%s] ACTION: %s | UserID: %d", time.Now().Format(time.RFC3339), action, userID)
	fmt.Println(logEntry)
}
