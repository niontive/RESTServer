package main

import (
	"errors"
	"net/mail"
	"net/url"
	"strings"
)

type maintainer struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type AppMetaData struct {
	Title       string `yaml:"title"`
	Version     string `yaml:"version"`
	Maintainers []maintainer
	Company     string `yaml:"company"`
	Website     string `yaml:"website"`
	Source      string `yaml:"source"`
	License     string `yaml:"license"`
	Description string `yaml:"description"`
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

func validateAppMetaData(data AppMetaData) (err error) {
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
