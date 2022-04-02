package accountClient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/patriciabonaldy/interview-accountapi/pkg/client"
)

var validBIC = regexp.MustCompile(`^([A-Z]{6}[A-Z0-9]{2}|[A-Z]{6}[A-Z0-9]{5})$`)

// Account defines an interface for account client.
type Account interface {
	Delete(ctx context.Context, accountID string, version int) error
	GetByAccountID(ctx context.Context, accountID string) (RequestAccount, error)
	Save(ctx context.Context, accountData AccountData) (RequestAccount, error)
}

type account struct {
	baseURL string
	client  client.Client
}

// New function return an instance of account client
func New(baseURL string) (Account, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL is empty")
	}

	return &account{baseURL: baseURL, client: client.New()}, nil
}

func (a *account) Delete(ctx context.Context, accountID string, version int) error {
	if accountID == "" {
		return errors.New("accountID is empty")
	}

	url := fmt.Sprintf("%s/%s?version=%d", a.baseURL, accountID, version)
	err := a.client.Delete(ctx, url)
	if err != nil {
		return fmt.Errorf("failed deleting account [%s]", err)
	}

	return nil
}

func (a *account) GetByAccountID(ctx context.Context, accountID string) (RequestAccount, error) {
	if accountID == "" {
		return RequestAccount{}, errors.New("accountID is empty")
	}

	url := fmt.Sprintf("%s/%s", a.baseURL, accountID)
	resp, err := a.client.Get(ctx, url)
	if err != nil {
		return RequestAccount{}, fmt.Errorf("failed fetching account [%s]", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RequestAccount{}, fmt.Errorf("failed fetching account [%s]", err)
	}

	defer resp.Body.Close() // nolint:errcheck
	var accountResp RequestAccount
	err = json.Unmarshal(data, &accountResp)
	if err != nil {
		return RequestAccount{}, fmt.Errorf("failed fetching account [%s]", err)
	}

	return accountResp, nil
}

func (a *account) Save(ctx context.Context, accountData AccountData) (RequestAccount, error) {
	request, err := createCustomAccount(accountData)
	if err != nil {
		return RequestAccount{}, fmt.Errorf("failed saving account [%s]", err)
	}

	body, err := json.Marshal(request)
	if err != nil {
		return RequestAccount{}, fmt.Errorf("failed saving account [%s]", err)
	}

	headers := []client.Header{{Key: "Content-Type", Value: "application/json"}}
	resp, err := a.client.Post(ctx, a.baseURL, body, headers...)
	if err != nil {
		return RequestAccount{}, fmt.Errorf("failed saving account [%s]", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RequestAccount{}, fmt.Errorf("failed fetching account [%s]", err)
	}

	defer resp.Body.Close() // nolint:errcheck

	var accountResp RequestAccount
	err = json.Unmarshal(data, &accountResp)
	if err != nil {
		return RequestAccount{}, fmt.Errorf("failed fetching account [%s]", err)
	}

	return accountResp, nil
}

func createCustomAccount(ac AccountData) (RequestAccount, error) {
	if ac.Attributes == nil {
		return RequestAccount{}, errors.New("attributes can not be null")
	}

	err := validation.Errors{
		"id":              validation.Validate(&ac.ID, validation.Required, is.UUID),
		"organisation_id": validation.Validate(&ac.OrganisationID, validation.Required, is.UUID),
		"type":            validation.Validate(&ac.Type, validation.Required),
		"country":         validation.Validate(&ac.Attributes.Country, validation.Required),
		"name":            validation.Validate(&ac.Attributes.Name, validation.Required),
	}.Filter()
	if err != nil {
		return RequestAccount{}, err
	}

	fn, ok := mapAttributes[*ac.Attributes.Country]
	if !ok {
		return RequestAccount{}, errors.New("country does not exist")
	}

	err = fn(ac.Attributes)
	if err != nil {
		return RequestAccount{}, err
	}

	if ac.Attributes.BankID != "" && isNotNumber(ac.Attributes.BankID) {
		return RequestAccount{}, errors.New("bank ID must be number")
	}

	if ac.Attributes.Bic != "" && !validBIC.MatchString(ac.Attributes.Bic) {
		return RequestAccount{}, errors.New("invalid BIC")
	}

	return RequestAccount{Account: ac}, nil
}

func isNotNumber(id string) bool {
	if _, err := strconv.Atoi(id); err != nil {
		return true
	}

	return false
}
