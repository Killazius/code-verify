package compilation

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type testCase struct {
	input    Lang
	expected bool
}

func TestIsValidLang(t *testing.T) {
	t.Run("check correct lang", func(t *testing.T) {
		testCases := []testCase{
			{
				input:    LangCpp,
				expected: true,
			},
			{
				input:    LangPy,
				expected: true,
			},
			{
				input:    LangGo,
				expected: true,
			},
			{
				input:    "",
				expected: false,
			},
			{
				input:    "lang",
				expected: false,
			},
		}
		t.Parallel()
		for _, tc := range testCases {
			assert.Equal(t, tc.expected, isValidLang(tc.input))
		}
	})

	t.Run("func is idempotent", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			assert.True(t, isValidLang(LangCpp))
		}
	})
}

func TestCreateFile(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		code      string
		lang      Lang
		expectErr bool
	}{
		{
			name:      "Valid language and successful file creation",
			filePath:  "testfile.cpp",
			code:      "#include <iostream>\nint main(){}",
			lang:      LangCpp,
			expectErr: false,
		},
		{
			name:      "Invalid language",
			filePath:  "testfile.txt",
			code:      "some code",
			lang:      "txt",
			expectErr: true,
		},
	}
	t.Parallel()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := CreateFile(tt.filePath, tt.code, tt.lang)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)

				content, readErr := os.ReadFile(tt.filePath)
				require.NoError(t, readErr)
				assert.Equal(t, tt.code, string(content))

				os.Remove(tt.filePath)
			}
		})
	}
}
