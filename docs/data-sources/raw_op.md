---
page_title: "tensorflow_raw_op Data Source - TensorFlow"
subcategory: ""
description: |-
  Generates a Python call to tf.raw_ops.
---

# `tensorflow_raw_op` Data Source

`tensorflow_raw_op` creates a call to `tf.raw_ops.<op>`.

## Example Usage

```terraform
data "tensorflow_raw_op" "sum" {
  op = "AddV2"
  kwargs = {
    x = "1"
    y = "2"
  }
}
```

Generated expression:

```python
tf.raw_ops.AddV2(x=1, y=2)
```

## Argument Reference

- `op` (Required) Raw op name under `tf.raw_ops`.
- `args` (Optional) Positional argument expressions.
- `kwargs` (Optional) Keyword argument expressions.

## Attribute Reference

- `expression` Generated Python call expression.
- `statement` Generated Python expression statement.
