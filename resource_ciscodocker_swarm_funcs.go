package main

import (
	//"errors"
	"fmt"
	//"strconv"
	"time"
	"log"

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

	var orchestrationConfig swarm.OrchestrationConfig
	if v, ok := d.GetOk("task_history_retention_limit"); ok {
		limit := int64(v.(int))
		orchestrationConfig.TaskHistoryRetentionLimit = &limit
	}

	var raftConfig swarm.RaftConfig
	if v, ok := d.GetOk("snapshot_interval"); ok {
		interval := uint64(v.(int))
		raftConfig.SnapshotInterval = interval
	}
	if v, ok := d.GetOk("keep_old_snapshots"); ok {
		snapshots := uint64(v.(int))
		raftConfig.KeepOldSnapshots = &snapshots
	}
	if v, ok := d.GetOk("log_entries_for_slow_followers"); ok {
		followers := uint64(v.(int))
		raftConfig.LogEntriesForSlowFollowers = followers
	}
	if v, ok := d.GetOk("election_tick"); ok {
		election := v.(int)
		raftConfig.ElectionTick = election
	}
	if v, ok := d.GetOk("heartbeat_tick"); ok {
		heartbeat := v.(int)
		raftConfig.HeartbeatTick = heartbeat
	}

	var dispatcherConfig swarm.DispatcherConfig
	if v, ok := d.GetOk("heartbeat_period"); ok {
		heartbeat := time.Duration(v.(int))
		dispatcherConfig.HeartbeatPeriod = heartbeat
	}


	var caConfig swarm.CAConfig
	if v, ok := d.GetOk("node_cert_expiry"); ok {
		expiry := time.Duration(v.(int))
		caConfig.NodeCertExpiry = expiry
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

	log.Println("initSwarmOptions............")

	var swarmResp string
	if swarmResp, err = client.InitSwarm(initSwarmOptions); err != nil {
		return fmt.Errorf("Unable to init swarm: %s", err)
	}
	if swarmResp == "" {
		return fmt.Errorf("Returned swarm response is nil")
	}


	inspectErr := resourceDockerSwarmInspect(d, meta)
	if inspectErr != nil {
		return fmt.Errorf("Returned swarm inpect: %s", inspectErr)
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

	d.Set("manager_token", initiedSwarm.JoinTokens.Manager)
	d.Set("worker_token", initiedSwarm.JoinTokens.Worker)
	d.Set("version", initiedSwarm.ClusterInfo.Meta.Version.Index)

	return nil
}

func resourceDockerSwarmUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Println("Update started....")
	d.Partial(true)

	if d.HasChange("task_history_retention_limit") {
		d.GetChange("task_history_retention_limit")
		d.SetPartial("task_history_retention_limit")
	}
	if d.HasChange("snapshot_interval") {
		d.GetChange("snapshot_interval")

		d.SetPartial("snapshot_interval")
	}
	if d.HasChange("keep_old_snapshots") {
		d.GetChange("keep_old_snapshots")

		d.SetPartial("keep_old_snapshots")
	}
	if d.HasChange("log_entries_for_slow_followers") {
		d.GetChange("log_entries_for_slow_followers")

		d.SetPartial("log_entries_for_slow_followers")
	}
	if d.HasChange("election_tick") {
		d.GetChange("election_tick")

		d.SetPartial("election_tick")
	}
	if d.HasChange("heartbeat_tick") {


		d.SetPartial("heartbeat_tick")
	}
	if d.HasChange("heartbeat_period") {
		d.GetChange("heartbeat_period")

		d.SetPartial("heartbeat_period")
	}
	if d.HasChange("node_cert_expiry") {
		d.GetChange("node_cert_expiry")
		d.SetPartial("node_cert_expiry")
	}

	if err := updateSwarm(d, meta); err != nil {
		return err
	}
	d.Partial(false)
	return nil
}

func resourceDockerSwarmLeave(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*dc.Client)

	force_leave := d.Get("force_leave").(bool)

	leaveSwarmOptions := dc.LeaveSwarmOptions{
		force_leave,
		ctx,
	}

	if leaveErr := client.LeaveSwarm(leaveSwarmOptions); leaveErr != nil {
		return fmt.Errorf("Unable to leave swarm: %s", leaveErr)
	}

	d.SetId("")
	return nil
}

func updateSwarm(d *schema.ResourceData, meta interface{}) error {
	var err error
	client := meta.(*dc.Client)

	var orchestrationConfig swarm.OrchestrationConfig
	if v, ok := d.GetOk("task_history_retention_limit"); ok {
		limit := int64(v.(int))
		orchestrationConfig.TaskHistoryRetentionLimit = &limit
	}

	var raftConfig swarm.RaftConfig
	if v, ok := d.GetOk("snapshot_interval"); ok {
		interval := uint64(v.(int))
		raftConfig.SnapshotInterval = interval
	}

	if v, ok := d.GetOk("keep_old_snapshots"); ok {
		snapshots := uint64(v.(int))
		raftConfig.KeepOldSnapshots = &snapshots
	}

	if v, ok := d.GetOk("log_entries_for_slow_followers"); ok {
		followers := uint64(v.(int))
		raftConfig.LogEntriesForSlowFollowers = followers
	}

	if v, ok := d.GetOk("election_tick"); ok {
		election := v.(int)
		raftConfig.ElectionTick = election
	}

	if v, ok := d.GetOk("heartbeat_tick"); ok {
		heartbeat := v.(int)
		raftConfig.HeartbeatTick = heartbeat
	}

	var dispatcherConfig swarm.DispatcherConfig
	if v, ok := d.GetOk("heartbeat_period"); ok {
		heartbeat := time.Duration(v.(int))
		dispatcherConfig.HeartbeatPeriod = heartbeat
	}

	var caConfig swarm.CAConfig
	if v, ok := d.GetOk("node_cert_expiry"); ok {
		expiry := time.Duration(v.(int))
		caConfig.NodeCertExpiry = expiry
	}


	taskDefaults := swarm.TaskDefaults{
		// TODO: LogDriver section
		//LogDriver: d.get("log_driver").(string),
	}

	swarmUpdateSpec := swarm.Spec{
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


	updateSwarmOptions := dc.UpdateSwarmOptions{
		Version : d.Get("version").(int),
		RotateWorkerToken: d.Get("rotate_worker_token").(bool),
		RotateManagerToken: d.Get("rotate_manager_token").(bool),
		Swarm: swarmUpdateSpec,
		Context: ctx,
	}

	if err = client.UpdateSwarm(updateSwarmOptions); err != nil {
		return fmt.Errorf("Unable to update swarm: %s", err)
	}

	inspectErr := resourceDockerSwarmInspect(d, meta)
	if inspectErr != nil {
		return fmt.Errorf("Returned swarm inition: %s", inspectErr)
	}

	d.Set("manager_token", initiedSwarm.JoinTokens.Manager)
	d.Set("worker_token", initiedSwarm.JoinTokens.Worker)
	d.Set("version", initiedSwarm.ClusterInfo.Meta.Version.Index)
	d.SetId(initiedSwarm.ClusterInfo.ID)

	return nil
}

