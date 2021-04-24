package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/urfave/cli/v2"
)

type Token struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
}

func CreateToken(uid uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = uid
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET"))); err != nil {
		return "", err
	} else {
		return token, nil
	}
}

func generateToken(ctx *cli.Context) error {
	return nil
}

func validateToken(ctx *cli.Context) error {
	return nil
}

func tokenCommands() *cli.Command {
	return &cli.Command{
		Name:  "jwt",
		Usage: "performs various jwt related operations, such as generating new jwt tokens, validating tokens, etc...",
		Subcommands: []*cli.Command{
			{
				Name:   "generate",
				Usage:  "generates a new jwt token",
				Action: generateToken,
			},
			{
				Name:   "validate",
				Usage:  "validates a token against the selected backend",
				Action: validateToken,
			},
		},
	}
}
