package lexer

import "fmt"

type TokenType int

/*
Test Operators (used with [ ] or [[ ]]):

-eq, -ne, -lt, -le, -gt, -ge - Numeric comparisons
=, != - String comparisons
-z, -n - String empty/not empty
-e, -f, -d, -s, etc. - File tests
*/

const (
	// Basic tokens
	TokenWord TokenType = iota // Command names and arguments
	TokenEof                   // End of input

	// Operators
	TokenPipe        // | (pipe)
	TokenRedirectIn  // < (input redirection)
	TokenRedirectOut // > (output redirection)
	TokenAppendOut   // >> (append output redirection)
	TokenHereDoc     // << (here document)
	TokenHereString  // <<< (Here string)
	TokenDupOutputFd // >&
	TokenDupInputFd  // <&

	// Control operators
	TokenBackground // & (background process)
	TokenSemicolon  // ; (command separator)
	TokenAnd        // && (logical AND)
	TokenOr         // || (logical OR)

	// Grouping
	TokenLeftParen  // ( (subshell start)
	TokenRightParen // ) (subshell end)
	TokenLeftBrace  // { (command group start)
	TokenRightBrace // } (command group end)

	// Special characters
	// TODO: Add:
	//          * $() for command sub
	// 			* $(()) for arithmetic expansion
	// 			* ${param} for parameter expansion "with various modifiers"
	TokenDollar      // $ (variable expansion)
	TokenBacktick    // ` (command substitution)
	TokenQuote       // " (double quote)
	TokenSingleQuote // ' (single quote)
	TokenEscape      // \ (escape character)
	TokenTilde       // ~ (home directory expansion)

	// Optional - for more complex shells
	TokenAssign        // = (variable assignment)
	TokenComment       // # (comment)
	TokenGlob          // * or ? (glob patterns)
	TokenNegExitStatus // !
	TokenJobSpec       // %n - Job specification
)

type Token struct {
	Type       TokenType
	TextualRep string
}

func New(tokenType TokenType, textualRep string) *Token {
	return &Token{tokenType, textualRep}
}

func Type(token *Token) TokenType {
	return token.Type
}

func TextualRep(token *Token) string {
	return token.TextualRep
}

func ToString(tokens []*Token) string {
	fmt.Printf("Tokens length: %d\n", len(tokens))
	asString := ""

	for _, val := range tokens {
		fmt.Println(Type(val))
		asString += TextualRep(val) + ", "
	}

	return asString[:len(asString)-2]
}
