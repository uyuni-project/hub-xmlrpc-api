package xmlrpc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/controller"
)

func Test_encodeFaultErrorToXML(t *testing.T) {
	tt := []struct {
		name          string
		faultError    controller.FaultError
		expectedError string
	}{
		{name: "encodeFaultErrorToXML Success",
			faultError: controller.FaultApplicationError,
		},
		{name: "encodeFaultErrorToXML Success",
			faultError: controller.FaultInternalError,
		},
		{name: "encodeFaultErrorToXML Success",
			faultError: controller.FaultInvalidCredentials,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name+" "+tc.faultError.Message, func(t *testing.T) {
			const faultErrorXMLPayload = `<methodResponse><fault><value><struct><member><name>faultCode</name><value><int>%d</int></value></member><member><name>faultString</name><value><string>%s</string></value></member></struct></value></fault></methodResponse>`

			/*	FORMATTED:

				<methodResponse>
					<fault>
						<value>
							<struct>
								<member>
									<name>faultCode</name>
									<value><int>%d</int></value>
								</member>
								<member>
									<name>faultString</name>
									<value><string>%s</string></value>
								</member>
							</struct>
						</value>
					</fault>
				</methodResponse> */

			expectedXML := fmt.Sprintf(faultErrorXMLPayload, tc.faultError.Code, tc.faultError.Message)
			encodedFaultError, err := encodeFaultErrorToXML(tc.faultError)
			if err != nil {
				t.Fatalf("Error ocurred when parsing fault to XML:%v", err)
			}

			if strings.TrimSpace(expectedXML) != string(encodedFaultError) {
				t.Fatalf("Unexpected Result. Expected: %s, Got:%s", expectedXML, encodedFaultError)
			}
		})
	}
}
