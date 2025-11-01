// Package expression provides template-based expression evaluation for workflow conditions.
// It supports referencing node outputs, variables, context values, and complex boolean logic.
package expression

import "fmt"

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	TokenEOF TokenType = iota
	TokenIllegal

	// Literals
	TokenNumber  // 123, 45.67
	TokenString  // "hello", 'world'
	TokenTrue    // true
	TokenFalse   // false
	TokenNull    // null

	// Identifiers and references
	TokenIdent // identifier
	TokenDot   // .

	// Template delimiters
	TokenTemplateStart // {{
	TokenTemplateEnd   // }}

	// Operators
	TokenPlus     // +
	TokenMinus    // -
	TokenStar     // *
	TokenSlash    // /
	TokenPercent  // %
	TokenEq       // ==
	TokenNeq      // !=
	TokenLt       // <
	TokenLte      // <=
	TokenGt       // >
	TokenGte      // >=
	TokenAnd      // &&
	TokenOr       // ||
	TokenNot      // !
	TokenContains // contains (function-like operator)
	TokenMatches  // matches (regex)

	// Delimiters
	TokenLParen   // (
	TokenRParen   // )
	TokenLBracket // [
	TokenRBracket // ]
	TokenComma    // ,
)

var tokenNames = map[TokenType]string{
	TokenEOF:           "EOF",
	TokenIllegal:       "ILLEGAL",
	TokenNumber:        "NUMBER",
	TokenString:        "STRING",
	TokenTrue:          "TRUE",
	TokenFalse:         "FALSE",
	TokenNull:          "NULL",
	TokenIdent:         "IDENT",
	TokenDot:           "DOT",
	TokenTemplateStart: "{{",
	TokenTemplateEnd:   "}}",
	TokenPlus:          "+",
	TokenMinus:         "-",
	TokenStar:          "*",
	TokenSlash:         "/",
	TokenPercent:       "%",
	TokenEq:            "==",
	TokenNeq:           "!=",
	TokenLt:            "<",
	TokenLte:           "<=",
	TokenGt:            ">",
	TokenGte:           ">=",
	TokenAnd:           "&&",
	TokenOr:            "||",
	TokenNot:           "!",
	TokenContains:      "contains",
	TokenMatches:       "matches",
	TokenLParen:        "(",
	TokenRParen:        ")",
	TokenLBracket:      "[",
	TokenRBracket:      "]",
	TokenComma:         ",",
}

func (t TokenType) String() string {
	if name, ok := tokenNames[t]; ok {
		return name
	}
	return fmt.Sprintf("TokenType(%d)", t)
}

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q) at %d:%d", t.Type, t.Literal, t.Line, t.Column)
}

// Keywords maps keyword strings to their token types
var keywords = map[string]TokenType{
	"true":     TokenTrue,
	"false":    TokenFalse,
	"null":     TokenNull,
	"contains": TokenContains,
	"matches":  TokenMatches,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TokenIdent
}
