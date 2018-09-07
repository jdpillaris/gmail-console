package main

import (
	"app/controllers"
	"bufio"
	"os"
	"fmt"
	"strings"
	"time"
)

func init() {

}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Start of the retrieval period (DD-MM-YYYY): ")
	startDateStr, _ := reader.ReadString('\n')
	startDateStr = strings.TrimSuffix(startDateStr, "\n")
	fmt.Print("End of the retrieval period (DD-MM-YYYY): ")
	endDateStr, _ := reader.ReadString('\n')
	endDateStr = strings.TrimSuffix(endDateStr, "\n")

	for !(isValid(startDateStr, endDateStr)) {
		fmt.Print("Start of the retrieval period (DD-MM-YYYY): ")
		startDateStr, _ = reader.ReadString('\n')
		startDateStr = strings.TrimSuffix(startDateStr, "\n")
		fmt.Print("End of the retrieval period (DD-MM-YYYY): ")
		endDateStr, _ = reader.ReadString('\n')
		endDateStr = strings.TrimSuffix(endDateStr, "\n")
	}

	startDate, _ := time.Parse("02-01-2006", startDateStr)
	endDate, _ := time.Parse("02-01-2006", endDateStr)

	inbox := new(controllers.Inbox)
	inbox.Create()
	inbox.StoreMessages(startDate, endDate)
}

// isValid returns true. Assume that the correct start/end dates are entered 
func isValid(begin, end string) bool {
	return true
}
