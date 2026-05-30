---
page_title: "tensorflow_call Data Source - TensorFlow"
subcategory: ""
description: |-
  Generates a Python function or method call expression.
---

# `tensorflow_call` Data Source

`tensorflow_call` creates a Python call expression. Use it for APIs or method
calls that do not have a dedicated wrapper.

## Example Usage

```terraform
data "tensorflow_call" "compile" {
  function = "${data.tensorflow_assign.model.expression}.compile"
  kwargs = {
    optimizer = data.tensorflow_keras_optimizers_adam.optimizer.expression
    loss      = data.tensorflow_keras_losses_sparse_categorical_crossentropy.loss.expression
  }
}
```

## Argument Reference

- `function` (Required) Python callable expression.
- `args` (Optional) Positional argument expressions.
- `kwargs` (Optional) Keyword argument expressions.

## Attribute Reference

- `expression` Generated Python call expression.
- `statement` Generated Python expression statement.
