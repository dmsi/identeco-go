package storage

// TODO logger

// RefreshToken separate struct
type UserData struct {
	Username string
	Password string
	Hash     string
}

type PrivateKeyData struct {
	Data []byte
}

type JWKSetsData struct {
	Data []byte
}

type UserDataStorage interface {
	ReadUserData(username string) (*UserData, error)
	WriteUserData(user UserData) error
}

type KeyDataStorage interface {
	ReadPrivateKey() (*PrivateKeyData, error)
	WritePrivateKey(key PrivateKeyData) error

	ReadJWKSets() (*JWKSetsData, error)
	WriteJWKSets(jwksSets JWKSetsData) error
}
