package controllers

import (
	"encoding/json"

	"github.com/dmsi/identeco-go/pkg/config"
	"github.com/dmsi/identeco-go/pkg/keys"
	"github.com/dmsi/identeco-go/pkg/storage"
)

func (c *Controller) RotateKeys() error {
	keysData := storage.Keys{}

	privateKey, err := keys.GenerateKey(config.Cfg.KeyLength)
	if err != nil {
		return err
	}

	keysData.PrivateKey, err = keys.EncodePEM(privateKey)
	if err != nil {
		return err
	}

	j, err := keys.PublicKeyToJWK(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	new := keys.JWKSets{
		Keys: []keys.JWK{
			*j,
		},
	}

	keysDataStored, err := c.KeyStorage.ReadKeys()
	if err == nil {
		current := keys.JWKSets{}
		err = json.Unmarshal(keysDataStored.JWKS, &current)
		if err != nil {
			return err
		}
		new.Keys = append(new.Keys, current.Keys[0])
	}

	jwksData, err := json.Marshal(new)
	if err != nil {
		return err
	}
	keysData.JWKS = jwksData

	err = c.KeyStorage.WriteKeys(keysData)
	if err != nil {
		return err
	}

	return nil
}
