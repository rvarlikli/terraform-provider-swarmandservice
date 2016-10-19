package main

import (
	
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
		},
	}
}
func resourceServerCreate(d *schema.ResourceData, m interface{}) error {

 address := d.Get("address").(string)
 user_name := d.Get("user_name").(string)
 password := d.Get("password").(string)
 port := d.Get("port").(string)

 d.SetId(address + "!"+user_name+ "!"+password+"!"+port)
   
	client := &http.Client{}

	/* Authenticate */
	req, err := http.NewRequest("GET", "http://"+address+":"+port+"/simengine/rest/launch", nil)
	req.SetBasicAuth(user_name,password)
	res, err := client.Do(req)
	if err != nil {
	    log.Fatal(err)
	}

	// res, err := http.Get("http://10.204.106.107:19400/rest/projects")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// read body
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatal("Unexpected status code", res.StatusCode)
	}
  
	log.Printf("Body: %s\n", body)
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
