package main

import "strings"

func isAlreadyFinished(err error) bool {
	return strings.EqualFold(err.Error(), "os: process already finished")
}
