package main

import (
	//"bytes"
	//"fmt"
	//"regexp"
	//"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCiscoDockerSwarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceDockerSwarmInit,
		Read:   resourceDockerSwarmInspect,
		Update: resourceDockerSwarmUpdate,
		Delete: resourceDockerSwarmLeave,

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
			"force_new_cluster": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default: true,
			},
			"auto_lock_managers": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default: false,
			},
			"task_history_retention_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default: 10,
			},
			"snapshot_interval": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default: 10000,
			},
			"log_entries_for_slow_followers": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default: 500,
			},
			"election_tick": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default: 3,
			},
			"heartbeat_tick": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default: 1,
			},
			"heartbeat_period": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default: 5000000000,
			},
			"node_cert_expiry": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default: 7776000000000000,
			},

		},
	}
}

