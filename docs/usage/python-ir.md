# Python Code-Generation IR

PR2 introduces the internal Go package that later Terraform data sources will use
to generate Python code.

The package intentionally has no Terraform dependency. It models the smallest
Python fragments needed by upcoming data sources:

- `Expression`
- `Statement`
- Python literals
- references and dotted references
- attribute access
- function calls
- assignment

Example:

```go
values, err := python.Literal([]int{1, 2, 3})
if err != nil {
	return err
}

dtype, err := python.Reference("tf.float32")
if err != nil {
	return err
}

constant, err := python.CallNamed("tf.constant", []python.Expression{values}, []python.KeywordArgument{
	{Name: "dtype", Value: dtype},
})
if err != nil {
	return err
}

statement, err := python.Assign("x", constant)
if err != nil {
	return err
}

fmt.Println(statement.Code())
```

Output:

```python
x = tf.constant([1, 2, 3], dtype=tf.float32)
```

This package is deliberately conservative. It validates Python identifiers,
rejects unsupported literal types, and rejects non-finite float values.
