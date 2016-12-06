package main

import (
	//"errors"
	"fmt"
	//"strconv"
	"log"

	dc "github.com/fsouza/go-dockerclient"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

var (
	nodeCtx    context.Context
)

func resourceDockerSwarmNodeJoin(d *schema.ResourceData, meta interface{}) error {
	var err error
	client := meta.(*dc.Client)

	swarmJoinRequest := swarm.JoinRequest{
		ListenAddr: d.Get("listen_address").(string),
		AdvertiseAddr: d.Get("advertise_address").(string),
	}

	is_manager := d.Get("is_manager").(bool)

	if is_manager {
		if v, ok := d.GetOk("manager_token"); ok {
			swarmJoinRequest.JoinToken = v.(string)
		}
	}

	if is_manager == false {
		if v, ok := d.GetOk("worker_token"); ok {
			swarmJoinRequest.JoinToken = v.(string)
		}
	}

	if v, ok := d.GetOk("remote_address"); ok {
		swarmJoinRequest.RemoteAddrs = stringListToStringSlice(v.([]interface{}))
	}


	joinSwarmOptions := dc.JoinSwarmOptions{
		swarmJoinRequest,
		ctx,
	}

	log.Println("joinSwarmOptions............")

	if err = client.JoinSwarm(joinSwarmOptions); err != nil {
		return fmt.Errorf("Unable to join swarm: %s", err)
	}
	//TODO: set swarmnode id
	d.SetId("swarmnode-"+d.Get("advertise_address").(string))
	return nil
}



func resourceDockerSwarmNodeLeave(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*dc.Client)

	force_leave := true

	leaveSwarmOptions := dc.LeaveSwarmOptions{
		force_leave,
		nodeCtx,
	}

	log.Println("LeaveSwarmOptions............")

	if leaveErr := client.LeaveSwarm(leaveSwarmOptions); leaveErr != nil {
		return fmt.Errorf("Unable to leave swarm: %s", leaveErr)
	}

	d.SetId("")
	return nil
}

func resourceDockerSwarmNodeInspect(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func stringSetToStringSlice(stringSet *schema.Set) []string {
	ret := []string{}
	if stringSet == nil {
		return ret
	}
	for _, envVal := range stringSet.List() {
		ret = append(ret, envVal.(string))
	}
	return ret
}

func stringListToStringSlice(stringList []interface{}) []string {
	ret := []string{}
	for _, v := range stringList {
		if v == nil {
			ret = append(ret, "")
			continue
		}
		ret = append(ret, v.(string))
	}
	return ret
}
