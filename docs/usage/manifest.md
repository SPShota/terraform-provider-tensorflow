# API Manifest

`tftf-manifest` generates a JSON manifest from an installed TensorFlow Python
module. Later PRs use this manifest to generate thin Terraform data source
wrappers for TensorFlow APIs.

Generate a manifest:

```sh
go run ./cmd/tftf-manifest generate \
  -python python3 \
  -module tensorflow \
  -root tf \
  -max-depth 3 \
  -output tf-manifest.json
```

The command imports TensorFlow from the selected Python environment and walks the
public `tf` namespace. TensorFlow must be installed for that Python interpreter.

The manifest records:

- symbol path, such as `tf.constant`
- kind, such as `module`, `class`, `function`, `callable`, or `value`
- Python module name
- best-effort Python signature
- TensorFlow documentation URL
- direct child symbols for namespaces and classes

The default documentation base is:

```text
https://www.tensorflow.org/api_docs/python
```

This PR adds the generator and schema only. Provider wrappers are generated in a
later PR.
