package main

import (
	"testing"
)

func TestGet(t *testing.T) {
	key := "name"
	value := "testing "
	SET(store, key, value)
	value2 := GET(store, key)
	if value2 != value {
		t.Error("the get function is not working ")
	}
}
func TestSet(t *testing.T) {
	key := "name"
	value := "aviral"
	SET(store, key, value)
	if GET(store, key) != value {
		t.Errorf("Result was incorrect, the set has some error ")
	}

}
