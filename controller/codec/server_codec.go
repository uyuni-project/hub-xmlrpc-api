package codec

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/rpc"
	"github.com/kolo/xmlrpc"
	"github.com/uyuni-project/hub-xmlrpc-api/controller"
)

type Codec struct {
	mappings                 map[string]string
	defaultMethodByNamespace map[string]string
	defaultMethod            string
	transformers             map[string]Transformer
}

type Transformer func(request *ServerRequest, output interface{}) error

func NewCodec() *Codec {
	return &Codec{
		mappings:                 make(map[string]string),
		defaultMethodByNamespace: make(map[string]string),
		defaultMethod:            "",
		transformers:             make(map[string]Transformer),
	}
}

func (c *Codec) RegisterMapping(mapping string, method string, transformer Transformer) {
	c.mappings[mapping] = method
	c.transformers[c.resolveServiceMethod(method)] = transformer
}

func (c *Codec) RegisterDefaultMethod(method string, transformer Transformer) {
	c.defaultMethod = method
	c.transformers[c.resolveServiceMethod(method)] = transformer
}

func (c *Codec) RegisterDefaultMethodForNamespace(namespace, method string, transformer Transformer) {
	c.defaultMethodByNamespace[namespace] = method
	c.transformers[c.resolveServiceMethod(method)] = transformer
}

func (c *Codec) NewRequest(r *http.Request) rpc.CodecRequest {
	rawxml, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &CodecRequest{err: err}
	}
	defer r.Body.Close()

	r.Body = ioutil.NopCloser(bytes.NewBuffer(rawxml))

	var serverRequest ServerRequest
	if err := xmlrpc.UnmarshalMethodCall(rawxml, &serverRequest); err != nil {
		return &CodecRequest{err: err}
	}

	userMethod := serverRequest.MethodName
	serviceMethod := c.resolveServiceMethod(userMethod)
	transformer := c.resolveTransformer(serviceMethod)

	return &CodecRequest{request: &serverRequest, serviceMethod: serviceMethod, transformer: transformer}
}

func (c *Codec) resolveTransformer(requestMethod string) Transformer {
	if transformer, ok := c.transformers[requestMethod]; ok {
		return transformer
	}
	return nil
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
	transformer   Transformer
	err           error
}

func (c *CodecRequest) Method() (string, error) {
	if c.err == nil {
		return c.serviceMethod, nil
	}
	return "", c.err
}

func (c *CodecRequest) ReadRequest(args interface{}) error {
	if c.transformer == nil {
		return controller.FaultInternalError
	}
	c.err = c.transformer(c.request, args)
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
		var fault controller.FaultError

		switch c.err.(type) {
		case controller.FaultError:
			fault = c.err.(controller.FaultError)
		default:
			fault = controller.FaultApplicationError
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
