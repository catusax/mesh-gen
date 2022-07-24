//https://rotemtam.com/2021/03/22/creating-a-protoc-plugin-to-gen-go-code/

package protobuf

import (
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

var (
	osPackage           = protogen.GoImportPath("os")
	syncPackage         = protogen.GoImportPath("sync")
	grpcPackage         = protogen.GoImportPath("google.golang.org/grpc")
	grpcInsecurePackage = protogen.GoImportPath("google.golang.org/grpc/credentials/insecure")
)

// GenerateClientFile generates Service client with go-kit endpoint
func GenerateClientFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + "_client.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-go-_mesh-gen. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	//create client
	for _, srv := range file.Services {
		g.P("var _", srv.GoName, "Once ", syncPackage.Ident("Once"))
		g.P("var _", srv.GoName, "client ", srv.GoName, "Client")

		g.P("func Get", srv.GoName, "Client() ", srv.GoName, "Client {")
		g.P(` var host string`)
		g.P("if ", osPackage.Ident("Getenv"), "(\"CONTAINER\") != \"\" {")
		g.P("host = \"", getServiceMeshHost(GetConfig().Namespace, GetConfig().Mesh, GetConfig().Namespace), "\"")
		g.P(`
	} else {
		host = "localhost"
	}`)
		g.P("_", srv.GoName, "Once.Do(func() {")

		g.P(" dial, err := ", grpcPackage.Ident("Dial"), "(host+\":", _port, "\",",
			grpcPackage.Ident("WithTransportCredentials"),
			"(", grpcInsecurePackage.Ident("NewCredentials"), "()))")

		g.P(`if err != nil {
		panic(err)`)
		g.P("_", srv.GoName, "Once = ", syncPackage.Ident("Once{}"))
		g.P("}")

		g.P("_", srv.GoName, "client = New", srv.GoName, "Client(dial)")

		g.P("})")

		g.P("return _", srv.GoName, "client ")
		g.P("}")
	}

	return g
}

func getServiceMeshHost(name, mesh, namespace string) string {
	name = strings.ReplaceAll(name, "_", "-") + "-grpc"
	switch mesh {
	case "traefik-_mesh":
		return name + "." + namespace + ".traefik._mesh"
	case "istio":
		return name
	default:
		return name
	}
}
