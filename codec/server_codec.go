package codec

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gorilla/rpc"
	"github.com/kolo/xmlrpc"
)

type Codec struct {
	methods                  map[string]string
	defaultMethodByNamespace map[string]string
	defaultMethod            string
	parsers                  map[string]Parser
	defaultParser            Parser
}

func NewCodec() *Codec {
	return &Codec{
		methods:                  make(map[string]string),
		defaultMethodByNamespace: make(map[string]string),
		defaultMethod:            "",
		parsers:                  make(map[string]Parser),
		defaultParser:            nil,
	}
}

func (c *Codec) RegisterDefaultParser(parser Parser) {
	c.defaultParser = parser
}

func (c *Codec) RegisterMethod(method string) {
	c.methods[method] = method
}

func (c *Codec) RegisterMethodWithParser(method string, parser Parser) {
	c.methods[method] = method
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

	var request interface{}
	if err := xmlrpc.UnmarshalServerRequest(rawxml, &request); err != nil {
		return &CodecRequest{err: err}
	}
	xmlRequest := request.(map[string]interface{})
	userMethod := xmlRequest["methodName"].(string)
	serviceMethod := c.resolveServiceMethod(userMethod)

	parser := c.resolveParser(serviceMethod)

	return &CodecRequest{request: &serverRequest{xmlRequest, serviceMethod}, parser: parser}
}

func (c *Codec) resolveParser(requestMethod string) Parser {
	if parser, ok := c.parsers[requestMethod]; ok {
		return parser
	}
	return c.defaultParser
}

func (c *Codec) resolveServiceMethod(requestMethod string) string {
	namespace, methodStr := c.getNamespaceAndMethod(requestMethod)
	if _, ok := c.methods[requestMethod]; ok {
		return c.toLowerCase(namespace, methodStr)
	} else if method, ok := c.defaultMethodByNamespace[namespace]; ok {
		return method
	} else if c.defaultMethod != "" {
		return c.defaultMethod
	}
	return requestMethod
}

func (c *Codec) getNamespaceAndMethod(requestMethod string) (string, string) {
	//TODO:
	if len(requestMethod) > 1 {
		parts := strings.Split(requestMethod, ".")
		slice := parts[1:len(parts)]
		return parts[0], strings.Join(slice, ".")
	}
	return "", ""
}

func (c *Codec) toLowerCase(namespace, method string) string {
	//TODO:
	if namespace != "" && method != "" {
		r, n := utf8.DecodeRuneInString(method)
		if unicode.IsLower(r) {
			return namespace + "." + string(unicode.ToUpper(r)) + method[n:]
		}
	}
	return namespace + "." + method
}

type serverRequest struct {
	request       map[string]interface{}
	serviceMethod string
}

type CodecRequest struct {
	request *serverRequest
	err     error
	parser  Parser
}

func (c *CodecRequest) Method() (string, error) {
	if c.err == nil {
		return c.request.serviceMethod, nil
	}
	return "", c.err
}

func (c *CodecRequest) ReadRequest(args interface{}) error {
	c.err = c.parser(c.request.request, args)
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
