//go:build ignore

package main

import (
	"embed"
	"fmt"
)

//go:embed ../../internal/postgres/embed/*
var testFS embed.FS

func main() {
	entries, err := testFS.ReadDir("internal/postgres/embed")
	if err != nil {
		fmt.Println("ReadDir error:", err)
		entries, err = testFS.ReadDir(".")
		if err != nil {
			fmt.Println("Root ReadDir error:", err)
			return
		}
	}
	for _, e := range entries {
		info, _ := e.Info()
		fmt.Printf("%s - %d bytes\n", e.Name(), info.Size())
	}
}
