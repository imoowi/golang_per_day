package main

import (
	"fmt"

	"golang.org/x/text/language"
)

func TriggerVuln() {
	tag, err := language.Parse("en-US")
	if err != nil {
		panic(err)
	}
	fmt.Println(tag)
}
