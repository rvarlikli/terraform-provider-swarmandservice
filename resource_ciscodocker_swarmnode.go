package main
import (
	"bytes"
	//"fmt"
	//"io"
	//"bufio"
	// "mime/multipart"
	//"os"
	//"io/ioutil"
	"log"
	"net/http"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
)

func resourceCiscoDockerSwarmNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceDockerSwarmNodeCreate,
		Read:   resourceDockerSwarmNodeRead,
		Update: resourceDockerSwarmNodeUpdate,
		Delete: resourceDockerSwarmNodeDelete,
		Schema: map[string]*schema.Schema{
			"api_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"api_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"listen_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"advertise_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"swarm_manager_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"is_swarm_manager": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"swarm_manager_token":&schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"swarm_worker_token":&schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
func resourceDockerSwarmNodeCreate(d *schema.ResourceData, m interface{}) error {
	api_address := d.Get("api_address").(string)
	api_port := d.Get("api_port").(int)
	listen_address := d.Get("listen_address").(string)
	advertise_address := d.Get("advertise_address").(string)
	swarm_manager_address := d.Get("swarm_manager_address").(string)
	is_swarm_manager:= d.Get("is_swarm_manager").(bool)
	swarm_worker_token:="SWMTKN-1-57krntzshq7h1jroa93kz82r74ugsn4eq0jakr876gzxtkhxgh-cd7k2n9npms3nkoif5oqp7gcn"

	if is_swarm_manager {
		targetUrl:= "http://"+api_address+":"+strconv.Itoa(api_port)+"/swarm/join"
		jsonValue:="{\"ListenAddr\":\""+listen_address+"\",\"AdvertiseAddr\":\""+advertise_address+"\",\"RemoteAddrs\":[\""+swarm_manager_address+"\"],\"JoinToken\":\""+swarm_worker_token+"\"}"
		joinDockerSwarm(targetUrl,jsonValue)
	}

	d.SetId(api_address+":"+strconv.Itoa(api_port)+"/swarm/join")
	return nil
}
func resourceDockerSwarmNodeRead(d *schema.ResourceData, m interface{}) error {
	//     := &Params{Count: 5}
	return nil
}
func resourceDockerSwarmNodeUpdate(d *schema.ResourceData, m interface{}) error {
	// Enable partial state mode
	d.Partial(true)
	if d.HasChange("image_name") {
		// Try updating the image_name
		d.SetPartial("image_name")
	}
	if d.HasChange("replica_count") {
		// Try updating the replica_count
		d.SetPartial("replica_count")
	}

	// If we were to return here, before disabling partial mode below,
	// then only the "address" field would be saved.
	// We succeeded, disable partial mode. This causes Terraform to save
	// save all fields again.
	d.Partial(false)
	return nil
}
func resourceDockerSwarmNodeDelete(d *schema.ResourceData, m interface{}) error {
	api_address := d.Get("api_address").(string)
	image_name := d.Get("image_name").(string)
	service_name := d.Get("service_name").(string)
	api_port := d.Get("api_port").(int)
	targetUrl:= "http://"+api_address+":"+strconv.Itoa(api_port)+"/services/"+service_name
	log.Println("Test.....", targetUrl)
	log.Println("Image Name.....", image_name)
	removeDockerService(targetUrl)
	return nil
}

func joinDockerSwarm(targetUrl string, jsonValue string) error {
	bodyBuf := &bytes.Buffer{}
	//bodyWriter := multipart.NewWriter(bodyBuf)
	//bodyWriter:=bufio.NewWriter(bodyBuf)
	bodyBuf.Write([]byte(jsonValue))
	client := &http.Client{}
	/* Authenticate */
	req, err := http.NewRequest("POST", targetUrl,bodyBuf)
	log.Println("Test.....", err)
	log.Println("jsonValue.....", jsonValue)
	req.Header.Set("Content-Type", "Content-Type:text/xml;charset=UTF-8;Accept: application/json")
	res, err := client.Do(req)
	if res.StatusCode == 406 {
		log.Println("Node is already part of a swarm", res)
	}
	if res.StatusCode == 400 {
		log.Println("Bad parameter", res)
	}
	if res.StatusCode != 200 {
		log.Fatal("Unexpected status code...", res.StatusCode)
		log.Fatal("Unexpected error....", res)
	}
	return nil
}

func removeDockerSwarmNode(targetUrl string) error {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", targetUrl,nil)
	//req.SetContentType(contentType)
	res, err := client.Do(req)
	log.Println("Response......", res)
	log.Println("Response Code.....", res.StatusCode)
	if res.StatusCode == 500 {
		log.Println("Server error", res)
	}
	if res.StatusCode == 404 {
		log.Println("No such service", res)
	}
	if res.StatusCode != 200 {
		log.Fatal("Unexpected error", res.StatusCode)
	}
	if err != nil {
		return err
	}
	return nil
}

func updateDockerSwarmNode(targetUrl string, jsonValue string) error {
	bodyBuf := &bytes.Buffer{}
	//bodyWriter := multipart.NewWriter(bodyBuf)
	//bodyWriter:=bufio.NewWriter(bodyBuf)
	bodyBuf.Write([]byte(jsonValue))
	client := &http.Client{}
	/* Authenticate */
	req, err := http.NewRequest("POST", targetUrl,bodyBuf)
	log.Println("Test.....", err)
	req.Header.Set("Content-Type", "Content-Type:text/xml;charset=UTF-8")
	//req.SetBasicAuth(userName,password)
	//req.SetContentType(contentType)
	res, err := client.Do(req)
	if res.StatusCode == 500 {
		log.Println("Server error", res)
	}
	if res.StatusCode == 404 {
		log.Fatal("No such service", res.StatusCode)
	}
	if res.StatusCode != 200 {
		log.Fatal("Unexpected status code", res.StatusCode)
	}
	return nil
}
