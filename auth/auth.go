package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

type User struct {
	ID       uint   `jsonutils:"id"`
	Username string `jsonutils:"username"`
	Password string `jsonutils:"password"`
}

var user = User{
	ID:       0,
	Username: "admin",
	Password: "password",
}

func Login(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid JSON")
		return
	}
	if user.Username != u.Username || u.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Email or password do not match")
	}
	if token, err := CreateToken(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, token)
	}
}

func Commands() []*cli.Command {
	return []*cli.Command{
		authServer(),
	}
}
