package lexer

import (
	"testing"
)

func TestLexerBasicToken(t *testing.T) {
	input := `var x: int16; 
	const y: float32;
	// This is a comment	
	// Another comment
	private add(int32 a, int32 b) int32 {
		return a + b;
	}
		
	class MyClass {
		private field1: int32;
		public field2: float64;
	public get() int32 {
			return field1;
		}
	}
	"Hello, World!"
	'b'`

	tests := []struct {
		Type   TokenType
		Lexeme string
	}{
		{KwVar, "var"},
		{Identifier, "x"},
		{Colon, ":"},
		{KwInt16, "int16"},
		{Semicolon, ";"},
		{KwConst, "const"},
		{Identifier, "y"},
		{Colon, ":"},
		{KwFloat32, "float32"},
		{Semicolon, ";"},
		{Comment, "// This is a comment"},
		{Comment, "// Another comment"},
		{KwPrivate, "private"},
		{Identifier, "add"},
		{LeftParen, "("},
		{KwInt32, "int32"},
		{Identifier, "a"},
		{Comma, ","},
		{KwInt32, "int32"},
		{Identifier, "b"},
		{RightParen, ")"},
		{KwInt32, "int32"},
		{LeftBrace, "{"},
		{KwReturn, "return"},
		{Identifier, "a"},
		{Plus, "+"},
		{Identifier, "b"},
		{Semicolon, ";"},
		{RightBrace, "}"},
		{KwClass, "class"},
		{Identifier, "MyClass"},
		{LeftBrace, "{"},
		{KwPrivate, "private"},
		{Identifier, "field1"},
		{Colon, ":"},
		{KwInt32, "int32"},
		{Semicolon, ";"},
		{KwPublic, "public"},
		{Identifier, "field2"},
		{Colon, ":"},
		{KwFloat64, "float64"},
		{Semicolon, ";"},
		{KwPublic, "public"},
		{Identifier, "get"},
		{LeftParen, "("},
		{RightParen, ")"},
		{KwInt32, "int32"},
		{LeftBrace, "{"},
		{KwReturn, "return"},
		{Identifier, "field1"},
		{Semicolon, ";"},
		{RightBrace, "}"},
		{RightBrace, "}"},
		{String, `"Hello, World!"`},
		{Byte, `'b'`},
		{EOF, ""},
	}
	lexer := New(input, "test.vs")
	for i, tt := range tests {
		tok := lexer.Next()
		if tok.Type != tt.Type {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, tt.Type.String(), tok.Type.String())
		}
		if tok.Lexeme != tt.Lexeme {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q", i, tt.Lexeme, tok.Lexeme)
		}
	}
}

func TestIdentifer(t *testing.T) {
	input := "foo bar x y z' tail"
	expected := []string{"foo", "bar", "x", "y", "z'", "tail"}
	lexer := New(input, "test.vs")
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != Identifier {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, Identifier.String(), tok.Type.String())
		}
		if tok.Lexeme != exp {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q", i, exp, tok.Lexeme)
		}
	}
}

func TestString(t *testing.T) {
	input := `"Hello, World!" "Line1\nLine2" "Quote: \"" "Backslash: \\"`
	expected := []string{`"Hello, World!"`, `"Line1\nLine2"`, `"Quote: \""`, `"Backslash: \\"`}
	lexer := New(input, "test.vs")
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != String {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, String.String(), tok.Type.String())
		}
		if tok.Lexeme != exp {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q", i, exp, tok.Lexeme)
		}
	}
}

func TestByte(t *testing.T) {
	input := `'a' '\n' '\'' '\\'`
	expected := []string{`'a'`, `'\n'`, `'\''`, `'\\'`}
	lexer := New(input, "test.vs")
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != Byte {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, Byte.String(), tok.Type.String())
		}
		if tok.Lexeme != exp {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q", i, exp, tok.Lexeme)
		}
	}
}

func TestComments(t *testing.T) {
	input := `// This is a comment
	// Another comment`
	expected := []string{"// This is a comment", "// Another comment"}
	lexer := New(input, "test.vs")
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != Comment {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, Comment.String(), tok.Type.String())
		}
		if tok.Lexeme != exp {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q", i, exp, tok.Lexeme)
		}
	}
}

func TestIllegalToken(t *testing.T) {
	input := "@"
	lexer := New(input, "test.vs")
	tok := lexer.Next()
	if tok.Type != Illegal {
		t.Fatalf("token type wrong. expected=%q, got=%q", Illegal.String(), tok.Type.String())
	}
}

func TestKeyword(t *testing.T) {
	input := "var const private public class return if else for"
	lexer := New(input, "test.vs")
	expected := []TokenType{KwVar, KwConst, KwPrivate, KwPublic, KwClass, KwReturn, KwIf, KwElse, KwFor}
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != exp {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, exp.String(), tok.Type.String())
		}
	}
}

func TestOperators(t *testing.T) {
	input := "+ - * / % = == != < <= > >= ! && ||"
	lexer := New(input, "test.vs")
	expected := []TokenType{Plus, Minus, Asterisk, Slash, Percent, Assign, Equal, NotEqual, LessThan, LessEqual, GreaterThan, GreaterEqual, Not, And, Or}
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != exp {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, exp.String(), tok.Type.String())
		}
	}
}

func TestDelimiters(t *testing.T) {
	input := "( ) { } [ ] , ; :"
	lexer := New(input, "test.vs")
	expected := []TokenType{LeftParen, RightParen, LeftBrace, RightBrace, LeftBracket, RightBracket, Comma, Semicolon, Colon}
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != exp {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, exp.String(), tok.Type.String())
		}
	}
}

func TestNumericLiterals(t *testing.T) {
	input := "123 45.67 0u 255u 32767 2147483647 9223372036854775807 3.4028235 1.7976931348623157"
	expected := []struct {
		Type   TokenType
		Lexeme string
	}{
		{Integer, "123"},
		{Float, "45.67"},
		{Unsigned, "0u"},
		{Unsigned, "255u"},
		{Integer, "32767"},
		{Integer, "2147483647"},
		{Integer, "9223372036854775807"},
		{Float, "3.4028235"},
		{Float, "1.7976931348623157"},
	}
	lexer := New(input, "test.vs")
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != exp.Type {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, exp.Type.String(), tok.Type.String())
		}
		if tok.Lexeme != exp.Lexeme {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q", i, exp.Lexeme, tok.Lexeme)
		}
	}
}

func TestEOF(t *testing.T) {
	input := ""
	lexer := New(input, "test.vs")
	tok := lexer.Next()
	if tok.Type != EOF {
		t.Fatalf("token type wrong. expected=%q, got=%q", EOF.String(), tok.Type.String())
	}
}

func TestWhitespaceHandling(t *testing.T) {
	input := "   \n\t  var   \n\t x  ;  "
	expected := []struct {
		Type   TokenType
		Lexeme string
	}{
		{KwVar, "var"},
		{Identifier, "x"},
		{Semicolon, ";"},
	}
	lexer := New(input, "test.vs")
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != exp.Type {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, exp.Type.String(), tok.Type.String())
		}
		if tok.Lexeme != exp.Lexeme {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q", i, exp.Lexeme, tok.Lexeme)
		}
	}
}

func TestComplexInput(t *testing.T) {
	input := `
	class Test {
		private value: int32;

		public setValue(int32 v) {
			value = v;
		}

		public getValue() int32 {
			return value;
		}
	}

	var t: Test;
	t = Test();
	t.setValue(42);
	println("Value: " + t.getValue());
	`
	lexer := New(input, "test.vs")
	for {
		tok := lexer.Next()
		if tok.Type == EOF {
			break
		}
		if tok.Type == Illegal {
			t.Fatalf("Illegal token encountered: %q", tok.Lexeme)
		}
	}
}

func TestIdentifierWithApostrophe(t *testing.T) {
	input := "data' value' test'"
	expected := []string{"data'", "value'", "test'"}
	lexer := New(input, "test.vs")
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != Identifier {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, Identifier.String(), tok.Type.String())
		}
		if tok.Lexeme != exp {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q", i, exp, tok.Lexeme)
		}
	}
}

func TestMultipleComments(t *testing.T) {
	input := `// First comment
	// Second comment
	// Third comment`
	expected := []string{"// First comment", "// Second comment", "// Third comment"}
	lexer := New(input, "test.vs")
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != Comment {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, Comment.String(), tok.Type.String())
		}
		if tok.Lexeme != exp {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q", i, exp, tok.Lexeme)
		}
	}
}

func TestStringWithEscapes(t *testing.T) {
	input := `"Line1\nLine2\tTabbed\"Quote\"\\Backslash"`
	expected := `"Line1\nLine2\tTabbed\"Quote\"\\Backslash"`
	lexer := New(input, "test.vs")
	tok := lexer.Next()
	if tok.Type != String {
		t.Fatalf("token type wrong. expected=%q, got=%q", String.String(), tok.Type.String())
	}
	if tok.Lexeme != expected {
		t.Fatalf("lexeme wrong. expected=%q, got=%q", expected, tok.Lexeme)
	}
}

func TestByteWithEscapes(t *testing.T) {
	input := `'\n' '\t' '\'' '\\'`
	expected := []string{`'\n'`, `'\t'`, `'\''`, `'\\'`}
	lexer := New(input, "test.vs")
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != Byte {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, Byte.String(), tok.Type.String())
		}
		if tok.Lexeme != exp {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q", i, exp, tok.Lexeme)
		}
	}
}

func TestMixedInput(t *testing.T) {
	input := `var count: int32 = 10; // Initialize count
	count = count + 1;
	println("Count is: " + count);`
	lexer := New(input, "test.vs")
	for {
		tok := lexer.Next()
		if tok.Type == EOF {
			break
		}
		if tok.Type == Illegal {
			t.Fatalf("Illegal token encountered: %q", tok.Lexeme)
		}
	}
}

func TestWhitespaceOnly(t *testing.T) {
	input := "   \n\t  \r\n   "
	lexer := New(input, "test.vs")
	tok := lexer.Next()
	if tok.Type != EOF {
		t.Fatalf("token type wrong. expected=%q, got=%q", EOF.String(), tok.Type.String())
	}
}

func TestNoInput(t *testing.T) {
	input := ""
	lexer := New(input, "test.vs")
	tok := lexer.Next()
	if tok.Type != EOF {
		t.Fatalf("token type wrong. expected=%q, got=%q", EOF.String(), tok.Type.String())
	}
}

func TestAdjacentOperators(t *testing.T) {
	input := "a+++b--*c"
	expected := []struct {
		Type   TokenType
		Lexeme string
	}{
		{Identifier, "a"},
		{Plus, "+"},
		{Plus, "+"},
		{Plus, "+"},
		{Identifier, "b"},
		{Minus, "-"},
		{Minus, "-"},
		{Asterisk, "*"},
		{Identifier, "c"},
	}
	lexer := New(input, "test.vs")
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != exp.Type || tok.Lexeme != exp.Lexeme {
			t.Fatalf("tests[%d] failed: expected=(%q,%q), got=(%q,%q)", i, exp.Type.String(), exp.Lexeme, tok.Type.String(), tok.Lexeme)
		}
	}
}

func TestNumbersWithLeadingZeros(t *testing.T) {
	input := "00123 00045.67"
	expected := []struct {
		Type   TokenType
		Lexeme string
	}{
		{Integer, "00123"},
		{Float, "00045.67"},
	}
	lexer := New(input, "test.vs")
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != exp.Type || tok.Lexeme != exp.Lexeme {
			t.Fatalf("tests[%d] failed: expected=(%q,%q), got=(%q,%q)", i, exp.Type.String(), exp.Lexeme, tok.Type.String(), tok.Lexeme)
		}
	}
}

func TestIdentifierStartingWithKeyword(t *testing.T) {
	input := "varName constVal privateMethod"
	expected := []string{"varName", "constVal", "privateMethod"}
	lexer := New(input, "test.vs")
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != Identifier {
			t.Fatalf("tests[%d] - expected Identifier, got=%q", i, tok.Type.String())
		}
		if tok.Lexeme != exp {
			t.Fatalf("tests[%d] - expected lexeme %q, got=%q", i, exp, tok.Lexeme)
		}
	}
}

func TestCommentWithSpecialChars(t *testing.T) {
	input := `// comment with symbols !@#$%^&*()_+{}:"<>?`
	lexer := New(input, "test.vs")
	tok := lexer.Next()
	if tok.Type != Comment {
		t.Fatalf("expected Comment, got=%q", tok.Type.String())
	}
	if tok.Lexeme != input {
		t.Fatalf("expected lexeme %q, got=%q", input, tok.Lexeme)
	}
}

func TestMixedUnusualWhitespace(t *testing.T) {
	input := "var\t x \n= 10 ;"
	expected := []struct {
		Type   TokenType
		Lexeme string
	}{
		{KwVar, "var"},
		{Identifier, "x"},
		{Assign, "="},
		{Integer, "10"},
		{Semicolon, ";"},
	}
	lexer := New(input, "test.vs")
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != exp.Type || tok.Lexeme != exp.Lexeme {
			t.Fatalf("tests[%d] failed: expected=(%q,%q), got=(%q,%q)", i, exp.Type.String(), exp.Lexeme, tok.Type.String(), tok.Lexeme)
		}
	}
}

func TestStringWithLineBreaks(t *testing.T) {
	input := `"This is a \n multi-line \n string"`
	lexer := New(input, "test.vs")
	tok := lexer.Next()
	if tok.Type != String {
		t.Fatalf("expected String, got=%q", tok.Type.String())
	}
	if tok.Lexeme != input {
		t.Fatalf("expected lexeme %q, got=%q", input, tok.Lexeme)
	}
}

func TestNumberEdgeCases(t *testing.T) {
	input := "0 0.0 0u 0.0001"
	expected := []struct {
		Type   TokenType
		Lexeme string
	}{
		{Integer, "0"},
		{Float, "0.0"},
		{Unsigned, "0u"},
		{Float, "0.0001"},
	}
	lexer := New(input, "test.vs")
	for i, exp := range expected {
		tok := lexer.Next()
		if tok.Type != exp.Type || tok.Lexeme != exp.Lexeme {
			t.Fatalf("tests[%d] failed: expected=(%q,%q), got=(%q,%q)", i, exp.Type.String(), exp.Lexeme, tok.Type.String(), tok.Lexeme)
		}
	}
}

func TestStringWithOnlyEscapes(t *testing.T) {
	input := `"\\n\\t\\\"\\'"`
	expected := input
	lexer := New(input, "test.vs")
	tok := lexer.Next()
	if tok.Type != String {
		t.Fatalf("expected String, got=%q", tok.Type.String())
	}
	if tok.Lexeme != expected {
		t.Fatalf("expected lexeme %q, got=%q", expected, tok.Lexeme)
	}
}
