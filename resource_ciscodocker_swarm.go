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
	"encoding/json"
	"github.com/docker/docker/api/types/swarm"
)

func resourceCiscoDockerSwarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceDockerSwarmCreate,
		Read:   resourceDockerSwarmRead,
		Update: resourceDockerSwarmUpdate,
		Delete: resourceDockerSwarmDelete,
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
			"force_new": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"task_history_retention_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},




		},
	}
}
func resourceDockerSwarmCreate(d *schema.ResourceData, m interface{}) error {
	api_address := d.Get("api_address").(string)
	api_port := d.Get("api_port").(int)
	listen_address := d.Get("listen_address").(string)
	advertise_address := d.Get("advertise_address").(string)
	force_new := d.Get("force_new").(bool)

	targetUrl:= "http://"+api_address+":"+strconv.Itoa(api_port)+"/swarm/init"
	jsonValue:="{\"ListenAddr\":\""+listen_address+"\",\"AdvertiseAddr\":\""+advertise_address+"\",\"ForceNewCluster\":"+strconv.FormatBool(force_new)+",\"Spec\":{\"Orchestration\":{},\"Raft\":{},\"Dispatcher\":{},\"CAConfig\":{}}}"
	initDockerSwarm(targetUrl,jsonValue)
	log.Println("Swarm created............")

	getDockerSwarm("http://"+api_address+":"+strconv.Itoa(api_port)+"/swarm")


	d.SetId(api_address+":"+strconv.Itoa(api_port)+"/swarm/init")
	return nil
}
func resourceDockerSwarmRead(d *schema.ResourceData, m interface{}) error {
	//     := &Params{Count: 5}
	return nil
}
func resourceDockerSwarmUpdate(d *schema.ResourceData, m interface{}) error {
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
func resourceDockerSwarmDelete(d *schema.ResourceData, m interface{}) error {
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

func initDockerSwarm(targetUrl string, jsonValue string) error {
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

func getDockerSwarm(targetUrl string) error {

	client := &http.Client{}
	/* Authenticate */
	req, err := http.NewRequest("GET", targetUrl,nil)
	log.Println("Request error.....", err)
	req.Header.Set("Content-Type", "Content-Type:text/xml;charset=UTF-8")

	res, err := client.Do(req)

	if res.StatusCode != 200 {
		log.Fatal("Unexpected status code...", res.StatusCode)
		log.Fatal("Unexpected error....", res)
	}
	if res.StatusCode == 200 {

		// TODO: get version number and service id from response


	}

	swarmInfo:= swarm.Swarm{}
	defer res.Body.Close()
	log.Println("KADİRRRRRRR......", json.NewDecoder(res.Body).Decode(&swarmInfo))

	if err != nil {
		return  err
	}
	return nil
}

func removeDockerSwarm(targetUrl string) error {
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

func updateDockerSwarm(targetUrl string, jsonValue string) error {
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

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	log.Println("Get json body", r.Body)
	log.Println("Get json error", err)

	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
