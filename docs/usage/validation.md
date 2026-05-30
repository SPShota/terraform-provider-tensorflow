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
added separately as an opt-in integration test.
