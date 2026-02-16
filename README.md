# Page

Page is a parser generator written in Go, designed to facilitate the creation of parsers, lexers, and abstract syntax trees (AST) for custom languages.

## Features

- **Grammar Definition**: Define your language grammar using a clear and concise syntax or load from files.
- **Lexer Generation**: Automatically generates a lexer based on your grammar rules.
- **Parser Generation**: Generates a parser to process your language's input.
- **Abstract Syntax Tree (AST)**: Automatically builds an AST for easy traversal and manipulation.
- **Visitor Pattern**: Includes support for the visitor pattern to walk the AST.

## Architecture

The project is divided into **Build Time** (Grammar Processing) and **Runtime** (Parsing).

```mermaid
graph TD
    subgraph Build Time
        G[Grammar File (.gp)] -->|Load| GP[Grammar Parser]
        GP -->|Generate| V[Vocabulary]
        GP -->|Generate| S[Syntax]
    end

    subgraph Runtime
        Src[Source Code] -->|Read| I[Input Stream]
        I -->|Feed| L[Lexer]
        V -.->|Config| L
        L -->|Tokens| P[Parser]
        S -.->|Config| P
        P -->|Build| AST[Abstract Syntax Tree]
        AST -->|Walk| Vis[Visitor]
    end
```

## Installation

```bash
go get github.com/fabiouggeri/page
```

## Usage

### Defining a Grammar

You can define your grammar programmatically or load it from a file.

Example of defining a complex rule for an identifier:
```go
func id() *rule.NonTerminalRule {
	// matches identifiers starting with letter or underscore
	return rule.New("id",
		rule.And(
			rule.Or(
				rule.And(rule.OneOrMore(rule.Char('_')), alphanum()),
				letter()
			),
			rule.ZeroOrMore(idChar())
		)
	)
}
```

### Parsing Input based on a Grammar File

```go
package main

import (
	"fmt"
	"log"

	"github.com/fabiouggeri/page/build/grammar"
	"github.com/fabiouggeri/page/build/syntax"
	"github.com/fabiouggeri/page/build/vocabulary"
	"github.com/fabiouggeri/page/runtime/input"
	"github.com/fabiouggeri/page/runtime/lexer"
	"github.com/fabiouggeri/page/runtime/parser"
)

func main() {
    // 1. Load Grammar
	g, err := grammar.FromFile("path/to/grammar.gp")
	if err != nil {
		log.Fatal(err)
	}

    // 2. Prepare Input Source
	input, err := input.NewFileInput("path/to/source.code")
	if err != nil {
		log.Fatal(err)
	}

    // 3. Create Lexer and Parser
	v := vocabulary.FromGrammar(g)
	lex := lexer.New(v, input)
	syn := syntax.FromGrammar(g, v)
	p := parser.New(lex, syn)

    // 4. Parse
	rootNode := p.Execute()

	if len(p.Errors()) > 0 {
		fmt.Println("Parse errors:")
		for _, e := range p.Errors() {
			fmt.Printf("   %s\n", e.Message())
		}
	} else {
		fmt.Println("Parse successful!")
        // Process rootNode...
	}
}
```

## Project Structure

- `build/`: Contains build-time logic (grammar, automata conversion).
- `runtime/`: Runtime components (lexer, parser, visitor, input handling).
- `util/`: Utility functions.
- `examples/`: Example usage and test cases.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
