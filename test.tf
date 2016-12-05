provider "ciscodocker" {
  host = "tcp://192.168.99.100:2376/"
}

resource "ciscodocker_swarm" "swarm" {
  listen_address = "0.0.0.0:2377"
  advertise_address = "192.168.99.100:2377"
  force_new_cluster = false
  force_leave = true
}
