package main
import (
	"bytes"
	"fmt"
	"io"
	"bufio"
	// "mime/multipart"
	"os"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/hashicorp/terraform/helper/schema"
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
			"port": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}


func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	api_address := d.Get("api_address").(string)
	//image_name := d.Get("image_name").(string)
	//service_name := d.Get("service_name").(string)
	//port := d.Get("port").(string)

	targetUrl:= "http://"+api_address+":2375/services/create"
	jsonValue:= "{\"Name\":\"hello\",\"TaskTemplate\":{\"ContainerSpec\":{\"Image\":\"kitematic/hello-world-nginx\"},\"Resources\":{\"Limits\":{},\"Reservations\":{}},\"RestartPolicy\":{\"Condition\":\"any\",\"MaxAttempts\":0},\"Placement\":{}},\"Mode\":{\"Replicated\":{\"Replicas\":1}},\"UpdateConfig\":{\"Parallelism\":1,\"FailureAction\":\"pause\"},\"EndpointSpec\":{\"Mode\":\"vip\"}}"

	startDockerService(targetUrl,jsonValue)

	return nil
}


func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	// params := &Params{Count: 5}
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	// Enable partial state mode
	d.Partial(true)
	if d.HasChange("address") {
		// Try updating the address
		d.SetPartial("address")
	}
	// If we were to return here, before disabling partial mode below,
	// then only the "address" field would be saved.
	// We succeeded, disable partial mode. This causes Terraform to save
	// save all fields again.
	d.Partial(false)
	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Raguuuuu")
	address := d.Get("address").(string)
	user_name := d.Get("user_name").(string)
	password := d.Get("password").(string)
	port := d.Get("port").(string)
	simulation_name := d.Get("simulation_name").(string)

	files, _ := ioutil.ReadDir("./"+simulation_name)

	for _, f := range files {

		targetUrl:= "http://"+address+":"+port+"/simengine/rest/stop/"+simulation_name+"!!"+f.Name()

		stopSimulation(targetUrl,user_name,password)
	}

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
	req.Header.Set("Content-Type", "Content-Type:text/xml;charset=UTF-8")
	//req.SetBasicAuth(userName,password)
	//req.SetContentType(contentType)
	res, err := client.Do(req)
	if res.StatusCode == 400 {
		log.Println("Unexpected status code1", res)
	}
	if res.StatusCode != 201 {
		log.Fatal("Unexpected status code2", res.StatusCode)
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