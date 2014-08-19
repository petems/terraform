package heroku

import (
	"fmt"
	"log"

	"github.com/bgentry/heroku-go"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func resourceHerokuDrain() *schema.Resource {
	return &schema.Resource{
		Create: resourceHerokuDrainCreate,
		Read:   resourceHerokuDrainRead,
		Delete: resourceHerokuDrainDelete,

		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"app": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"token": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceHerokuDrainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*heroku.Client)

	app := d.Get("app").(string)
	url := d.Get("url").(string)

	log.Printf("[DEBUG] Drain create configuration: %#v, %#v", app, url)

	dr, err := client.LogDrainCreate(app, url)
	if err != nil {
		return err
	}

	d.SetId(dr.Id)
	d.Set("url", dr.URL)
	d.Set("token", dr.Token)
	d.SetDependencies([]terraform.ResourceDependency{
		terraform.ResourceDependency{ID: app},
	})

	log.Printf("[INFO] Drain ID: %s", d.Id())
	return nil
}

func resourceHerokuDrainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*heroku.Client)

	log.Printf("[INFO] Deleting drain: %s", d.Id())

	// Destroy the drain
	err := client.LogDrainDelete(d.Get("app").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting drain: %s", err)
	}

	return nil
}

func resourceHerokuDrainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*heroku.Client)

	dr, err := client.LogDrainInfo(d.Get("app").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving drain: %s", err)
	}

	d.Set("url", dr.URL)
	d.Set("token", dr.Token)

	return nil
}
