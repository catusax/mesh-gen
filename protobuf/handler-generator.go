package protobuf

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"google.golang.org/protobuf/compiler/protogen"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// https://github.com/grpc/grpc-go/blob/master/cmd/protoc-gen-go-grpc/grpc.go#L215

// GenerateHandlerFile generates _handler to be implemented.
func GenerateHandlerFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + "_handler.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	//prevFile, err := os.Open(filename)
	entry, _ := os.ReadDir(GetConfig().Handler)

	//defer prevFile.Close()
	exist := len(entry) != 0
	if exist {
		//src, err := os.ReadFile(filename)
		//if err != nil {
		//	panic(err)
		//}

		g.Skip()

		addHandlerToFile(gen, file)

		return nil

	} else {

		return generateNewHandlerFile(g, file)
	}

}

type Decl struct {
	File string
	Decl ast.Decl
}

var files = make(map[string][]byte)

func addHandlerToFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	fset := token.NewFileSet()
	asts, err := ParsePackageDir(fset, GetConfig().Handler)
	if err != nil {
		panic(err)
	}

	astfiles := asts.Files

	var decls []Decl
	for handlerFileName, astfile := range astfiles {
		for i := range astfile.Decls {
			decls = append(decls, Decl{
				File: handlerFileName,
				Decl: astfile.Decls[i],
			})
		}

	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "===panic====\n\n%s\n\n====panic===")
		panic(err)
	}

	//TODO: remember src index ,so we don't need to parse ast again
	//var index = make([]int, len(src))
	//for i, _ := range index {
	//	index[i] = i
	//}

	//create client
	for _, srv := range file.Services {

	MethodLoop:
		for methodIndex, method := range srv.Methods {
			for _, d := range decls { //range ast decls，find if any match
				if genDecl, ok := d.Decl.(*ast.FuncDecl); ok &&
					isMethod(genDecl, method.GoName, srv.GoName) { //method already exists

					inputs := getFuncParamTypes(genDecl)
					outputs := getFuncResultTypes(genDecl)

					if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() { //normal method
						if len(inputs) > 0 && inputs[len(inputs)-1] == "*"+method.Input.GoIdent.GoName { //params not match

							if len(outputs) > 0 && outputs[0] == "*"+method.Output.GoIdent.GoName {
								continue MethodLoop
							}
						}

						changeMethod(d.File, file, gen, genDecl, method)
						return nil
					}

					if method.Desc.IsStreamingClient() { // client only have xxxServer param
						if len(inputs) > 0 && inputs[len(inputs)-1] == method.Parent.GoName+"_"+method.GoName+"Server" {
							continue MethodLoop
						}
						changeMethod(d.File, file, gen, genDecl, method)
						return nil
					}

					if method.Desc.IsStreamingServer() || !method.Desc.IsStreamingClient() { //server have two param, xxxrequest and xxxServer
						if len(inputs) > 1 &&
							inputs[len(inputs)-2] == "*"+method.Input.GoIdent.GoName &&
							inputs[len(inputs)-1] == method.Parent.GoName+"_"+method.GoName+"Server" {
							continue MethodLoop
						}
					}

					//method exist but not correct, try to change method signature
					//continue MethodLoop
					// replace signature
					changeMethod(d.File, file, gen, genDecl, method)
					return nil

				}
			}

			//method do not exist, generate and insert

			//fmt.Fprintf(os.Stderr, "\n\n===generate====%s==generate===== \n\n", method.GoName)

			var src []byte
			if v, ok := files[filepath.Join(GetConfig().Handler, srv.GoName+".go")]; ok {
				src = v
			} else {
				src, _ = os.ReadFile(filepath.Join(GetConfig().Handler, srv.GoName+".go"))
			}

			g := gen.NewGeneratedFile(srv.GoName+".go", file.GoImportPath)

			//find next method position
			var insertPosition = len(src)
			if len(srv.Methods) != methodIndex+1 {
			SearchNext:
				for _, d := range decls {
					if genDecl, ok := d.Decl.(*ast.FuncDecl); ok {
						if isMethod(genDecl, srv.Methods[methodIndex+1].GoName, srv.GoName) {
							if genDecl.Doc != nil {
								//fmt.Fprintf(os.Stderr, "\n\n$$$found doc %d $$$$ \n\n", genDecl.Doc.Pos())
								insertPosition = int(genDecl.Doc.Pos())
								break SearchNext
							}

							//fmt.Fprintf(os.Stderr, "\n\n$$$not found doc %d $$$$ \n\n", genDecl.Pos()-1)
							insertPosition = int(genDecl.Pos() - 1)
							break SearchNext
						}
					}
				}
			}

			// insert into insertPosition

			buff := bytes.NewBuffer(make([]byte, 0))

			buff.Write(src[:insertPosition])

			generateMethod(buff, g, srv, method)

			buff.Write(src[insertPosition:])

			// 修改数据后长度会变，需要重新解析语法树 , 否则需要记录每次插入的长度
			files[filepath.Join(GetConfig().Handler, srv.GoName+".go")] = buff.Bytes()
			addHandlerToFile(gen, file)
			g.Skip()
			return nil
		}
	}

	for k, v := range files {
		g := gen.NewGeneratedFile(k, file.GoImportPath)
		if len(v) >= 0 {
			_, err := g.Write(v)
			if err != nil {
				panic(err)
			}
		}
	}

	if !GetConfig().Test {
		return nil
	}

	for handlerFileName, astfile := range astfiles {
		for i := range astfile.Decls {
			decls = append(decls, Decl{
				File: handlerFileName,
				Decl: astfile.Decls[i],
			})
		}

		if !strings.HasSuffix(handlerFileName, "_test.go") {
			GenerateTestFile(gen, file, handlerFileName, astfile)
		}

	}

	return nil

}

func generateNewHandlerFile(g *protogen.GeneratedFile, file *protogen.File) *protogen.GeneratedFile {
	g.P("// Code generated by protoc-gen-go-_mesh-gen.")
	g.P()
	g.P("package _handler")
	g.P()

	g.P(`import(
	"context"
`)
	g.P("pb ", g.QualifiedGoIdent(file.GoImportPath.Ident(strings.Replace(file.GoImportPath.String(), ".", string(file.GoPackageName), 1))))
	g.P(")")

	//create client
	for _, srv := range file.Services {
		g.P("type ", srv.GoName, " struct {")
		g.P("	pb.Unimplemented", srv.GoName, "Server")
		g.P("}")

		for _, method := range srv.Methods {

			var buf bytes.Buffer
			generateMethod(&buf, g, srv, method)
			buf.WriteTo(g)

		}
	}

	return g
}

func generateMethod(buf *bytes.Buffer, g *protogen.GeneratedFile, srv *protogen.Service, method *protogen.Method) {
	g.Annotate(srv.GoName+"."+method.GoName, method.Location)

	FPrint(buf, g, method.Comments.Leading,
		"func ", serverSignature(g, method), " {")

	FPrint(buf, g, "//TODO:implement me")
	FPrint(buf, g, "panic(\"implement me\")")

	FPrint(buf, g, "}")

}

func changeMethod(fileName string, file *protogen.File, gen *protogen.Plugin, funcDecl *ast.FuncDecl, method *protogen.Method) {
	//method exist but not correct, try to change method signature
	//continue MethodLoop

	// replace signature

	var src []byte
	if v, ok := files[fileName]; ok {
		src = v
	} else {
		src, _ = os.ReadFile(fileName)
	}
	g := gen.NewGeneratedFile(fileName, file.GoImportPath)

	buff := bytes.NewBuffer(make([]byte, 0))

	buff.Write(src[:funcDecl.Pos()-1]) //position before "func"

	// find line end
	var end = int(funcDecl.Pos())
	for ; end < len(src); end++ {
		if src[end] == byte('\n') {
			break
		}
	}

	//fmt.Fprintf(os.Stderr, "\n\n===replace====%s==replace===== \n\n", method.GoName)

	FPrint(buff, g, "func ", serverSignature(g, method), " { //signature replaced by _mesh-gen due to proto change")
	buff.Write(src[end:])

	// 修改数据后长度会变，需要重新解析语法树 , 否则需要记录每次插入的长度
	files[fileName] = buff.Bytes()
	addHandlerToFile(gen, file)
	g.Skip()

}

func serverSignature(g *protogen.GeneratedFile, method *protogen.Method) string {
	var reqArgs []string
	ret := "error"
	if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
		reqArgs = append(reqArgs, "ctx context.Context")
		ret = "(*pb." + g.QualifiedGoIdent(method.Output.GoIdent) + ", error)"
	}
	if !method.Desc.IsStreamingClient() {
		reqArgs = append(reqArgs, "req "+"*pb."+g.QualifiedGoIdent(method.Input.GoIdent))
	}
	if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
		reqArgs = append(reqArgs, "stream pb."+method.Parent.GoName+"_"+method.GoName+"Server")
	}

	g.P()
	prefix := "(e *" + method.Parent.GoName + ")" + method.GoName

	return prefix + "(" + strings.Join(reqArgs, ", ") + ") " + ret
}

func FPrint(buf io.Writer, g *protogen.GeneratedFile, v ...interface{}) {
	for _, x := range v {
		switch x := x.(type) {
		case protogen.GoIdent:
			fmt.Fprint(buf, g.QualifiedGoIdent(x))
		default:
			fmt.Fprint(buf, x)
		}
	}

	fmt.Fprintln(buf)
}

func ParsePackageDir(fset *token.FileSet, path string) (pkg *ast.Package, first error) {
	pkgs, err := parser.ParseDir(fset, path, func(info fs.FileInfo) bool {
		if _, ok := files[info.Name()]; ok {
			return false
		}
		return true
	}, 0)
	if err != nil {
		return nil, err
	}

	for k, v := range files {
		f, err := parser.ParseFile(fset, "", v, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		pkgs[GetConfig().Handler].Files[k] = f
	}

	return pkgs[GetConfig().Handler], nil

}
