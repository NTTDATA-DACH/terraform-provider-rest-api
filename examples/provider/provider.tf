terraform {
  required_providers {
    hashicups = {
      source = "edu/hashicups"
#      version = "1.0.0"
    }
  }
}

provider "hashicups" {
    base_url = "https://jsonplaceholder.typicode.com"
    #auth_token = "token"
}

resource "hashicups_apiresource" "edu" {
    endpoint_path = "posts/1"
    payload = jsonencode({
        title = "foo"
        body = "bar"
        userId = 1
    })
}

output "response" {
    value = hashicups_apiresource.edu.response
}
