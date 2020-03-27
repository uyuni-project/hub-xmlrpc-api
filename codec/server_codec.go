package codec

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/rpc"
	"github.com/kolo/xmlrpc"
)

type Codec struct {
	mappings                 map[string]string
	defaultMethodByNamespace map[string]string
	defaultMethod            string
	parsers                  map[string]Parser
	defaultParser            Parser
}

func NewCodec() *Codec {
	return &Codec{
		mappings:                 make(map[string]string),
		defaultMethodByNamespace: make(map[string]string),
		defaultMethod:            "",
		parsers:                  make(map[string]Parser),
		defaultParser:            nil,
	}
}

func (c *Codec) RegisterDefaultParser(parser Parser) {
	c.defaultParser = parser
}

func (c *Codec) RegisterMapping(mapping string, method string) {
	c.mappings[mapping] = method
}

func (c *Codec) RegisterMappingWithParser(mapping string, method string, parser Parser) {
	c.mappings[mapping] = method
	c.parsers[c.resolveServiceMethod(method)] = parser
}

func (c *Codec) RegisterDefaultMethod(method string, parser Parser) {
	c.defaultMethod = method
	c.parsers[c.resolveServiceMethod(method)] = parser
}

func (c *Codec) RegisterDefaultMethodForNamespace(namespace, method string, parser Parser) {
	c.defaultMethodByNamespace[namespace] = method
	c.parsers[c.resolveServiceMethod(method)] = parser
}

func (c *Codec) NewRequest(r *http.Request) rpc.CodecRequest {
	rawxml, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &CodecRequest{err: err}
	}
	defer r.Body.Close()

	r.Body = ioutil.NopCloser(bytes.NewBuffer(rawxml))

	var serverRequest ServerRequest
	if err := xmlrpc.UnmarshalServerRequest(rawxml, &serverRequest); err != nil {
		return &CodecRequest{err: err}
	}

	userMethod := serverRequest.MethodName
	serviceMethod := c.resolveServiceMethod(userMethod)
	parser := c.resolveParser(serviceMethod)

	return &CodecRequest{request: &serverRequest, serviceMethod: serviceMethod, parser: parser}
}

func (c *Codec) resolveParser(requestMethod string) Parser {
	if parser, ok := c.parsers[requestMethod]; ok {
		return parser
	}
	return c.defaultParser
}

func (c *Codec) resolveServiceMethod(requestMethod string) string {
	namespace := c.getNamespace(requestMethod)
	if method, ok := c.mappings[requestMethod]; ok {
		return method
	} else if method, ok := c.defaultMethodByNamespace[namespace]; ok {
		return method
	} else if c.defaultMethod != "" {
		return c.defaultMethod
	}
	return requestMethod
}

func (c *Codec) getNamespace(requestMethod string) string {
	if len(requestMethod) > 1 {
		parts := strings.Split(requestMethod, ".")
		return parts[0]
	}
	return ""
}

type ServerRequest struct {
	MethodName string        `xmlrpc:"methodName"`
	Params     []interface{} `xmlrpc:"params"`
}

type CodecRequest struct {
	serviceMethod string
	request       *ServerRequest
	parser        Parser
	err           error
}

func (c *CodecRequest) Method() (string, error) {
	if c.err == nil {
		return c.serviceMethod, nil
	}
	return "", c.err
}

func (c *CodecRequest) ReadRequest(args interface{}) error {
	c.err = c.parser(c.request, args)
	if c.err != nil {
		return c.err
	}
	return nil
}

func (c *CodecRequest) WriteResponse(w http.ResponseWriter, response interface{}, methodErr error) error {
	var encodedResponse []byte
	err := c.err
	if err == nil {
		err = methodErr
	}
	if err != nil {
		var fault FaultError

		switch c.err.(type) {
		case FaultError:
			fault = c.err.(FaultError)
		default:
			fault = FaultApplicationError
			fault.Message += fmt.Sprintf(": %v", err)
		}
		encodedResponse, err = encodeFaultErrorToXML(fault)
		if err != nil {
			return err
		}
	} else {
		encodedResponse, err = encodeResponseToXML(response)
		if err != nil {
			return err
		}
	}
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Write(encodedResponse)
	return nil
}
