package main

import (
	"testing"

	"umbrella.github.com/go-kit-example/cha1/util"
)

func TestGenKey(t *testing.T) {
	err := util.GenPubandPriKey(1024, "./pem/")
	if err != nil {
		t.Fatal(err)
	}
}
