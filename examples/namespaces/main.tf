terraform {
  required_providers {
    tensorflow = {
      source = "SPShota/tensorflow"
    }
  }
}

provider "tensorflow" {}

data "tensorflow_random_set_seed" "seed" {
  args = ["42"]
}

data "tensorflow_random_normal" "noise" {
  args = ["[2, 3]"]
}

data "tensorflow_assign" "noise" {
  name  = "noise"
  value = data.tensorflow_random_normal.noise.expression
}

data "tensorflow_literal" "words" {
  value_json = jsonencode(["1", "2", "3"])
}

data "tensorflow_strings_to_number" "numbers" {
  args = [data.tensorflow_literal.words.expression]
}

data "tensorflow_assign" "numbers" {
  name  = "numbers"
  value = data.tensorflow_strings_to_number.numbers.expression
}

data "tensorflow_ragged_constant" "ragged" {
  args = ["[[1, 2], [3]]"]
}

data "tensorflow_assign" "ragged" {
  name  = "ragged"
  value = data.tensorflow_ragged_constant.ragged.expression
}

data "tensorflow_program" "main" {
  statements = [
    data.tensorflow_random_set_seed.seed.statement,
    data.tensorflow_assign.noise.statement,
    data.tensorflow_assign.numbers.statement,
    data.tensorflow_assign.ragged.statement,
    "print(noise)",
    "print(numbers)",
    "print(ragged)",
  ]
}

output "python" {
  value = data.tensorflow_program.main.content
}
