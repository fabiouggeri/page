package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fabiouggeri/page/build/automata"
	"github.com/fabiouggeri/page/build/grammar"
	"github.com/fabiouggeri/page/build/rule"
	"github.com/fabiouggeri/page/build/syntax"
	"github.com/fabiouggeri/page/build/vocabulary"
	"github.com/fabiouggeri/page/runtime/ast"
	"github.com/fabiouggeri/page/runtime/input"
	"github.com/fabiouggeri/page/runtime/lexer"
	"github.com/fabiouggeri/page/runtime/parser"
	"github.com/fabiouggeri/page/util"
)

func main() {
	//g, err := grammar.FromFile("C:\\Users\\fabio\\dev\\gitlab\\yapp-2.0\\yapp-java-runtime\\src\\main\\java\\org\\uggeri\\yapp\\runtime\\java\\test\\Harbour.gy")
	//g, err := grammar.FromFile("/home/fabio_uggeri/dev/yapp/yapp-java-runtime/src/main/java/org/uggeri/yapp/runtime/java/test/Harbour.gy")
	//cg := go.NewGenerator()
	//err := cg.Generate("c:\\temp\\test-parser", cg.BuildParser())
	//testCode()
	testLexer()
}

func upperLetter() *rule.RangeRule {
	return rule.Range('A', 'Z')
}

func lowerLetter() *rule.RangeRule {
	return rule.Range('a', 'z')
}

func letter() *rule.OrRule {
	return rule.Or(lowerLetter(), upperLetter())
}

func digit() *rule.RangeRule {
	return rule.Range('0', '9')
}

func alphanum() *rule.OrRule {
	return rule.Or(lowerLetter(), upperLetter(), digit())
}

func idChar() *rule.OrRule {
	return rule.Or(lowerLetter(), upperLetter(), digit(), rule.Char('_'))
}

func id() *rule.NonTerminalRule {
	//rules.Build("((_+[a-zA-Z0-9])|([A-Za-z]))[a-zA-Z0-9_]*")
	return rule.New("id",
		rule.And(rule.Or(rule.And(rule.OneOrMore(rule.Char('_')), alphanum()), letter()), rule.ZeroOrMore(idChar())))
}

func functionKeyword() *rule.NonTerminalRule {
	return rule.New("function", rule.StringI("function"))
}

func openPar() *rule.NonTerminalRule {
	return rule.New("open_par", rule.String("("))
}

func closePar() *rule.NonTerminalRule {
	return rule.New("close_par", rule.String(")"))
}

func parameters() *rule.NonTerminalRule {
	return rule.New("parameters", id())
}

func functionName() *rule.NonTerminalRule {
	return rule.New("function_name", id())
}

func functionSyntax() *rule.NonTerminalRule {
	return rule.New("function_declaration", rule.And(rule.String("private"), functionKeyword(), functionName(), openPar(), parameters(), closePar(), rule.Char('{')))
}

func private() *rule.NonTerminalRule {
	return rule.New("private", rule.String("private"))
}

func plus() *rule.NonTerminalRule {
	return rule.New("plus", rule.Char('+'))
}

func minus() *rule.NonTerminalRule {
	return rule.New("minus", rule.Char('-'))
}

func space() *rule.NonTerminalRule {
	return rule.New("spaces", rule.OneOrMore(rule.Or(rule.Char(' '), rule.Char('\t'), rule.Char('\n'))))
}

func lineComment() *rule.NonTerminalRule {
	return rule.New("line_comment", rule.And(rule.String("//"), rule.ZeroOrMore(rule.Not(rule.Or(rule.Char('\n'), rule.EOI)))))
}

func blockComment() *rule.NonTerminalRule {
	return rule.New("block_comment", rule.And(rule.String("/*"), rule.ZeroOrMore(rule.Or(rule.Not(rule.Char('*')), rule.And(rule.Char('*'), rule.Not(rule.Char('/'))))), rule.String("*/")))
}

func testLexer() {
	g1, errGrammar := grammar.FromFile("C:\\Users\\fabio\\dev\\github\\page\\examples\\HarbourPP.gp")
	if errGrammar != nil {
		fmt.Print(errGrammar)
		return
	}
	fmt.Println("Grammar name: ", g1.Name())
	// g1 := grammar.New("teste")
	// err := g1.Rules(id(), functionKeyword(), private(), plus(), minus(), openPar(), closePar(), functionSyntax(), space(), lineComment(), blockComment())
	// if err != nil {
	// 	fmt.Print(err)
	// 	return
	// }
	w := util.NewStringTextWriter()
	g1.ToText(w)
	w.NewLine()
	// programs := builder.Build(g)
	// for _, p := range programs {
	// 	p.Generate(gg, w)
	// }

	lexerRules := g1.LexerRules()
	if g1.HasError() {
		fmt.Printf("Grammar errors:\n")
		for _, e := range g1.Errors() {
			fmt.Printf("   %s\n", e.Error())
		}
		return
	}
	nfa := vocabulary.RulesToNFA(lexerRules...)
	w.WriteString("==================== NFA ======================\n")
	w.WriteString(nfa.String())
	dfa := automata.NFAToDFA(nfa)
	w.WriteString("==================== DFA ======================\n")
	w.WriteString(dfa.String())

	//v := vocabulary.FromDFA(d)
	v := vocabulary.FromGrammar(g1)
	w.WriteString("==================== VOCABULARY ======================\n")
	v.Write(w)
	w.NewLine()
	w.WriteString("====================== SYNTAX ========================\n")
	s := syntax.FromGrammar(g1, v)
	s.Write(w)

	os.WriteFile("C:\\Users\\fabio\\temp\\teste.txt", []byte(w.String()), 0644)

	// i := input.NewStringInput("private function teste(a, b)\n {")
	//i, inputErr := input.NewFileInput("C:\\Users\\fabio\\dev\\github\\sdbrdd\\sdb-med-lib\\src\\main\\clipper\\sdb_api_med.prg")
	i, inputErr := input.NewFileInput("C:\\Users\\fabio\\temp\\teste.prg")
	if inputErr != nil {
		fmt.Print(inputErr)
		return
	}
	l := lexer.New(v, i)
	// d := automata.NFAToDFA(vocabulary.RulesToNFA(g.ParserRules()...))
	// fmt.Print(d.String())
	p := parser.New(l, s)
	ast := p.Execute()
	if ast != nil {
		printTree(ast, i, s, 0)
	}
	//printTokens(l, v)
}

func printTokens(l *lexer.Lexer, v *lexer.Vocabulary) {
	fmt.Print("==================== TOKENS ======================\n")
	// for _, token := range tokens {
	// 	fmt.Printf("Row: %d, Col: %d, Types: %v\n", token.Row(), token.Col(), tokensNames(v, token))
	// }
	// for _, lexError := range l.Errors() {
	// 	fmt.Printf("Error: %s\n", lexError)
	// }
	token, lexError := l.NextToken()
	for token != nil && !token.IsType(lexer.TKN_EOF) {
		if lexError == nil {
			fmt.Printf("Row: %d, Col: %d, Types: %v\n", token.Row(), token.Col(), tokensNames(v, token))
		} else {
			fmt.Printf("Error: %s\n", lexError)
		}
		token, lexError = l.NextToken()
	}
}

func printTree(node *ast.Node, in input.Input, s *parser.Syntax, i int) {
	printNode(node, in, s, i)
	child := node.FirstChild()
	for child != nil {
		printTree(child, in, s, i+1)
		child = child.Sibling()
	}
}

func printNode(node *ast.Node, in input.Input, s *parser.Syntax, i int) {
	for range i {
		fmt.Print("   ")
	}
	fmt.Print("[")
	fmt.Print(s.RuleName(node.RuleType()))
	fmt.Print("] : '")
	fmt.Print(formatText(in.GetText(node.Start(), node.End())))
	fmt.Println("'")
}

func formatText(text string) string {
	sb := strings.Builder{}
	for _, c := range text {
		switch c {
		case '\n':
			sb.WriteString("\\n")
		case '\r':
			sb.WriteString("\\r")
		case '\t':
			sb.WriteString("\\t")
		case '\f':
			sb.WriteString("\\f")
		default:
			sb.WriteRune(rune(c))
		}
	}
	return sb.String()
}

func tokensNames(v *lexer.Vocabulary, t *lexer.Token) string {
	str := strings.Builder{}
	for i, tokenType := range t.Types() {
		if i > 0 {
			str.WriteRune(',')
		}
		str.WriteString(v.TokenName(tokenType))
	}
	return str.String()
}
