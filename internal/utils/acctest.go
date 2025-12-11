package utils

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func AtLeastOne(value string) error {
	if value == "0" {
		return fmt.Errorf("列表应该至少有一个元素")
	}
	return nil
}

func GetSubdirectories(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var subdirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			subdirs = append(subdirs, entry.Name())
		}
	}
	return subdirs, nil
}

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
	return generateRandomStringWithRetry(0)
}

func generateRandomStringWithRetry(retryCount int) string {
	if retryCount > 1 { // 最多重试1次
		return ""
	}

	length := 10
	builder := strings.Builder{}
	builder.Grow(length)
	r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(retryCount)))

	// 生成字符串
	for i := 0; i < length; i++ {
		randomIndex := r.Intn(len(charset))
		builder.WriteByte(charset[randomIndex])
	}

	result := builder.String()

	// 检查是否有连续3个递增字符
	for i := 0; i < length-2; i++ {
		if result[i]+1 == result[i+1] && result[i+1]+1 == result[i+2] {
			// 如果有，重新生成一次
			return generateRandomStringWithRetry(retryCount + 1)
		}
	}

	return result
}
func GenerateRandomPort(min int, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// 生成一个 min 到 max 之间的随机数
	randomNum := rand.Intn(max-min) + min // Intn(n) 返回一个范围是 [0, n) 的随机数
	return randomNum
}
