package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type File struct {
	PackageName string
	Imports     []string
	Methods     []*Method
}

type Method struct {
	Name string
	Code string
	Docs []string
}

func ParseFile(src []byte, structName string) (*File, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	file := &File{
		PackageName: f.Name.String(),
		Imports:     make([]string, 0),
		Methods:     make([]*Method, 0),
	}

	for _, i := range f.Imports {
		if i.Name != nil {
			if i.Name.String() == "_" {
				continue
			}
			file.Imports = append(file.Imports, fmt.Sprintf("%s %s", i.Name.String(), i.Path.Value))
		} else {
			file.Imports = append(file.Imports, i.Path.Value)
		}
	}

	for _, d := range f.Decls {
		if !isExportedMethod(d) {
			continue
		}
		fd := d.(*ast.FuncDecl)
		if strings.TrimPrefix(ParseExpr(fd.Recv.List[0].Type), "*") != structName {
			continue
		}
		m := &Method{
			Name: fd.Name.String(),
		}
		_ = ParseFieldList(fd.Type.TypeParams)
		params := ParseFieldList(fd.Type.Params)
		ret := ParseFieldList(fd.Type.Results)

		m.Code = fmt.Sprintf("%s(%s) (%s)", m.Name, strings.Join(params, ", "), strings.Join(ret, ", "))

		if fd.Doc != nil {
			for _, doc := range fd.Doc.List {
				m.Docs = append(m.Docs, doc.Text)
			}
		}
		file.Methods = append(file.Methods, m)
	}
	return file, nil
}

func ParseFieldList(fields *ast.FieldList) []string {
	if fields == nil {
		return nil
	}
	results := make([]string, 0)
	for _, f := range fields.List {
		names := make([]string, 0)
		for _, name := range f.Names {
			names = append(names, name.Name)
		}
		results = append(results, fmt.Sprintf("%s %s", strings.Join(names, ", "), ParseExpr(f.Type)))
	}
	return results
}

func isExportedMethod(d ast.Decl) bool {
	fd, ok := d.(*ast.FuncDecl)
	if !ok {
		return false
	}
	return fd.Recv != nil && fd.Name.IsExported()
}

func ParseExpr(x ast.Expr) string {
	switch t := x.(type) {
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", ParseExpr(t.X))
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", ParseExpr(t.X), t.Sel.Name)
	case *ast.ArrayType:
		return fmt.Sprintf("[]%s", ParseExpr(t.Elt))
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", ParseExpr(t.Key), ParseExpr(t.Value))
	default:
		panic(fmt.Errorf("tt:  %T\n", t))
	}
}

func ParsePackage(dir string, structName string) (results []*File, _ error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), "_test.go") || strings.HasSuffix(file.Name(), ".iface.go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		file, err := ParseFile(data, structName)
		if err != nil {
			return nil, err
		}
		results = append(results, file)
	}
	return results, nil
}

func getStructName(src []byte) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return ""
	}
	l, _ := strconv.Atoi(os.Getenv("GOLINE"))
	for _, d := range f.Decls {
		st, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}
		switch t := st.Specs[0].(type) {
		case *ast.TypeSpec:
			pos := fset.Position(t.Pos())
			if pos.Line < l {
				continue
			}
			return t.Name.String()
		}
	}
	return ""
}
