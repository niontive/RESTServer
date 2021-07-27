package main

import (
	"reflect"
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

func TestValidateAppMetaData(t *testing.T) {
	var data appMetaData
	var err error

	data.Title = "Test Title"
	data.Version = "1"
	data.Maintainers = append(data.Maintainers, maintainer{Name: "Test Name", Email: "test@test.com"})
	data.Company = "Test Company"
	data.Website = "https://www.test.com"
	data.Source = "https://github.com/test"
	data.License = "TestLicense"
	data.Description = "Test Description"

	err = validateAppMetaData(data)
	if err != nil {
		t.Errorf("Error validating ok metadata: %v", err)
	}

	data.Website = "test"
	err = validateAppMetaData(data)
	if err == nil {
		t.Error("Error validating incorrect metadata: website field")
	}

	data = appMetaData{}
	err = validateAppMetaData(data)
	if err == nil {
		t.Error("Error validting incorrect metadata: empty struct")
	}
}

func TestAppMetaDataStore(t *testing.T) {
	var dataStore = appMetaDataStore{store: make([]appMetaData, 0), dupTracker: make(map[string]bool)}
	var data appMetaData
	var err error

	data.Title = "Test Title"
	data.Version = "1"
	data.Maintainers = append(data.Maintainers, maintainer{Name: "Test Name", Email: "test@test.com"})
	data.Company = "Test Company"
	data.Website = "https://www.test.com"
	data.Source = "https://github.com/test"
	data.License = "TestLicense"
	data.Description = "Test Description"

	err = dataStore.Add(data)
	if err != nil {
		t.Errorf("Unable to add MD to data store: %v", err)
	}

	err = dataStore.Add(data)
	if err == nil {
		t.Error("Added duplicate metadata to data store")
	}

	if dataStore.TotalEntries() != 1 {
		t.Error("Incorrect number of data store entries")
	}

	searchMD := make([]string, 1)
	searchMD = append(searchMD, "Test Title")
	retrievedMD, err := dataStore.Search("title", searchMD)
	if err != nil {
		t.Errorf("Error in data store search: %v", err)
	}
	if !reflect.DeepEqual(retrievedMD[0], data) {
		t.Error("Retrieved MD does not match added MD")
	}
}
