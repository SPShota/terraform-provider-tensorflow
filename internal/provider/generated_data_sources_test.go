package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestGeneratedDataSources(t *testing.T) {
	t.Parallel()

	dataSources := GeneratedDataSources()
	if len(dataSources) != 114 {
		t.Fatalf("len(GeneratedDataSources()) = %d, want 114", len(dataSources))
	}

	got := make(map[string]struct{}, len(dataSources))
	for _, newDataSource := range dataSources {
		ds := newDataSource()
		var resp datasource.MetadataResponse
		ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "tensorflow"}, &resp)
		got[resp.TypeName] = struct{}{}
	}

	for _, name := range []string{
		"tensorflow_constant",
		"tensorflow_convert_to_tensor",
		"tensorflow_cast",
		"tensorflow_reshape",
		"tensorflow_concat",
		"tensorflow_stack",
		"tensorflow_math_reduce_sum",
		"tensorflow_math_reduce_mean",
		"tensorflow_math_add",
		"tensorflow_variable",
		"tensorflow_gradient_tape",
		"tensorflow_module",
		"tensorflow_keras_sequential",
		"tensorflow_keras_model",
		"tensorflow_keras_input",
		"tensorflow_keras_layers_dense",
		"tensorflow_keras_layers_dropout",
		"tensorflow_keras_optimizers_adam",
		"tensorflow_keras_losses_sparse_categorical_crossentropy",
		"tensorflow_keras_metrics_sparse_categorical_accuracy",
		"tensorflow_data_dataset_from_tensor_slices",
		"tensorflow_data_dataset_from_tensors",
		"tensorflow_data_dataset_range",
		"tensorflow_data_dataset_zip",
		"tensorflow_data_dataset_list_files",
		"tensorflow_data_tf_record_dataset",
		"tensorflow_random_normal",
		"tensorflow_random_uniform",
		"tensorflow_image_resize",
		"tensorflow_image_decode_image",
		"tensorflow_io_read_file",
		"tensorflow_io_parse_example",
		"tensorflow_audio_decode_wav",
		"tensorflow_strings_split",
		"tensorflow_strings_to_number",
		"tensorflow_sparse_to_dense",
		"tensorflow_sparse_reorder",
		"tensorflow_ragged_constant",
		"tensorflow_ragged_map_flat_values",
		"tensorflow_zeros",
		"tensorflow_ones",
	} {
		if _, ok := got[name]; !ok {
			t.Fatalf("expected generated data source %q; got %#v", name, got)
		}
	}
}

func TestGeneratedWrapperSchema(t *testing.T) {
	t.Parallel()

	ds := GeneratedDataSources()[3]()
	var resp datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	if resp.Schema.MarkdownDescription == "" {
		t.Fatalf("expected generated wrapper schema description")
	}

	for _, name := range []string{"args", "kwargs", "expression", "statement"} {
		if _, ok := resp.Schema.Attributes[name]; !ok {
			t.Fatalf("schema is missing %q attribute", name)
		}
	}
}
