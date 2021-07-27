package main

import (
	"errors"
	"net/mail"
	"net/url"
	"strings"
	"sync"
)

type maintainer struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type appMetaData struct {
	Title       string `yaml:"title"`
	Version     string `yaml:"version"`
	Maintainers []maintainer
	Company     string `yaml:"company"`
	Website     string `yaml:"website"`
	Source      string `yaml:"source"`
	License     string `yaml:"license"`
	Description string `yaml:"description"`
}

type appMetaDataStore struct {
	store      []appMetaData
	dupTracker map[string]bool
	storeLock  sync.RWMutex
	dupLock    sync.RWMutex
}

//
// Check if metadata already exists in store
//
func (mdStore *appMetaDataStore) checkDuplicate(md appMetaData) (err error) {
	mdStore.dupLock.RLock()
	if _, value := mdStore.dupTracker[md.Title]; !value {
		mdStore.dupLock.RUnlock()
		mdStore.dupLock.Lock()
		mdStore.dupTracker[md.Title] = true
		mdStore.dupLock.Unlock()
	} else {
		err = errors.New("Duplicate value")
		mdStore.dupLock.RUnlock()
	}
	return
}

//
// Add new metadata entry to store
//
func (mdStore *appMetaDataStore) Add(md appMetaData) (err error) {
	mdStore.storeLock.Lock()
	defer mdStore.storeLock.Unlock()
	err = mdStore.checkDuplicate(md)
	if err != nil {
		return
	}
	mdStore.store = append(mdStore.store, md)
	return
}

//
// Retrieve total number of entries in store
//
func (mdStore *appMetaDataStore) TotalEntries() int {
	mdStore.storeLock.RLock()
	defer mdStore.storeLock.RUnlock()
	return len(mdStore.store)
}

//
// Retrieve metadata app titles
//
func (mdStore *appMetaDataStore) GetAppTitles() (titles []string) {
	mdStore.storeLock.RLock()
	defer mdStore.storeLock.RUnlock()

	for _, element := range mdStore.store {
		titles = append(titles, element.Title)
	}

	return
}

//
// Search store and return entires that match the key/value pairs
//
func (mdStore *appMetaDataStore) Search(key string, values []string) (md []appMetaData, err error) {
	mdStore.storeLock.RLock()
	defer mdStore.storeLock.RUnlock()

	if mdStore.TotalEntries() == 0 {
		err = errors.New("MD store is empty")
		return md, err
	}

	//
	// Iterate through each store entry and check if
	// value matches for given key
	//
	for _, element := range mdStore.store {
		for _, v := range values {
			switch {
			case key == "title":
				if element.Title == v {
					md = append(md, element)
				}
			case key == "version":
				if element.Version == v {
					md = append(md, element)
				}
			case key == "name":
				for _, m := range element.Maintainers {
					if m.Name == v {
						md = append(md, element)
					}
				}
			case key == "email":
				for _, m := range element.Maintainers {
					if m.Email == v {
						md = append(md, element)
					}
				}
			case key == "company":
				if element.Company == v {
					md = append(md, element)
				}
			case key == "website":
				if element.Website == v {
					md = append(md, element)
				}
			case key == "source":
				if element.Source == v {
					md = append(md, element)
				}
			case key == "license":
				if element.License == v {
					md = append(md, element)
				}
			case key == "description":
				if element.License == v {
					md = append(md, element)
				}
			default:
				err = errors.New("Invalid key")
			}
		}
	}

	if len(md) == 0 {
		err = errors.New("No metadata found")
	}

	return md, err
}

func checkEmptyString(str string) error {
	if len(strings.TrimSpace(str)) == 0 {
		return errors.New("Empty string")
	}
	return nil
}

func checkEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func checkURL(URL string) error {
	_, err := url.ParseRequestURI(URL)
	return err
}

func validateTitle(title string) error {
	if checkEmptyString(title) != nil {
		return errors.New("Empty title")
	}
	return nil
}

func validateVersion(version string) error {
	if checkEmptyString(version) != nil {
		return errors.New("Empty version")
	}
	return nil
}

func validateMaintainers(maintainers []maintainer) error {
	if len(maintainers) == 0 {
		return errors.New("Empty maintainers list")
	}

	for _, element := range maintainers {
		if checkEmptyString(element.Name) != nil {
			return errors.New("Empty name in maintainer")
		}
		if checkEmptyString(element.Email) != nil {
			return errors.New("Empty email in maintainer")
		}
		if checkEmail(element.Email) != nil {
			return errors.New("Invalid email")
		}
	}
	return nil
}

func validateCompany(company string) error {
	if checkEmptyString(company) != nil {
		return errors.New("Empty company")
	}
	return nil
}

func validateWebsite(website string) error {
	if checkEmptyString(website) != nil {
		return errors.New("Empty website")
	}
	if checkURL(website) != nil {
		return errors.New("Invalid website URL")
	}
	return nil
}

func validateSource(source string) error {
	if checkEmptyString(source) != nil {
		return errors.New("Empty source")
	}
	if checkURL(source) != nil {
		return errors.New("Invalid source URL")
	}
	return nil
}

func validateLicense(license string) error {
	if checkEmptyString(license) != nil {
		return errors.New("Empty license")
	}
	return nil
}

func validateDescription(description string) error {
	if checkEmptyString(description) != nil {
		return errors.New("Empty description")
	}
	return nil
}

//
// Main function for validating app metadata
//
func validateAppMetaData(data appMetaData) (err error) {
	err = validateTitle(data.Title)
	if err != nil {
		return err
	}

	err = validateVersion(data.Version)
	if err != nil {
		return err
	}

	err = validateMaintainers(data.Maintainers)
	if err != nil {
		return err
	}

	err = validateCompany(data.Company)
	if err != nil {
		return err
	}

	err = validateWebsite(data.Website)
	if err != nil {
		return err
	}

	err = validateSource(data.Source)
	if err != nil {
		return err
	}

	err = validateLicense(data.Source)
	if err != nil {
		return err
	}

	err = validateDescription(data.Description)
	if err != nil {
		return err
	}

	return err
}

//
// Replace underscore for spaces
//
func replaceUnderscore(k string, v []string) []string {
	switch k {
	case "title":
		fallthrough
	case "name":
		fallthrough
	case "company":
		fallthrough
	case "license":
		fallthrough
	case "description":
		for idx, value := range v {
			v[idx] = strings.Replace(value, "_", " ", -1)
		}
		return v
	default:
		return v
	}
}
