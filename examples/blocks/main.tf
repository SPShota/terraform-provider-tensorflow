terraform {
  required_providers {
    tensorflow = {
      source = "SPShota/tensorflow"
    }
  }
}

provider "tensorflow" {}

data "tensorflow_literal" "initial_weight" {
  value_json = jsonencode(3.0)
}

data "tensorflow_variable" "weight" {
  args = [data.tensorflow_literal.initial_weight.expression]
}

data "tensorflow_assign" "weight" {
  name  = "weight"
  value = data.tensorflow_variable.weight.expression
}

data "tensorflow_gradient_tape" "tape" {}

data "tensorflow_with" "tape" {
  context = data.tensorflow_gradient_tape.tape.expression
  alias   = "tape"

  statements = [
    "loss = weight * weight",
  ]
}

data "tensorflow_return" "loss" {
  value = "loss"
}

data "tensorflow_function" "train_step" {
  name       = "train_step"
  decorators = ["tf.function"]

  statements = [
    data.tensorflow_with.tape.statement,
    data.tensorflow_return.loss.statement,
  ]
}

data "tensorflow_program" "main" {
  statements = [
    data.tensorflow_assign.weight.statement,
    data.tensorflow_function.train_step.statement,
    "print(train_step())",
  ]
}

output "python" {
  value = data.tensorflow_program.main.content
}
