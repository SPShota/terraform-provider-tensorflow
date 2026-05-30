# Examples

The examples show progressively richer ways to generate TensorFlow Python code.

| Example | Purpose |
| --- | --- |
| [provider](provider) | Minimal provider configuration. |
| [program](program) | Direct `tensorflow_program` usage with raw statements. |
| [expressions](expressions) | Expression data sources feeding `tensorflow_program`. |
| [basic](basic) | End-to-end source generation into `generated.py`. |

The `basic` example uses `hashicorp/local` to write the generated source file.
The TensorFlow provider must be available locally or from a registry before
running Terraform.
