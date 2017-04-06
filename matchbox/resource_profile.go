package matchbox

import (
	"context"

	matchbox "github.com/coreos/matchbox/matchbox/client"
	"github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfileCreate,
		Read:   resourceProfileRead,
		Delete: resourceProfileDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"kernel": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"initrd": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
			"args": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceProfileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*matchbox.Client)
	ctx := context.TODO()

	name := d.Get("name").(string)

	var initrds []string
	for _, initrd := range d.Get("initrd").([]interface{}) {
		initrds = append(initrds, initrd.(string))
	}

	var args []string
	for _, arg := range d.Get("args").([]interface{}) {
		args = append(args, arg.(string))
	}

	profile := &storagepb.Profile{
		Id:         name,
		IgnitionId: d.Get("config").(string),
		Boot: &storagepb.NetBoot{
			Kernel: d.Get("kernel").(string),
			Initrd: initrds,
			Args:   args,
		},
	}

	_, err := client.Profiles.ProfilePut(ctx, &serverpb.ProfilePutRequest{
		Profile: profile,
	})
	if err != nil {
		return err
	}

	d.SetId(profile.GetId())
	return err
}

func resourceProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*matchbox.Client)
	ctx := context.TODO()

	name := d.Get("name").(string)
	_, err := client.Profiles.ProfileGet(ctx, &serverpb.ProfileGetRequest{
		Id: name,
	})
	if err != nil {
		// resource doesn't exist anymore
		d.SetId("")
		return nil
	}
	return err
}

func resourceProfileDelete(d *schema.ResourceData, meta interface{}) error {
	// TODO: Delete API is not yet implemented
	return nil
}
