package main

import (
	"testing"
)

func TestReplaceUnderscore(t *testing.T) {
	var testSlice []string

	testSlice = append(testSlice, "hello_world")
	titleKey := "title"
	testSlice = replaceUnderscore(titleKey, testSlice)
	if testSlice[0] != "hello world" {
		t.Errorf("replaceUnderscore str = '%v'; want 'hello world'", testSlice[0])
	}

	testSlice = nil
	testSlice = append(testSlice, "hello_world")
	emailKey := "email"
	testSlice = replaceUnderscore(emailKey, testSlice)
	if testSlice[0] != "hello_world" {
		t.Errorf("replaceUnderscore str = '%v'; want 'hello_world'", testSlice[0])
	}
}
