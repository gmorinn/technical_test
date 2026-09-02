package utils

import (
	"fmt"
	"strings"
	"testing"
)

// Testing ErrorResponse function
func TestErrorResponse(t *testing.T) {
	errMsg := "Some error"
	errCode := fmt.Errorf("ErrorCode")
	err := ErrorResponse(errMsg, errCode)
	if err == nil || !strings.Contains(err.Error(), errMsg) || !strings.Contains(err.Error(), errCode.Error()) {
		t.Error("Error response not as expected")
	}
}
