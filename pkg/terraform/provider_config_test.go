package terraform

import (
	"fmt"
	"runtime"
	"testing"
)

func TestProviderConfig_GetBinaryName(t *testing.T) {
	type fields struct {
		Key     string
		Version string
		Postfix string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test for aws provider",
			fields: fields{
				Key:     "aws",
				Version: "3.24.1",
				Postfix: "x5",
			},
			want: "terraform-provider-aws_v3.24.1_x5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ProviderConfig{
				Key:     tt.fields.Key,
				Version: tt.fields.Version,
				Postfix: tt.fields.Postfix,
			}
			if got := c.GetBinaryName(); got != tt.want {
				t.Errorf("GetBinaryName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProviderConfig_GetDownloadUrl(t *testing.T) {
	type fields struct {
		Key     string
		Version string
		Postfix string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test for aws provider",
			fields: fields{
				Key:     "aws",
				Version: "3.24.1",
			},
			want: fmt.Sprintf(
				"https://releases.hashicorp.com/terraform-provider-aws/3.24.1/terraform-provider-aws_3.24.1_%s_%s.zip",
				runtime.GOOS,
				runtime.GOARCH,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ProviderConfig{
				Key:     tt.fields.Key,
				Version: tt.fields.Version,
				Postfix: tt.fields.Postfix,
			}
			if got := c.GetDownloadUrl(); got != tt.want {
				t.Errorf("GetDownloadUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}