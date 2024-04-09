package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func PanicOnError(err error, messages ...string) {
	if err != nil {
		message := ""
		if len(messages) > 0 {
			message = messages[0]
		}
		log.Panicf("%s\n%+v\n", message, err)
	}
}

func FatalError(err error, messages ...string) {
	if err != nil {
		message := ""
		if len(messages) > 0 {
			message = messages[0]
		}
		log.Fatalf("%s\n%+v\n", message, err)
	}
}

func PrintError(err error, messages ...string) string {
	if err != nil {
		message := ""
		if len(messages) > 0 {
			message = messages[0]
		}
		errorMessage := fmt.Sprintf("%s\n%+v\n", message, err)
		log.Print(errorMessage)
		return errorMessage
	}
	return ""
}

func HttpResponseError(res http.ResponseWriter, err error, statusCode int, txn *sql.Tx, messages ...string) {
	errorMessage := PrintError(err, messages...)
	errorMessage = strings.TrimSpace(errorMessage)
	message := ""
	if len(messages) > 0 {
		message = messages[0]
	} else {
		message = errorMessage
	}

	errorResponse := struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Message: message,
		Status:  statusCode,
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	err = json.NewEncoder(res).Encode(errorResponse)
	if err != nil {
		http.Error(res, "Internal Server Error!", http.StatusInternalServerError)
	}

	if txn != nil {
		txn.Rollback()
	}

}

func HttpResponseJson(res http.ResponseWriter, response interface{}, statusCode int) {
	/* Setting the headers before encoding the response body is important because the headers need to be sent first,
	before the response body is written. This allows the client to properly interpret the response.
	If we try to set the headers after encoding the response body, the headers will be ignored,
	and the client may not be able to properly interpret the response.
	*/
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	err := json.NewEncoder(res).Encode(response)
	if err != nil {
		HttpResponseError(res, err, http.StatusInternalServerError, nil)
		return
	}
}

func DBInsertUpdateDeleteHelper(res http.ResponseWriter, txn *sql.Tx, statement string, args ...interface{}) (int, int, error) {
	// prepare the statement
	query, err := txn.Prepare(statement)
	if err != nil {
		return 0, 0, err
	}
	defer query.Close()

	//execute the query
	result, err := query.Exec(args...)
	if err != nil {
		return 0, 0, err
	}

	lastRow, err := result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	numRows, err := result.RowsAffected()
	if err != nil {
		return 0, 0, err
	}

	return int(lastRow), int(numRows), nil
}
