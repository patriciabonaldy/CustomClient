package accountClient

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	auBankID = "AUBSB"
	beBankID = "BE"
	caBankID = "CACPA"
	gbBankID = "GBDSC"
)

var mapAttributes = map[string]func(data *Attributes) error{
	"AU": createAUAttributes,
	"BE": createBEAttributes,
	"CA": createCAAttributes,
	"GB": createGBAttributes,
	"-": func(attributes *Attributes) error {
		return nil
	},
}

func createAUAttributes(atrb *Attributes) error {
	err := validation.Errors{
		"bic":          validation.Validate(&atrb.Bic, validation.Required),
		"bank_id_code": validation.Validate(&atrb.BankIDCode, validation.Required),
	}.Filter()
	if err != nil {
		return err
	}

	if atrb.BankIDCode != auBankID {
		return errors.New("invalid Bank ID Code")
	}

	if atrb.Iban != "" {
		return errors.New("IBAN has to be empty")
	}

	if len(atrb.AccountNumber) > 0 && atrb.AccountNumber[0] == '0' {
		return errors.New("error AccountNumber - first character cannot be 0")
	}

	const i = 5
	if len(atrb.BankID) > i {
		return errors.New("error BankID, has more than 6 digits")
	}

	return nil
}

func createBEAttributes(atrb *Attributes) error {
	const min = 3
	err := validation.Errors{
		"bank_id":      validation.Validate(&atrb.BankID, validation.Required, validation.Length(min, min)),
		"bank_id_code": validation.Validate(&atrb.BankIDCode, validation.Required),
	}.Filter()
	if err != nil {
		return err
	}

	if atrb.BankIDCode != beBankID {
		return errors.New("invalid Bank ID Code")
	}

	const i = 7
	if len(atrb.AccountNumber) > i {
		return errors.New("error AccountNumber - has to have only 7 digits")
	}

	return nil
}

func createCAAttributes(atrb *Attributes) error {
	err := validation.Errors{
		"bic": validation.Validate(&atrb.Bic, validation.Required),
	}.Filter()
	if err != nil {
		return err
	}

	if atrb.BankIDCode != caBankID {
		return errors.New("invalid Bank ID Code")
	}

	if len(atrb.BankID) > 0 && atrb.BankID[0] != '0' {
		return errors.New("error BankID, first character must be 0")
	}

	if len(atrb.Iban) > 0 {
		return errors.New("iban not supported, has to be empty")
	}

	return nil
}

func createGBAttributes(atrb *Attributes) error {
	const min = 6
	err := validation.Errors{
		"bank_id":      validation.Validate(&atrb.BankID, validation.Required, validation.Length(min, min)),
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
