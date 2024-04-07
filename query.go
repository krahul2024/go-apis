package main

import (
	"fmt"
	"log"
)

func DisplayQueryResults(query string) {
	rows, err := DB.Query(query)
	if err != nil {
		log.Fatal("Error executing query:", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal("Error getting column names:", err)
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	fmt.Println("Query Results:")
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal("Error scanning row values:", err)
		}

		for i, column := range columns {
			fmt.Printf("%s: %v\t", column, values[i])
		}
		fmt.Println()
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating over rows:", err)
	}
}
