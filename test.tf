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

resource "ciscodocker_service" "gitlab-ce" {
  depends_on = ["ciscodocker_swarm.manager", "ciscodocker_swarmnode.node1"]
  provider = "ciscodocker.manager"
  service_name="gitlab"
  image_name="gitlab/gitlab-ce"
  restart_policy_condition = "any"
  service_replica_count=1
  env = ["GITLAB_OMNIBUS_CONFIG=\"external_url 'http://nesillocal.com:8088/'; gitlab_rails['gitlab_shell_ssh_port'] = 8922;\""]
  ports = {
    published_port = 8088
    target_port = 80
    protocol = "tcp"
  }
  ports = {
    published_port = 8443
    target_port = 443
    protocol = "tcp"
  }
  ports = {
    published_port = 8922
    target_port = 22
    protocol = "tcp"
  }
  resolution_mode = "vip"

}

resource "ciscodocker_service" "java-pie-app" {
  depends_on = ["ciscodocker_swarm.manager", "ciscodocker_swarmnode.node1"]
  provider = "ciscodocker.manager"
  service_name="java-pie-app"
  image_name="cloudnesil/openjdk:8u111-jdk"
  restart_policy_condition = "any"
  service_replica_count=1
  ports = {
    published_port = 8090
    target_port = 8080
    protocol = "tcp"
  }
  resolution_mode = "vip"
  command = ["startapp.sh"]
  args = ["master"]

}

resource "ciscodocker_service" "java-pie-appv2" {
  depends_on = ["ciscodocker_swarm.manager", "ciscodocker_swarmnode.node1"]
  provider = "ciscodocker.manager"
  service_name="java-pie-appv2"
  image_name="cloudnesil/openjdk:8u111-jdk"
  restart_policy_condition = "any"
  service_replica_count=1
  ports = {
    published_port = 8091
    target_port = 8080
    protocol = "tcp"
  }
  resolution_mode = "vip"
  command = ["startapp.sh"]
  args = ["v2"]

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
