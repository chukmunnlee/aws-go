package main

import "log"

func CheckError(err error) {
	if nil == err {
		return
	}
	log.Fatalf("Error: %v\n", err)
}
