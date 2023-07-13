package jwksets

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/dmsi/identeco-go/pkg/controllers"
	e "github.com/dmsi/identeco-go/pkg/lib/err"
)

func wrap(name string, err error) error {
	return e.Wrap("controllers.jwksets."+name, err)
}

type JWKSetsController struct {
	controllers.Controller
}

func (j *JWKSetsController) GetJWKSets() (*string, error) {
	jwkSetsData, err := j.KeyStorage.ReadJWKSets()
	if err != nil {
		return nil, wrap("GetJWKSets", err)
	}

	return aws.String(string(jwkSetsData.Data)), nil
}
