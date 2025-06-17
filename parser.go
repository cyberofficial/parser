package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type TokenType string

const (
	// Special tokens
	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"

	// Literals
	IDENTIFIER TokenType = "IDENTIFIER"
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"

	// Operators
	EQ       TokenType = "EQ"       // =
	NE       TokenType = "NE"       // !=
	LT       TokenType = "LT"       // <
	GT       TokenType = "GT"       // >
	AND      TokenType = "AND"      // AND
	OR       TokenType = "OR"       // OR
	CONTAINS TokenType = "CONTAINS" // CONTAINS
	LPAREN   TokenType = "LPAREN"   // (
	RPAREN   TokenType = "RPAREN"   // )
)

type Token struct {
	Type    TokenType
	Literal string
}

type Expression interface {
	Evaluate(item reflect.Value) (bool, error)
}

type ComparisonExpression struct {
	Field    string
	Operator TokenType
	Value    string
}

type ConjunctionExpression struct {
	Expressions []Expression
}

// OrExpression supports logical OR
type OrExpression struct {
	Expressions []Expression
}

func Parse[T any](query string, data []T) ([]T, error) {
	l := NewLexer(query)
	p := NewParser(l)

	ast, err := p.ParseQuery()
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}
	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("parsing errors: %s", strings.Join(p.Errors(), "; "))
	}

	filteredData := []T{}

	for _, item := range data {
		val := reflect.ValueOf(item)
		if val.Kind() == reflect.Ptr && val.IsNil() {
			continue
		}
		if val.Kind() == reflect.Ptr {
			val = val.Elem() // Dereference if it's a pointer to a struct
		}

		if val.Kind() != reflect.Struct {
			return nil, fmt.Errorf("expected slice of structs, got %s in data", val.Kind())
		}

		match, err := ast.Evaluate(val)
		if err != nil {
			fmt.Printf("Warning: Skipping item due to evaluation error: %v\n", err)
			continue
		}

		if match {
			filteredData = append(filteredData, item)
		}
	}

	return filteredData, nil
}

// Enhanced getFieldValue: returns a slice of reflect.Value if a slice is encountered in the path
func getFieldValues(item reflect.Value, fieldPath string) ([]reflect.Value, error) {
	parts := strings.Split(fieldPath, ".")
	currentValues := []reflect.Value{item}

	for _, part := range parts {
		nextValues := []reflect.Value{}
		for _, val := range currentValues {
			fmt.Printf("DEBUG: part=%q, val=%#v, kind=%v\n", part, val.Interface(), val.Kind())
			if val.Kind() == reflect.Ptr {
				if val.IsNil() {
					continue
				}
				val = val.Elem()
			}
			if val.Kind() == reflect.Slice {
				for j := 0; j < val.Len(); j++ {
					elem := val.Index(j)
					if elem.Kind() == reflect.Ptr {
						if elem.IsNil() {
							continue
						}
						elem = elem.Elem()
					}
					if elem.Kind() == reflect.Struct {
						field := elem.FieldByName(part)
						if field.IsValid() {
							nextValues = append(nextValues, field)
						}
					} else {
						nextValues = append(nextValues, elem)
					}
				}
				continue
			}
			if val.Kind() == reflect.Struct {
				field := val.FieldByName(part)
				if !field.IsValid() {
					continue
				}
				nextValues = append(nextValues, field)
				continue
			}
			// For non-struct, non-slice, just append (should only happen at leaf)
			nextValues = append(nextValues, val)
		}
		currentValues = nextValues
		if len(currentValues) == 0 {
			return nil, fmt.Errorf("field %q not found in path %q", part, fieldPath)
		}
	}
	// Flatten any slices at the leaf
	flat := []reflect.Value{}
	for _, v := range currentValues {
		if v.Kind() == reflect.Slice {
			for i := 0; i < v.Len(); i++ {
				flat = append(flat, v.Index(i))
			}
		} else {
			flat = append(flat, v)
		}
	}
	fmt.Printf("DEBUG getFieldValues(%q): %v\n", fieldPath, flat)
	return flat, nil
}

// The core Evaluate method for ComparisonExpression
func (ce *ComparisonExpression) Evaluate(item reflect.Value) (bool, error) {
	fieldValues, err := getFieldValues(item, ce.Field)
	if err != nil || len(fieldValues) == 0 {
		fmt.Printf("DEBUG Evaluate: field %q not found or empty\n", ce.Field)
		return false, nil // If field is missing or not found, do not match
	}
	for _, fieldValue := range fieldValues {
		match, _ := ce.compareValue(fieldValue)
		fmt.Printf("DEBUG Compare: %v %v %q => %v\n", fieldValue, ce.Operator, ce.Value, match)
		if match {
			return true, nil
		}
	}
	return false, nil
}

// compareValue handles the actual comparison for a single value
func (ce *ComparisonExpression) compareValue(fieldValue reflect.Value) (bool, error) {
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			return false, nil
		}
		fieldValue = fieldValue.Elem()
	}
	switch fieldValue.Kind() {
	case reflect.String:
		s := fieldValue.Interface().(string)
		switch ce.Operator {
		case EQ:
			return s == ce.Value, nil
		case NE:
			return s != ce.Value, nil
		case LT:
			return s < ce.Value, nil
		case GT:
			return s > ce.Value, nil
		case CONTAINS:
			return strings.Contains(s, ce.Value), nil
		}
	case reflect.Bool:
		b, _ := strconv.ParseBool(ce.Value)
		switch ce.Operator {
		case EQ:
			return fieldValue.Interface().(bool) == b, nil
		case NE:
			return fieldValue.Interface().(bool) != b, nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, _ := strconv.ParseInt(ce.Value, 10, 64)
		fv := fieldValue.Int()
		switch ce.Operator {
		case EQ:
			return fv == v, nil
		case NE:
			return fv != v, nil
		case LT:
			return fv < v, nil
		case GT:
			return fv > v, nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, _ := strconv.ParseUint(ce.Value, 10, 64)
		fv := fieldValue.Uint()
		switch ce.Operator {
		case EQ:
			return fv == v, nil
		case NE:
			return fv != v, nil
		case LT:
			return fv < v, nil
		case GT:
			return fv > v, nil
		}
	case reflect.Float32, reflect.Float64:
		v, _ := strconv.ParseFloat(ce.Value, 64)
		fv := fieldValue.Float()
		switch ce.Operator {
		case EQ:
			return fv == v, nil
		case NE:
			return fv != v, nil
		case LT:
			return fv < v, nil
		case GT:
			return fv > v, nil
		}
	case reflect.Slice:
		if ce.Operator == CONTAINS {
			for i := 0; i < fieldValue.Len(); i++ {
				item := fieldValue.Index(i)
				if item.Kind() == reflect.Ptr && !item.IsNil() {
					item = item.Elem()
				}
				if item.Kind() == reflect.String {
					if item.String() == ce.Value || strings.Contains(item.String(), ce.Value) {
						return true, nil
					}
				} else if item.Kind() == reflect.Interface {
					if s, ok := item.Interface().(string); ok && (s == ce.Value || strings.Contains(s, ce.Value)) {
						return true, nil
					}
				}
			}
			return false, nil
		}
	}
	return false, nil
}

type Parser struct {
	l *Lexer

	currentToken Token
	peekToken    Token
	errors       []string
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) currentTokenIs(t TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) ParseQuery() (Expression, error) {
	expr := p.parseOrExpression()
	// Skip any trailing RPAREN tokens after parsing the main expression
	for p.currentToken.Type == RPAREN {
		p.nextToken()
	}
	if p.currentToken.Type != EOF {
		p.errors = append(p.errors, "unexpected token after end of query")
	}
	return expr, nil
}

// parseOrExpression handles OR precedence
func (p *Parser) parseOrExpression() Expression {
	expr := p.parseAndExpression()
	for p.currentTokenIs(OR) {
		p.nextToken() // move to right expr
		right := p.parseAndExpression()
		if orExpr, ok := expr.(*OrExpression); ok {
			orExpr.Expressions = append(orExpr.Expressions, right)
			expr = orExpr
		} else {
			expr = &OrExpression{Expressions: []Expression{expr, right}}
		}
		if !p.currentTokenIs(OR) {
			break
		}
	}
	return expr
}

// parseAndExpression handles AND precedence
func (p *Parser) parseAndExpression() Expression {
	expr := p.parsePrimary()
	for p.currentTokenIs(AND) {
		p.nextToken() // move to right expr
		right := p.parsePrimary()
		if andExpr, ok := expr.(*ConjunctionExpression); ok {
			andExpr.Expressions = append(andExpr.Expressions, right)
			expr = andExpr
		} else {
			expr = &ConjunctionExpression{Expressions: []Expression{expr, right}}
		}
		if !p.currentTokenIs(AND) {
			break
		}
	}
	return expr
}

func (p *Parser) parsePrimary() Expression {
	if p.currentTokenIs(LPAREN) {
		p.nextToken()                 // move to first expr inside parens
		expr := p.parseOrExpression() // Use parseOrExpression for full precedence inside parens
		if !p.currentTokenIs(RPAREN) && !p.currentTokenIs(EOF) {
			fmt.Printf("DEBUG: parsePrimary expected RPAREN or EOF, got %s (%q)\n", p.currentToken.Type, p.currentToken.Literal)
			p.errors = append(p.errors, "expected )")
			return expr
		}
		if p.currentTokenIs(RPAREN) {
			p.nextToken() // Advance past RPAREN so parent sees next token
		}
		return expr
	}
	// Otherwise, parse a comparison
	expr, err := p.parseComparison()
	if err != nil {
		p.errors = append(p.errors, err.Error())
		return nil
	}
	return expr
}

// parseComparison is now only used by parsePrimary
func (p *Parser) parseComparison() (*ComparisonExpression, error) {
	expr := &ComparisonExpression{}

	if !p.currentTokenIs(IDENTIFIER) {
		return nil, fmt.Errorf("expected identifier, got %s (%q)", p.currentToken.Type, p.currentToken.Literal)
	}
	expr.Field = p.currentToken.Literal

	p.nextToken()

	switch p.currentToken.Type {
	case EQ, NE, LT, GT, CONTAINS:
		expr.Operator = p.currentToken.Type
	default:
		return nil, fmt.Errorf("expected operator (=, !=, <, >, CONTAINS), got %s (%q)", p.currentToken.Type, p.currentToken.Literal)
	}

	p.nextToken()
	expr.Value = p.currentToken.Literal
	p.nextToken() // Always advance after reading the value

	return expr, nil
}

// parseExpression is now an alias for parseOrExpression for compatibility
func (p *Parser) parseExpression() Expression {
	return p.parseOrExpression()
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(EQ, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: NE, Literal: literal}
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	case '<':
		tok = newToken(LT, l.ch)
	case '>':
		tok = newToken(GT, l.ch)
	case '(': // Add LPAREN
		tok = newToken(LPAREN, l.ch)
	case ')': // Add RPAREN
		tok = newToken(RPAREN, l.ch)
	case '\'':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0: // EOF
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			typeFromIdent := LookupIdentifier(tok.Literal)
			tok.Type = typeFromIdent
			return tok
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '.' || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '\'' || l.ch == 0 {
			break
		}
	}
	literal := l.input[position:l.position]
	l.readChar()
	return literal
}

// newToken creates a new Token.
func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func LookupIdentifier(identifier string) TokenType {
	switch strings.ToUpper(identifier) { // Case-insensitive for keywords
	case "AND":
		return AND
	case "OR":
		return OR
	case "CONTAINS":
		return CONTAINS
	default:
		return IDENTIFIER
	}
}

// Evaluate for ConjunctionExpression
func (ce *ConjunctionExpression) Evaluate(item reflect.Value) (bool, error) {
	if len(ce.Expressions) == 1 {
		return ce.Expressions[0].Evaluate(item)
	}
	// Check if all conditions are on the same field (including nested fields)
	allCmp := true
	var field string
	for _, expr := range ce.Expressions {
		cmp, ok := expr.(*ComparisonExpression)
		if !ok {
			allCmp = false
			break
		}
		if field == "" {
			field = cmp.Field
		} else if cmp.Field != field {
			allCmp = false
			break
		}
	}
	if allCmp {
		fieldValues, err := getFieldValues(item, field)
		if err != nil || len(fieldValues) == 0 {
			return false, nil
		}
		if fieldValues[0].Kind() == reflect.Slice {
			for i := 0; i < fieldValues[0].Len(); i++ {
				elem := fieldValues[0].Index(i)
				allTrue := true
				for _, expr := range ce.Expressions {
					cmp := expr.(*ComparisonExpression)
					if match, _ := cmp.compareValue(elem); !match {
						allTrue = false
						break
					}
				}
				if allTrue {
					return true, nil
				}
			}
			return false, nil
		} else {
			val := fieldValues[0]
			allTrue := true
			for _, expr := range ce.Expressions {
				cmp := expr.(*ComparisonExpression)
				if match, _ := cmp.compareValue(val); !match {
					allTrue = false
					break
				}
			}
			return allTrue, nil
		}
	}
	// Fallback: for AND over different fields, all must be true for the same item
	for _, expr := range ce.Expressions {
		match, err := expr.Evaluate(item)
		if err != nil || !match {
			return false, nil
		}
	}
	return true, nil
}

// Evaluate for OrExpression
func (oe *OrExpression) Evaluate(item reflect.Value) (bool, error) {
	if len(oe.Expressions) == 1 {
		return oe.Expressions[0].Evaluate(item)
	}
	// Check if all conditions are on the same field (including nested fields)
	allCmp := true
	var field string
	for _, expr := range oe.Expressions {
		cmp, ok := expr.(*ComparisonExpression)
		if !ok {
			allCmp = false
			break
		}
		if field == "" {
			field = cmp.Field
		} else if cmp.Field != field {
			allCmp = false
			break
		}
	}
	if allCmp {
		fieldValues, err := getFieldValues(item, field)
		if err != nil || len(fieldValues) == 0 {
			return false, nil
		}
		if fieldValues[0].Kind() == reflect.Slice {
			for i := 0; i < fieldValues[0].Len(); i++ {
				elem := fieldValues[0].Index(i)
				for _, expr := range oe.Expressions {
					cmp := expr.(*ComparisonExpression)
					if match, _ := cmp.compareValue(elem); match {
						return true, nil
					}
				}
			}
			return false, nil
		} else {
			val := fieldValues[0]
			for _, expr := range oe.Expressions {
				cmp := expr.(*ComparisonExpression)
				if match, _ := cmp.compareValue(val); match {
					return true, nil
				}
			}
			return false, nil
		}
	}
	// Fallback: for OR over different fields, return true if any condition is true for the same item
	for _, expr := range oe.Expressions {
		match, err := expr.Evaluate(item)
		if err == nil && match {
			return true, nil
		}
	}
	return false, nil
}
