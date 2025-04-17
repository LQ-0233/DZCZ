package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordEncrypt(t *testing.T) {
	password := "123456"
	hashedPassword, err := PasswordEncrypt(password)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hashedPassword)
	match := PasswordMatch(password, hashedPassword)
	fmt.Println(match)
	assert.True(t, match)
}
