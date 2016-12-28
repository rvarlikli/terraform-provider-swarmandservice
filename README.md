# Terraform Provider for Docker Swarm and Service
The provider creates Docker Swarm clusters and services.

### Prerequisites
-[Docker]

-[Terraform]

-[Go]

#### Steps
1-) Clone the repo and cd to redo dir

2-) Set GOPATH

`export GOPATH=/repo/dir`

2-) Get provider dependencies

`go get github.com/hashicorp/terraform/helper/schema`

`go get github.com/docker/docker/api/types/swarm`

`go get github.com/fsouza/go-dockerclient`

github.com/docker

3-) Build the provider with name

`go build -o terraform-provider-ciscodocker`

4-) Create Docker host machine(s) for Docker Swarm Cluster

`docker-machine create --driver virtualbox manager1`

`docker-machine create --driver virtualbox node1`

5-) Create Swarm Cluster and Docker Services with terraform

`terraform apply`

##### About Docker Images
###### Java Image(cloudnesil/openjdk:8u111-jdk)
Java image runs a sample jar file that is a java application with [spring-boot-rest-example]
Image has a simple bash script to download the jar file from github repo and run it. The bash script gets branch value as arguments.

`docker run -d --name java-sample -p 8080:8090 cloudnesil/openjdk:8u111-jdk startapp.sh master`

###### GitLab-ce(docker pull gitlab/gitlab-ce)
Official [GitLab] image


[//]: # (These are reference links used in the body of this note and get stripped out when the markdown processor does its job. There is no need to format nicely because it shouldn't be seen. Thanks SO - http://stackoverflow.com/questions/4823468/store-comments-in-markdown-syntax)


   [docker]: <https://www.docker.com/products/overview>
   [terraform]: <https://www.terraform.io/downloads.html>
   [spring-boot-rest-example]: <https://github.com/rvarlikli/spring-boot-rest-example>
   [Gitlab]: <https://hub.docker.com/r/gitlab/gitlab-ce/>
   [Go]: <https://golang.org/dl/>
   
