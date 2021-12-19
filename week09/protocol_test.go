package week09

import (
	"fmt"
	"testing"
)

func Test_decoder(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "测试 hello",
			data: "hello",
		},
		{
			name: "测试空字符",
			data: "",
		},
		{
			name: "测试空字符",
			data: "hellohellohellohellohellohellohellohellohellohellohellohellohello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := encode(tt.data)
			p := decode(buf)
			body := string(p.Body)
			fmt.Printf("body: %s\n", body)
			if body != tt.data {
				t.Errorf("解码失败: %s", tt.data)
			}
		})
	}
}
