package main

import (
	"log"
)

func main() {
	entrypoint("ap-southeast-2")
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
