package accountClient

import (
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/patriciabonaldy/interview-accountapi/pkg/genericClient"
)

type mockClient struct {
	wantError bool
}

func (m mockClient) Delete(_ context.Context, _ string, _ ...genericClient.Header) error {
	if m.wantError {
		return errors.New("unknown error")
	}

	return nil
}

func (m mockClient) Get(_ context.Context, _ string) (resp *http.Response, err error) {
	if m.wantError {
		return nil, errors.New("unknown error")
	}

	reader := io.NopCloser(strings.NewReader(`{
		"data": {
			"type": "accounts",
				"id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
				"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
				"attributes": {
				"country": "GB",
					"base_currency": "GBP",
					"account_number": "41426819",
					"bank_id": "400300",
					"bank_id_code": "GBDSC",
					"bic": "NWBKGB22",
					"iban": "GB11NWBK40030041426819",
					"status": "confirmed"
			}}}
		`))

	return &http.Response{
		StatusCode:    http.StatusCreated,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: 0,
		Body:          reader,
	}, nil
}

func (m mockClient) Post(_ context.Context, _ string, _ []byte, _ ...genericClient.Header) (resp *http.Response, err error) {
	if m.wantError {
		return nil, errors.New("unknown error")
	}

	reader := io.NopCloser(strings.NewReader(`{
		"data": {
			"type": "accounts",
				"id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
				"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
				"version": 0,
				"attributes": {
				"country": "GB",
					"base_currency": "GBP",
					"account_number": "41426819",
					"bank_id": "400300",
					"bank_id_code": "GBDSC",
					"bic": "NWBKGB22",
					"iban": "GB11NWBK40030041426819",
					"status": "confirmed"
			},
			"relationships": {
				"account_events": {
					"data": [{
						"type": "account_events",
						"id": "c1023677-70ee-417a-9a6a-e211241f1e9c"
						},{"type": "account_events",
						"id": "aca32528-d4cf-4d54-93fe-5d80d27ab773"
						}]}}}}
		`))

	return &http.Response{
		StatusCode:    http.StatusCreated,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: 0,
		Body:          reader,
	}, nil
}

var _ genericClient.Client = &mockClient{}

func Test_account_Save(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		client      genericClient.Client
		accountData AccountData
		wantErr     bool
	}{
		{
			name: "error creating CustomAccount",
			accountData: AccountData{
				ID:             "b8fc26d7-ca23-4b32-a5ad-5a5e39b047be",
				OrganisationID: "b8fc26d7-ca23-4b32-a5ad-5a5e39b047be",
				Type:           "accounts",
				Attributes: &Attributes{
					Country:    getString("GB"),
					Name:       []string{"John"},
					BankID:     "ZZ",
					Bic:        "ZZZ",
					BankIDCode: "zzzz",
				},
			},
			wantErr: true,
		},
		{
			name:   "error in post",
			client: &mockClient{wantError: true},
			accountData: AccountData{
				ID:             "b8fc26d7-ca23-4b32-a5ad-5a5e39b047be",
				OrganisationID: "b8fc26d7-ca23-4b32-a5ad-5a5e39b047be",
				Type:           "accounts",
				Attributes: &Attributes{
					Country:    getString("GB"),
					Name:       []string{"John"},
					BankID:     "ZZ",
					Bic:        "ZZZ",
					BankIDCode: "GBDSC",
				},
			},
			wantErr: true,
		},
		{
			name:   "success",
			client: &mockClient{},
			accountData: AccountData{
				ID:             "b8fc26d7-ca23-4b32-a5ad-5a5e39b047be",
				OrganisationID: "b8fc26d7-ca23-4b32-a5ad-5a5e39b047be",
				Type:           "accounts",
				Attributes: &Attributes{
					Country:    getString("GB"),
					Name:       []string{"John"},
					BankID:     "400302",
					Bic:        "NWBKGB42",
					BankIDCode: "GBDSC",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &account{
				baseURL: tt.baseURL,
				client:  tt.client,
			}
			if _, err := a.Save(context.Background(), tt.accountData); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createCustomAccount(t *testing.T) {
	tests := []struct {
		name    string
		acc     AccountData
		want    AccountData
		wantErr bool
	}{
		{
			name:    "attributes is nil",
			acc:     AccountData{},
			wantErr: true,
		},
		{
			name:    "error some fields are required",
			acc:     AccountData{Attributes: &Attributes{}},
			wantErr: true,
		},
		{
			name: "error country fn does not exist",
			acc: AccountData{
				ID:             "b8fc26d7-ca23-4b32-a5ad-5a5e39b047be",
				OrganisationID: "b8fc26d7-ca23-4b32-a5ad-5a5e39b047be",
				Type:           "accounts",
				Attributes: &Attributes{
					Country: getString("ARG"),
					Name:    []string{"John"},
				},
			},
			wantErr: true,
		},
		{
			name: "error fields are required in fn",
			acc: AccountData{
				ID:             "b8fc26d7-ca23-4b32-a5ad-5a5e39b047be",
				OrganisationID: "b8fc26d7-ca23-4b32-a5ad-5a5e39b047be",
				Type:           "accounts",
				Attributes: &Attributes{
					Country: getString("GB"),
					Name:    []string{"John"},
				},
			},
			wantErr: true,
		},
		{
			name: "error invalid BankIDCode",
			acc: AccountData{
				ID:             "b8fc26d7-ca23-4b32-a5ad-5a5e39b047be",
				OrganisationID: "b8fc26d7-ca23-4b32-a5ad-5a5e39b047be",
				Type:           "accounts",
				Attributes: &Attributes{
					Country:    getString("GB"),
					Name:       []string{"John"},
					BankID:     "ZZ",
					Bic:        "ZZZ",
					BankIDCode: "zzzz",
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			acc: AccountData{
				ID:             "b8fc26d7-ca23-4b32-a5ad-5a5e39b047be",
				OrganisationID: "b8fc26d7-ca23-4b32-a5ad-5a5e39b047be",
				Type:           "accounts",
				Attributes: &Attributes{
					Country:    getString("GB"),
					Name:       []string{"John"},
					BankID:     "400302",
					Bic:        "NWBKGB42",
					BankIDCode: "GBDSC",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := createCustomAccount(tt.acc)
			if (err != nil) != tt.wantErr {
				t.Errorf("createCustomAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func getString(s string) *string {
	return &s
}

func Test_account_GetByAccountID(t *testing.T) {
	type args struct {
		accountID string
		baseURL   string
		client    genericClient.Client
	}
	tests := []struct {
		name    string
		args    args
		want    RequestAccount
		wantErr bool
	}{
		{
			name: "error accountID is empty",
			args: args{
				accountID: "",
				baseURL:   "http://localhost:8080/v1/organisation/accounts",
			},
			wantErr: true,
		},
		{
			name: "error getting account",
			args: args{
				accountID: "b8fc26d7-ca23-4b32-a5ad-5a5e39b057be",
				baseURL:   "http://localhost:8080/v1/organisation/accounts",
				client:    &mockClient{wantError: true},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				accountID: "b8fc26d7-ca23-4b32-a5ad-5a5e39b057be",
				baseURL:   "http://localhost:8080/v1/organisation/accounts",
				client:    &mockClient{},
			},
			want: RequestAccount{
				Account: AccountData{
					Attributes: &Attributes{
						AccountClassification:   nil,
						AccountMatchingOptOut:   nil,
						AccountNumber:           "41426819",
						AlternativeNames:        nil,
						BankID:                  "400300",
						BankIDCode:              "GBDSC",
						BaseCurrency:            "GBP",
						Bic:                     "NWBKGB22",
						Country:                 getString("GB"),
						Iban:                    "GB11NWBK40030041426819",
						JointAccount:            nil,
						Name:                    nil,
						SecondaryIdentification: "",
						Status:                  getString("confirmed"),
						Switched:                nil,
					},
					ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
					OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
					Type:           "accounts",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &account{
				baseURL: tt.args.baseURL,
				client:  tt.args.client,
			}
			got, err := a.GetByAccountID(context.Background(), tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByAccountID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByAccountID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_account_Delete(t *testing.T) {
	type args struct {
		client    genericClient.Client
		baseURL   string
		accountID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "error accountID is empty",
			args: args{
				baseURL: "http://localhost:8080/v1/organisation/accounts",
			},
			wantErr: true,
		},
		{
			name: "error accountID is empty",
			args: args{
				accountID: "b8fc26d7-ca23-4b32-a5ad-5a5e39b057be",
				baseURL:   "http://localhost:8080/v1/organisation/accounts",
				client:    mockClient{wantError: true},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				accountID: "b8fc26d7-ca23-4b32-a5ad-5a5e39b057be",
				baseURL:   "http://localhost:8080/v1/organisation/accounts",
				client:    mockClient{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &account{
				baseURL: tt.args.baseURL,
				client:  tt.args.client,
			}
			if err := a.Delete(context.Background(), tt.args.accountID, 0); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
