package main

import (
	"wxbot/bot"
)

//打包命令:GOOS=linux GOARCH=amd64 go build -o wxbot

func main() {
	bot.Bot()
}
