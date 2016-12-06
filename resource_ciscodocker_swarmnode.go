package main

import (
	//"bytes"
	//"fmt"
	//"regexp"
	//"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCiscoDockerSwarmNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceDockerSwarmNodeJoin,
		Read:   resourceDockerSwarmNodeInspect,
		Delete: resourceDockerSwarmNodeLeave,

		Schema: map[string]*schema.Schema{
			"listen_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"advertise_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_manager": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: false,
				ForceNew: true,
			},
			"manager_token": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"worker_token": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remote_address": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

