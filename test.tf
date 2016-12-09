provider "ciscodocker" {
  alias ="manager"
  host = "tcp://192.168.99.100:2376/"
  cert_path = "/Users/kadirtaskiran/.docker/machine/machines/manager1"
}

provider "ciscodocker" {
  alias = "node1"
  host = "tcp://192.168.99.101:2376/"
  cert_path = "/Users/kadirtaskiran/.docker/machine/machines/node1"
}

resource "ciscodocker_swarm" "manager" {
  provider = "ciscodocker.manager"
  listen_address = "0.0.0.0:2377"
  advertise_address = "192.168.99.100:2377"
  force_new_cluster = false
  force_leave = true
  rotate_manager_token = true
  heartbeat_tick = 2
}

resource "ciscodocker_swarmnode" "node1" {
  depends_on = ["ciscodocker_swarm.manager"]
  provider = "ciscodocker.node1"
  listen_address = "0.0.0.0:2377"
  advertise_address = "192.168.99.101:2377"
  is_manager = false
  manager_token = "${ciscodocker_swarm.manager.manager_token}"
  worker_token = "${ciscodocker_swarm.manager.worker_token}"
  remote_address = ["192.168.99.100:2377"]
}

resource "ciscodocker_service" "hello" {
  depends_on = ["ciscodocker_swarm.manager", "ciscodocker_swarmnode.node1"]
  provider = "ciscodocker.manager"
  service_name = "hello"
  image_name = "kitematic/hello-world-nginx"
  restart_policy_condition = "any"
  service_replica_count = 1
  ports = {
    published_port = 8081
    target_port = 80
    protocol = "tcp"
  }

}