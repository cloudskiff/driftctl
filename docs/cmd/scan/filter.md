---
id: filtering
title: Filtering Resources
---

driftctl offers two ways to filter resources:

- .driftignore
- Filter rules

**.driftignore** is a simple way to ignore resources, you put resources in a `.driftignore` file like a `.gitignore`.

**Filter rules** allow you to build complex expression to include and exclude a set of resources in your workflow.
Powered by expression language [JMESPath](https://jmespath.org/) you could build a complex include and exclude expression.

If you need only to exclude a set of resources you should use .driftignore, if you need something more advanced, check filter rules.

## .driftignore

Create the .driftignore file where you launch driftctl (usually the root of your IaC repo).

Each line must be of kind

- `resource_type.resource_id`, resource_id could be a wildcard to exclude all resources of a given type.
- `resource_type.resource_id.path.to.FieldName`, resource_id can be wildcard to ignore a drift on given field for a given type, path could also contain wildcards.

> **N.B.** Fields are not case-sensitive.

If your resource id or the path of a field contains dot or backslash you can escape them with backslashes:

```ignore
resource_type.resource\.id\.containing\.dots.path.to.dotted\.FieldName
resource_type.resource_id_containing\\backslash.path.to.backslash\\FieldName
```

### Examples

```ignore
# Will ignore S3 bucket called my-bucket
aws_s3_bucket.my-buckey
# Will ignore every aws_instance resource
aws_instance.*
# Will ignore environement for all lambda functions
aws_lambda_function.*.Environment
# Will ignore lastModified for my-lambda-name lambda function
aws_lambda_function.my-lambda-name.LastModified
```

## Filter rules

Filter rules could be passed to `scan` cmd with `--filter` flag.
You could also use the environment variable `DCTL_FILTER`.
Filter rules syntax in use is actually [JMESPath](https://jmespath.org/specification.html).

Filter are applied on a normalized struct which contains the following fields:

- **Type**: Type of the resource, e.g. `aws_s3_bucket`
- **Id**: Id of the resource, e.g. `my-bucket-name`
- **Attr**: Contains every resource attributes (check `pkg/resource/aws/aws_s3_bucket.go` for a full list of supported attributes of a bucket)

### Examples

```shell
# Will include only S3 bucket in the search
$ driftctl scan --filter "Type=='aws_s3_bucket'"
# OR (beware of escape your shell special chars between double quotes)
$ driftctl scan --filter $'Type==\'aws_s3_bucket\''

# Excludes only s3 bucket named 'my-bucket-name'
$ driftctl scan --filter $'Type==\'aws_s3_bucket\' && Id!=\'my-bucket-name\''

# Ignore buckets that have tags terraform equal to 'false'
$ driftctl scan --filter $'!(Type==\'aws_s3_bucket\' && Attr.Tags.terraform==\'false\')'

# Ignore buckets that don't have tag terraform
$ driftctl scan --filter $'!(Type==\'aws_s3_bucket\' && Attr.Tags != null && !contains(keys(Attr.Tags), \'terraform\'))'

# Ignore buckets with an ID prefix of 'terraform-'
$ driftctl scan --filter $'!(Type==\'aws_s3_bucket\' && starts_with(Id, \'terraform-\'))'

# Ignore buckets with an ID suffix of '-test'
$ driftctl scan --filter $'!(Type==\'aws_s3_bucket\' && ends_with(Id, \'-test\'))'

```
