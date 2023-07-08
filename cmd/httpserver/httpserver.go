package main

import (
	"os"

	"github.com/dmsi/identeco/pkg/storage"
	"github.com/dmsi/identeco/pkg/storage/mongodb/keysmongodb"
	"github.com/dmsi/identeco/pkg/storage/mongodb/usersmongodb"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/exp/slog"
)

func main() {
	lg := slog.Default()

	u, err := usersmongodb.New(lg, os.Getenv("MONGODB_URL"), "main", "users")
	if err != nil {
		panic(err)
	}

	err = u.WriteUserData(storage.UserData{
		Username: "boss",
		Hash:     "$l**$$$sh256",
	})
	if err != nil {
		panic(err)
	}

	user, err := u.ReadUserData("boss")
	if err != nil {
		panic(err)
	}
	lg.Info("read user", "user", *user)

	k, err := keysmongodb.New(lg, os.Getenv("MONGODB_URL"), "main", "keys")
	if err != nil {
		panic(err)
	}

	err = k.WritePrivateKey(storage.PrivateKeyData{
		Data: []byte("hello this is private key"),
	})
	if err != nil {
		panic(err)
	}

	pk, err := k.ReadPrivateKey()
	if err != nil {
		panic(err)
	}
	lg.Info("read private key", "key", string(pk.Data))

	err = k.WriteJWKSets(storage.JWKSetsData{
		Data: []byte("hoy! this is JWKS"),
	})
	if err != nil {
		panic(err)
	}

	j, err := k.ReadJWKSets()
	if err != nil {
		panic(err)
	}
	lg.Info("read jwk sets", "jwks", string(j.Data))
}
