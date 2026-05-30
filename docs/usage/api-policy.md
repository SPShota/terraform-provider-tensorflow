# API Coverage Policy

TF.tf uses three layers to cover TensorFlow APIs:

1. Generic expression data sources make arbitrary Python calls possible.
2. Generated wrappers provide ergonomic data sources for selected TensorFlow
   symbols.
3. `tensorflow_raw_op` covers the large `tf.raw_ops` namespace.

## Dedicated Wrappers

Dedicated wrappers are added by namespace and kept thin. A wrapper fixes the
callable, accepts `args` and `kwargs`, and returns:

- `expression`
- `statement`

This keeps generated wrappers predictable and easy to review.

## Manifest Flow

Generate a TensorFlow manifest from the Python environment you want to target:

```sh
go run ./cmd/tftf-manifest generate -output tf-manifest.json
```

Generate wrapper source from that manifest:

```sh
go run ./cmd/tftf-wrappergen generate \
  -input tf-manifest.json \
  -output internal/provider/generated_data_sources.go
```

Compare the manifest with registered wrappers:

```sh
go run ./cmd/tftf-coverage report \
  -input tf-manifest.json \
  -output coverage.md
```

## Raw Ops

`tf.raw_ops` is treated as covered by `tensorflow_raw_op` by default. To report
raw ops as individual missing symbols, run coverage with:

```sh
go run ./cmd/tftf-coverage report \
  -input tf-manifest.json \
  -include-raw-ops=false
```
