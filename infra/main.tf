terraform {
  required_version = ">= 1.5"
  required_providers {
    aiven = { source = "aiven/aiven", version = "~> 4.0" }
    render = { source = "render-oss/render", version = "1.9.1" }
  }
}
provider "aiven" {}
provider "render" {}

data "aiven_mysql" "qa_service" {
  project      = var.aiven_project
  service_name = var.aiven_service
}

resource "aiven_mysql_database" "preprod" {
  project       = var.aiven_project
  service_name  = var.aiven_service
  database_name = "app_preprod"
}

resource "aiven_mysql_user" "preprod" {
  project       = var.aiven_project
  service_name  = var.aiven_service
  username      = "preprod_user"
}

locals {
  conn = "${aiven_mysql_user.preprod.username}:${aiven_mysql_user.preprod.password}@tcp(${data.aiven_mysql.qa_service.service_host}:${data.aiven_mysql.qa_service.service_port})/${aiven_mysql_database.preprod.database_name}?tls=skip-verify"
}

resource "render_web_service" "api" {
  name   = "miapp-api-preprod"
  plan   = "free"
  region = var.region
  runtime_source = {
    image = { image_url = var.image_repo_api, tag = var.image_tag }
  }
  env_vars = {
    DB_DSN = { value = local.conn }
  }
}

resource "render_web_service" "front" {
  name   = "miapp-front-preprod"
  plan   = "free"
  region = var.region
  runtime_source = {
    image = { image_url = var.image_repo_front, tag = var.image_tag }
  }
  env_vars = {
    BACKEND_URL  = { value = render_web_service.api.url }
    DNS_RESOLVER = { value = "8.8.8.8" }
  }
}
