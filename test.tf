resource "ciscodocker_service" "dev" {
  api_address = "192.168.99.100"
  api_port=2375
  image_name="kitematic/hello-world-nginx"
  service_name="hello"
  replica_count=1
  env = ["SERVICE=elastic", "PROJECT=stage", "ENVIRONMENT=operations"]
  ports = {
    published = 8089
    target = 80
    protocol = "tcp"
  }
  ports = {
    published = 9090
    target = 8080
    protocol = "tcp"
  }
}
