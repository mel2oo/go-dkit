package xfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name string
		data []byte

		expectedName string
		expectedExt  string
	}{
		{"", nil, "default.bin", "bin"},
		{"", []byte{22}, "default.unknown", "unknown"},
		{"aa", nil, "aa.bin", "bin"},
		{"aa", []byte{22}, "aa.unknown", "unknown"},
		{"aa.exe", []byte{22}, "aa.exe", "exe"},
	}

	for _, tc := range testCases {
		f := New(tc.name, tc.data)
		assert.Equal(t, tc.expectedName, f.Name)
		assert.Equal(t, tc.expectedExt, f.Extension)
	}
}

func TestLoadFromRawID(t *testing.T) {
	testCases := []struct {
		str string

		expectedID   string
		expectedName string
		expectedExt  string
	}{
		{
			"file:aa.exe<e7c1a9b1-a18a-41b6-954b-6f2888cc0f66>",
			"e7c1a9b1-a18a-41b6-954b-6f2888cc0f66", "aa.exe", "exe",
		},
		{
			"file:aa.exe.bin<e7c1a9b1-a18a-41b6-954b-6f2888cc0f66>",
			"e7c1a9b1-a18a-41b6-954b-6f2888cc0f66", "aa.exe.bin", "bin",
		},
	}

	for _, tc := range testCases {
		f, err := LoadFromRawID(tc.str)
		assert.NoError(t, err)
		assert.Equal(t, tc.expectedID, f.ID)
		assert.Equal(t, tc.expectedName, f.Name)
		assert.Equal(t, tc.expectedExt, f.Extension)
	}
}
