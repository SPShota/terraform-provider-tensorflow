# Wrapper Generation

`tftf-wrappergen` turns an API manifest into Go source that registers generated
Terraform data sources.

The generated data sources are thin wrappers over the generic expression call
shape:

```hcl
data "tensorflow_constant" "x" {
  args = [data.tensorflow_literal.values.expression]
  kwargs = {
    dtype = data.tensorflow_ref.float32.expression
  }
}
```

The wrapper emits the same outputs as `tensorflow_call`:

- `expression`
- `statement`

Generate wrapper source from a manifest:

```sh
go run ./cmd/tftf-wrappergen generate \
  -input tf-manifest.json \
  -output internal/provider/generated_data_sources.go
```

Generated wrappers for selected TensorFlow namespaces are checked in under the
provider package. Use [API Coverage](coverage.md) to see which manifest symbols
have dedicated wrappers.
