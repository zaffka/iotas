# Iota enums generator for Golang

This little tool helps you to create `String()`, `Marshal()` and `Unmarshal()` methods for any enumerated type in Go. Thus your code will satisfy the `fmt.Stringer` interface and also marshal-unmarshal enums as strings to and from JSON payloads.

The app inspired by Google's [Sringer](https://pkg.go.dev/golang.org/x/tools/cmd/stringer) tool.

:exclamation: The tool's functionality is a little bit reduced as of yet.  
See the TODO list below to see what's on the go.

As of yet, **it will work only** for enums declared like this:

```
const (
	Unknown MatrixType = iota
	OLED
	AMOLED
	TFT
)
```

**It does not support** enums declared with iota-shifted first constant, aliased constants, etc.:

```
const (
	Unknown MatrixType = iota + 1
	OLED
	AMOLED = OLED
	TFT MatrixType = iota
)
```

## Installation

```
go install github.com/zaffka/iotas@latest
```

## How it works

The code is searching for specific code block within the ast-tree.  
It extracts constant names starting with typed zero-valued iota constant.  
The search stops if any value or type attached to the constant found.

```
*ast.File {
     Name: *ast.Ident {
       Name: "examples"                 <-- loaded at the Parser initialization stage (parse.NewParser func)
     }
     Decls: []ast.Decl {
       0: *ast.GenDecl {                <-- finding a general ast declaration
           Tok: const                   <-- finding a token with constants
           Specs: []ast.Spec {
            0: *ast.ValueSpec {         <-- finding a first value specification
               Names: []*ast.Ident {
                    0: *ast.Ident {
                      Name: "Unknown"   <-- getting a name of the constant
                    }
               }
               Type: *ast.Ident {
                    Name: "MatrixType"  <-- getting a type name of the iota constant sequence
               }
               Values: []ast.Expr {
                    0: *ast.Ident {
                      Name: "iota"      <-- checking if the sequence starts with zero-valued iota
                    }
               }
            }
            1: *ast.ValueSpec {         <-- get next constant
               Names: []*ast.Ident {
                    0: *ast.Ident {
                      Name: "OLED"      <-- get its name
                    }
               }
               Type: nil                <-- check it has no type
               Values: nil              <-- ...and no values (overwise parsing stops)
            }
            2: *ast.ValueSpec {
               Names: []*ast.Ident {
                    0: *ast.Ident {
                      Name: "AMOLED"
                    }
               }
               Type: nil
               Values: nil
               Comment: nil
            }
           }
```

## TODOs

- [ ] add test cases for current functionality
- [ ] add automation build and package publishing
- [ ] add format support for marshalling ("SuperAMOLED"-> "super_amoled", "SuperAMOLED"-> "SUPER_AMOLED", etc.)
- [ ] add support for shifted iota-declaration
