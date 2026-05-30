# Limitations

TF.tf is a code generator. Terraform evaluates data sources and produces Python
source; it does not execute TensorFlow during a normal plan or apply.

## Python Semantics

TF.tf validates basic identifier shapes and checks generated Python syntax in
tests, but it does not type-check TensorFlow programs. Runtime errors such as
shape mismatches, invalid dtypes, missing files, or unsupported TensorFlow
arguments are reported by Python/TensorFlow when the generated script runs.

## Expressions Are Strings

Most data sources exchange Python expressions as strings. This keeps TF.tf
flexible enough to represent the TensorFlow API surface, but Terraform cannot
understand Python semantics inside those strings.

Use `tensorflow_literal` for HCL values when possible, and use raw expression
strings when you need Python constructs such as tuples, method calls, or
constants like `tf.data.AUTOTUNE`.

## Wrapper Coverage

Dedicated generated wrappers are intentionally incremental. If an API does not
have a wrapper yet, use:

- `tensorflow_call` for normal `tf.*` APIs
- `tensorflow_raw_op` for `tf.raw_ops.*`

Use [API Coverage](coverage.md) to compare a TensorFlow manifest against the
registered wrappers.

## Terraform Provider Installation

Examples assume the provider is available locally or from a Terraform registry.
During development, use a Terraform CLI `dev_overrides` block as shown in the
[Usage Guide](index.md).
