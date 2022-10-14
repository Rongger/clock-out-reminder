package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/chengxuncc/shutdownhook"
)

var timeLayout = "2006-01-02 15:04:05  -0700 MST"

func main() {
	f, err := os.Create("clock-out-reminder.log")
	if err != nil {
		panic(err)
	}

	err = shutdownhook.New(func() {
		fmt.Fprintf(f, "%s shutdown\n", time.Now())
		f.Close()

		key := os.Getenv("BARK_KEY")

		ok := isTimeToClockOut()
		if ok == true {
			http.Get(fmt.Sprintf("https://api.day.app/%s/通知/您似乎关电脑下班了，记得打卡~", key))
		}
	})
	if err != nil {
		panic(err)
	}
}

func isTimeToClockOut() bool {
	now := time.Now()
	y, m, d := now.Date()
	startTime, err1 := time.Parse(timeLayout, fmt.Sprintf("%d-%d-%d 19:00:00 +0800 CST", y, m, d))
	endTime, err2 := time.Parse(timeLayout, fmt.Sprintf("%d-%d-%d 23:59:59 +0800 CST", y, m, d))

	if err1 == nil && err2 == nil && (now.After(startTime) && now.Before(endTime)) {
		return true
	}
	return false
}
