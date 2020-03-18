package server

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"reflect"
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
			httpRequest:          buildSuccessHTTPRequest(),
			expectedCodecRequest: buildSuccessCodecRequest(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			codec := NewCodec()
			codec.RegisterDefaultParser(new(StructParser))
			codecRequest := codec.NewRequest(tc.httpRequest)

			if !reflect.DeepEqual(codecRequest, &tc.expectedCodecRequest) {
				t.Fatalf("Expected and actual don't match")
			}
		})
	}
}

func buildSuccessCodecRequest() CodecRequest {
	var request ServerRequest
	rawxml := []byte(successRequestBody)
	xml.Unmarshal(rawxml, &request)
	request.rawxml = rawxml

	return CodecRequest{request: &request, parser: new(StructParser)}
}

func buildSuccessHTTPRequest() *http.Request {
	httpRequest := httptest.NewRequest(http.MethodPost, "http://localhost:8888/hub/rpc/api", strings.NewReader(successRequestBody))
	httpRequest.Header.Set("Content-Type", "text/xml")

	return httpRequest
}

const successRequestBody = `
			<methodName>hub.login</methodName>
				<params>
					<param>
						<value><string>admin</string></value>
					</param>
					<param>
						<value><string>admin</string></value>
					</param>
				</params>
			</methodCall>`
