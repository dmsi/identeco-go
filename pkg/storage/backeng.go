package storage

// TODO PROTOTYPE

// RefreshToken separate struct
type UserData struct {
	Username     string
	Password     string
	Hash         string
	RefreshToken string
}

type PrivateKeyData struct {
	Data []byte
}

type JWKSetsData struct {
	Data []byte
}

type UserDataStorage interface {
	ReadUserData(username string) (*UserData, error)
	WriteUserData(username string, user UserData) error
}

type KeyDataStorage interface {
	// TODO - keey the name in args?
	// Do I need to pass name?
	// It suggests that I can save two or more private keys
	// When in fact I can have only one at a time
	// But this interface allows to make a backup!
	ReadPrivateKey() (*PrivateKeyData, error)
	WritePrivateKey() error

	ReadJWKSets() (*JWKSetsData, error)
	WriteJWKSets() error
}

// IDO_DDB_TABLE_NAME
// IDO_S3_BUCKET_NAME
// IDO_PRIVATE_KEY_NAME
// IDO_JWKS_JSON_NAME
// IDO_PRIVATE_KEY_BITS
// IDO_PRIVATE_KEY_LIFETIME
// IDO_ACCESS_TOKEN_LIFETIME
// IDO_REFRESH_TOKEN_LIFETIME
// IDO_CLAIM_ISS
// func GetUserData(username string) (interface{}, error) {
// 	return nil, nil
// }

// func SetUserData(userdata interface{}) error {
// 	return nil
// }

// func GetPrivateKey() (interface{}, error) {
// 	return nil, nil
// }

// func SetPrivateKey(key interface{}) error {
// 	return nil
// }

// func GetJWKS() (interface{}, error) {
// 	return nil, nil
// }

// func SetJWKS(j interface{}) error {
// 	return nil
// }
