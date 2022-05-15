package protobuf

import "go/ast"

func isMethod(genDecl *ast.FuncDecl, method, receiver string) bool {

	var funcReceiver = getFuncReceiverType(genDecl)

	if (funcReceiver == receiver || funcReceiver == "*"+receiver) && genDecl.Name.Name == method {
		return true
	}
	return false
}

func isFunc(genDecl *ast.FuncDecl, method string) bool {
	return genDecl.Name.Name == method
}

func getFuncReceiverType(genDecl *ast.FuncDecl) string {
	if genDecl.Recv.List == nil { // nil or len() == 1
		return ""
	}
	return getTypeName(genDecl.Recv.List[0].Type)

}

func getFuncParamTypes(genDecl *ast.FuncDecl) []string {
	var res []string
	for _, asd := range genDecl.Type.Params.List {
		res = append(res, getTypeName(asd.Type))
	}
	return res
}

func getFuncResultTypes(genDecl *ast.FuncDecl) []string {
	var res []string
	for _, asd := range genDecl.Type.Results.List {
		res = append(res, getTypeName(asd.Type))
	}
	return res

}

func getTypeName(expr ast.Expr) string {

	star, ok := expr.(*ast.StarExpr)
	if ok {
		return "*" + getTypeName(star.X)
	} else {

		switch expr.(type) {
		case *ast.Ident:
			return expr.(*ast.Ident).String()
		case *ast.SelectorExpr:
			selc := expr.(*ast.SelectorExpr)
			return selc.Sel.String()
		}

	}
	//asd := expr.(*ast.Ident)

	return ""

}
