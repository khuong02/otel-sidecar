resource "helm_release" "sidecar-injector" {
  name       = "sidecar-injector"

  chart      = "${path.module}/../../../charts/sidecar-injector"
  namespace  = "sidecar-injector"

  set {
    name  = "image.repository"
    value = lookup(var.sidecar_injector_image, "repository")
  }
  
  set {
    name  = "image.tag"
    value = lookup(var.sidecar_injector_image, "tag")
  }
}