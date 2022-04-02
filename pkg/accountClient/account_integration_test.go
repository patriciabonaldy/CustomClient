package accountClient

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var accounAddr string

func TestMain(m *testing.M) {
	accounAddr = fmt.Sprintf("%s/v1/organisation/accounts", os.Getenv("ACCOUNT_ADDR"))
	if os.Getenv("ACCOUNT_ADDR") == "" {
		accounAddr = "http://localhost:8080/v1/organisation/accounts"
	}

	m.Run()
	os.Exit(0)
}

func Test_Save(t *testing.T) {
	tests := []struct {
		name        string
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
			name: "error in post",
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
			name: "invalid BankID",
			accountData: AccountData{
				ID:             "b8fc26d7-ca23-4b32-a5ad-5a5e39b057be",
				OrganisationID: "b8fc26d7-ca23-4b32-a5ad-5a5e39b097be",
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
			name: "invalid BIC",
			accountData: AccountData{
				ID:             "b8fc26d7-ca23-4b32-a5ad-5a5e39b057be",
				OrganisationID: "b8fc26d7-ca23-4b32-a5ad-5a5e39b097be",
				Type:           "accounts",
				Attributes: &Attributes{
					Country:    getString("GB"),
					Name:       []string{"John"},
					BankID:     "400302",
					Bic:        "ZZZ",
					BankIDCode: "GBDSC",
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			accountData: AccountData{
				ID:             generateID(),
				OrganisationID: generateID(),
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
			a, err := New(accounAddr)
			require.NoError(t, err)

			if _, err := a.Save(context.Background(), tt.accountData); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_Delete(t *testing.T) {
	type args struct {
		accountID string
		fn        func(baseURL string) Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "error accountID is empty",
			args: args{
				fn: func(baseURL string) Account {
					a, err := New(baseURL)
					require.NoError(t, err)

					return a
				},
			},
			wantErr: true,
		},
		{
			name: "error accountID is empty",
			args: args{
				fn: func(baseURL string) Account {
					a, err := New(baseURL)
					require.NoError(t, err)

					return a
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				accountID: "b8fc26d7-ca23-4b32-a5ad-3a5e39b048be",
				fn: func(baseURL string) Account {
					a, err := New(baseURL)
					require.NoError(t, err)

					accountData := AccountData{
						ID:             "b8fc26d7-ca23-4b32-a5ad-3a5e39b048be",
						OrganisationID: "b8fc26d7-ca23-4b32-a5ad-5a5e39b047be",
						Type:           "accounts",
						Attributes: &Attributes{
							Country:    getString("GB"),
							Name:       []string{"John"},
							BankID:     "400302",
							Bic:        "NWBKGB42",
							BankIDCode: "GBDSC",
						},
					}
					a.Save(context.Background(), accountData) // nolint:errcheck
					return a
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.args.fn(accounAddr)
			if err := a.Delete(context.Background(), tt.args.accountID, 0); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func generateID() string {
	uid, _ := uuid.NewUUID()
	return uid.String()
}
