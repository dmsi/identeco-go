package rotatekeys

import (
	"encoding/json"

	"github.com/dmsi/identeco-go/pkg/controllers"
	e "github.com/dmsi/identeco-go/pkg/lib/err"
	"github.com/dmsi/identeco-go/pkg/services/keys"
	"github.com/dmsi/identeco-go/pkg/storage"
)

// TODO rename to rotatekeys

func wrap(name string, err error) error {
	return e.Wrap("controllers.rotatekeys."+name, err)
}

type RotateController struct {
	controllers.Controller
}

func (r *RotateController) RotateKeys() error {
	keyData := storage.PrivateKeyData{}
	// jwkSetsData := &storage.JWKSetsData{}

	privateKey, err := r.KeyService.GenerateKey()
	if err != nil {
		return wrap("RotateKeys", err)
	}

	keyData.Data, err = r.KeyService.PrivateKeyEncodePEM(privateKey)
	if err != nil {
		return wrap("RotateKeys", err)
	}

	j, err := r.KeyService.PublicKeyToJWK(privateKey.PublicKey)
	if err != nil {
		return wrap("RotateKeys", err)
	}
	new := keys.JWKSets{
		Keys: []keys.JWK{
			*j,
		},
	}

	jwkSetsData, err := r.KeyStorage.ReadJWKSets()
	if err == nil {
		current := keys.JWKSets{}
		err = json.Unmarshal(jwkSetsData.Data, &current)
		if err != nil {
			return wrap("RotateKeys", err)
		}
		new.Keys = append(new.Keys, current.Keys[0])
	}

	// Write data
	// TODO atomic
	err = r.KeyStorage.WritePrivateKey(keyData)
	if err != nil {
		return wrap("RotateKeys", err)
	}

	data, err := json.Marshal(&new)
	if err != nil {
		return wrap("RotateKeys", err)
	}
	err = r.KeyStorage.WriteJWKSets(storage.JWKSetsData{
		Data: data,
	})
	if err != nil {
		return wrap("RotateKeys", err)
	}

	return nil
}
