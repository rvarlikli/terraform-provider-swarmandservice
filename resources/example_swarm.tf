resource "ciscodocker_swarm" "swarm" {
  api_address = "192.168.99.100"
  api_port = 2375
  listen_address = "0.0.0.0:2377"
  advertise_address = "192.168.99.100:2377"
  force_new = true
}