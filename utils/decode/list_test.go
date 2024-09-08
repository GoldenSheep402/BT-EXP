package decode

import (
	"reflect"
	"testing"
)

func TestToList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	tests := []struct {
		input    string
		expected []string
		hasError bool
	}{
		// 正常的列表输入
		{"l4:spam4:eggse", []string{"spam", "eggs"}, false},
		// 无效的输入 - 缺少冒号
		{"l4spam4eggse", nil, true},
		// 无效的输入 - 长度错误
		{"l5:spam4:eggse", nil, true},
		// 无效的输入 - 不以 'e' 结尾
		{"l4:spam4:eggs", nil, true},
		// 边界情况 - 空列表
		{"le", []string{}, false},
		// 无效的输入 - 非法长度
		{"l-4:spam4:eggse", nil, true},
		// 无效的输入 - 缺少 'l' 开头
		{"4:spam4:eggse", nil, true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := ToList(test.input)

			if (err != nil) != test.hasError {
				t.Errorf("expected error: %v, got: %v", test.hasError, err)
			}

			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}
