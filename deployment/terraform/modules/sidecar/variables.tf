variable "sidecar_injector_image" {
  description = "sidecar_injector_image"
  type = map
  default = {
    repository = "khuong02/sidecar-injector"
    tag = "0.1.0"
  }
}