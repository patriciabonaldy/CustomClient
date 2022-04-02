package account

import (
	"context"
	"github.com/patriciabonaldy/interview-accountapi/pkg/client"
	"testing"
)

func Test_account_Save(t *testing.T) {
	type fields struct {
		baseURL string
		client  client.Client
	}
	type args struct {
		ctx         context.Context
		accountData AccountData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &account{
				baseURL: tt.fields.baseURL,
				client:  tt.fields.client,
			}
			if err := a.Save(tt.args.ctx, tt.args.accountData); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
