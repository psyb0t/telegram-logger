package badgerdb

import (
	"bytes"
	"testing"
)

func TestGetUserKey(t *testing.T) {
	testCases := []struct {
		name     string
		userID   string
		expected []byte
	}{
		{
			name:     "user ID 12345",
			userID:   "12345",
			expected: []byte("user-12345"),
		},
		{
			name:     "user ID abcdef",
			userID:   "abcdef",
			expected: []byte("user-abcdef"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := getUserKey(tc.userID)
			if !bytes.Equal(actual, tc.expected) {
				t.Errorf("got %v, want %v", actual, tc.expected)
			}
		})
	}
}
