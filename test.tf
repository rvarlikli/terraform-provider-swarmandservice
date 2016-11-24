resource "ciscodocker_service" "dev" {
  api_address = "192.168.99.100"
  api_port=2375
  image_name="kitematic/hello-world-nginx"
  service_name="hello"
  replica_count=1
  published_port=8089
  target_port=80
}
