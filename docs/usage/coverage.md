# API Coverage

`tftf-coverage` compares a TensorFlow API manifest with the provider's generated
wrappers and reports covered and missing symbols.

Generate a manifest:

```sh
go run ./cmd/tftf-manifest generate -output tf-manifest.json
```

Generate a Markdown coverage report:

```sh
go run ./cmd/tftf-coverage report \
  -input tf-manifest.json \
  -format markdown \
  -output coverage.md
```

Generate JSON:

```sh
go run ./cmd/tftf-coverage report \
  -input tf-manifest.json \
  -format json \
  -output coverage.json
```

By default, `tf.raw_ops.*` symbols are treated as covered by
`tensorflow_raw_op`, because raw ops are exposed through a generic data source.
Disable that behavior with:

```sh
go run ./cmd/tftf-coverage report \
  -input tf-manifest.json \
  -include-raw-ops=false
```
