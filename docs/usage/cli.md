# CLI Tools

TF.tf includes small development CLIs for API discovery, wrapper generation, and
coverage reporting.

## `tftf-manifest`

Generate a TensorFlow API manifest from an installed Python environment:

```sh
go run ./cmd/tftf-manifest generate \
  -python python3 \
  -module tensorflow \
  -root tf \
  -max-depth 3 \
  -output tf-manifest.json
```

TensorFlow must be importable from the selected Python executable.

## `tftf-wrappergen`

Generate provider wrapper source from a manifest:

```sh
go run ./cmd/tftf-wrappergen generate \
  -input tf-manifest.json \
  -output internal/provider/generated_data_sources.go
```

Generated wrappers are thin data sources that accept `args` and `kwargs` and
return `expression` and `statement`.

## `tftf-coverage`

Generate an API coverage report:

```sh
go run ./cmd/tftf-coverage report \
  -input tf-manifest.json \
  -format markdown \
  -output coverage.md
```

Use `-format json` for machine-readable output.
