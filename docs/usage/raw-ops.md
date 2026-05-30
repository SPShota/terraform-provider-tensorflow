# Raw Ops

`tf.raw_ops` is a very large TensorFlow namespace. TF.tf exposes it through a
generic data source instead of checking in thousands of dedicated wrappers.

Use `tensorflow_raw_op` with an op name under `tf.raw_ops`:

```hcl
data "tensorflow_raw_op" "sum" {
  op = "AddV2"
  kwargs = {
    x = data.tensorflow_literal.left.expression
    y = data.tensorflow_literal.right.expression
  }
}
```

Generated Python:

```python
tf.raw_ops.AddV2(x=1, y=2)
```

The `op` value must be a Python identifier. It is always rendered under
`tf.raw_ops`, so it cannot escape into arbitrary Python code.

See [examples/raw-ops](../../examples/raw-ops).
