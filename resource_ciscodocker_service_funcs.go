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
	serviceCtx    context.Context
	createdService *swarm.Service
)

func resourceDockerServiceCreate(d *schema.ResourceData, meta interface{}) error {
	var err error
	client := meta.(*dc.Client)

	var serviceSpec swarm.ServiceSpec
	if v, ok := d.GetOk("service_name"); ok {
		serviceSpec.Annotations.Name=v.(string)
	}

	if v, ok := d.GetOk("service_labels"); ok {
		serviceSpec.Annotations.Labels = mapTypeMapValsToString(v.(map[string]interface{}))
	}

	var taskTemplate swarm.TaskSpec
	var containerSpec swarm.ContainerSpec
	if v, ok := d.GetOk("image_name"); ok {
		containerSpec.Image = v.(string)
	}

	if v, ok := d.GetOk("container_labels"); ok {
		containerSpec.Labels = mapTypeMapValsToString(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("command"); ok {
		containerSpec.Command = stringListToStringSlice(v.([]interface{}))
		for _, v := range containerSpec.Command {
			if v == "" {
				return fmt.Errorf("values for command may not be empty")
			}
		}
	}

	if v, ok := d.GetOk("args"); ok {
		containerSpec.Args = stringListToStringSlice(v.([]interface{}))
	}

	if v, ok := d.GetOk("hostname"); ok {
		containerSpec.Hostname = v.(string)
	}

	if v, ok := d.GetOk("env"); ok {
		containerSpec.Env = stringListToStringSlice(v.([]interface{}))
	}

	if v, ok := d.GetOk("dir"); ok {
		containerSpec.Dir = v.(string)
	}

	if v, ok := d.GetOk("user"); ok {
		containerSpec.User = v.(string)
	}

	if v, ok := d.GetOk("groups"); ok {
		containerSpec.Env = stringListToStringSlice(v.([]interface{}))
	}

	if v, ok := d.GetOk("tty"); ok {
		containerSpec.TTY = v.(bool)
	}

	if v, ok := d.GetOk("open_stdin"); ok {
		containerSpec.OpenStdin = v.(bool)
	}

	//TODO: mounts section

	if v, ok := d.GetOk("stop_grace_period"); ok {
		period := time.Duration(v.(int))
		containerSpec.StopGracePeriod = &period
	}

	//TODO: healthcheck section

	//TODO: hosts section

	//TODO: dnsconfig section

	//TODO: secrets section

	taskTemplate.ContainerSpec = containerSpec

	log.Println("container spec............")

	//resources limits and reservations
	if v, ok := d.GetOk("resources_limits_nano_cpus"); ok {
		limits_cpu := int64(v.(int))
		taskTemplate.Resources.Limits.NanoCPUs = limits_cpu
	}
	log.Println("resources_limits_nano_cpus")

	if v, ok := d.GetOk("resources_limits_memory_bytes"); ok {
		limits_memory := int64(v.(int))
		taskTemplate.Resources.Limits.MemoryBytes = limits_memory
	}
	log.Println("resources_limits_memory_bytes")

	if v, ok := d.GetOk("resources_reservations_nano_cpus"); ok {
		reservations_cpu := int64(v.(int))
		taskTemplate.Resources.Reservations.NanoCPUs = reservations_cpu
	}
	log.Println("resources_reservations_nano_cpus")


	if v, ok := d.GetOk("resources_reservations_memory_bytes"); ok {
		reservations_memory := int64(v.(int))
		taskTemplate.Resources.Reservations.MemoryBytes = reservations_memory
	}
	log.Println("resources_reservations_memory_bytes")


	 //restart policy
	var restartPolicy swarm.RestartPolicy

	if v, ok := d.GetOk("restart_policy_condition"); ok {
		var condition swarm.RestartPolicyCondition
		if v.(string) == "none" {
			condition = swarm.RestartPolicyConditionNone
		}
		if v.(string) == "on-failure" {
			condition = swarm.RestartPolicyConditionOnFailure
		}
		if v.(string) == "any" {
			condition = swarm.RestartPolicyConditionAny
		}
		restartPolicy.Condition = condition
	}

	log.Println("restart_policy_condition")


	if v, ok := d.GetOk("restart_policy_delay"); ok {
		delay := time.Duration(v.(int))
		restartPolicy.Delay = &delay
	}


	log.Println("restart_policy_delay")


	if v, ok := d.GetOk("restart_policy_attempts"); ok {
		attempts := uint64(v.(int))
		restartPolicy.MaxAttempts = &attempts
	}

	log.Println("restart_policy_attempts")


	if v, ok := d.GetOk("restart_policy_window"); ok {
		window := time.Duration(v.(int))
		restartPolicy.Window = &window
	}
	taskTemplate.RestartPolicy = &restartPolicy

	//placement

	if v, ok := d.GetOk("placement"); ok {
		taskTemplate.Placement.Constraints = stringListToStringSlice(v.([]interface{}))

	}

	//TODO: task networks section

	//TODO: logdriver section

	//TODO: forceupdate section

	serviceSpec.TaskTemplate = taskTemplate

	log.Println("Task Template............")

	//serviceSpec Mode

	var serviceSpecMode swarm.ServiceMode

	if v, ok := d.GetOk("is_global_service"); ok {
		if v.(bool) {
			globalService := swarm.GlobalService{}
			serviceSpecMode.Global = &globalService
		}
	}
	log.Println("is_global_service............")

	var serviceSpecModeReplicated swarm.ReplicatedService

	if v, ok := d.GetOk("service_replica_count"); ok {
		replica := uint64(v.(int))
		serviceSpecModeReplicated.Replicas = &replica
	}
	serviceSpecMode.Replicated = &serviceSpecModeReplicated
	serviceSpec.Mode = serviceSpecMode

	log.Println("service_replica_count............")


	//updateconfig

	var serviceSpecUpdateConfig swarm.UpdateConfig

	if v, ok := d.GetOk("update_parallelism_count"); ok {
		parallelism := uint64(v.(int))
		serviceSpecUpdateConfig.Parallelism = parallelism
	}

	if v, ok := d.GetOk("update_delay"); ok {
		delay := time.Duration(v.(int))
		serviceSpecUpdateConfig.Delay = delay
	}

	if v, ok := d.GetOk("update_failure_action"); ok {
		serviceSpecUpdateConfig.FailureAction = v.(string)
	}

	serviceSpec.UpdateConfig = &serviceSpecUpdateConfig
	log.Println("service spec update config............")


	//TODO: updateconfig Monitor and MaxFailureRatio

	//TODO: ServiceSpec Networks is deprecated

	//EndpointSpec

	var endpointSpec swarm.EndpointSpec
	if v, ok := d.GetOk("resolution_mode"); ok {
		var mode swarm.ResolutionMode
		if v.(string) == "vip" {
			mode = swarm.ResolutionModeVIP
		}
		if v.(string) == "dnsrr" {
			mode = swarm.ResolutionModeDNSRR
		}
		endpointSpec.Mode = mode
	}

	portConfigs := []swarm.PortConfig{}

	if v, ok := d.GetOk("ports"); ok {

		portConfigs = portSetToDockerPortConfig(v.(*schema.Set))
	}

	if len(portConfigs) != 0 {
		endpointSpec.Ports = portConfigs
	}

	serviceSpec.EndpointSpec = &endpointSpec

	log.Println("Service spec............")

	createServiceOptions := dc.CreateServiceOptions{
		serviceSpec,
		serviceCtx,
	}

	log.Println("create service options............")

	if createdService, err = client.CreateService(createServiceOptions); err != nil {
		return fmt.Errorf("Unable to create service: %s", err)
	}
	if createdService == nil {
		return fmt.Errorf("Returned service is nil")
	}

	d.SetId(createdService.ID)

	inspectErr := resourceDockerServiceInspect(d, meta)
	if inspectErr != nil {
		return fmt.Errorf("Returned service inpect: %s", inspectErr)
	}



	return nil
}

func resourceDockerServiceInspect(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*dc.Client)

	var service *swarm.Service
	var err error

	if service, err = client.InspectService(d.Id()); err != nil {
		return fmt.Errorf("Error inspecting service: %s", err)
	}

	d.Set("service_version", service.Meta.Version.Index)
	return nil
}

func resourceDockerServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDockerServiceRemove(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*dc.Client)

	removeServiceOptions := dc.RemoveServiceOptions{
		d.Id(),
		serviceCtx,
	}

	if err := client.RemoveService(removeServiceOptions); err != nil {
		return fmt.Errorf("Unable to remove service: %s", err)
	}

	d.SetId("")
	return nil
}

func updateService(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func mapTypeMapValsToString(typeMap map[string]interface{}) map[string]string {
	mapped := make(map[string]string, len(typeMap))
	for k, v := range typeMap {
		mapped[k] = v.(string)
	}
	return mapped
}


func portSetToDockerPortConfig(ports *schema.Set) ([]swarm.PortConfig) {
	retPortConfigs := []swarm.PortConfig{}

	for _, portInt := range ports.List() {
		port := portInt.(map[string]interface{})
		target := uint32(port["target_port"].(int))
		published := uint32(port["published_port"].(int))
		protocol_string := port["protocol"].(string)
		var protocol swarm.PortConfigProtocol
		if protocol_string == "tcp" {
			protocol = swarm.PortConfigProtocolTCP
		}
		if protocol_string == "udp" {
			protocol = swarm.PortConfigProtocolUDP
		}
		publish_mode_string := port["publish_mode"].(string)
		var publish_mode swarm.PortConfigPublishMode
		if publish_mode_string == "ingress" {
			publish_mode = swarm.PortConfigPublishModeIngress
		}
		if publish_mode_string == "host" {
			publish_mode = swarm.PortConfigPublishModeHost
		}

		portConfig := swarm.PortConfig{
			Protocol: protocol,
			TargetPort: target,
			PublishedPort: published,
			PublishMode: publish_mode,

		}
		retPortConfigs = append(retPortConfigs, portConfig)
	}

	return retPortConfigs
}




