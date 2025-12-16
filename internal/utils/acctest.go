package utils

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
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

	// 👇 仅新增这1行（定义字符集长度）
	charsetLen := big.NewInt(int64(len(charset)))

	// 👇 循环逻辑和你完全一样，仅改随机索引生成
	for i := 0; i < length; i++ {
		// 👇 核心改动（替代你原来的 r.Intn()）
		// rand.Int 第一个参数固定传 rand.Reader，第二个传字符集长度（范围）
		randomIndex, _ := crand.Int(crand.Reader, charsetLen)
		builder.WriteByte(charset[randomIndex.Int64()])
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
