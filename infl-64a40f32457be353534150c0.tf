module "lots_of_vars" {
  source = "github.com/infralight/hcl-indexing-test-repo//modules/lots_of_vars"

  with_description_no_default = "foo"
  number_without_default      = 1
  nullable_string             = "foo"
  with_ugly_validation        = "y"
  with_validation             = "string"
  list_map_any_long_default = [
    {
      foo = "bar"
      bin = "baz"
      baz = "bin"
      bar = "foo"
    },
    {
      foo = "bar"
      bin = "baz"
      baz = "bin"
      bar = "foo"
    },
    {
      foo = "bar"
      bin = "baz"
      baz = "bin"
      bar = "foo"
    },
    {
      foo = "bar"
      bin = "baz"
      baz = "bin"
      bar = "foo"
    },
  ]
}
