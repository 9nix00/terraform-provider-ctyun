package utils

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

func LoadTestCase(filename string, parameters ...interface{}) string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	fullPath := filepath.Join(pwd, "testdata", filename)
	f, err := os.ReadFile(fullPath)
	if err != nil {
		return ""
	}

	return fmt.Sprintf(string(f), parameters...)
}

const charset = "abcdefghijklmnopqrstuvwxyz"

func GenerateRandomString() string {
	length := 10
	builder := strings.Builder{}
	builder.Grow(length)
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		builder.WriteByte(charset[randomIndex])
	}
	return builder.String()
}
