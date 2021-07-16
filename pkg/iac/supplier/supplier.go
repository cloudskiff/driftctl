package supplier

import (
	"fmt"

	"github.com/cloudskiff/driftctl/pkg/alerter"
	"github.com/cloudskiff/driftctl/pkg/iac/terraform/state/backend"
	"github.com/cloudskiff/driftctl/pkg/output"
	"github.com/cloudskiff/driftctl/pkg/terraform"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/cloudskiff/driftctl/pkg/iac/config"

	"github.com/cloudskiff/driftctl/pkg/iac/terraform/state"

	"github.com/cloudskiff/driftctl/pkg/resource"
)

var supportedSuppliers = []string{
	state.TerraformStateReaderSupplier,
}

func IsSupplierSupported(supplierKey string) bool {
	for _, s := range supportedSuppliers {
		if s == supplierKey {
			return true
		}
	}
	return false
}

func GetIACSupplier(configs []config.SupplierConfig,
	library *terraform.ProviderLibrary,
	backendOpts *backend.Options,
	progress output.Progress,
	alerter *alerter.Alerter,
	factory resource.ResourceFactory) (resource.Supplier, error) {

	chainSupplier := NewIacChainSupplier()
	for _, config := range configs {
		if !IsSupplierSupported(config.Key) {
			alerter.SendAlert("", state.NewStateEnumerationAlertAlert(errors.Errorf("Unsupported supplier '%s'", config.Key)))
			continue
		}

		deserializer := resource.NewDeserializer(factory)

		var supplier resource.Supplier
		var err error
		switch config.Key {
		case state.TerraformStateReaderSupplier:
			supplier, err = state.NewReader(config, library, backendOpts, progress, deserializer)
		default:
			alerter.SendAlert("", state.NewStateEnumerationAlertAlert(errors.Errorf("Unsupported supplier '%s'", config.Key)))
			continue
		}

		if err != nil {
			alerter.SendAlert("", state.NewStateEnumerationAlertAlert(err))
			continue
		}

		logrus.WithFields(logrus.Fields{
			"supplier": config.Key,
			"backend":  config.Backend,
			"path":     config.Path,
		}).Debug("Found IAC supplier")

		chainSupplier.AddSupplier(supplier)
	}

	if chainSupplier.CountSuppliers() <= 0 {
		return nil, errors.New("none of given from where supported")
	}

	return chainSupplier, nil
}

func GetSupportedSuppliers() []string {
	return supportedSuppliers
}

func GetSupportedSchemes() []string {
	schemes := []string{
		"tfstate://",
	}
	for _, supplier := range supportedSuppliers {
		for _, backend := range backend.GetSupportedBackends() {
			schemes = append(schemes, fmt.Sprintf("%s+%s://", supplier, backend))
		}
	}
	return schemes
}
