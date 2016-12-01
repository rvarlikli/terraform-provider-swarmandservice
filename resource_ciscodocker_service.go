package main
import (
	"bytes"
	"fmt"
	//"io"
	//"bufio"
	// "mime/multipart"
	//"os"
	//"io/ioutil"
	"log"
	"net/http"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/hashcode"
	"strconv"
	"strings"
)
//type Params struct {
//	Count int `url:"count,omitempty"`
//}
func resourceCiscoDockerService() *schema.Resource {
	return &schema.Resource{
		Create: resourceDockerServiceCreate,
		Read:   resourceDockerServiceRead,
		Update: resourceDockerServiceUpdate,
		Delete: resourceDockerServiceDelete,
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
			"ports": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"published": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},

						"target": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},

						"protocol": &schema.Schema{
							Type:     schema.TypeString,
							Default:  "tcp",
							Optional: true,
							ForceNew: true,
						},
					},
				},
				Set: resourceServicePortsHash,

			},
			"env": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"service_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_number": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}
func resourceDockerServiceCreate(d *schema.ResourceData, m interface{}) error {
	api_address := d.Get("api_address").(string)
	image_name := d.Get("image_name").(string)
	service_name := d.Get("service_name").(string)
	replica_count := d.Get("replica_count").(int)
	api_port := d.Get("api_port").(int)
	start_command := d.Get("start_command").(string)
	start_command_args := d.Get("start_command_args").(string)
	portsString := ""
	envString := ""
	d.SetId(api_address+":"+strconv.Itoa(api_port)+"/service/"+service_name)

	if v, ok := d.GetOk("ports"); ok {
		//log.Println("Ports.....", v)
		portsString = portsConvertToString(v.(*schema.Set))
		//log.Println("Ports.....", portsString)
	}
	if v, ok := d.GetOk("env"); ok {
		envString = stringSetToStringSlice(v.(*schema.Set))
	}

	if start_command !="" {
		log.Println("Start command", start_command)
		start_command="\""+start_command+"\""
	}
	if start_command_args != "" {
		log.Println("Start command args", start_command_args)
		start_command_args="\""+start_command_args+"\""
	}


	targetUrl:= "http://"+api_address+":"+strconv.Itoa(api_port)+"/services/create"
	jsonValue:="{\"Name\":\""+service_name+"\",\"TaskTemplate\":{\"ContainerSpec\":{\"Image\":\""+image_name+"\",\"Command\":["+start_command+"],\"Args\":["+start_command_args+"],\"Env\":["+envString+"]},\"Resources\":{\"Limits\":{},\"Reservations\":{}},\"RestartPolicy\":{\"Condition\":\"any\",\"MaxAttempts\":0},\"Placement\":{}},\"Mode\":{\"Replicated\":{\"Replicas\":"+strconv.Itoa(replica_count)+"}},\"UpdateConfig\":{\"Parallelism\":1,\"FailureAction\":\"pause\"},\"EndpointSpec\":{\"Mode\":\"vip\",\"Ports\":["+portsString+"]}}"
	startDockerService(targetUrl,jsonValue)
	return nil
}
func resourceDockerServiceRead(d *schema.ResourceData, m interface{}) error {
	api_address := d.Get("api_address").(string)
	api_port := d.Get("api_port").(int)
	service_name := d.Get("service_name").(string)
	targetUrl:= "http://"+api_address+":"+strconv.Itoa(api_port)+"/services/"+service_name
	getDockerService(targetUrl)

	d.Set("service_id", "test")
	d.Set("version_number", 538)


	return nil
}
func resourceDockerServiceUpdate(d *schema.ResourceData, m interface{}) error {
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
func resourceDockerServiceDelete(d *schema.ResourceData, m interface{}) error {
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
func getDockerService(targetUrl string) error {
	client := &http.Client{}
	/* Authenticate */
	req, err := http.NewRequest("GET", targetUrl,nil)
	log.Println("Request error.....", err)
	req.Header.Set("Content-Type", "Content-Type:text/xml;charset=UTF-8")

	res, err := client.Do(req)
	if res.StatusCode == 409 {
		log.Println("Existing service", res)
	}
	if res.StatusCode == 406 {
		log.Println("Server error or node is not part of swarm", res)
	}
	if res.StatusCode != 200 {
		log.Fatal("Unexpected status code...", res.StatusCode)
		log.Fatal("Unexpected error....", res)
	}
	if res.StatusCode == 200 {
		log.Println("Response......", res)
		// TODO: get version number and service id from response

	}
	if err != nil {
		return  err
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

func resourceServicePortsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%v-", m["published"].(int)))

	if v, ok := m["target"]; ok {
		buf.WriteString(fmt.Sprintf("%v-", v.(int)))
	}

	if v, ok := m["protocol"]; ok {
		buf.WriteString(fmt.Sprintf("%v-", v.(string)))
	}
	//log.Println("Start command", hashcode.String(buf.String()))
	return hashcode.String(buf.String())
}


func portsConvertToString(ports *schema.Set) (string) {
	portString :=""
	forCount := 0
	for _, portInt := range ports.List() {
		forCount++
		if forCount > 1{
			portString=portString+","
		}
		port := portInt.(map[string]interface{})
		published := port["published"].(int)
		target := port["target"].(int)
		protocol := port["protocol"].(string)

		portString=portString+"{"+
			"\"Protocol\": \""+protocol+"\","+
			"\"TargetPort\":"+strconv.Itoa(target)+","+
			"\"PublishedPort\":"+strconv.Itoa(published)+
			"}"
	}

	return portString
}

func stringSetToStringSlice(stringSet *schema.Set) string {
	envString := ""
	envCount := 0
	if stringSet == nil {
		return envString
	}
	for _, envVal := range stringSet.List() {
		envCount++
		if envCount > 1 {
			envString+=","
		}
		tmpString:= strings.Replace(envVal.(string),"\"","\\\"",-1)
		envString += "\""+tmpString+"\""
		log.Println("Env....", envString)
	}
	return envString
}

