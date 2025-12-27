package lexer

import "unicode"

type Lexer struct {
	source   []rune
	position int
	line     int
	column   int
	file     string
}

func Tokenize(source string, file string) ([]Token, error) {
	lexer := New(source, file)
	out := []Token{}
	for {
		tok := lexer.Next()
		out = append(out, tok)
		if tok.Type == EOF {
			break
		}
		if tok.Type == Illegal {
			lexer.error("Illegal token")
		}
	}
	return out, nil
}

func New(source string, file string) *Lexer {
	return &Lexer{
		source:   []rune(source),
		position: 0,
		line:     1,
		column:   1,
		file:     file,
	}
}

func (l *Lexer) Next() Token {
	l.consumeWhitespace()
	line, column := l.line, l.column
	ch := l.peekRuneAt(0)

	if ch == 0 {
		return l.makeToken(EOF, "", line, column)
	}

	if unicode.IsLetter(ch) || ch == '_' {
		lex := l.scanIdentifier()
		tokType := lookupIdentifier(lex)
		return l.makeToken(tokType, lex, line, column)
	}

	if unicode.IsDigit(ch) {
		lex, tokType := l.scanNumber()
		return l.makeToken(tokType, lex, line, column)
	}

	if ch == '"' {
		lex := l.scanStringLiteral()
		return l.makeToken(String, lex, line, column)
	}

	if ch == '\'' {
		lex := l.scanByteLiteral()
		return l.makeToken(Byte, lex, line, column)
	}

	if ch == '/' && l.peekRuneAt(1) == '/' {
		lex := l.scanLineComment()
		return l.makeToken(Comment, lex, line, column)
	}

	switch ch {
	case '=':
		l.advanceRune()
		if l.peekRuneAt(0) == '=' {
			l.advanceRune()
			return l.makeToken(Equal, "==", line, column)
		}
		return l.makeToken(Assign, "=", line, column)
	case '!':
		l.advanceRune()
		if l.peekRuneAt(0) == '=' {
			l.advanceRune()
			return l.makeToken(NotEqual, "!=", line, column)
		}
		return l.makeToken(Not, "!", line, column)
	case '<':
		l.advanceRune()
		if l.peekRuneAt(0) == '=' {
			l.advanceRune()
			return l.makeToken(LessEqual, "<=", line, column)
		}
		return l.makeToken(LessThan, "<", line, column)
	case '>':
		l.advanceRune()
		if l.peekRuneAt(0) == '=' {
			l.advanceRune()
			return l.makeToken(GreaterEqual, ">=", line, column)
		}
		return l.makeToken(GreaterThan, ">", line, column)
	case '&':
		l.advanceRune()
		if l.peekRuneAt(0) == '&' {
			l.advanceRune()
			return l.makeToken(And, "&&", line, column)
		}
		return l.makeToken(Illegal, "&", line, column)
	case '|':
		l.advanceRune()
		if l.peekRuneAt(0) == '|' {
			l.advanceRune()
			return l.makeToken(Or, "||", line, column)
		}
		return l.makeToken(Illegal, "|", line, column)
	case '+':
		l.advanceRune()
		return l.makeToken(Plus, "+", line, column)
	case '-':
		l.advanceRune()
		return l.makeToken(Minus, "-", line, column)
	case '*':
		l.advanceRune()
		return l.makeToken(Asterisk, "*", line, column)
	case '/':
		l.advanceRune()
		return l.makeToken(Slash, "/", line, column)
	case '%':
		l.advanceRune()
		return l.makeToken(Percent, "%", line, column)
	case '(':
		l.advanceRune()
		return l.makeToken(LeftParen, "(", line, column)
	case ')':
		l.advanceRune()
		return l.makeToken(RightParen, ")", line, column)
	case '{':
		l.advanceRune()
		return l.makeToken(LeftBrace, "{", line, column)
	case '}':
		l.advanceRune()
		return l.makeToken(RightBrace, "}", line, column)
	case '[':
		l.advanceRune()
		return l.makeToken(LeftBracket, "[", line, column)
	case ']':
		l.advanceRune()
		return l.makeToken(RightBracket, "]", line, column)
	case ',':
		l.advanceRune()
		return l.makeToken(Comma, ",", line, column)
	case ';':
		l.advanceRune()
		return l.makeToken(Semicolon, ";", line, column)
	case ':':
		l.advanceRune()
		return l.makeToken(Colon, ":", line, column)
	case '.':
		l.advanceRune()
		return l.makeToken(Dot, ".", line, column)
	default:
		l.advanceRune()
		return l.makeToken(Illegal, string(ch), line, column)
	}
}
