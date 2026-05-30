terraform {
  required_providers {
    tensorflow = {
      source = "SPShota/tensorflow"
    }
  }
}

provider "tensorflow" {}

data "tensorflow_literal" "values" {
  value_json = jsonencode([1, 2, 3])
}

data "tensorflow_ref" "float32" {
  name = "tf.float32"
}

data "tensorflow_constant" "x" {
  args = [data.tensorflow_literal.values.expression]
  kwargs = {
    dtype = data.tensorflow_ref.float32.expression
  }
}

data "tensorflow_assign" "x" {
  name  = "x"
  value = data.tensorflow_constant.x.expression
}

data "tensorflow_math_reduce_sum" "total" {
  args = [data.tensorflow_assign.x.expression]
}

data "tensorflow_program" "main" {
  statements = [
    data.tensorflow_assign.x.statement,
    "print(${data.tensorflow_math_reduce_sum.total.expression})",
  ]
}

output "python" {
  value = data.tensorflow_program.main.content
}
