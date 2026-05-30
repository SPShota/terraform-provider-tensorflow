terraform {
  required_providers {
    tensorflow = {
      source = "SPShota/tensorflow"
    }
  }
}

provider "tensorflow" {}

data "tensorflow_literal" "left" {
  value_json = jsonencode(1)
}

data "tensorflow_literal" "right" {
  value_json = jsonencode(2)
}

data "tensorflow_raw_op" "sum" {
  op = "AddV2"
  kwargs = {
    x = data.tensorflow_literal.left.expression
    y = data.tensorflow_literal.right.expression
  }
}

data "tensorflow_assign" "sum" {
  name  = "sum_value"
  value = data.tensorflow_raw_op.sum.expression
}

data "tensorflow_program" "main" {
  statements = [
    data.tensorflow_assign.sum.statement,
    "print(sum_value)",
  ]
}

output "python" {
  value = data.tensorflow_program.main.content
}
