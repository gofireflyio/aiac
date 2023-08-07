module "ha_secured_instances" {
  source = "github.com/infralight/well-architected-env//modules/ha_secured_instances"

  instance_name = "eranbibi"
  cluster_size  = 3
  instance_type = "R4_LARGE"
}
