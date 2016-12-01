package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"ciscodocker_service": resourceCiscoDockerService(),
			"ciscodocker_swarm": resourceCiscoDockerSwarm(),
			"ciscodocker_swarmnode": resourceCiscoDockerSwarmNode(),
		},
	}
}
