terraform {
  required_providers {
    tensorflow = {
      source = "SPShota/tensorflow"
    }
  }
}

provider "tensorflow" {}

data "tensorflow_literal" "features" {
  value_json = jsonencode([[1.0, 2.0], [3.0, 4.0], [5.0, 6.0], [7.0, 8.0]])
}

data "tensorflow_literal" "labels" {
  value_json = jsonencode([0, 1, 0, 1])
}

data "tensorflow_data_dataset_from_tensor_slices" "dataset" {
  args = [
    "(${data.tensorflow_literal.features.expression}, ${data.tensorflow_literal.labels.expression})",
  ]
}

data "tensorflow_call" "cache" {
  function = "${data.tensorflow_data_dataset_from_tensor_slices.dataset.expression}.cache"
}

data "tensorflow_call" "batch" {
  function = "${data.tensorflow_call.cache.expression}.batch"
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

data "tensorflow_program" "main" {
  statements = [
    data.tensorflow_assign.dataset.statement,
    "print(dataset)",
  ]
}

output "python" {
  value = data.tensorflow_program.main.content
}
