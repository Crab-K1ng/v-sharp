package lexer

type TokenType int

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Column int

	File   string
	Source []rune
}

const (
	Illegal TokenType = iota
	Comment

	Identifier
	Integer  // Represents all integer types (untyped)
	Float    // Represents all float types (untyped)
	Unsigned // Represents all unsigned integer types (untyped)
	Int8
	Int16
	Int32
	Int64
	UInt8
	UInt16
	UInt32
	UInt64
	Float32
	Float64
	Boolean
	String
	Byte
	Void

	Plus     // +
	Minus    // -
	Asterisk // *
	Slash    // /
	Percent  // %

	LeftParen    // (
	RightParen   // )
	LeftBrace    // {
	RightBrace   // }
	LeftBracket  // [
	RightBracket // ]

	Comma     // ,
	Semicolon // ;
	Colon     // :
	Dot       // .

	Assign       // =
	Equal        // ==
	NotEqual     // !=
	LessThan     // <
	GreaterThan  // >
	LessEqual    // <=
	GreaterEqual // >=

	Not // !
	And // &&
	Or  // ||

	KwPublic
	KwPrivate
	KwVirtual
	KwOverride
	KwStatic
	KwConst
	KwVar
	KwIf
	KwElse
	KwMatch
	KwFor
	KwReturn
	KwStructure
	KwEnumeration
	KwDefine
	KwTypedef
	KwClass
	KwInt8
	KwInt16
	KwInt32
	KwInt64
	KwUInt8
	KwUInt16
	KwUInt32
	KwUInt64
	KwFloat32
	KwFloat64
	KwBoolean
	KwString
	KwByte
	KwVoid

	EOF
)

var keywordsMap = map[string]TokenType{
	"public":      KwPublic,
	"private":     KwPrivate,
	"virtual":     KwVirtual,
	"override":    KwOverride,
	"static":      KwStatic,
	"const":       KwConst,
	"var":         KwVar,
	"if":          KwIf,
	"else":        KwElse,
	"match":       KwMatch,
	"for":         KwFor,
	"return":      KwReturn,
	"structure":   KwStructure,
	"enumeration": KwEnumeration,
	"define":      KwDefine,
	"typedef":     KwTypedef,
	"class":       KwClass,
	"int8":        KwInt8,
	"int16":       KwInt16,
	"int32":       KwInt32,
	"int64":       KwInt64,
	"uint8":       KwUInt8,
	"uint16":      KwUInt16,
	"uint32":      KwUInt32,
	"uint64":      KwUInt64,
	"float32":     KwFloat32,
	"float64":     KwFloat64,
	"boolean":     KwBoolean,
	"string":      KwString,
	"byte":        KwByte,
	"void":        KwVoid,
}

var precedence = map[TokenType]int{
	Or:  1,
	And: 2,

	Equal:        3,
	NotEqual:     3,
	LessThan:     4,
	GreaterThan:  4,
	LessEqual:    4,
	GreaterEqual: 4,

	Plus:     5,
	Minus:    5,
	Asterisk: 6,
	Slash:    6,
	Percent:  6,
}

func (t TokenType) String() string {
	switch t {
	case Illegal:
		return "Illegal"
	case Comment:
		return "Comment"
	case Identifier:
		return "Identifier"
	case Int8:
		return "Int8"
	case Int16:
		return "Int16"
	case Int32:
		return "Int32"
	case Int64:
		return "Int64"
	case UInt8:
		return "UInt8"
	case UInt16:
		return "UInt16"
	case UInt32:
		return "UInt32"
	case UInt64:
		return "UInt64"
	case Float32:
		return "Float32"
	case Float64:
		return "Float64"
	case Boolean:
		return "Boolean"
	case String:
		return "String"
	case Byte:
		return "Byte"
	case Void:
		return "Void"
	case Plus:
		return "Plus"
	case Minus:
		return "Minus"
	case Asterisk:
		return "Asterisk"
	case Slash:
		return "Slash"
	case Percent:
		return "Percent"
	case LeftParen:
		return "LeftParen"
	case RightParen:
		return "RightParen"
	case LeftBrace:
		return "LeftBrace"
	case RightBrace:
		return "RightBrace"
	case LeftBracket:
		return "LeftBracket"
	case RightBracket:
		return "RightBracket"
	case Comma:
		return "Comma"
	case Semicolon:
		return "Semicolon"
	case Colon:
		return "Colon"
	case Assign:
		return "Assign"
	case Equal:
		return "Equal"
	case NotEqual:
		return "NotEqual"
	case LessThan:
		return "LessThan"
	case GreaterThan:
		return "GreaterThan"
	case LessEqual:
		return "LessEqual"
	case GreaterEqual:
		return "GreaterEqual"
	case Not:
		return "Not"
	case And:
		return "And"
	case Or:
		return "Or"
	case KwPublic:
		return "KwPublic"
	case KwPrivate:
		return "KwPrivate"
	case KwVirtual:
		return "KwVirtual"
	case KwOverride:
		return "KwOverride"
	case KwStatic:
		return "KwStatic"
	case KwConst:
		return "KwConst"
	case KwVar:
		return "KwVar"
	case KwIf:
		return "KwIf"
	case KwElse:
		return "KwElse"
	case KwMatch:
		return "KwMatch"
	case KwFor:
		return "KwFor"
	case KwReturn:
		return "KwReturn"
	case KwStructure:
		return "KwStructure"
	case KwEnumeration:
		return "KwEnumeration"
	case KwDefine:
		return "KwDefine"
	case KwTypedef:
		return "KwTypedef"
	case KwClass:
		return "KwClass"
	case Integer:
		return "Integer"
	case Float:
		return "Float"
	case Unsigned:
		return "Unsigned"
	case KwInt8:
		return "KwInt8"
	case KwInt16:
		return "KwInt16"
	case KwInt32:
		return "KwInt32"
	case KwInt64:
		return "KwInt64"
	case KwUInt8:
		return "KwUInt8"
	case KwUInt16:
		return "KwUInt16"
	case KwUInt32:
		return "KwUInt32"
	case KwUInt64:
		return "KwUInt64"
	case KwFloat32:
		return "KwFloat32"
	case KwFloat64:
		return "KwFloat64"
	case KwBoolean:
		return "KwBoolean"
	case KwString:
		return "KwString"
	case KwByte:
		return "KwByte"
	case KwVoid:
		return "KwVoid"
	case EOF:
		return "EOF"
	default:
		return "Unknown"
	}
}
