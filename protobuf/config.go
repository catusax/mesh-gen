package protobuf

import (
	"flag"
)

var _test *bool
var _handler *string
var _mesh *string
var _namespace *string
var _port *string
var _name *string
var _options *string
var Flags = new(flag.FlagSet)

type Flag struct {
	Test      bool
	Handler   string
	Mesh      string
	Namespace string
	Port      string
	Name      string
	Options   string
}

func init() {
	_test = Flags.Bool("test", true, "generate testing code")
	_handler = Flags.String("handler", "handler", "generate handler code")
	_mesh = Flags.String("mesh", "", "service-mesh type for generated client")
	_namespace = Flags.String("namespace", "default", "k8s namespace of your service")
	_port = Flags.String("port", "8080", "grpc port of your service")
	_name = Flags.String("name", "localhost", "name your service")
	_options = Flags.String("options", "", "grpc options")
}

var _config *Flag = nil

func GetConfig() *Flag {
	if _config != nil {
		return _config
	}
	_config = &Flag{
		Test:      *_test,
		Handler:   *_handler,
		Mesh:      *_mesh,
		Namespace: *_namespace,
		Port:      *_port,
		Name:      *_name,
		Options:   *_options,
	}
	return _config
}
