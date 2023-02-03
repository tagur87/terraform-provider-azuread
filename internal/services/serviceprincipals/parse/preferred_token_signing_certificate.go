package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
)

type PreferredTokenSigningCertificateId struct {
	ServicePrincipalId string
	Thumbprint         string
}

func NewPreferredTokenSigningCertificateID(servicePrincipalId, thumbprint string) PreferredTokenSigningCertificateId {
	return PreferredTokenSigningCertificateId{
		ServicePrincipalId: servicePrincipalId,
		Thumbprint:         thumbprint,
	}
}

func (id PreferredTokenSigningCertificateId) String() string {
	return id.ServicePrincipalId + "/tokenSigningCertificate/" + id.Thumbprint
}

func PreferredTokenSigningCertificateID(idString string) (*PreferredTokenSigningCertificateId, error) {
	parts := strings.Split(idString, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Object Resource ID should be in the format {servicePrincipalId}/{type}/{thumbprint} - but got %q", idString)
	}

	id := PreferredTokenSigningCertificateId{
		ServicePrincipalId: parts[0],
		Thumbprint:         parts[2],
	}

	if _, err := uuid.ParseUUID(id.ServicePrincipalId); err != nil {
		return nil, fmt.Errorf("ServicePrincipalId isn't a valid UUID (%q): %+v", id.ServicePrincipalId, err)
	}

	if parts[1] == "" {
		return nil, fmt.Errorf("Type in {servicePrincipalId}/{type}/{subID} should not be empty")
	}

	return &id, nil
}
