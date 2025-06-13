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
	EQ  TokenType = "EQ"  // =
	NE  TokenType = "NE"  // !=
	LT  TokenType = "LT"  // <
	GT  TokenType = "GT"  // >
	AND TokenType = "AND" // AND
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

func (ce *ComparisonExpression) Evaluate(item reflect.Value) (bool, error) {
	fieldValue, err := getFieldValue(item, ce.Field)
	if err != nil {
		return false, nil
	}

	// Handle slice values
	if fieldValue.Kind() == reflect.Slice {
		for i := 0; i < fieldValue.Len(); i++ {
			elem := fieldValue.Index(i)
			// If slice of pointers, dereference
			if elem.Kind() == reflect.Ptr && !elem.IsNil() {
				elem = elem.Elem()
			}
			// If comparing to a field of struct (e.g., Tags.Name)
			if elem.Kind() == reflect.Struct && strings.Contains(ce.Field, ".") {
				// Remove the slice field from ce.Field
				parts := strings.SplitN(ce.Field, ".", 2)
				if len(parts) == 2 {
					subField := parts[1]
					subExpr := &ComparisonExpression{Field: subField, Operator: ce.Operator, Value: ce.Value}
					match, _ := subExpr.Evaluate(elem)
					if match {
						return true, nil
					}
				}
			} else {
				// Compare element directly (e.g., Interests = 'music')
				cmpExpr := &ComparisonExpression{Operator: ce.Operator, Value: ce.Value}
				// Set Field to empty so getFieldValue returns elem
				match, _ := cmpExpr.compareValue(elem)
				if match {
					return true, nil
				}
			}
		}
		return false, nil
	}

	return ce.compareValue(fieldValue)
}

// compareValue handles the actual comparison for a single value
func (ce *ComparisonExpression) compareValue(fieldValue reflect.Value) (bool, error) {
	if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
		fieldValue = reflect.Zero(fieldValue.Type().Elem())
	}
	// if fieldValue.Kind() == reflect.Invalid || (fieldValue.Kind() != reflect.Ptr && fieldValue.IsZero()) {
	// }

	switch fieldValue.Kind() {
	case reflect.String:
		switch ce.Operator {
		case EQ:
			return fieldValue.Interface().(string) == ce.Value, nil
		case NE:
			return fieldValue.Interface().(string) != ce.Value, nil
		case LT:
			return fieldValue.Interface().(string) < ce.Value, nil
		case GT:
			return fieldValue.Interface().(string) > ce.Value, nil
		}
	case reflect.Bool:
		sb, err := strconv.ParseBool(ce.Value)
		if err != nil {
			return false, nil
		}
		switch ce.Operator {
		case EQ:
			return fieldValue.Interface().(bool) == sb, nil
		case NE:
			return fieldValue.Interface().(bool) != sb, nil
		}
	case reflect.Int:
		v, err := strconv.Atoi(ce.Value)
		if err != nil {
			return false, nil
		}
		fv, ok := fieldValue.Interface().(int)
		if !ok {
			return false, nil
		}
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
	case reflect.Int8:
		vx, err := strconv.ParseInt(ce.Value, 10, 8)
		if err != nil {
			return false, nil
		}
		v := int8(vx)
		fv, ok := fieldValue.Interface().(int8)
		if !ok {
			return false, nil
		}
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
	case reflect.Uint8:
		vx, err := strconv.ParseUint(ce.Value, 10, 8)
		if err != nil {
			return false, nil
		}
		v := uint8(vx)
		fv, ok := fieldValue.Interface().(uint8)
		if !ok {
			return false, nil
		}
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
	case reflect.Int32:
		vx, err := strconv.ParseInt(ce.Value, 10, 32)
		if err != nil {
			return false, nil
		}
		v := int32(vx)
		fv, ok := fieldValue.Interface().(int32)
		if !ok {
			return false, nil
		}
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
	case reflect.Uint32:
		xv, err := strconv.ParseUint(ce.Value, 10, 32)
		if err != nil {
			return false, nil
		}
		v := uint32(xv)
		fv, ok := fieldValue.Interface().(uint32)
		if !ok {
			return false, nil
		}
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
	case reflect.Int64:
		v, err := strconv.ParseInt(ce.Value, 10, 64)
		if err != nil {
			return false, nil
		}
		fv, ok := fieldValue.Interface().(int64)
		if !ok {
			return false, nil
		}
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
	case reflect.Uint64:
		v, err := strconv.ParseUint(ce.Value, 10, 64)
		if err != nil {
			return false, nil
		}
		fv, ok := fieldValue.Interface().(uint64)
		if !ok {
			return false, nil
		}
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
	case reflect.Float32:
		xv, err := strconv.ParseFloat(ce.Value, 32)
		if err != nil {
			return false, nil
		}
		v := float32(xv)
		fv, ok := fieldValue.Interface().(float32)
		if !ok {
			return false, nil
		}
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
	case reflect.Float64:
		v, err := strconv.ParseFloat(ce.Value, 64)
		if err != nil {
			return false, nil
		}
		fv, ok := fieldValue.Interface().(float64)
		if !ok {
			return false, nil
		}
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
	}

	return false, fmt.Errorf("unsupported operator %q", ce.Operator)
}

func (ce *ConjunctionExpression) Evaluate(item reflect.Value) (bool, error) {
	for _, expr := range ce.Expressions {
		result, err := expr.Evaluate(item)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}
	return true, nil
}

func getFieldValue(item reflect.Value, fieldPath string) (reflect.Value, error) {
	parts := strings.Split(fieldPath, ".")
	currentValue := item

	for i, part := range parts {
		// Dereference pointers
		if currentValue.Kind() == reflect.Ptr {
			if currentValue.IsNil() {
				return reflect.Value{}, fmt.Errorf("nil pointer at path %q", strings.Join(parts[:i+1], "."))
			}
			currentValue = currentValue.Elem()
		}

		if currentValue.Kind() == reflect.Slice {
			// If this is the last part, return the slice
			if i == len(parts)-1 {
				return currentValue, nil
			}
			// Otherwise, we want to access a field of each element in the slice
			return currentValue, nil // Let Evaluate handle the rest
		}

		if currentValue.Kind() != reflect.Struct {
			return reflect.Value{}, fmt.Errorf("field %q is not a struct, cannot access nested field %q", strings.Join(parts[:i], "."), part)
		}

		field := currentValue.FieldByName(part)
		if !field.IsValid() {
			return reflect.Value{}, fmt.Errorf("field %q not found in struct %s", part, currentValue.Type().Name())
		}
		currentValue = field
	}

	return currentValue, nil
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
	expressions := []Expression{}
	expr, err := p.parseComparison()
	if err != nil {
		return nil, err
	}
	expressions = append(expressions, expr)

	for p.peekTokenIs(AND) {
		p.nextToken() // Consume AND
		p.nextToken() // Move to the start of the next comparison

		expr, err = p.parseComparison()
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	if len(p.errors) > 0 {
		return nil, fmt.Errorf("parsing errors: %s", strings.Join(p.errors, "; "))
	}

	if len(expressions) == 1 {
		return expressions[0], nil
	}

	return &ConjunctionExpression{Expressions: expressions}, nil
}

func (p *Parser) parseComparison() (*ComparisonExpression, error) {
	expr := &ComparisonExpression{}

	if !p.currentTokenIs(IDENTIFIER) {
		return nil, fmt.Errorf("expected identifier, got %s (%q)", p.currentToken.Type, p.currentToken.Literal)
	}
	expr.Field = p.currentToken.Literal

	p.nextToken()

	switch p.currentToken.Type {
	case EQ, NE, LT, GT:
		expr.Operator = p.currentToken.Type
	default:
		return nil, fmt.Errorf("expected operator (=, !=, <, >), got %s (%q)", p.currentToken.Type, p.currentToken.Literal)
	}

	p.nextToken()
	expr.Value = p.currentToken.Literal

	return expr, nil
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
	case '\'':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0: // EOF
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdentifier(tok.Literal)
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
	default:
		return IDENTIFIER
	}
}
