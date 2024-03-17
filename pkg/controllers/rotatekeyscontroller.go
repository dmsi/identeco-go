package controllers

import (
	"encoding/json"

	"github.com/dmsi/identeco-go/pkg/services/keys"
	"github.com/dmsi/identeco-go/pkg/storage"
)

type RotateController struct {
	Controller
}

func (c *RotateController) RotateKeys() error {
	keyData := storage.PrivateKeyData{}

	privateKey, err := c.KeyService.GenerateKey()
	if err != nil {
		return wrap("RotateKeys", err)
	}

	keyData.Data, err = c.KeyService.PrivateKeyEncodePEM(privateKey)
	if err != nil {
		return wrap("RotateKeys", err)
	}

	j, err := c.KeyService.PublicKeyToJWK(privateKey.PublicKey)
	if err != nil {
		return wrap("RotateKeys", err)
	}
	new := keys.JWKSets{
		Keys: []keys.JWK{
			*j,
		},
	}

	jwkSetsData, err := c.KeyStorage.ReadJWKSets()
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
	err = c.KeyStorage.WritePrivateKey(keyData)
	if err != nil {
		return wrap("RotateKeys", err)
	}

	data, err := json.Marshal(&new)
	if err != nil {
		return wrap("RotateKeys", err)
	}
	err = c.KeyStorage.WriteJWKSets(storage.JWKSetsData{
		Data: data,
	})
	if err != nil {
		return wrap("RotateKeys", err)
	}

	return nil
}
