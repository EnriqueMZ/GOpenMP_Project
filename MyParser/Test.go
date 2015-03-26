package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	//"reflect"
)

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	// Parse the file containing this very example
	// but stop after processing the imports.
	f, err := parser.ParseFile(fset, "<stdin>", os.Stdin, parser.DeclarationErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the AST.
	//ast.Print(fset, f)

	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		var tokStr string
		var specs []ast.Spec
		switch x := n.(type) {
		case *ast.GenDecl:
			//fmt.Println("GenDecl: ", x)
			tokStr = x.Tok.String()
			specs = x.Specs
			if tokStr == "var" {
				if len(specs) == 0 {
					fmt.Println("Var sin declaraciones")
				} else {
					for j := range specs {
						fmt.Println("Specs: ", specs[j])
						switch y := specs[j].(type) {
						case *ast.ValueSpec:
							if y.Values != nil { // Variables inicializadas
								for i := range y.Names { 
									id := y.Names[i].Obj.Name
									typ := y.Type
									values := y.Values
									fmt.Println("Contenido de Values:", values[i])
									switch z := values[i].(type) {
									case *ast.BasicLit:
										val := z.Value
										fmt.Println("Variable inicializada BasicLit :", i, id, typ, val)
									case *ast.Ident:
										val := z.Name
										fmt.Println("Variable inicializada Ident:", i, id, typ, val)
									}
								}
							} else { // Variables sin inicializar
								for i := range y.Names {
									id := y.Names[i]
									typ := y.Type
									fmt.Println("Variable no inicializada:", i, id, typ, nil)
								}
							}
						}
					}
				}
			}
		}
		return true
	})
}
