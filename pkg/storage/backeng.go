package storage

// TODO PROTOTYPE

// RefreshToken separate struct
type UserData struct {
	Username     string `dynamodbav:"pk"`
	Password     string `dynamodbav:"password"`
	Hash         string `dynamodbav:"hash"`
	RefreshToken string `dynamodbav:"sk"`
}

// PK - token -> get user by token
type TokenData struct {
	Token    string
	Username string
}

type PrivateKeyData struct {
	PK   string
	SK   string
	Data string
}

type JWKSetsData struct {
	PK   string
	SK   string
	Data string
}

type UserDataStorage interface {
	GetUserData(username string) (*UserData, error)
	SetUserData(username string, user UserData) error
}

type KeyDataStorage interface {
	GetPrivateKey(name string) (*PrivateKeyData, error)
	SetPrivateKey(name string, key PrivateKeyData) error
	GetJWKSets(name string) (*JWKSetsData, error)
	SetJWKSets(name string, jwkSets JWKSetsData) error
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
func GetUserData(username string) (interface{}, error) {
	return nil, nil
}

func SetUserData(userdata interface{}) error {
	return nil
}

func GetPrivateKey() (interface{}, error) {
	return nil, nil
}

func SetPrivateKey(key interface{}) error {
	return nil
}

func GetJWKS() (interface{}, error) {
	return nil, nil
}

func SetJWKS(j interface{}) error {
	return nil
}
