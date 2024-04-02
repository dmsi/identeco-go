package storage

type UserData struct {
	Username string
	Password string
	Hash     string
}

type Keys struct {
	PrivateKey []byte
	JWKS       []byte
}

type UsersStorage interface {
	ReadUserData(username string) (*UserData, error)
	WriteUserData(user UserData) error
}

type KeysStorage interface {
	ReadKeys() (*Keys, error)
	WriteKeys(k Keys) error
}
