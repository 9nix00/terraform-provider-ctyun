package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
)

type ChangeLog struct {
	Entries map[string][]string
}

const changelogDir = ".changelog"

func main() {
	cl := &ChangeLog{
		Entries: make(map[string][]string),
	}

	// 读取所有PR文件
	files, err := ioutil.ReadDir(changelogDir)
	if err != nil {
		panic(fmt.Sprintf("Error reading changelog directory: %v", err))
	}

	// 解析每个文件内容
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".txt" {
			processFile(cl, filepath.Join(changelogDir, f.Name()))
		}
	}

	// 生成最终CHANGELOG.md
	generateMarkdown(cl)
}

func processFile(cl *ChangeLog, path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Error reading file %s: %v", path, err))
	}

	lines := strings.Split(string(content), "\n")
	var currentType string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "release-note:") {
			currentType = strings.TrimPrefix(line, "release-note:")
			continue
		}

		if line != "" && currentType != "" {
			// 去重处理
			if !contains(cl.Entries[currentType], line) {
				cl.Entries[currentType] = append(cl.Entries[currentType], line)
			}
		}
	}
}

func generateMarkdown(cl *ChangeLog) {
	// 定义输出顺序和标题
	categories := []struct {
		Key   string
		Title string
	}{
		{"new-resource", "New Resources"},
		{"new-data-source", "New Data Sources"},
		{"enhancement", "Enhancements"},
		{"bug", "Bug Fixes"},
		{"deprecation", "Deprecations"},
	}

	// 创建模板
	tmpl := `# Changelog

{{range .Categories}}
## {{.Title}}
{{range $index, $entry := (index $.Entries .Key)}}
* {{$entry}}{{end}}
{{end}}`

	// 准备模板数据
	data := struct {
		Entries    map[string][]string
		Categories []struct {
			Key   string
			Title string
		}
	}{
		Entries:    cl.Entries,
		Categories: categories,
	}

	// 执行模板渲染
	t := template.Must(template.New("changelog").Parse(tmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		panic(fmt.Sprintf("Error executing template: %v", err))
	}

	// 写入文件
	if err := ioutil.WriteFile("CHANGELOG.md", buf.Bytes(), 0644); err != nil {
		panic(fmt.Sprintf("Error writing CHANGELOG.md: %v", err))
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
