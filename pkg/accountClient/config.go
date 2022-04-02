package accountClient

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	gbBankID = "GBDSC"
)

var mapAttributes = map[string]func(data *Attributes) error{
	"AU": createAUAttributes,
	"BE": createBEAttributes,
	"CA": createCAAttributes,
	"EE": createEEAttributes,
	"GB": createGBAttributes,
	"FR": createFRAttributes,
}

func createAUAttributes(attributes *Attributes) error {
	return nil
}

func createBEAttributes(attributes *Attributes) error {
	return nil
}

func createCAAttributes(attributes *Attributes) error {
	return nil
}

func createEEAttributes(attributes *Attributes) error {
	return nil
}

func createFRAttributes(attributes *Attributes) error {
	return nil
}

func createGBAttributes(atrb *Attributes) error {
	err := validation.Errors{
		"bank_id":      validation.Validate(&atrb.BankID, validation.Required),
		"bic":          validation.Validate(&atrb.Bic, validation.Required),
		"bank_id_code": validation.Validate(&atrb.BankIDCode, validation.Required),
	}.Filter()
	if err != nil {
		return err
	}

	if atrb.BankIDCode != gbBankID {
		return errors.New("invalid Bank ID Code")
	}

	atrb.BaseCurrency = "GBP"

	return nil
}
