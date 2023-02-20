package scaleway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cockpit "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayCockpitToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayCockpitTokenCreate,
		ReadContext:   resourceScalewayCockpitTokenRead,
		DeleteContext: resourceScalewayCockpitTokenDelete,
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(defaultCockpitTimeout),
			Read:    schema.DefaultTimeout(defaultCockpitTimeout),
			Delete:  schema.DefaultTimeout(defaultCockpitTimeout),
			Default: schema.DefaultTimeout(defaultCockpitTimeout),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the token",
			},
			"scopes": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "Endpoints",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query_metrics": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
							Description: "Query metrics",
						},
						"write_metrics": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							ForceNew:    true,
							Description: "Write metrics",
						},
						"setup_metrics_rules": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
							Description: "Setup metrics rules",
						},
						"query_logs": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
							Description: "Query logs",
						},
						"write_logs": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							ForceNew:    true,
							Description: "Write logs",
						},
						"setup_logs_rules": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
							Description: "Setup logs rules",
						},
						"setup_alerts": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
							Description: "Setup alerts",
						},
					},
				},
			},
			"secret_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret key of the token",
				Sensitive:   true,
			},
			"project_id": projectIDSchema(),
		},
	}
}

func resourceScalewayCockpitTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api, err := cockpitAPI(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)
	scopes := expandCockpitTokenScopes(d.Get("scopes"))

	res, err := api.CreateToken(&cockpit.CreateTokenRequest{
		Name:      name,
		Scopes:    scopes,
		ProjectID: projectID,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("secret_key", res.SecretKey)
	d.SetId(res.ID)
	return resourceScalewayCockpitTokenRead(ctx, d, meta)
}

func resourceScalewayCockpitTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api, err := cockpitAPI(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := api.GetToken(&cockpit.GetTokenRequest{
		TokenID: d.Id(),
	}, scw.WithContext(ctx))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("name", res.Name)
	_ = d.Set("scopes", flattenCockpitTokenScopes(res.Scopes))
	_ = d.Set("project_id", res.ProjectID)

	return nil
}

func resourceScalewayCockpitTokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api, err := cockpitAPI(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	err = api.DeleteToken(&cockpit.DeleteTokenRequest{
		TokenID: d.Id(),
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}