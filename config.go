package main

import (
	"flag"
	"os"
	"strconv"
)

func GetConfig(key string) string {
	return os.Getenv(key);
}

func SetConfig(key string, value string) {
	os.Setenv(key, value);
}

func HandleFlags() {
	portF := flag.Int("p", 0, "port number")
	userNameF := flag.String("n", "", "username")
	chatRoomNameF := flag.String("r", "general", "chatroom name")
	logsF := flag.Bool("logs", false, "chatroom name")

	flag.Parse()
	
	SetConfig("port", strconv.Itoa(*portF))
	SetConfig("userName", *userNameF)
	SetConfig("chatRoomName", *chatRoomNameF)
	SetConfig("logs", strconv.FormatBool(*logsF))
}