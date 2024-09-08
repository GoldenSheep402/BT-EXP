package decode

import (
	"reflect"
	"testing"
)

func TestToDictionaries(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]interface{}
		hasError bool
	}{
		// 简单字典测试
		{
			input:    "d3:cow3:moo4:spam4:eggse",
			expected: map[string]interface{}{"cow": "moo", "spam": "eggs"},
			hasError: false,
		},
		// 嵌套字典测试
		{
			input:    "d3:cowd3:mooe4:spam4:eggse",
			expected: map[string]interface{}{"cow": map[string]interface{}{"moo": ""}, "spam": "eggs"},
			hasError: false,
		},
		// 列表嵌套测试
		{
			input:    "d4:spaml1:a1:bee",
			expected: map[string]interface{}{"spam": []string{"a", "b"}},
			hasError: false,
		},
		// 空字典测试
		{
			input:    "de",
			expected: map[string]interface{}{},
			hasError: false,
		},
		// 无效字典测试：缺少结束符 'e'
		{
			input:    "d3:cow3:moo4:spam4:eggs",
			expected: nil,
			hasError: true,
		},
		// 无效字典测试：缺少冒号
		{
			input:    "d3cow3mooe",
			expected: nil,
			hasError: true,
		},
		// 无效长度
		{
			input:    "d5:cow3:moo4:spam4:eggse",
			expected: nil,
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := ToDictionaries(test.input)

			// 检查是否有错误
			if (err != nil) != test.hasError {
				t.Errorf("expected error: %v, got: %v", test.hasError, err)
			}

			// 检查结果
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}
