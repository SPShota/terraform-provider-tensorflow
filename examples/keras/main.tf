terraform {
  required_providers {
    tensorflow = {
      source = "SPShota/tensorflow"
    }
  }
}

provider "tensorflow" {}

data "tensorflow_keras_layers_dense" "hidden" {
  args = ["64"]
  kwargs = {
    activation  = "\"relu\""
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

data "tensorflow_call" "summary" {
  function = "${data.tensorflow_assign.model.expression}.summary"
}

data "tensorflow_program" "main" {
  statements = [
    data.tensorflow_assign.model.statement,
    data.tensorflow_call.compile.statement,
    data.tensorflow_call.summary.statement,
  ]
}

output "python" {
  value = data.tensorflow_program.main.content
}
