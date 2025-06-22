terraform {
  required_providers {
    hashicups = {
      source = "edu/hashicups"
#      version = "1.0.0"
    }
  }
}

provider "hashicups" {}

data "hashicups_coffees" "example" {}
