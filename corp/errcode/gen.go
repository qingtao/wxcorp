// 根据微信企业号全局错误码描述信息生产go的错误文档
//	https://work.weixin.qq.com/api/doc#90000/90139/90313
//	取出微信企业号的错误代码markdown文本保存到相同目录下的${file}.md
// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	file = flag.String("file", "", "微信企业号全局错误码描述信息")
)

func main() {
	flag.Parse()
	wd, _ := os.Getwd()

	if *file == "" {
		fmt.Println("请提供全局错误码描述文件名")
		os.Exit(1)
	}

	f := filepath.Join(wd, *file)

	b, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Println("读取文件错误:", err.Error())
		os.Exit(1)
	}
	if len(b) < 1 {
		fmt.Println("读取文件错误")
		os.Exit(1)
	}

	// errorCode 错误码
	type errorCode struct {
		Data map[string]interface{}
	}

	var code errorCode

	if err = json.Unmarshal(b, &code); err != nil {
		fmt.Println("解析json错误:", err.Error())
		os.Exit(1)
	}

	if code.Data == nil {
		fmt.Println("全局错误码文本为空")
		os.Exit(1)
	}

	// fmt.Printf("%s\n", code)
	// fmt.Printf("%s\n", code.Data)
	doc, ok := code.Data["document"]
	if !ok || doc == nil {
		fmt.Println("全局错误码文本为空")
		os.Exit(1)
	}
	items, ok := doc.(map[string]interface{})
	if !ok || items == nil {
		fmt.Println("1.全局错误文本格式错误")
		os.Exit(1)
	}

	text, ok := items["content_txt"].(string)
	if !ok || text == "" {
		fmt.Println("2.全局错误文本格式错误")
		os.Exit(1)
	}
	// text = strings.Replace(text, "\r\n", "\n", -1)
	start := strings.Index(text, "-1")
	end := strings.LastIndex(text, "排查方法")
	text = text[start:end]

	re := regexp.MustCompile(`(-?\d+)\s+([^\n\r]+)\s*`)
	all := re.FindAllStringSubmatch(text, -1)
	// fmt.Println(text)

	var buf bytes.Buffer
	buf.WriteString(`// Package errcode 错误代码错误信息对照表
// 注意：此文件由代码生成工具生成，不要直接编辑
package errcode

import (
	"github.com/pkg/errors"
)

var (
	// ErrInvalidAccessToken 无效的访问令牌
	ErrInvalidAccessToken = errors.New("无效的访问令牌")
	// ErrUnknown 未知错误
	ErrUnknown = errors.New("未知错误")
)

// Error 取得错误信息
func Error(i int) error {
	if i == 0 {
		return nil
	}
	if i == 40014 || i == 40021 {
		return ErrInvalidAccessToken
	}
	msg, ok := errCode[i]
	if !ok {
		return ErrUnknown
	}
	return errors.New(msg)
}

var errCode = map[int]string{
`)
	for _, t := range all {
		if len(t) == 3 {
			buf.WriteByte('\t')
			buf.WriteString(t[1])
			buf.WriteString(": `")
			buf.WriteString(strings.TrimSpace(t[2]))
			buf.WriteString("`,\n")
		}
	}
	buf.WriteString("}\n\n")
	target := filepath.Join(wd, "global_errcode.go")
	err = ioutil.WriteFile(target, buf.Bytes(), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
