package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/chengxuncc/shutdownhook"
)

var timeLayout = "2006-01-02 15:04:05  -0700 MST"

func main() {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("无法获取当前用户信息：", err)
		return
	}
	filePath := filepath.Join(currentUser.HomeDir, "clock-out-reminder.log")

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("无法打开文件：", err)
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Printf("%s run\n", time.Now().Format(timeLayout))
	err = shutdownhook.New(func() {
		ok := isTimeToClockOut()
		key := os.Getenv("BARK_KEY")
		log.Printf("%s shutdown\n", time.Now().Format(timeLayout))

		if ok == true {
			http.Get(fmt.Sprintf("https://api.day.app/%s/通知/您似乎关电脑下班了，记得打卡~?isArchive=0&level=timeSensitive", key))
		}
		f.Close()
	})
	if err != nil {
		log.Fatal("shutdownhook err:", err)
	}
}

func isTimeToClockOut() bool {
	now := time.Now()
	y, m, d := now.Date()
	startTime, err1 := time.Parse(timeLayout, fmt.Sprintf("%04d-%02d-%02d 17:00:00 +0800 CST", y, m, d))
	endTime, err2 := time.Parse(timeLayout, fmt.Sprintf("%04d-%02d-%02d 23:59:59 +0800 CST", y, m, d))

	if err1 == nil && err2 == nil && (now.After(startTime) && now.Before(endTime)) {
		return true
	}
	return false
}
