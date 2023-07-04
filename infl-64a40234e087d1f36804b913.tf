module "lots_of_vars" {
  source = "github.com/infralight/hcl-indexing-test-repo//modules/lots_of_vars"

  number_without_default      = 3
  with_description_no_default = "X"
  with_validation             = "X"
  with_ugly_validation        = "Z"
  nullable_string             = "X"
}
