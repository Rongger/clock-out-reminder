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
	filePath := "clock-out-reminder.log"
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		panic(err)
	}

	fmt.Fprintf(f, "%s run\n", time.Now().Format(timeLayout))
	err = shutdownhook.New(func() {
		ok := isTimeToClockOut()
		key := os.Getenv("BARK_KEY")
		fmt.Fprintf(f, "%s shutdown\n", time.Now().Format(timeLayout))

		if ok == true {
			http.Get(fmt.Sprintf("https://api.day.app/%s/通知/您似乎关电脑下班了，记得打卡~", key))
		}
		f.Close()
	})
	if err != nil {
		fmt.Fprintf(f, "shutdownhook err: %s\n", err)
		panic(err)
	}
}

func isTimeToClockOut() bool {
	now := time.Now()
	y, m, d := now.Date()
	startTime, err1 := time.Parse(timeLayout, fmt.Sprintf("%04d-%02d-%02d 13:00:00 +0800 CST", y, m, d))
	endTime, err2 := time.Parse(timeLayout, fmt.Sprintf("%04d-%02d-%02d 23:59:59 +0800 CST", y, m, d))

	if err1 == nil && err2 == nil && (now.After(startTime) && now.Before(endTime)) {
		return true
	}
	return false
}
