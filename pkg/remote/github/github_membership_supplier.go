package github

import (
	remoteerror "github.com/cloudskiff/driftctl/pkg/remote/error"
	"github.com/cloudskiff/driftctl/pkg/resource"
	resourcegithub "github.com/cloudskiff/driftctl/pkg/resource/github"
	"github.com/cloudskiff/driftctl/pkg/terraform"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
)

type GithubMembershipSupplier struct {
	reader       terraform.ResourceReader
	deserializer *resource.Deserializer
	repository   GithubRepository
	runner       *terraform.ParallelResourceReader
}

func NewGithubMembershipSupplier(provider *GithubTerraformProvider, repository GithubRepository, deserializer *resource.Deserializer) *GithubMembershipSupplier {
	return &GithubMembershipSupplier{
		provider,
		deserializer,
		repository,
		terraform.NewParallelResourceReader(provider.Runner().SubRunner()),
	}
}

func (s *GithubMembershipSupplier) SuppliedType() resource.ResourceType {
	return resourcegithub.GithubMembershipResourceType
}

func (s *GithubMembershipSupplier) Resources() ([]resource.Resource, error) {

	resourceList, err := s.repository.ListMembership()

	if err != nil {
		return nil, remoteerror.NewResourceEnumerationError(err, s.SuppliedType())
	}

	for _, id := range resourceList {
		id := id
		s.runner.Run(func() (cty.Value, error) {
			completeResource, err := s.reader.ReadResource(terraform.ReadResourceArgs{
				Ty: s.SuppliedType(),
				ID: id,
			})
			if err != nil {
				logrus.Warnf("Error reading %s[%s]: %+v", id, s.SuppliedType(), err)
				return cty.NilVal, err
			}
			return *completeResource, nil
		})
	}

	results, err := s.runner.Wait()
	if err != nil {
		return nil, err
	}

	return s.deserializer.Deserialize(s.SuppliedType(), results)
}
