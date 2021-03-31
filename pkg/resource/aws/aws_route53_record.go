// GENERATED, DO NOT EDIT THIS FILE
package aws

const AwsRoute53RecordResourceType = "aws_route53_record"

type AwsRoute53Record struct {
	AllowOverwrite                *bool     `cty:"allow_overwrite" diff:"-" computed:"true"`
	Fqdn                          *string   `cty:"fqdn" computed:"true"`
	HealthCheckId                 *string   `cty:"health_check_id"`
	Id                            string    `cty:"id" computed:"true"`
	MultivalueAnswerRoutingPolicy *bool     `cty:"multivalue_answer_routing_policy"`
	Name                          *string   `cty:"name" diff:"-"`
	Records                       *[]string `cty:"records"`
	SetIdentifier                 *string   `cty:"set_identifier"`
	Ttl                           *int      `cty:"ttl"`
	Type                          *string   `cty:"type"`
	ZoneId                        *string   `cty:"zone_id"`
	Alias                         *[]struct {
		EvaluateTargetHealth *bool   `cty:"evaluate_target_health"`
		Name                 *string `cty:"name"`
		ZoneId               *string `cty:"zone_id"`
	} `cty:"alias"`
	FailoverRoutingPolicy *[]struct {
		Type *string `cty:"type"`
	} `cty:"failover_routing_policy"`
	GeolocationRoutingPolicy *[]struct {
		Continent   *string `cty:"continent"`
		Country     *string `cty:"country"`
		Subdivision *string `cty:"subdivision"`
	} `cty:"geolocation_routing_policy"`
	LatencyRoutingPolicy *[]struct {
		Region *string `cty:"region"`
	} `cty:"latency_routing_policy"`
	WeightedRoutingPolicy *[]struct {
		Weight *int `cty:"weight"`
	} `cty:"weighted_routing_policy"`
}

func (r *AwsRoute53Record) TerraformId() string {
	return r.Id
}

func (r *AwsRoute53Record) TerraformType() string {
	return AwsRoute53RecordResourceType
}
