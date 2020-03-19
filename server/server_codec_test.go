package server

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func Test_RegisterDefaultMethodForNamespace(t *testing.T) {
	tt := []struct {
		name      string
		namespace string
		method    string
		parser    Parser
	}{
		{name: "RegisterDefaultMethodForNamespace",
			namespace: "multicast",
			method:    "MulticastService.DefaultMethod",
			parser:    new(MulticastArgsParser),
		},
		{name: "RegisterDefaultMethodForNamespace",
			namespace: "multicast",
			method:    "MulticastService.DefaultMethod",
			parser:    nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			codec := NewCodec()
			codec.RegisterDefaultMethodForNamespace(tc.namespace, tc.method, tc.parser)
			if codec.defaultMethodByNamespace[tc.namespace] != tc.method {
				t.Fatalf("defaultMethodByNamespace doesn't match, Expected value was: %v", tc.method)
			}
			if codec.parsers["MulticastService.DefaultMethod"] != tc.parser {
				t.Fatalf("parser for method doesn't match, Expected was: %v", reflect.TypeOf(tc.parser).String())
			}
		})
	}
}

func Test_RegisterDefaultMethod(t *testing.T) {
	tt := []struct {
		name   string
		method string
		parser Parser
	}{
		{name: "RegisterDefaultMethod",
			method: "MulticastService.DefaultMethod",
			parser: new(MulticastArgsParser),
		},
		{name: "RegisterDefaultMethod",
			method: "MulticastService.DefaultMethod",
			parser: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			codec := NewCodec()
			codec.RegisterDefaultMethod(tc.method, tc.parser)
			if codec.defaultMethod != tc.method {
				t.Fatalf("defaultMethod doesn't match, Expected value was: %v", tc.method)
			}
			if codec.parsers["MulticastService.DefaultMethod"] != tc.parser {
				t.Fatalf("parser for method doesn't match, Expected was: %v", reflect.TypeOf(tc.parser).String())
			}
		})
	}
}

func Test_resolveMethod(t *testing.T) {
	tt := []struct {
		name                      string
		namespace                 string
		defaultMethodForNamespace string
		defaultMethod             string
		method                    string
		expectedMethod            string
	}{
		{name: "CodecRequest resolveMethod success",
			defaultMethod:             "DefaultService.DefaultMethod",
			namespace:                 "multicast",
			defaultMethodForNamespace: "MulticastService.DefaultMethod",
			method:                    "multicastService.method",
			expectedMethod:            "multicastService.Method",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			codec := NewCodec()
			codec.RegisterDefaultMethod(tc.defaultMethod, new(StructParser))
			codec.RegisterDefaultMethodForNamespace(tc.namespace, tc.defaultMethodForNamespace, new(StructParser))
			codec.RegisterMethod(tc.method)

			method := codec.resolveMethod(tc.method)
			if method != tc.expectedMethod {
				t.Fatalf("Method doesn't match. Method value was: %v, expected value was: %v", method, tc.expectedMethod)
			}

			method = codec.resolveMethod(tc.namespace + "." + "unregistered_method")
			if method != tc.defaultMethodForNamespace {
				t.Fatalf("Method doesn't match with the defaultMethodForNamespace. Method value was: %v, expected value was: %v", method, tc.defaultMethodForNamespace)
			}

			method = codec.resolveMethod("unregistered_method")
			if method != tc.defaultMethod {
				t.Fatalf("Method doesn't match with the defaultMethod. Method value was: %v, expected value was: %v", method, tc.defaultMethod)
			}
		})
	}
}

func Test_resolveParser(t *testing.T) {
	tt := []struct {
		name          string
		method        string
		defaultParser Parser
	}{
		{name: "CodecRequest resolveParser Success",
			method:        "multicastService.method",
			defaultParser: new(ListParser),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			codec := NewCodec()
			codec.RegisterDefaultParser(tc.defaultParser)

			parser := codec.resolveParser("unregistered_method")
			if parser != tc.defaultParser {
				t.Fatalf("Parser doesn't match with the defaultParser. Parser value was: %v, expected value was: %v", reflect.TypeOf(parser).String(), reflect.TypeOf(tc.defaultParser).String())
			}

			parser = codec.resolveParser(tc.method)
			if parser != tc.defaultParser {
				t.Fatalf("Parser doesn't match with the defaultParser. Parser value was: %v, expected value was: %v", reflect.TypeOf(parser).String(), reflect.TypeOf(tc.defaultParser).String())
			}
		})
	}
}

func Test_NewRequest(t *testing.T) {
	tt := []struct {
		name                 string
		httpRequest          *http.Request
		expectedCodecRequest CodecRequest
	}{
		{name: "Create new CodecRequest Success",
			httpRequest:          buildSuccessHTTPRequestWithBody(successRequestBody),
			expectedCodecRequest: buildSuccessCodecRequestWithBody(successRequestBody),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			codec := NewCodec()
			codec.RegisterDefaultParser(new(MulticastArgsParser))
			codecRequest := codec.NewRequest(tc.httpRequest)

			if !reflect.DeepEqual(codecRequest, &tc.expectedCodecRequest) {
				t.Fatalf("Expected and actual don't match")
			}
		})
	}
}

const sessionKey = "39x0d8f9d78559bf45bc41d5145444fc4e05489549f2f99a8d30768b478fe219dc2"

var serverIds = []int64{1000010001, 1000010002}

var serverArgs = [][]interface{}{{"arg1_Server1", "arg1_Server2"}, {"arg2_Server1", "arg2_Server2"}}

func Test_CodecRequest_ReadRequest(t *testing.T) {
	tt := []struct {
		name                 string
		httpRequest          *http.Request
		expectedCodecRequest CodecRequest
		structToHydrate      interface{}
		expectedStruct       MulticastArgs
		expectedError        bool
	}{
		{name: "Create new CodecRequest Success",
			httpRequest:          buildSuccessHTTPRequestWithBody(successRequestBodyForMulticastArgs),
			expectedCodecRequest: buildSuccessCodecRequestWithBody(successRequestBodyForMulticastArgs),
			structToHydrate:      &MulticastArgs{},
			expectedStruct:       MulticastArgs{HubSessionKey: sessionKey, ServerIDs: serverIds, ServerArgs: serverArgs},
			expectedError:        false,
		},
		{name: "Create new CodecRequest no_pointer_structToHydrate_pass Failed",
			httpRequest:          buildSuccessHTTPRequestWithBody(successRequestBodyForMulticastArgs),
			expectedCodecRequest: buildSuccessCodecRequestWithBody(successRequestBodyForMulticastArgs),
			structToHydrate:      MulticastArgs{},
			expectedStruct:       MulticastArgs{},
			expectedError:        true,
		},
		{name: "Create new CodecRequest error_when_unmarshalling Failed",
			httpRequest:          buildSuccessHTTPRequestWithBody(brokenRequestBodyForMulticastArgs),
			expectedCodecRequest: buildSuccessCodecRequestWithBody(brokenRequestBodyForMulticastArgs),
			structToHydrate:      &MulticastArgs{},
			expectedStruct:       MulticastArgs{},
			expectedError:        true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			codec := NewCodec()
			codec.RegisterDefaultParser(new(MulticastArgsParser))

			codecRequest := codec.NewRequest(tc.httpRequest)

			if !reflect.DeepEqual(codecRequest, &tc.expectedCodecRequest) {
				t.Fatalf("Expected and actual don't match")
			}

			err := codecRequest.ReadRequest(tc.structToHydrate)

			if err != nil && !tc.expectedError {
				t.Fatalf("Error ocurred when reading request. %v", err)
			}

			if err == nil && !reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) {
				t.Fatalf("Expected and actual don't match")
			}
		})
	}
}

func buildSuccessCodecRequestWithBody(body string) CodecRequest {
	var request ServerRequest
	rawxml := []byte(body)
	xml.Unmarshal(rawxml, &request)
	request.rawxml = rawxml

	return CodecRequest{request: &request, parser: new(MulticastArgsParser)}
}

func buildSuccessHTTPRequestWithBody(body string) *http.Request {
	httpRequest := httptest.NewRequest(http.MethodPost, "http://localhost:8888/hub/rpc/api", strings.NewReader(body))
	httpRequest.Header.Set("Content-Type", "text/xml")

	return httpRequest
}

const successRequestBody = `
<?xml version='1.0'?>
<methodCall>
	<methodName>multicast.system.listUserSystems</methodName>
	<params>
		<param>
			<value><string>39x0d8f9d78559bf45bc41d5145444fc4e05489549f2f99a8d30768b478fe219dc2</string></value>
		</param><methodCall>
		<methodName>hub.login</methodName>
		<params>
			<param>
				<value><string>admin</string></value>
			</param>
			<param>
				<value><string>admin</string></value>
			</param>
		</params>
	</methodCall>
		<param>
			<value><array><data>
			<value><int>1000010000</int></value>
			<value><int>1000010001</int></value>
			</data></array></value>
		</param>
		<param>
			<value><array><data>
			<value><string>admin</string></value>
			<value><string>admin</string></value>
			</data></array></value>
		</param>
	</params>
</methodCall>`

var successRequestBodyForMulticastArgs = `
<?xml version='1.0'?>
<methodCall>
	<methodName>multicast.system.listUserSystems</methodName>
	<params>
		<param>
			<value><string>` + sessionKey + `</string></value>
		</param>
		<param>
			<value><array><data>
			<value><int>` + strconv.FormatInt(serverIds[0], 10) + `</int></value>
			<value><int>` + strconv.FormatInt(serverIds[1], 10) + `</int></value>
			</data></array></value>
		</param>
		<param>
			<value><array><data>
			<value><string>` + serverArgs[0][0].(string) + `</string></value>
			<value><string>` + serverArgs[0][1].(string) + `</string></value>
			</data></array></value>
		</param>
		<param>
			<value><array><data>
			<value><string>` + serverArgs[1][0].(string) + `</string></value>
			<value><string>` + serverArgs[1][1].(string) + `</string></value>
			</data></array></value>
		</param>
	</params>
</methodCall>`

var brokenRequestBodyForMulticastArgs = `
<?xml version='1.0'?>
<methodCall>
	<methodName>multicast.system.listUserSystems</methodName>
	<params>
		<param>
			<value><string>` + sessionKey + `</string></value>
		</param>
		<param>
			<value><array><data>
			<value><int>abd</int></value>
			<value><int>` + strconv.FormatInt(serverIds[1], 10) + `</int></value>
			</data></array></value>
		</param>
		<param>
			<value><array><data>
			<value><string>` + serverArgs[0][0].(string) + `</string></value>
			<value><string>` + serverArgs[0][1].(string) + `</string></value>
			</data></array></value>
		</param>
		<param>
			<value><array><data>
			<value><string>` + serverArgs[1][0].(string) + `</string></value>
			<value><string>` + serverArgs[1][1].(string) + `</string></value>
			</data></array></value>
		</param>
	</params>
</methodCall>`
