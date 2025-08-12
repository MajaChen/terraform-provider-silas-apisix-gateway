terraform {
  required_providers {
    apisix = {
      source = "silas.com/ssf/apisix-gateway"
    }
  }
  required_version = ">= 1.1.0"
}

provider "apisix" {
  env = "local"
}

resource "apisix_route" "ssf-java-sdk-springboot3-demo-dynLoggingLevel" {
  id          = "ssf-java-sdk-springboot3-demo-dynLoggingLevel"
  uris = ["/api/v1/demo/dynLoggingLevel"]
  upstream_id = "1"
  plugins = {
    openid_connect = {
      client_id = "client-id"
      discovery = "https://domain.authing.cn/oidc/.well-known/jwks.json"
      required_scopes = ["admin", "book:read", "book:write", "book:read"]
    }
  }
  name     = "ssf-java-sdk-springboot3-demo-dynLoggingLevel"
  desc     = "ssf-java-sdk-springboot3-demo dynLoggingLevel"
  methods = ["GET"]
  priority = 10
  vars = [["http_user", "==", "ios"]]
  timeout = {
    connect = 10
    send    = 30
    read    = 30
  }
  status = 1
}
