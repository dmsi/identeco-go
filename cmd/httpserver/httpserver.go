package main

import (
	"os"

	"github.com/dmsi/identeco/pkg/storage"
	"github.com/dmsi/identeco/pkg/storage/mongodb/userdata"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/exp/slog"
)

func main() {
	lg := slog.Default()

	u, err := userdata.New(lg, os.Getenv("MONGODB_URL"), "main", "users")
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
}
