package main

import (
	"fmt"
	"testing"
)

func TestFooBar(t *testing.T) {
	in := map[string]string {
		"first": "1",
		"second": "2",
	}
	out := writeRespMap(in)
	fmt.Println(out)
}