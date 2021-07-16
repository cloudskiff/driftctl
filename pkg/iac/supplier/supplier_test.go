package supplier

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/cloudskiff/driftctl/pkg/alerter"
	"github.com/cloudskiff/driftctl/pkg/iac/config"
	"github.com/cloudskiff/driftctl/pkg/iac/terraform/state"
	"github.com/cloudskiff/driftctl/pkg/iac/terraform/state/backend"
	"github.com/cloudskiff/driftctl/pkg/output"
	"github.com/cloudskiff/driftctl/pkg/terraform"
	"github.com/cloudskiff/driftctl/test/resource"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetIACSupplier(t *testing.T) {
	type args struct {
		config  []config.SupplierConfig
		options *backend.Options
	}
	tests := []struct {
		name       string
		args       args
		wantErr    error
		wantAlerts alerter.Alerts
	}{
		{
			name: "test unknown supplier",
			args: args{
				config: []config.SupplierConfig{
					{
						Key: "foobar",
					},
				},
				options: &backend.Options{
					Headers: map[string]string{},
				},
			},
			wantErr: fmt.Errorf("none of given from where supported"),
			wantAlerts: map[string][]alerter.Alert{
				"": {state.NewStateEnumerationAlertAlert(errors.Errorf("Unsupported supplier 'foobar'"))},
			},
		},
		{
			name: "test unknown supplier in multiples states",
			args: args{
				config: []config.SupplierConfig{
					{
						Key: "foobar",
					},
					{
						Key:     "tfstate",
						Backend: "",
						Path:    "terraform.tfstate",
					},
				},
				options: &backend.Options{
					Headers: map[string]string{},
				},
			},
			wantErr: nil,
			wantAlerts: map[string][]alerter.Alert{
				"": {state.NewStateEnumerationAlertAlert(errors.Errorf("Unsupported supplier 'foobar'"))},
			},
		},
		{
			name: "test valid tfstate://terraform.tfstate",
			args: args{
				config: []config.SupplierConfig{
					{Key: "tfstate", Backend: "", Path: "terraform.tfstate"},
				},
				options: &backend.Options{
					Headers: map[string]string{},
				},
			},
			wantErr:    nil,
			wantAlerts: alerter.Alerts{},
		},
		{
			name: "test valid multiples states",
			args: args{
				config: []config.SupplierConfig{
					{Key: "tfstate", Backend: "", Path: "terraform.tfstate"},
					{Key: "tfstate", Backend: "s3", Path: "terraform.tfstate"},
					{Key: "tfstate", Backend: "", Path: "terraform2.tfstate"},
				},
				options: &backend.Options{
					Headers: map[string]string{},
				},
			},
			wantErr:    nil,
			wantAlerts: alerter.Alerts{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := &output.MockProgress{}
			progress.On("Start").Return().Times(1)

			alerter := alerter.NewAlerter()

			repo := resource.InitFakeSchemaRepository("aws", "3.19.0")
			factory := terraform.NewTerraformResourceFactory(repo)

			_, err := GetIACSupplier(tt.args.config, terraform.NewProviderLibrary(), tt.args.options, progress, alerter, factory)
			if tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetIACSupplier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantAlerts, alerter.Retrieve())
		})
	}
}

func TestGetSupportedSchemes(t *testing.T) {

	want := []string{
		"tfstate://",
		"tfstate+s3://",
		"tfstate+http://",
		"tfstate+https://",
		"tfstate+tfcloud://",
	}

	if got := GetSupportedSchemes(); !reflect.DeepEqual(got, want) {
		t.Errorf("GetSupportedSchemes() = %v, want %v", got, want)
	}
}
