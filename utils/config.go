package utils

import (
	"flag"
	"os"
	"strconv"
)

// Gets the environment variable with [key]
func GetConfig(key string) string {
	return os.Getenv(key);
}

// Sets the environment variable with [key] and [value]
func SetConfig(key string, value string) {
	os.Setenv(key, value);
}

// Handles the flags passed when the program starts
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