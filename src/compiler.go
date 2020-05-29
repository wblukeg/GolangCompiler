package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type tokenizer interface {
	tokenize(*[]byte) []token
}

type sourceTokenizer struct{}

type token struct {
	tokenType string
	value     string
}

type tokenDefinition struct {
	tokenType string
	regEx     string
}

var tokenTypes = [7]tokenDefinition{
	{tokenType: "def", regEx: `\bdef\b`},
	{tokenType: "end", regEx: `\bend\b`},
	{tokenType: "identifier", regEx: `\b[a-zA-Z]+\b`},
	{tokenType: "integer", regEx: `\b[0-9]+\b`},
	{tokenType: "oparen", regEx: `\(`},
	{tokenType: "cparen", regEx: `\)`},
	{tokenType: "comma", regEx: `,`},
}

type parser interface {
	parse() defNode
}

type tokenParser struct{}

type defNode struct {
	name     string
	argNames []string
	body     treeNode
}

type treeNode struct {
	nodeType    string
	integerNode integerNode
	callNode    callNode
	varRefNode  varRefNode
}

type integerNode struct {
	value int64
}

type callNode struct {
	name      string
	argExpers []treeNode
}

type varRefNode struct {
	value string
}

type generator interface {
	generate(node defNode) string
}

type codeGenerator struct{}

var tokens []token
var tree defNode

func main() {
	//Read Source
	data, err := ioutil.ReadFile("code.lg")
	check(err)

	//Convert to tokens
	var t tokenizer = sourceTokenizer{}
	////Check Tokens
	tokens = t.tokenize(&data)
	// for _, value := range tokens {
	// 	fmt.Printf("%v\n", value)
	// }

	//Parse Tokens
	var p parser = tokenParser{}
	tree = p.parse()
	//	fmt.Printf("Name: %v\nArg Names: %v\nBody: %v\n", tree.name, tree.argNames, tree.body)

	//Generate Code
	var g generator = codeGenerator{}
	c := g.generate(tree)

	fmt.Println("function add(x, y) { return x + y };")
	fmt.Println(c)
	fmt.Println("console.log(f(1,2));")

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (st sourceTokenizer) tokenize(code *[]byte) []token {

	var tokens []token
	for len(*code) > 0 {
		token, err := tokenizeOneToken(code)
		check(err)
		tokens = append(tokens, token)
	}
	return tokens
}

func tokenizeOneToken(code *[]byte) (token, error) {

	var foundReturn token
	for i := 0; i < len(tokenTypes); i++ {

		re := regexp.MustCompile(`\A(` + tokenTypes[i].regEx + `)`)
		found := re.Find([]byte(*code))
		if len(found) > 0 {
			*code = bytes.TrimSpace([]byte((*code)[len(found):]))
			foundReturn = token{tokenType: tokenTypes[i].tokenType, value: string(found)}
			break
		}

	}
	if foundReturn == (token{}) {
		return token{}, fmt.Errorf("Token Type: %v not defined", string(*code))

	}
	return foundReturn, nil
}

func (tp tokenParser) parse() defNode {

	return parseDef()
}

func parseDef() defNode {
	consume("def")
	name := consume("identifier")
	argNames := parseArgNames()
	body := parseExpr()
	return defNode{name: name.value, argNames: argNames, body: body}
}

func parseArgNames() []string {
	var args []string
	consume("oparen")
	if peek("identifier") {
		args = append(args, consume("identifier").value)
		for peek("comma") == true {
			consume("comma")
			args = append(args, consume("identifier").value)
		}
	}
	consume("cparen")
	return args
}

func peek(expectedType string, offset ...int) bool {
	os := 0
	if len(offset) > 0 {
		os = offset[0]
	}
	return tokens[os].tokenType == expectedType
}

func parseExpr() treeNode {
	if peek("integer") {
		return parseInteger()
	} else if peek("identifier") && peek("oparen", 1) {
		return parseCall()
	}

	return parseVarRef()

}

func parseInteger() treeNode {
	i, err := strconv.ParseInt(consume("integer").value, 10, 64)
	check(err)
	return treeNode{nodeType: "int", integerNode: integerNode{value: i}}
}

func parseCall() treeNode {
	callName := consume("identifier").value

	return treeNode{nodeType: "call", callNode: callNode{name: callName, argExpers: parseArgExpers()}}
}

func parseArgExpers() []treeNode {
	var argExpers []treeNode
	consume("oparen")
	if !peek("cparen") {
		argExpers = append(argExpers, parseExpr())
		for peek("comma") == true {
			consume("comma")
			argExpers = append(argExpers, parseExpr())
		}
	}
	//fmt.Println("CONSUME THIS PAREN")
	consume("cparen")
	return argExpers
}

func parseVarRef() treeNode {
	return treeNode{nodeType: "varRef", varRefNode: varRefNode{value: consume("identifier").value}}
}

func consume(expectedType string) token {
	token := (tokens)[0]
	(tokens) = (tokens)[1:]
	if token.tokenType == expectedType {
		return token
	}
	panic(fmt.Errorf("Expected token type %v but got %v", expectedType, token.tokenType))
}

func (gc codeGenerator) generate(node defNode) string {

	body := generateBody(node.body)
	code := fmt.Sprintf("function %v(%v) { return %v };", node.name, strings.Join(node.argNames, ","), body)
	return code
}

func generateBody(bodyNode treeNode) string {

	bodyCode := ""
	switch bodyNode.nodeType {
	case "call":
		var callArgs []string
		for _, value := range bodyNode.callNode.argExpers {
			callArgs = append(callArgs, generateBody(value))
		}
		bodyCode = fmt.Sprintf("%v(%v)", bodyNode.callNode.name, strings.Join(callArgs, ","))
	case "varRef":
		bodyCode = bodyNode.varRefNode.value
	case "int":

		bodyCode = fmt.Sprint(bodyNode.integerNode.value)
	default:
		panic(fmt.Errorf("Unknown Body Node Type %v", bodyNode.nodeType))
	}
	return bodyCode
}
