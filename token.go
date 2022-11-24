package sky

type tokenType string

const (
	EOF        = "EOF"
	ILLEGAL    = "ILLEGAL"
	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"
	FLOAT      = "FLOAT"

	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	BANG   = "!"
	STAR   = "*"
	SLASH  = "/"
	MOD    = "%"

	LT     = "<"
	LEQ    = "<="
	GT     = ">"
	GEQ    = ">="
	EQ     = "=="
	NOT_EQ = "!="

	COMMA = ","

	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"

	LBRACE = "{"
	RBRACE = "}"

	FUNC   = "func"
	VAR    = "var"
	TRUE   = "true"
	FALSE  = "false"
	IF     = "if"
	ELIF   = "elif"
	ELSE   = "else"
	RETURN = "return"
	WHILE  = "while"
	FOR    = "for"
)

type token struct {
	tokenType tokenType
	lexeme    string
	line      int
	column    int
}

func newToken(tokenType tokenType, lexeme string, line, column int) token {
	return token{
		tokenType: tokenType,
		lexeme:    lexeme,
		line:      line,
		column:    column,
	}
}

var keywords = map[tokenType]tokenType{
	FUNC:   FUNC,
	VAR:    VAR,
	TRUE:   TRUE,
	FALSE:  FALSE,
	IF:     IF,
	ELIF:   ELIF,
	ELSE:   ELSE,
	RETURN: RETURN,
	WHILE:  WHILE,
	FOR:    FOR,
}

func toKeyword(key tokenType) (bool, tokenType) {
	if e, ok := keywords[key]; ok {
		return true, e
	}
	return false, ""
}
