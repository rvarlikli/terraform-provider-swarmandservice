package main

import (
	"bytes"
	"fmt"
	"regexp"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCiscoDockerService() *schema.Resource {
	return &schema.Resource{
		Create: resourceDockerServiceCreate,
		Read:   resourceDockerServiceInspect,
		Update: resourceDockerServiceUpdate,
		Delete: resourceDockerServiceRemove,

		Schema: map[string]*schema.Schema{
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"image_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"container_labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"command": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"args": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"env": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dir": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"groups": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tty": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: false,
			},
			"open_stdin": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: false,
			},
			"stop_grace_period": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default: 5000000000,
			},
			"resources_limits_nano_cpus": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"resources_limits_memory_bytes": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"resources_reservations_nano_cpus": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"resources_reservations_memory_bytes": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"restart_policy_condition": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "any",
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					if !regexp.MustCompile(`^(none|on-failure|any)$`).MatchString(value) {
						es = append(es, fmt.Errorf(
							"%q must be one of \"none\", \"on-failure\" or \"any\"", k))
					}
					return
				},
			},

			"restart_policy_delay": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default: 10000000000.0,
			},
			"restart_policy_attempts": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default: 0,
			},
			"restart_policy_window": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default: 0,
			},
			"placement": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_global_service": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ConflictsWith: []string{"service_replica_count"},

			},
			"service_replica_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:1,
			},
			"update_parallelism_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default: 1,
			},
			"update_delay": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"update_failure_action": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "continue",
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					if !regexp.MustCompile(`^(continue|pause)$`).MatchString(value) {
						es = append(es, fmt.Errorf(
							"%q must be one of \"continue\" or \"pause\"", k))
					}
					return
				},
			},
			"resolution_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "vip",
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					if !regexp.MustCompile(`^(vip|dnsrr)$`).MatchString(value) {
						es = append(es, fmt.Errorf(
							"%q must be one of \"vip\" or \"dnsrr\"", k))
					}
					return
				},
			},
			"ports": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"published_port": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"protocol": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "tcp",
							ForceNew: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
								value := v.(string)
								if !regexp.MustCompile(`^(tcp|udp)$`).MatchString(value) {
									es = append(es, fmt.Errorf(
										"%q must be one of \"tcp\" or \"udp\"", k))
								}
								return
							},
						},
						"publish_mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ingress",
							ForceNew: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
								value := v.(string)
								if !regexp.MustCompile(`^(ingress|host)$`).MatchString(value) {
									es = append(es, fmt.Errorf(
										"%q must be one of \"ingress\" or \"host\"", k))
								}
								return
							},
						},
					},
				},
				Set: resourceDockerPortsHash,
			},

		},
	}
}

func resourceDockerPortsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%v-", m["target_port"].(int)))

	if v, ok := m["published_port"]; ok {
		buf.WriteString(fmt.Sprintf("%v-", v.(int)))
	}

	if v, ok := m["publish_mode"]; ok {
		buf.WriteString(fmt.Sprintf("%v-", v.(string)))
	}

	if v, ok := m["protocol"]; ok {
		buf.WriteString(fmt.Sprintf("%v-", v.(string)))
	}

	return hashcode.String(buf.String())
}

