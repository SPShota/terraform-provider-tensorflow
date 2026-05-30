# Validation

TF.tf validates generated Python syntax in unit tests with:

```sh
python3 -m py_compile generated.py
```

This check verifies Python syntax only. It does not import TensorFlow or execute
the generated program, so it stays fast and does not require TensorFlow to be
installed.

For local examples, generate a file and run the same command:

```sh
terraform apply
python3 -m py_compile generated.py
```

Runtime validation that imports TensorFlow and executes generated code will be
is available as an opt-in integration test:

```sh
TF_TF_INTEGRATION=1 go test ./internal/integration/...
```

The integration test skips when TensorFlow is not importable from `python3`.
Install TensorFlow in the Python environment used by `python3` before running it
when you want runtime coverage.
