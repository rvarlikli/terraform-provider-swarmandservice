package main

import (
	"bytes"
    "fmt"
    "io"
   "bufio"
    // "mime/multipart"
    "os"

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
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"virl_file": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"simulation_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

 


func resourceServerCreate(d *schema.ResourceData, m interface{}) error {

 address := d.Get("address").(string)
 user_name := d.Get("user_name").(string)
 password := d.Get("password").(string)
 port := d.Get("port").(string)
 virl_file := d.Get("virl_file").(string)
 simulation_name := d.Get("simulation_name").(string)
  
 d.SetId(address + "!"+user_name+ "!"+password+"!"+port)

bodyBuf := &bytes.Buffer{}
//bodyWriter := multipart.NewWriter(bodyBuf)
bodyWriter:=bufio.NewWriter(bodyBuf)

filename:= "./"+virl_file
targetUrl:= "http://"+address+":"+port+"/simengine/rest/launch?session="+simulation_name

    // this step is very important
//fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
 
   
    // open file handle
    fh, err := os.Open(filename)
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
  //bodyWriter.Close()

 
   
	client := &http.Client{}

	//virl_template, err := ioutil.ReadFile("templ-lndc-scan-ff.virl")
	//mody := &virl_template.Buffer{}

	/* Authenticate */
	req, err := http.NewRequest("POST", targetUrl,bodyBuf)
	req.Header.Set("Content-Type", "Content-Type:text/xml;charset=UTF-8")
	req.SetBasicAuth(user_name,password)
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
	return nil
}
