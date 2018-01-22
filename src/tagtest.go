package main

import (
"github.com/dhowden/tag"
"fmt"
"os"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", m.Format())

	tags := m.Raw()
	for k, v := range tags {
	fmt.Printf("%v: %v\n", k, v)
	}
	flac, _ := tag.ReadFLACTags(f)
	flactags := flac.Raw()
	for k, v := range flactags {
	fmt.Println("%s", k, v) 
	}
}