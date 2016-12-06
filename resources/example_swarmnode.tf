provider "ciscodocker" {
  alias = "node1"
  host = "tcp://192.168.99.101:2376/"
  cert_path = "/Users/kadirtaskiran/.docker/machine/machines/node1"
}

resource "ciscodocker_swarmnode" "node1" {
  provider = "ciscodocker.node1"
  listen_address = "0.0.0.0:2377"
  advertise_address = "192.168.99.101:2377"
  is_manager = false
  //manager_token = comes from swarm manager / string
  //worker_token = comes from swarm manager / string
  remote_address = ["192.168.99.100:2377"]
}