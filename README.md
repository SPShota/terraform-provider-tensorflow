# TF.tf

TF.tf is a Terraform provider for generating TensorFlow Python programs from HCL.

This repository is currently in the provider scaffold phase. The first milestone
establishes the Terraform Plugin Framework entrypoint, provider metadata, CI, and
test structure. TensorFlow code-generation data sources will be added in later
PRs.

## Development

Requirements:

- Go 1.24+
- Terraform 1.8+

Run the unit tests:

```sh
go test ./...
```

Check formatting:

```sh
test -z "$(gofmt -l .)"
```

Run the provider directly to verify that it compiles:

```sh
go run .
```

Terraform provider binaries are not meant to be executed directly, so a successful
compile exits with Terraform plugin startup guidance.

## Roadmap

The planned implementation proceeds in small PRs:

1. Provider scaffold, CI, and README.
2. Python code-generation IR.
3. `tf_program` data source.
4. Generic expression and call data sources.
5. Usage documentation and examples.

Later PRs will add generated wrappers for TensorFlow APIs from the official
`tf` namespace documentation.
