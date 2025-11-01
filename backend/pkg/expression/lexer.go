package expression

import (
	"unicode"
)

// Lexer tokenizes expression strings
type Lexer struct {
	input   string
	pos     int  // current position
	readPos int  // next position
	ch      byte // current character
	line    int
	column  int
}

// NewLexer creates a new lexer for the given input
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar advances the lexer position
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
	l.column++
	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

// peekChar looks at the next character without advancing
func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case 0:
		tok.Type = TokenEOF
		tok.Literal = ""
	case '{':
		if l.peekChar() == '{' {
			ch := l.ch
			l.readChar()
			tok.Type = TokenTemplateStart
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = TokenIllegal
			tok.Literal = string(l.ch)
		}
	case '}':
		if l.peekChar() == '}' {
			ch := l.ch
			l.readChar()
			tok.Type = TokenTemplateEnd
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = TokenIllegal
			tok.Literal = string(l.ch)
		}
	case '.':
		tok.Type = TokenDot
		tok.Literal = string(l.ch)
	case '+':
		tok.Type = TokenPlus
		tok.Literal = string(l.ch)
	case '-':
		tok.Type = TokenMinus
		tok.Literal = string(l.ch)
	case '*':
		tok.Type = TokenStar
		tok.Literal = string(l.ch)
	case '/':
		tok.Type = TokenSlash
		tok.Literal = string(l.ch)
	case '%':
		tok.Type = TokenPercent
		tok.Literal = string(l.ch)
	case '(':
		tok.Type = TokenLParen
		tok.Literal = string(l.ch)
	case ')':
		tok.Type = TokenRParen
		tok.Literal = string(l.ch)
	case '[':
		tok.Type = TokenLBracket
		tok.Literal = string(l.ch)
	case ']':
		tok.Type = TokenRBracket
		tok.Literal = string(l.ch)
	case ',':
		tok.Type = TokenComma
		tok.Literal = string(l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = TokenNeq
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = TokenNot
			tok.Literal = string(l.ch)
		}
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = TokenEq
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = TokenIllegal
			tok.Literal = string(l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = TokenLte
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = TokenLt
			tok.Literal = string(l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = TokenGte
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = TokenGt
			tok.Literal = string(l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok.Type = TokenAnd
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = TokenIllegal
			tok.Literal = string(l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok.Type = TokenOr
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = TokenIllegal
			tok.Literal = string(l.ch)
		}
	case '"', '\'':
		tok.Type = TokenString
		tok.Literal = l.readString(l.ch)
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = TokenNumber
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok.Type = TokenIllegal
			tok.Literal = string(l.ch)
		}
	}

	l.readChar()
	return tok
}

// skipWhitespace skips whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier reads an identifier
func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

// readNumber reads a number (integer or float)
func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	// Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // skip '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[pos:l.pos]
}

// readString reads a string literal
func (l *Lexer) readString(quote byte) string {
	pos := l.pos + 1 // skip opening quote
	for {
		l.readChar()
		if l.ch == quote || l.ch == 0 {
			break
		}
		// Handle escape sequences
		if l.ch == '\\' {
			l.readChar()
		}
	}
	return l.input[pos:l.pos]
}

// isLetter checks if a character is a letter
func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}

// isDigit checks if a character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
