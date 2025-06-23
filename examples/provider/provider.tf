terraform {
  required_providers {
    nttdata-rest-api = {
      source = "nttdata-rest-api"
      version = "~> 1.0.0"
    }
  }
}

provider "nttdata-rest-api" {
    base_url = "https://jsonplaceholder.typicode.com"
    #auth_token = "token"
}

resource "nttdata-rest-api_apiresource" "edu" {
    endpoint_path = "posts/1"
    payload = jsonencode({
        title = "foo"
        body = "bar"
        userId = 1
    })
}

output "response" {
    value = nttdata-rest-api_apiresource.edu.response
}
