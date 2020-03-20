package server

import (
	"fmt"
	"strings"
	"testing"
)

func TestFault2XML(t *testing.T) {
	const faultXMLPayload = `
	<methodResponse><fault><value><struct><member><name>faultCode</name><value><int>%d</int></value></member><member><name>faultString</name><value><string>%s</string></value></member></struct></value></fault></methodResponse>`
	tt := []struct {
		name Fault
	}{
		{name: FaultApplicationError},
		{name: FaultInternalError},
		{name: FaultInvalidCredentials},
	}

	for _, tc := range tt {
		t.Run(tc.name.String, func(t *testing.T) {

			expected := fmt.Sprintf(faultXMLPayload, tc.name.Code, tc.name.String)
			actual := fault2XML(tc.name)
			if strings.TrimSpace(expected) != actual {
				t.Fatalf("Unexpected Result. Expected: %s, Got:%s", expected, actual)
			}

		})
	}
}
