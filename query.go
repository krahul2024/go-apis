package main

import (
	"api/utils"
	"fmt"
)

func DisplayQueryResults(query string) {
	rows, err := DB.Query(query)
	utils.FatalError(err)
	defer rows.Close()

	columns, err := rows.Columns()
	utils.FatalError(err)

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	fmt.Println("Query Results:")
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		utils.FatalError(err)

		for i, column := range columns {
			fmt.Printf("%s: %v\t", column, values[i])
		}
		fmt.Println()
	}

	utils.FatalError(err)
}
