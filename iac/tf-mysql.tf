resource "helm_release" "nvr-mysql" {
  name       = "nvr-mysql"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "mariadb"
  version    = "11.4.4"
  force_update = "true"

  set {
    name  = "cluster.enabled"
    value = "true"
  }

  set {
    name  = "auth.database"
    value = "frigate_db"
  }
  set {
    name  = "MARIADB_ROOT_PASSWORD"
    value = "yourfrigate_pass"
  }

}
