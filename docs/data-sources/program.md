---
page_title: "tensorflow_program Data Source - TensorFlow"
subcategory: ""
description: |-
  Generates a TensorFlow Python program from Python statements.
---

# `tensorflow_program` Data Source

`tensorflow_program` joins Python statements into a complete TensorFlow Python
program.

## Example Usage

```terraform
data "tensorflow_program" "main" {
  statements = [
    "x = tf.constant([1, 2, 3])",
    "print(tf.reduce_sum(x))",
  ]
}

output "python" {
  value = data.tensorflow_program.main.content
}
```

## Argument Reference

- `statements` (Required) Python statements to append to the generated program.
- `header` (Optional) File header.
- `imports` (Optional) Python import specs.

## Attribute Reference

- `content` Generated Python program content.
