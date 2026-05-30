# Additional TensorFlow Namespaces

TF.tf includes generated wrappers for initial subsets of several TensorFlow
namespaces beyond tensor, Keras, and `tf.data` APIs.

## `tf.random`

- `tensorflow_random_categorical`
- `tensorflow_random_normal`
- `tensorflow_random_set_seed`
- `tensorflow_random_shuffle`
- `tensorflow_random_uniform`

```hcl
data "tensorflow_random_normal" "noise" {
  args = ["[2, 3]"]
}
```

## `tf.image`

- `tensorflow_image_adjust_brightness`
- `tensorflow_image_convert_image_dtype`
- `tensorflow_image_crop_to_bounding_box`
- `tensorflow_image_decode_image`
- `tensorflow_image_flip_left_right`
- `tensorflow_image_random_flip_left_right`
- `tensorflow_image_resize`

```hcl
data "tensorflow_io_read_file" "image_bytes" {
  args = ["\"image.png\""]
}

data "tensorflow_image_decode_image" "image" {
  args = [data.tensorflow_io_read_file.image_bytes.expression]
}

data "tensorflow_image_resize" "resized" {
  args = [
    data.tensorflow_image_decode_image.image.expression,
    "[224, 224]",
  ]
}
```

## `tf.io`

- `tensorflow_io_decode_csv`
- `tensorflow_io_decode_json_example`
- `tensorflow_io_decode_raw`
- `tensorflow_io_encode_base64`
- `tensorflow_io_parse_example`
- `tensorflow_io_read_file`
- `tensorflow_io_write_file`

## `tf.audio`

- `tensorflow_audio_decode_wav`
- `tensorflow_audio_encode_wav`

## `tf.strings`

- `tensorflow_strings_join`
- `tensorflow_strings_length`
- `tensorflow_strings_lower`
- `tensorflow_strings_regex_replace`
- `tensorflow_strings_split`
- `tensorflow_strings_strip`
- `tensorflow_strings_to_number`
- `tensorflow_strings_upper`

```hcl
data "tensorflow_literal" "words" {
  value_json = jsonencode(["1", "2", "3"])
}

data "tensorflow_strings_to_number" "numbers" {
  args = [data.tensorflow_literal.words.expression]
}
```

## `tf.sparse`

- `tensorflow_sparse_add`
- `tensorflow_sparse_concat`
- `tensorflow_sparse_reorder`
- `tensorflow_sparse_reshape`
- `tensorflow_sparse_reset_shape`
- `tensorflow_sparse_to_dense`

## `tf.ragged`

- `tensorflow_ragged_boolean_mask`
- `tensorflow_ragged_constant`
- `tensorflow_ragged_map_flat_values`

See [examples/namespaces](../../examples/namespaces).
