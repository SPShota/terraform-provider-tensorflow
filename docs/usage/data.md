# `tf.data`

TF.tf includes generated wrappers for an initial `tf.data` dataset creation
subset:

- `tensorflow_data_dataset_from_tensor_slices`
- `tensorflow_data_dataset_from_tensors`
- `tensorflow_data_dataset_from_generator`
- `tensorflow_data_dataset_range`
- `tensorflow_data_dataset_zip`
- `tensorflow_data_dataset_list_files`
- `tensorflow_data_tf_record_dataset`

Dataset transformation methods such as `batch`, `cache`, `map`, `prefetch`,
`repeat`, and `shuffle` are method calls, so use `tensorflow_call` with a method
expression.

## Tensor Slices Pipeline

```hcl
data "tensorflow_literal" "features" {
  value_json = jsonencode([[1.0, 2.0], [3.0, 4.0]])
}

data "tensorflow_data_dataset_from_tensor_slices" "dataset" {
  args = [data.tensorflow_literal.features.expression]
}

data "tensorflow_call" "batch" {
  function = "${data.tensorflow_data_dataset_from_tensor_slices.dataset.expression}.batch"
  args     = ["2"]
}

data "tensorflow_call" "prefetch" {
  function = "${data.tensorflow_call.batch.expression}.prefetch"
  args     = ["tf.data.AUTOTUNE"]
}

data "tensorflow_assign" "dataset" {
  name  = "dataset"
  value = data.tensorflow_call.prefetch.expression
}
```

Generated Python:

```python
dataset = tf.data.Dataset.from_tensor_slices([[1, 2], [3, 4]]).batch(2).prefetch(tf.data.AUTOTUNE)
```

## Method Chain Pattern

Each pipeline step can consume the previous step's `expression`:

```hcl
data "tensorflow_call" "cache" {
  function = "${data.tensorflow_assign.dataset.expression}.cache"
}

data "tensorflow_call" "shuffle" {
  function = "${data.tensorflow_call.cache.expression}.shuffle"
  args     = ["1024"]
}
```

Use `tensorflow_assign` when you want to bind the final pipeline expression to a
name before passing it to Keras `fit` or another generated program statement.

See [examples/data](../../examples/data).
