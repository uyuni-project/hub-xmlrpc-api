package codec

import (
	"fmt"
	"strings"
	"testing"
)

func TestFault2XML(t *testing.T) {
	const faultXMLPayload = `
	<methodResponse><fault><value><struct><member><name>faultCode</name><value><int>%d</int></value></member><member><name>faultString</name><value><string>%s</string></value></member></struct></value></fault></methodResponse>`

	//Formatted:
	//
	//<methodResponse>
	//	<fault>
	//		<value>
	//			<struct>
	//				<member>
	//					<name>faultCode</name>
	//					<value><int>%d</int></value>
	//				</member>
	//				<member>
	//					<name>faultString</name>
	//					<value><string>%s</string></value>
	//				</member>
	//			</struct>
	//		</value>
	//	</fault>
	//</methodResponse>

	tt := []struct {
		name Fault
	}{
		{name: FaultApplicationError},
		{name: FaultInternalError},
		{name: FaultInvalidCredentials},
	}

	for _, tc := range tt {
		t.Run(tc.name.Message, func(t *testing.T) {
			expected := fmt.Sprintf(faultXMLPayload, tc.name.Code, tc.name.Message)
			actual, err := encodeFaultToXML(tc.name)
			if err != nil {
				t.Fatalf("Error ocurred when parsing fault to XML:%v", err)
			}

			if strings.TrimSpace(expected) != actual {
				t.Fatalf("Unexpected Result. Expected: %s, Got:%s", expected, actual)
			}
		})
	}
}
