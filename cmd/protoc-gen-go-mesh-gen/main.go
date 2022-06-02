package main

import (
	"flag"
	"github.com/catusax/mesh-gen/protobuf"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var test *bool
var handler *bool
var mesh *string
var namespace *string
var port *string
var name *string

func main() {

	var flags flag.FlagSet
	test = flags.Bool("test", true, "generate testing code")
	handler = flags.Bool("handler", true, "generate handler code")
	mesh = flags.String("mesh", "", "service-mesh type for generated client")
	namespace = flags.String("namespace", "default", "k8s namespace of your service")
	port = flags.String("port", "8080", "grpc port of your service")
	name = flags.String("name", "localhost", "name your service")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			protobuf.GenerateClientFile(gen, f, *mesh, *namespace, *port, *name)

			if *test {
				protobuf.GenerateTestFile(gen, f)
			}

			if *handler {
				protobuf.GenerateHandlerFile(gen, f)
			}
		}
		return nil
	})
}
