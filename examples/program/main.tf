terraform {
  required_providers {
    tensorflow = {
      source = "SPShota/tensorflow"
    }
  }
}

provider "tensorflow" {}

data "tensorflow_program" "main" {
  statements = [
    "x = tf.constant([1, 2, 3])",
    "print(tf.reduce_sum(x))",
  ]
}

output "python" {
  value = data.tensorflow_program.main.content
}
