package main

import (
	"api/utils"
	"fmt"
	"os"
	"strings"
)

func SetEnvs() {
	data, err := os.ReadFile(".env")
	utils.PrintError(err)

	line := ""

	for i := range data {
		if data[i] != 10 && data[i] != 32 {
			line += string(data[i])
		}
		if data[i] == 10 || i == len(data)-1 {
			envs := strings.Split(line, "=")
			if len(envs) == 1 {
				envs = append(envs, "")
			}
			key := strings.TrimSpace(envs[0])
			value := strings.TrimSpace(envs[1])
			os.Setenv(key, value)
			// fmt.Println(key, value)
			line = ""
		}
	}
	fmt.Println("Value of the env vars set successfully!")
}
