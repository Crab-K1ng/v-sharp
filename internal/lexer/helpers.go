package lexer

import (
	"fmt"
	"unicode"
)

func (l *Lexer) advanceRune() rune {
	if l.position >= len(l.source)-1 {
		l.position = len(l.source)
		return 0
	}
	l.position++
	if l.source[l.position] == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return l.source[l.position]
}

func (l *Lexer) peekRuneAt(offset int) rune {
	pos := l.position + offset
	if pos >= len(l.source) {
		return 0
	}
	return l.source[pos]
}

func (l *Lexer) consumeWhitespace() {
	for {
		ch := l.peekRuneAt(0)
		if ch == 0 {
			return
		}
		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			l.advanceRune()
			continue
		}
		break
	}
}

func (l *Lexer) makeToken(t TokenType, lexeme string, line int, column int) Token {
	return Token{
		Type:   t,
		Lexeme: lexeme,
		Line:   line,
		Column: column,
		File:   l.file,
		Source: l.source,
	}
}

func lookupIdentifier(lexeme string) TokenType {
	if k, ok := keywordsMap[lexeme]; ok {
		return k
	}
	return Identifier
}

func (l *Lexer) scanIdentifier() string {
	startPos := l.position
	for {
		ch := l.peekRuneAt(0)
		if ch == 0 || !(unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '\'') {
			break
		}
		l.advanceRune()
	}
	return string(l.source[startPos:l.position])
}

func (l *Lexer) scanNumber() (string, TokenType) {
	start := l.position
	isFloat := false
	isUnsigned := false
	for {
		ch := l.peekRuneAt(0)
		if ch == '.' {
			if isFloat {
				l.error("invalid number format: multiple decimal points")
			}
			isFloat = true
			l.advanceRune()
			continue
		}
		if unicode.IsDigit(ch) {
			l.advanceRune()
			continue
		}
		break
	}
	if l.peekRuneAt(0) == 'u' {
		isUnsigned = true
		l.advanceRune()
	}
	if isUnsigned {
		return string(l.source[start:l.position]), Unsigned
	}
	if isFloat {
		return string(l.source[start:l.position]), Float
	}
	return string(l.source[start:l.position]), Integer
}

func (l *Lexer) scanStringLiteral() string {
	start := l.position
	if l.peekRuneAt(0) != '"' {
		l.advanceRune()
		return string(l.source[start:l.position])
	}
	l.advanceRune()
	for {
		ch := l.peekRuneAt(0)
		if ch == 0 {
			l.error("unterminated string literal")
			return string(l.source[start:l.position])
		}
		if ch == '\\' {
			l.advanceRune()
			esc := l.peekRuneAt(0)
			if esc == 0 {
				l.error("unterminated escape sequence in string literal")
				return string(l.source[start:l.position])
			}
			switch esc {
			case 'n', 't', 'r', '\\', '\'', '"', '0':
			default:
				l.error(fmt.Sprintf("invalid escape character: \\%c", esc))
				return string(l.source[start:l.position])
			}
			l.advanceRune()
			continue
		}
		if ch == '"' {
			l.advanceRune()
			break
		}
		l.advanceRune()
	}
	return string(l.source[start:l.position])
}

func (l *Lexer) scanByteLiteral() string {
	start := l.position
	if l.peekRuneAt(0) != '\'' {
		l.advanceRune()
		return string(l.source[start:l.position])
	}
	l.advanceRune()
	ch := l.peekRuneAt(0)
	if ch == 0 {
		l.error("unterminated character literal")
		return string(l.source[start:l.position])
	}
	if ch == '\\' {
		l.advanceRune()
		esc := l.peekRuneAt(0)
		if esc == 0 {
			l.error("unterminated escape sequence in character literal")
			return string(l.source[start:l.position])
		}
		switch esc {
		case 'n', 't', 'r', '\\', '\'', '"', '0':
		default:
			l.error(fmt.Sprintf("invalid escape character: \\%c", esc))
			return string(l.source[start:l.position])
		}
		l.advanceRune()
	} else {
		l.advanceRune()
	}
	if l.peekRuneAt(0) == 0 {
		l.error("unterminated character literal")
		return string(l.source[start:l.position])
	}
	if l.peekRuneAt(0) != '\'' {
		l.error("extra characters in character literal (expected closing ')")
		return string(l.source[start:l.position])
	}
	l.advanceRune()
	return string(l.source[start:l.position])
}

func (l *Lexer) scanLineComment() string {
	start := l.position
	l.advanceRune()
	l.advanceRune()
	for {
		ch := l.peekRuneAt(0)
		if ch == 0 || ch == '\n' {
			break
		}
		l.advanceRune()
	}
	end := l.position
	for end > start {
		last := l.source[end-1]
		if last == ' ' || last == '\t' || last == '\r' {
			end--
			continue
		}
		break
	}
	return string(l.source[start:end])
}
