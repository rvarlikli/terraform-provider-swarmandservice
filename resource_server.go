package main
import (
	"bytes"
	"fmt"
	"io"
	"bufio"
	// "mime/multipart"
	"os"
	//"io/ioutil"
	"log"
	"net/http"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
)
type Params struct {
	Count int `url:"count,omitempty"`
}
func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,
		Schema: map[string]*schema.Schema{
			"api_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"image_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"api_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"replica_count": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"start_command": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_command_args": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"published_port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"target_port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}
func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	api_address := d.Get("api_address").(string)
	image_name := d.Get("image_name").(string)
	service_name := d.Get("service_name").(string)
	replica_count := d.Get("replica_count").(int)
	api_port := d.Get("api_port").(int)
	start_command := d.Get("start_command").(string)
	start_command_args := d.Get("start_command_args").(string)
	published_port := d.Get("published_port").(int)
	target_port := d.Get("target_port").(int)
	d.SetId(api_address+":"+strconv.Itoa(api_port)+"/service/"+service_name)

	if start_command !="" {
		log.Println("Start command", start_command)
		start_command="\""+start_command+"\""
	}
	if start_command_args != "" {
		log.Println("Start command args", start_command_args)
		start_command_args="\""+start_command_args+"\""
	}


	targetUrl:= "http://"+api_address+":"+strconv.Itoa(api_port)+"/services/create"
	jsonValue:="{\"Name\":\""+service_name+"\",\"TaskTemplate\":{\"ContainerSpec\":{\"Image\":\""+image_name+"\",\"Command\":["+start_command+"],\"Args\":["+start_command_args+"]},\"Resources\":{\"Limits\":{},\"Reservations\":{}},\"RestartPolicy\":{\"Condition\":\"any\",\"MaxAttempts\":0},\"Placement\":{}},\"Mode\":{\"Replicated\":{\"Replicas\":"+strconv.Itoa(replica_count)+"}},\"UpdateConfig\":{\"Parallelism\":1,\"FailureAction\":\"pause\"},\"EndpointSpec\":{\"Mode\":\"vip\",\"Ports\":[{\"Protocol\":\"tcp\",\"PublishedPort\":"+strconv.Itoa(published_port)+",\"TargetPort\":"+strconv.Itoa(target_port)+"}]}}"
	startDockerService(targetUrl,jsonValue)
	return nil
}
func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	//     := &Params{Count: 5}
	return nil
}
func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
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
	api_address := d.Get("api_address").(string)
	image_name := d.Get("image_name").(string)
	service_name := d.Get("service_name").(string)
	replica_count := d.Get("replica_count").(int)
	api_port := d.Get("api_port").(int)
	start_command := d.Get("start_command").(string)
	start_command_args := d.Get("start_command_args").(string)
	published_port := d.Get("published_port").(int)
	target_port := d.Get("target_port").(int)
	d.SetId(api_address+":"+strconv.Itoa(api_port)+"/service/"+service_name)

	if start_command !="" {
		log.Println("Start command", start_command)
		start_command="\""+start_command+"\""
	}
	if start_command_args != "" {
		log.Println("Start command args", start_command_args)
		start_command_args="\""+start_command_args+"\""
	}


	targetUrl:= "http://"+api_address+":"+strconv.Itoa(api_port)+"/services/"+service_name+"/update?version=1278"
	jsonValue:="{\"Name\":\""+service_name+"\",\"TaskTemplate\":{\"ContainerSpec\":{\"Image\":\""+image_name+"\",\"Command\":["+start_command+"],\"Args\":["+start_command_args+"]},\"Resources\":{\"Limits\":{},\"Reservations\":{}},\"RestartPolicy\":{\"Condition\":\"any\",\"MaxAttempts\":0},\"Placement\":{}},\"Mode\":{\"Replicated\":{\"Replicas\":"+strconv.Itoa(replica_count)+"}},\"UpdateConfig\":{\"Parallelism\":1,\"FailureAction\":\"pause\"},\"EndpointSpec\":{\"Mode\":\"vip\",\"Ports\":[{\"Protocol\":\"tcp\",\"PublishedPort\":"+strconv.Itoa(published_port)+",\"TargetPort\":"+strconv.Itoa(target_port)+"}]}}"
	updateDockerService(targetUrl,jsonValue)
	// If we were to return here, before disabling partial mode below,
	// then only the "address" field would be saved.
	// We succeeded, disable partial mode. This causes Terraform to save
	// save all fields again.
	d.Partial(false)
	return nil
}
func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
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
func startDockerService(targetUrl string, jsonValue string) error {
	bodyBuf := &bytes.Buffer{}
	//bodyWriter := multipart.NewWriter(bodyBuf)
	//bodyWriter:=bufio.NewWriter(bodyBuf)
	bodyBuf.Write([]byte(jsonValue))
	client := &http.Client{}
	/* Authenticate */
	req, err := http.NewRequest("POST", targetUrl,bodyBuf)
	log.Println("Test.....", err)
	log.Println("jsonValue.....", jsonValue)
	req.Header.Set("Content-Type", "Content-Type:text/xml;charset=UTF-8")
	//req.SetBasicAuth(userName,password)
	//req.SetContentType(contentType)
	res, err := client.Do(req)
	if res.StatusCode == 409 {
		log.Println("Existing service", res)
	}
	if res.StatusCode == 406 {
		log.Println("Server error or node is not part of swarm", res)
	}
	if res.StatusCode != 201 {
		log.Fatal("Unexpected status code...", res.StatusCode)
		log.Fatal("Unexpected error....", res)
	}
	return nil
}
func removeDockerService(targetUrl string) error {
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
func updateDockerService(targetUrl string, jsonValue string) error {
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
func startSimulation(localFileName string, targetUrl string, userName string, password string) error{
	bodyBuf := &bytes.Buffer{}
	//bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter:=bufio.NewWriter(bodyBuf)
	fh, err := os.Open(localFileName)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	//iocopy
	_, err = io.Copy(bodyWriter, fh)
	if err != nil {
		return err
	}
	log.Println(bodyBuf)
	client := &http.Client{}
	/* Authenticate */
	req, err := http.NewRequest("POST", targetUrl,bodyBuf)
	req.Header.Set("Content-Type", "Content-Type:text/xml;charset=UTF-8")
	req.SetBasicAuth(userName,password)
	//req.SetContentType(contentType)
	res, err := client.Do(req)
	if res.StatusCode == 400 {
		log.Println("Unexpected status code1", res)
	}
	if res.StatusCode != 200 {
		log.Fatal("Unexpected status code2", res.StatusCode)
	}
	return nil
}
func stopSimulation(targetUrl string, userName string , password string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", targetUrl,nil)
	req.SetBasicAuth(userName,password)
	//req.SetContentType(contentType)
	res, err := client.Do(req)
	if res.StatusCode == 400 {
		log.Println("Unexpected status code1", res)
	}
	if res.StatusCode != 200 {
		log.Fatal("Unexpected status code2", res.StatusCode)
	}
	if err != nil {
		return err
	}
	return nil
}