package client

import (
	"bytes"
	"testing"
)

func TestNewClientAuth(t *testing.T) {
	testCookies := `uluru_user=redacted%3D%3D;XSRF-TOKEN=redacted;`
	input := bytes.NewBufferString(testCookies)
	c, err := NewClientAuth(input)
	if err != nil {
		t.Errorf("NewClientAuth failed: %v", err)
	}

	output := new(bytes.Buffer)
	if _, err := c.WriteTo(output); err != nil {
		t.Errorf("NewClientAuth: write to failed: %v", err)
	}
	if output.String() != testCookies {
		t.Errorf("NewClientAuth: expected: %s, unexpected: %s", testCookies, output.String())
	}
}
