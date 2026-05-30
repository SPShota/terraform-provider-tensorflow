# Blocks And Functions

TF.tf supports Python block statements for TensorFlow workflows that need
functions, decorators, or context managers.

The block-oriented data sources are:

- `tensorflow_return`
- `tensorflow_with`
- `tensorflow_function`

Related wrappers:

- `tensorflow_variable` wraps `tf.Variable`
- `tensorflow_gradient_tape` wraps `tf.GradientTape`
- `tensorflow_module` wraps `tf.Module`

## Function With `tf.function`

Use `tensorflow_function` to generate a Python function definition. Decorators
are expressions without the leading `@`.

```hcl
data "tensorflow_return" "loss" {
  value = "loss"
}

data "tensorflow_function" "train_step" {
  name       = "train_step"
  decorators = ["tf.function"]
  statements = [
    "loss = weight * weight",
    data.tensorflow_return.loss.statement,
  ]
}
```

Generated Python:

```python
@tf.function
def train_step():
	loss = weight * weight
	return loss
```

## GradientTape

Use `tensorflow_gradient_tape` to create the context manager expression, then
wrap statements with `tensorflow_with`.

```hcl
data "tensorflow_gradient_tape" "tape" {}

data "tensorflow_with" "tape" {
  context = data.tensorflow_gradient_tape.tape.expression
  alias   = "tape"

  statements = [
    "loss = weight * weight",
  ]
}
```

Generated Python:

```python
with tf.GradientTape() as tape:
	loss = weight * weight
```

## End-To-End Example

See [examples/blocks](../../examples/blocks).
