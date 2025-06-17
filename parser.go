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
	GE       TokenType = "GE"       // >=
	LE       TokenType = "LE"       // <=
	AND      TokenType = "AND"      // AND
	OR       TokenType = "OR"       // OR
	CONTAINS TokenType = "CONTAINS" // CONTAINS
	LPAREN   TokenType = "LPAREN"   // (
	RPAREN   TokenType = "RPAREN"   // )
	IS       TokenType = "IS"       // IS
	NULL     TokenType = "NULL"     // NULL
	NOT      TokenType = "NOT"      // NOT
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
	// Use the enhanced lexer that supports negative numbers
	l := NewEnhancedLexer(query)
	p := NewParser(l)

	ast, err := p.ParseQuery()
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}
	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("parsing errors: %s", strings.Join(p.Errors(), "; "))
	}
	if ast == nil {
		return nil, fmt.Errorf("failed to parse query: AST is nil")
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
						field := getFieldByNameCaseInsensitive(elem, part)
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
				field := getFieldByNameCaseInsensitive(val, part)
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
	return flat, nil
}

// getFieldByNameCaseInsensitive returns the struct field with a name matching 'name' (case-insensitive), or an invalid reflect.Value if not found
func getFieldByNameCaseInsensitive(val reflect.Value, name string) reflect.Value {
	typeOfVal := val.Type()
	for i := 0; i < typeOfVal.NumField(); i++ {
		field := typeOfVal.Field(i)
		if strings.EqualFold(field.Name, name) {
			return val.Field(i)
		}
	}
	return reflect.Value{}
}

// The core Evaluate method for ComparisonExpression
func (ce *ComparisonExpression) Evaluate(item reflect.Value) (bool, error) {
	fieldValues, err := getFieldValues(item, ce.Field)
	if err != nil || len(fieldValues) == 0 {
		return false, nil // If field is missing or not found, do not match
	}
	for _, fieldValue := range fieldValues {
		match, _ := ce.compareValue(fieldValue)
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
		case LE:
			return s <= ce.Value, nil
		case GE:
			return s >= ce.Value, nil
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
		case LE:
			return fv <= v, nil
		case GE:
			return fv >= v, nil
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
		case LE:
			return fv <= v, nil
		case GE:
			return fv >= v, nil
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
		case LE:
			return fv <= v, nil
		case GE:
			return fv >= v, nil
		}
	case reflect.Slice:
		if fieldValue.IsNil() {
			return false, nil
		}
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
		} else if ce.Operator == EQ {
			for i := 0; i < fieldValue.Len(); i++ {
				item := fieldValue.Index(i)
				if item.Kind() == reflect.Ptr && !item.IsNil() {
					item = item.Elem()
				}
				if item.Kind() == reflect.String {
					if item.String() == ce.Value {
						return true, nil
					}
				} else if item.Kind() == reflect.Interface {
					if s, ok := item.Interface().(string); ok && s == ce.Value {
						return true, nil
					}
				}
			}
			return false, nil
		} else if ce.Operator == NE {
			for i := 0; i < fieldValue.Len(); i++ {
				item := fieldValue.Index(i)
				if item.Kind() == reflect.Ptr && !item.IsNil() {
					item = item.Elem()
				}
				if item.Kind() == reflect.String {
					if item.String() == ce.Value {
						return false, nil
					}
				} else if item.Kind() == reflect.Interface {
					if s, ok := item.Interface().(string); ok && s == ce.Value {
						return false, nil
					}
				}
			}
			return true, nil
		}
	}
	return false, nil
}

// LexerInterface defines the interface for any lexer that can be used with the parser
type LexerInterface interface {
	NextToken() Token
}

type Parser struct {
	l LexerInterface

	currentToken Token
	peekToken    Token
	errors       []string
}

func NewParser(l LexerInterface) *Parser {
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

func (p *Parser) ParseQuery() (Expression, error) {
	// Handle empty query
	if p.currentToken.Type == EOF {
		return nil, nil
	}

	// Skip leading AND/OR tokens for user-friendly SQL-like queries
	for p.currentToken.Type == AND || p.currentToken.Type == OR {
		p.nextToken()
	}
	expr := p.parseOrExpression()

	// Check for unexpected trailing RPAREN tokens after parsing the main expression
	if p.currentToken.Type == RPAREN {
		p.errors = append(p.errors, "unbalanced parenthesis: unexpected closing )")
		// Skip any trailing RPAREN tokens
		for p.currentToken.Type == RPAREN {
			p.nextToken()
		}
	}

	if p.currentToken.Type != EOF && len(p.errors) == 0 {
		p.errors = append(p.errors, "unexpected token after end of query")
	}

	if len(p.errors) > 0 {
		return nil, fmt.Errorf("%s", strings.Join(p.errors, "; "))
	}
	return expr, nil
}

// parseOrExpression handles OR precedence
func (p *Parser) parseOrExpression() Expression {
	expr := p.parseAndExpression()
	for p.currentTokenIs(OR) {
		p.nextToken() // move to right expr
		right := p.parseAndExpression()
		if right == nil {
			// If there's an error in the right side, stop parsing this expression
			p.errors = append(p.errors, "invalid expression after OR")
			return expr
		}
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
	if expr == nil {
		return nil
	}
	for p.currentTokenIs(AND) {
		p.nextToken() // move to right expr
		right := p.parsePrimary()
		if right == nil {
			// If there's an error in the right side, stop parsing this expression
			p.errors = append(p.errors, "invalid expression after AND")
			return expr
		}
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
		// We're starting a parenthesized expression
		p.nextToken()

		// Handle empty parentheses
		if p.currentTokenIs(RPAREN) {
			p.nextToken()
			// Return a special "always false" expression
			return &ConjunctionExpression{Expressions: []Expression{}}
		}

		expr := p.parseOrExpression() // Use parseOrExpression for full precedence inside parens

		if expr == nil {
			p.errors = append(p.errors, "invalid expression inside parentheses")
			// Skip to matching parenthesis or EOF
			for !p.currentTokenIs(EOF) && !p.currentTokenIs(RPAREN) {
				p.nextToken()
			}
			if p.currentTokenIs(RPAREN) {
				p.nextToken()
			}
			return nil
		}

		if p.currentTokenIs(RPAREN) {
			p.nextToken()
		} else { // Handle special case where we might have hit EOF after string
			// We can no longer check the last character of the input
			if p.currentToken.Type == EOF {
				// Just assume there's a missing closing parenthesis
				p.errors = append(p.errors, "unbalanced parenthesis: missing closing parenthesis at end of input")
			}

			p.errors = append(p.errors, "unbalanced parenthesis: missing closing )")
			return nil // Return nil to prevent cascading errors
		}
		return expr
	}

	if p.currentTokenIs(IDENTIFIER) {
		field := p.currentToken.Literal
		p.nextToken()

		// Handle IS NULL / IS NOT NULL
		if p.currentTokenIs(IS) {
			p.nextToken()
			not := false
			if p.currentTokenIs(NOT) {
				not = true
				p.nextToken()
			}
			if p.currentTokenIs(NULL) {
				p.nextToken()
				return &IsNullExpression{Field: field, Not: not}
			} else {
				p.errors = append(p.errors, "expected NULL after IS")
				return nil
			}
		}

		expr, err := p.parseComparisonWithField(field)
		if err != nil {
			p.errors = append(p.errors, err.Error())
			return nil
		}
		return expr
	}

	// If we get here, it's an unexpected token
	if !p.currentTokenIs(EOF) {
		p.errors = append(p.errors, fmt.Sprintf("unexpected token: %s", p.currentToken.Literal))
		p.nextToken() // Skip over this token to try to continue parsing
	}
	return nil
}

func (p *Parser) parseComparisonWithField(field string) (*ComparisonExpression, error) {
	expr := &ComparisonExpression{Field: field}

	switch p.currentToken.Type {
	case EQ, NE, LT, GT, GE, LE, CONTAINS:
		expr.Operator = p.currentToken.Type
	default:
		return nil, fmt.Errorf("expected operator (=, !=, <, >, <=, >=, CONTAINS), got %s (%q)", p.currentToken.Type, p.currentToken.Literal)
	}

	p.nextToken()
	expr.Value = p.currentToken.Literal
	p.nextToken() // Always advance after reading the value

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
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: LE, Literal: literal}
		} else {
			tok = newToken(LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: GE, Literal: literal}
		} else {
			tok = newToken(GT, l.ch)
		}
	case '(': // Add LPAREN
		tok = newToken(LPAREN, l.ch)
	case ')': // Add RPAREN
		tok = newToken(RPAREN, l.ch)
	case '\'':
		tok.Type = STRING
		tok.Literal = l.readString()
		return tok
	case '-':
		// Check if it's a negative number
		if isDigit(l.peekChar()) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
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
	// Handle negative numbers
	if l.ch == '-' {
		l.readChar()
	}
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
		// Handle escape sequences
		if l.ch == '\\' && l.peekChar() == '\'' {
			l.readChar() // Skip the backslash and include the quote
		}
	}
	literal := l.input[position:l.position]
	if l.ch != 0 { // Only advance if we're not at EOF
		l.readChar()
	}
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
	case "IS":
		return IS
	case "NULL":
		return NULL
	case "NOT":
		return NOT
	default:
		return IDENTIFIER
	}
}

// Evaluate for ConjunctionExpression
func (ce *ConjunctionExpression) Evaluate(item reflect.Value) (bool, error) {
	if len(ce.Expressions) == 0 {
		return false, nil // Empty conjunction is always false (for empty parentheses case)
	}

	if len(ce.Expressions) == 1 {
		return ce.Expressions[0].Evaluate(item)
	}

	// Check if all conditions are on the same field (including nested fields)
	allCmp := true
	var field string
	for _, expr := range ce.Expressions {
		if expr == nil {
			allCmp = false
			break
		}
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

	// Special case for AND conditions on the same field
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
			// For scalar values, check all conditions against each value
			for _, val := range fieldValues {
				allTrue := true
				for _, expr := range ce.Expressions {
					cmp := expr.(*ComparisonExpression)
					if match, _ := cmp.compareValue(val); !match {
						allTrue = false
						break
					}
				}
				if !allTrue {
					return false, nil
				}
			}
			return true, nil
		}
	}

	// Fallback: for AND over different fields, all must be true for the same item
	for _, expr := range ce.Expressions {
		if expr == nil {
			return false, nil
		}
		match, err := expr.Evaluate(item)
		if err != nil || !match {
			return false, nil
		}
	}
	return true, nil
}

// Evaluate for OrExpression
func (oe *OrExpression) Evaluate(item reflect.Value) (bool, error) {
	for _, expr := range oe.Expressions {
		match, err := expr.Evaluate(item)
		if err == nil && match {
			return true, nil
		}
	}
	return false, nil
}

// Add IsNullExpression type

type IsNullExpression struct {
	Field string
	Not   bool
}

func (e *IsNullExpression) Evaluate(item reflect.Value) (bool, error) {
	fieldValues, err := getFieldValues(item, e.Field)
	if err != nil || len(fieldValues) == 0 {
		return !e.Not, nil // IS NULL: true if not found; IS NOT NULL: false if not found
	}
	for _, v := range fieldValues {
		if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			if v.IsNil() {
				return !e.Not, nil
			}
		} else if v.Kind() == reflect.Slice {
			if v.IsNil() || v.Len() == 0 {
				return !e.Not, nil
			}
		} else if !v.IsValid() {
			return !e.Not, nil
		} else if v.IsZero() {
			return !e.Not, nil
		}
	}
	return e.Not, nil // IS NOT NULL: true if found and not nil/zero
}
