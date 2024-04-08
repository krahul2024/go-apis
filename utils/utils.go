package utils

import (
	"log"
	"net/http"
)

func PanicOnError(err error, messages ...string) {
	if err != nil {
		message := ""
		if len(messages) > 0 {
			message = messages[0]
		}
		log.Panicf("%s\n%v\n", message, err)
	}
}

func FatalError(err error, messages ...string) {
	if err != nil {
		message := ""
		if len(messages) > 0 {
			message = messages[0]
		}
		log.Fatalf("%s\n%v\n", message, err)
	}
}

func PrintError(err error, messages ...string) {
	if err != nil {
		message := ""
		if len(messages) > 0 {
			message = messages[0]
		}
		log.Printf("%s\n%v\n", message, err)
	}
}

func HttpResponseError(res http.ResponseWriter, err error, statusCode int, messages ...string) {
	if err != nil {
		message := ""
		if len(messages) > 0 {
			message = messages[0]
		} else {
			message = http.StatusText(statusCode)
		}

		PrintError(err, message)
		http.Error(res, message, statusCode)
		return
	}
}
