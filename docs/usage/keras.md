# Keras

TF.tf includes generated wrappers for an initial `tf.keras` subset:

- `tensorflow_keras_sequential`
- `tensorflow_keras_model`
- `tensorflow_keras_input`
- `tensorflow_keras_layers_dense`
- `tensorflow_keras_layers_dropout`
- `tensorflow_keras_layers_flatten`
- `tensorflow_keras_layers_conv2_d`
- `tensorflow_keras_layers_max_pooling2_d`
- `tensorflow_keras_layers_global_average_pooling2_d`
- `tensorflow_keras_layers_embedding`
- `tensorflow_keras_layers_lstm`
- `tensorflow_keras_layers_batch_normalization`
- `tensorflow_keras_layers_activation`
- `tensorflow_keras_optimizers_adam`
- `tensorflow_keras_optimizers_sgd`
- `tensorflow_keras_optimizers_rmsprop`
- `tensorflow_keras_losses_sparse_categorical_crossentropy`
- `tensorflow_keras_losses_categorical_crossentropy`
- `tensorflow_keras_losses_binary_crossentropy`
- `tensorflow_keras_losses_mean_squared_error`
- `tensorflow_keras_metrics_sparse_categorical_accuracy`
- `tensorflow_keras_metrics_categorical_accuracy`
- `tensorflow_keras_metrics_binary_accuracy`
- `tensorflow_keras_metrics_accuracy`
- `tensorflow_keras_metrics_auc`
- `tensorflow_keras_metrics_mean`

Method calls such as `model.compile(...)`, `model.fit(...)`, and
`model.summary()` use `tensorflow_call` with a method expression.

## Sequential Model

```hcl
data "tensorflow_keras_layers_dense" "hidden" {
  args = ["64"]
  kwargs = {
    activation = "\"relu\""
    input_shape = "(784,)"
  }
}

data "tensorflow_keras_layers_dropout" "dropout" {
  args = ["0.2"]
}

data "tensorflow_keras_layers_dense" "output" {
  args = ["10"]
  kwargs = {
    activation = "\"softmax\""
  }
}

data "tensorflow_keras_sequential" "model" {
  args = [
    "[${join(", ", [
      data.tensorflow_keras_layers_dense.hidden.expression,
      data.tensorflow_keras_layers_dropout.dropout.expression,
      data.tensorflow_keras_layers_dense.output.expression,
    ])}]",
  ]
}

data "tensorflow_assign" "model" {
  name  = "model"
  value = data.tensorflow_keras_sequential.model.expression
}
```

Generated Python:

```python
model = tf.keras.Sequential([
  tf.keras.layers.Dense(64, activation="relu", input_shape=(784,)),
  tf.keras.layers.Dropout(0.2),
  tf.keras.layers.Dense(10, activation="softmax"),
])
```

## Compile

```hcl
data "tensorflow_keras_optimizers_adam" "optimizer" {
  kwargs = {
    learning_rate = "0.001"
  }
}

data "tensorflow_keras_losses_sparse_categorical_crossentropy" "loss" {}

data "tensorflow_keras_metrics_sparse_categorical_accuracy" "accuracy" {}

data "tensorflow_call" "compile" {
  function = "${data.tensorflow_assign.model.expression}.compile"
  kwargs = {
    optimizer = data.tensorflow_keras_optimizers_adam.optimizer.expression
    loss      = data.tensorflow_keras_losses_sparse_categorical_crossentropy.loss.expression
    metrics   = "[${data.tensorflow_keras_metrics_sparse_categorical_accuracy.accuracy.expression}]"
  }
}
```

Use `data.tensorflow_call.compile.statement` in `tensorflow_program`.

See [examples/keras](../../examples/keras).
