package newrelic

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	"github.com/newrelic/newrelic-client-go/pkg/obfuscation"
	"log"
)

func resourceNewRelicObfuscationExpresion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicObfuscationExpresionCreate,
		ReadContext:   resourceNewRelicObfuscationExpresionRead,
		UpdateContext: resourceNewRelicObfuscationExpresionUpdate,
		DeleteContext: resourceNewRelicObfuscationExpresionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Account with the NRQL drop rule will be put.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of expression.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of expression.",
			},
			"regex": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description of expression.",
			},
		},
	}
}

func resourceNewRelicObfuscationExpresionCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, data)
	createInp := obfuscation.LogConfigurationsCreateObfuscationExpressionInput{
		Name:  data.Get("name").(string),
		Regex: data.Get("regex").(string),
	}
	if v, ok := data.GetOk("description"); ok {
		createInp.Description = v.(string)
	}

	log.Printf("[INFO] Creating New Relic One obfuscation expression %s", createInp.Name)

	created, err := client.Obfuscation.LogConfigurationsCreateObfuscationExpressionWithContext(ctx, accountID, createInp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := created.id

	log.Printf("[INFO] New Obfuscation Expression: %s", id)
	data.SetId(id)

	return resourceNewRelicObfuscationExpresionRead(ctx, data, meta)
}

func resourceNewRelicObfuscationExpresionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic obfuscation expression %s", d.Id())

	expression, err := client.Obfuscation.GetObfuscationExpressionsWithContext(ctx, accountID)
	if err != nil && expression == nil {
		d.SetId("")
		return nil
	}

	return diag.FromErr(flattenObfuscationExpression(expression, d))
}

func resourceNewRelicObfuscationExpresionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	updateInp := obfuscation.LogConfigurationsUpdateObfuscationExpressionInput{
		Name:  d.Get("name").(string),
		Regex: d.Get("regex").(string),
	}

	if v, ok := d.GetOk("description"); ok {
		updateInp.Description = v.(string)
	}

	log.Printf("[INFO] Updating New Relic One obfuscation expression %s", d.Id())
	updated, err := client.Obfuscation.LogConfigurationsUpdateObfuscationExpressionWithContext(ctx, accountID, updateInp)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicObfuscationExpresionRead(ctx, d, meta)

}

func resourceNewRelicObfuscationExpresionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Deleting New Relic obfuscation expression %s", d.Id())

	if _, err := client.Obfuscation.LogConfigurationsDeleteObfuscationExpressionWithContext(ctx, accountID, common.EntityGUID(d.Id())); err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			return nil
		}
		return diag.FromErr(err)
	}
	return nil
}

func flattenObfuscationExpression(e *obfuscation.LogConfigurationsObfuscationExpression, d *schema.ResourceData) error {
	_ = d.Set("name", e.Name)
	_ = d.Set("regex", e.Regex)
	if e.Description != "" {
		_ = d.Set("description", e.Description)
	}
	return nil
}
