package main

import (
	//"errors"
	"fmt"
	//"strconv"
	"time"

	dc "github.com/fsouza/go-dockerclient"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

//var (
//	creationTime time.Time
//)

var (
	initiedSwarm swarm.Swarm
	ctx    context.Context
)

func resourceDockerSwarmInit(d *schema.ResourceData, meta interface{}) error {
	var err error
	client := meta.(*dc.Client)

	orchestrationConfig := swarm.OrchestrationConfig{
		TaskHistoryRetentionLimit: d.Get("task_history_retention_limit").(*int64),
	}

	raftConfig := swarm.RaftConfig{
		SnapshotInterval: d.Get("snapshot_interval").(uint64),
		KeepOldSnapshots: d.Get("keep_old_snapshots").(*uint64),
		LogEntriesForSlowFollowers: d.Get("log_entries_for_slow_followers").(uint64),
		ElectionTick: d.Get("election_tick").(int),
		HeartbeatTick: d.Get("heartbeat_tick").(int),
	}

	dispatcherConfig := swarm.DispatcherConfig{
		HeartbeatPeriod: d.Get("heartbeat_period").(time.Duration),
	}

	caConfig := swarm.CAConfig{
		NodeCertExpiry: d.Get("node_cert_expiry").(time.Duration),
		// TODO: externalCAs section
		//ExternalCAs:
	}

	taskDefaults := swarm.TaskDefaults{
		// TODO: LogDriver section
		//LogDriver: d.get("log_driver").(string),
	}

	swarmInitRequestSpec := swarm.Spec{
		Orchestration: orchestrationConfig,
		Raft: raftConfig,
		Dispatcher: dispatcherConfig,
		CAConfig: caConfig,
		TaskDefaults: taskDefaults,
		//  TODO: EncryptionConfig section
		//EncryptionConfig: &swarm.EncryptionConfig{
		//	AutoLockManagers: d.get("auto_lock_managers").(bool),
		//},
	}

	swarInitRequest := swarm.InitRequest{
		ListenAddr: d.Get("listen_address").(string),
		AdvertiseAddr: d.Get("advertise_address").(string),
		ForceNewCluster: d.Get("force_new_cluster").(bool),
		Spec: swarmInitRequestSpec,
		//AutoLockManagers: d.Get("auto_lock_managers").(bool),
	}

	//if v, ok := d.GetOk("env"); ok {
	//	createOpts.Config.Env = stringSetToStringSlice(v.(*schema.Set))
	//}

	//var ctx := context.Context
	initSwarmOptions := dc.InitSwarmOptions{
		swarInitRequest,
		ctx,
	}

	var swarmResp string
	if swarmResp, err = client.InitSwarm(initSwarmOptions); err != nil {
		return fmt.Errorf("Unable to init swarm: %s", err)
	}
	if swarmResp == "" {
		return fmt.Errorf("Returned swarm response is nil")
	}


	initionErr := resourceDockerSwarmInspect(d, meta)
	if initionErr != nil {
		return fmt.Errorf("Returned swarm inition: %s", initionErr)
	}

	d.SetId(initiedSwarm.ClusterInfo.ID)

	return nil
}

func resourceDockerSwarmInspect(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*dc.Client)
	//ctx := context.Context{}

	var swarm swarm.Swarm

	swarm, err := client.InspectSwarm(ctx)

	if err != nil {
		return fmt.Errorf("Error inspecting swarm: %s", err)
	}

	initiedSwarm = swarm

	return nil
}

func resourceDockerSwarmUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDockerSwarmLeave(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*dc.Client)

	// Stop the container before removing if destroy_grace_seconds is defined
	if d.Get("destroy_grace_seconds").(int) > 0 {
		var timeout = uint(d.Get("destroy_grace_seconds").(int))
		if err := client.StopContainer(d.Id(), timeout); err != nil {
			return fmt.Errorf("Error stopping container %s: %s", d.Id(), err)
		}
	}

	removeOpts := dc.RemoveContainerOptions{
		ID:            d.Id(),
		RemoveVolumes: true,
		Force:         true,
	}

	if err := client.RemoveContainer(removeOpts); err != nil {
		return fmt.Errorf("Error deleting container %s: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

