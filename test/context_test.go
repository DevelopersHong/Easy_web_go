package test1

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestConText
func TestContext(t *testing.T) {
	testCases := []struct {
		name      string
		input     context.Context
		field     string
		outputVal interface{}
		outputErr error
	}{
		{
			name:  "test1",
			input: context.Background(),
		},
	}
	for _, tc := range testCases {
		timeoutCtx, cancel := context.WithTimeout(tc.input, time.Second)
		defer cancel()
		time.Sleep(time.Second)
		err := timeoutCtx.Err()
		assert.Equal(t, tc.outputErr, err)
		fmt.Println(err)
	}

}
