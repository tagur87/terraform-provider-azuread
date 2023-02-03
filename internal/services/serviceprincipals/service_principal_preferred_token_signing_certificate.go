package serviceprincipals

import (
	"context"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/services/serviceprincipals/parse"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func servicePrincipalPreferredTokenSigningCertificateResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: servicePrincipalPreferredTokenSigningCertificateResourceCreate,
		ReadContext:   servicePrincipalPreferredTokenSigningCertificateResourceRead,
		DeleteContext: servicePrincipalPreferredTokenSigningCertificateResourceDelete,

		Importer: tf.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ObjectSubResourceID(id, "claimsMappingPolicy")
			return err
		}),

		Schema: map[string]*schema.Schema{
			"preferred_thumbprint": {
				Description: "The thumbprint of the preferred token signing certificate",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},

			"service_principal_id": {
				Description: "Object ID of the service principal for which to configure the preferred certificate",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
		},
	}
}

func servicePrincipalPreferredTokenSigningCertificateResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).ServicePrincipals.ServicePrincipalsClient

	thumbprint := d.Get("preferred_thumbprint").(string)
	servicePrincipal := d.Get("service_principal_id").(string)

	_, err := client.Update(ctx, msgraph.ServicePrincipal{
		DirectoryObject: msgraph.DirectoryObject{
			Id: utils.String(servicePrincipal),
		},
		PreferredTokenSigningKeyThumbprint: utils.NullableString(thumbprint),
	})

	if err != nil {
		return tf.ErrorDiagF(
			err,
			"Could not set PreferredTokenSigningCertificate, service_principal_id: %q, preferred_thumbprint: %q",
			servicePrincipal,
			thumbprint,
		)
	}

	id := parse.NewPreferredTokenSigningCertificateID(
		servicePrincipal,
		thumbprint,
	)

	d.SetId(id.String())

	return servicePrincipalPreferredTokenSigningCertificateResourceRead(ctx, d, meta)
}

func servicePrincipalPreferredTokenSigningCertificateResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).ServicePrincipals.ServicePrincipalsClient

	id, err := parse.PreferredTokenSigningCertificateID(d.Id())
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Parsing Preferred Token Signing Certificate ID %q", d.Id())
	}

	spID := id.ServicePrincipalId

	servicePrincipal, status, err := client.Get(ctx, spID, odata.Query{})
	if err != nil {
		if status == http.StatusNotFound {
			log.Printf("[DEBUG] Service Principal with Object ID %q was not found - removing preferred token signing certificate from state!", spID)
			d.SetId("")
			return nil
		}
		return tf.ErrorDiagF(err, "retrieving service principal with object ID: %q", spID)
	}

	tf.Set(d, "service_principal_id", spID)
	tf.Set(d, "preferred_thumbprint", servicePrincipal.PreferredTokenSigningKeyThumbprint)

	return nil
}

func servicePrincipalPreferredTokenSigningCertificateResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).ServicePrincipals.ServicePrincipalsClient

	id, err := parse.PreferredTokenSigningCertificateID(d.Id())
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Parsing Preferred Token Signing Certificate ID %q", d.Id())
	}

	_, err = client.Update(ctx, msgraph.ServicePrincipal{
		DirectoryObject: msgraph.DirectoryObject{
			Id: utils.String(id.ServicePrincipalId),
		},
		PreferredTokenSigningKeyThumbprint: utils.NullableString(""),
	})

	if err != nil {
		return tf.ErrorDiagF(
			err,
			"Could not remove PreferredTokenSigningCertificate, service_principal_id: %q",
			id.ServicePrincipalId,
		)
	}

	return servicePrincipalPreferredTokenSigningCertificateResourceRead(ctx, d, meta)
}
