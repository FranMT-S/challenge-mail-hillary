terraform {
  required_providers {
    cockroach = {
      source  = "cockroachdb/cockroach"
      version = "~> 1.10.0"
    }
  }
}

provider "cockroach" {
  apikey = var.cockroach_apikey
}

resource "cockroach_cluster" "serverless" {
  name           = var.cluster_name
  cloud_provider = "AWS"
  plan           = "BASIC"
  serverless     = {}
  regions = [
    {
      name = var.region
    }
  ]
  delete_protection = false
}


resource "cockroach_sql_user" "cockroach" {
  name       = var.sql_user_name
  password   = var.sql_user_password
  cluster_id = cockroach_cluster.serverless.id
}

resource "cockroach_database" "cockroach" {
  count      = var.db_name == "defaultdb" ? 0 : 1 # not create if defaultdb is the name
  name       = var.db_name
  cluster_id = cockroach_cluster.serverless.id
}


data "cockroach_connection_string" "db" {
  id       = cockroach_cluster.serverless.id
  sql_user = cockroach_sql_user.cockroach.name
  password = cockroach_sql_user.cockroach.password
  database = var.db_name != "" ? var.db_name : "defaultdb"
  os       = "LINUX"
}


