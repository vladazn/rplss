package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

type JwtPkg struct {
	key string
}

type Player struct {
	jwt.RegisteredClaims
	Player string
}

func NewJwtPkg(key string) *JwtPkg {
	return &JwtPkg{key: key}
}

func (j JwtPkg) NewKey(playerName string) string {
	k, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Player{Player: playerName}).
		SignedString([]byte(j.key))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return k
}

func (j JwtPkg) ParseKey(token string) string {
	p := &Player{}
	t, err := jwt.ParseWithClaims(token, p,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.key), nil
		})
	if err != nil || !t.Valid {
		return ""
	}

	return p.Player
}
