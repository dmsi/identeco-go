package jwksets

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/dmsi/identeco-go/pkg/controllers"
)

type JWKSetsController struct {
	controllers.Controller
}

func (j *JWKSetsController) GetJWKSets() (*string, error) {
	jwkSetsData, err := j.KeyStorage.ReadJWKSets()
	if err != nil {
		return nil, err
	}

	return aws.String(string(jwkSetsData.Data)), nil
}
