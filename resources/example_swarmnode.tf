resource "ciscodocker_swarmnode" "swarm-node" {
  api_address = "192.168.99.101"
  api_port = 2375
  listen_address = "0.0.0.0:2377"
  advertise_address = "192.168.99.101:2377"
  is_swarm_manager = true
  swarm_manager_address = "192.168.99.100:2377"
}