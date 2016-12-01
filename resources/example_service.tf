resource "ciscodocker_service" "gitlab-ce" {
  api_address = "192.168.99.100"
  api_port=2375
  image_name="gitlab/gitlab-ce"
  service_name="gitlab"
  replica_count=1
  env = ["GITLAB_OMNIBUS_CONFIG=\"external_url 'http://nesillocal.com:8081/'; gitlab_rails['gitlab_shell_ssh_port'] = 8922;\""]
  ports = {
    published = 8081
    target = 80
    protocol = "tcp"
  }
  ports = {
    published = 8443
    target = 443
    protocol = "tcp"
  }
  ports = {
    published = 8922
    target = 22
    protocol = "tcp"
  }

}