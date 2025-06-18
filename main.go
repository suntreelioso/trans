package main

import (
	"log"
	"trans/cmd"
)

func main() {
	if err := cmd.Cmd.Execute(); err != nil {
		log.Printf("error: %v", err.Error())
		return
	}
}
