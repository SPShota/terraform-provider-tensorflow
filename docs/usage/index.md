# Usage Guide

TF.tf generates TensorFlow Python source from Terraform data sources.

The recommended flow is:

1. Build Python expressions with `tensorflow_literal`, `tensorflow_ref`,
   `tensorflow_attr`, and `tensorflow_call`.
2. Turn expressions into statements with `tensorflow_assign` or the
   `statement` output from `tensorflow_call`.
3. Join statements with `tensorflow_program`.
4. Write the generated source with a Terraform resource such as `local_file`.

## Start Here

- [Basic Example](basic.md)
- [Expressions](expressions.md)
- [Blocks And Functions](blocks.md)
- [`tf_program`](program.md)

## TensorFlow APIs

- [Generated Wrappers](generated-wrappers.md)
- [Keras](keras.md)
- [`tf.data`](data.md)
- [Additional TensorFlow Namespaces](namespaces.md)
- [Raw Ops](raw-ops.md)

## Tooling

- [Validation](validation.md)
- [API Manifest](manifest.md)
- [Wrapper Generation](wrapper-generation.md)
- [API Coverage](coverage.md)
- [CLI Tools](cli.md)

## Project Notes

- [Limitations](limitations.md)
- [API Coverage Policy](api-policy.md)
- [Internal Python IR](python-ir.md)

## Data Sources

| Data source | Purpose |
| --- | --- |
| `tensorflow_literal` | Converts JSON/HCL values into Python literal expressions. |
| `tensorflow_ref` | Emits a Python reference, such as `tf.float32`. |
| `tensorflow_attr` | Emits attribute access, such as `tf.keras.Sequential`. |
| `tensorflow_call` | Emits a Python call expression and expression statement. |
| `tensorflow_assign` | Emits an assignment statement and reference expression. |
| `tensorflow_return` | Emits a Python return statement. |
| `tensorflow_with` | Emits a Python `with` block. |
| `tensorflow_function` | Emits a Python function definition. |
| `tensorflow_raw_op` | Emits a `tf.raw_ops.<op>(...)` call. |
| `tensorflow_program` | Emits the final Python program content. |

## Local Provider Development

During local development, configure Terraform to use the provider binary built
from this repository. One common approach is a Terraform CLI development
override in `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "SPShota/tensorflow" = "/path/to/terraform-provider-tensorflow/bin"
  }

  direct {}
}
```

Then build the provider binary into that directory:

```sh
go build -o bin/terraform-provider-tensorflow
```

After that, examples can be evaluated with `terraform plan` or
`terraform apply`.

## Validation Commands

```sh
go test ./...
TF_TF_INTEGRATION=1 go test ./internal/integration/...
```

The integration command requires TensorFlow to be installed for `python3`.
