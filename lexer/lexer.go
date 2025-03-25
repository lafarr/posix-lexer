package lexer

import (
	"errors"
	"strings"
)

// \$, \\, \"
func Tokenize(command string) ([]*Token, error) {
	tokens := []*Token{}

	inSingleQuotes := false
	inDoubleQuotes := false
	currToken := ""

	for i := 0; i < len(command); i++ {
		ch := string(command[i])

		switch ch {
		case ` `:
			if !inSingleQuotes && !inDoubleQuotes && len(currToken) > 0 {
				tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
				currToken = ""
			} else {
				currToken += ch
			}
		case `<`:
			if !inSingleQuotes && !inDoubleQuotes {
				if i < len(command)-1 && command[i+1] == '&' {
					// Dup file descriptor
					if len(currToken) > 0 {
						tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
					}
					tokens = append(tokens, New(TokenDupInputFd, "<&"))
					currToken = ""
				} else if i < len(command)-2 && command[i+1] == '<' && command[i+2] == '<' {
					// Here string
					if len(currToken) > 0 {
						tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
					}
					tokens = append(tokens, New(TokenHereString, "<<<"))
					currToken = ""
				} else if i < len(command)-1 && command[i+1] == '<' {
					// Here document
					if len(currToken) > 0 {
						tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
					}
					tokens = append(tokens, New(TokenHereDoc, "<<"))
					currToken = ""

				} else {
					// Redirect input
					if len(currToken) > 0 {
						tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
					}
					tokens = append(tokens, New(TokenRedirectIn, "<"))
					currToken = ""
				}
			} else {
				currToken += ch
			}
		// Have to also handle >&, >>
		case `>`:
			if !inSingleQuotes && !inDoubleQuotes {
				if i < len(command)-1 && command[i+1] == '&' {
					// Dup file descriptor
					if len(currToken) > 0 {
						tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
					}
					tokens = append(tokens, New(TokenDupOutputFd, ">&"))
					currToken = ""
				} else if i < len(command)-1 && command[i+1] == '>' {
					// Append document
					if len(currToken) > 0 {
						tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
					}
					tokens = append(tokens, New(TokenAppendOut, ">>"))
					currToken = ""

				} else {
					// Redirect output
					if len(currToken) > 0 {
						tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
					}
					tokens = append(tokens, New(TokenRedirectOut, ">"))
					currToken = ""
				}
			} else {
				currToken += ch
			}

		// Have to also handle &&
		case `&`:
			if !inSingleQuotes && !inDoubleQuotes {
				if len(currToken) > 0 {
					tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
				}

				if i != len(command)-1 && command[i+1] == '&' {
					tokens = append(tokens, New(TokenAnd, "&&"))
				} else {
					tokens = append(tokens, New(TokenBackground, "&"))
				}
			} else {
				currToken += ch
			}
		case `;`:
			if !inSingleQuotes && !inDoubleQuotes {
				if len(currToken) > 0 {
					tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
					currToken = ""
				}
			} else {
				currToken += ch
			}
		// Have to also handle ||
		case `|`:
			if !inSingleQuotes && !inDoubleQuotes {
				if len(currToken) > 0 {
					tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
					currToken = ""
				}
				if i < len(command)-1 && command[i+1] == '|' {
					tokens = append(tokens, New(TokenOr, "||"))
				} else {
					tokens = append(tokens, New(TokenPipe, "|"))
				}
			} else {
				currToken += ch
			}
		case `(`:
			if !inSingleQuotes && !inDoubleQuotes {
				if len(currToken) > 0 {
					tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
					currToken = ""
				}
				tokens = append(tokens, New(TokenLeftParen, "("))
			} else {
				currToken += ch
			}
		case `)`:
			if !inSingleQuotes && !inDoubleQuotes {
				if len(currToken) > 0 {
					tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
					currToken = ""
				}
				tokens = append(tokens, New(TokenRightParen, ")"))
			} else {
				currToken += ch
			}
		case `{`:
			if !inSingleQuotes && !inDoubleQuotes {
				if len(currToken) > 0 {
					tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
					currToken = ""
				}
				tokens = append(tokens, New(TokenLeftBrace, "{"))
			} else {
				currToken += ch
			}
		case `}`:
			if !inSingleQuotes && !inDoubleQuotes {
				if len(currToken) > 0 {
					tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
					currToken = ""
				}
				tokens = append(tokens, New(TokenRightBrace, "}"))
			} else {
				currToken += ch
			}
			// can be escaped, but we handle the escaping in the case for '\'

		case `'`:
			if inDoubleQuotes {
				currToken += ch
			} else {
				inSingleQuotes = !inSingleQuotes
			}
		// variable replacement
		case `$`:
			// TODO: Also handle $(()) for math
			// don't need to do this here, will do this in an expansion phase
			currToken += ch
		// command sub
		case "`":
			// don't need to do this here, will do this in an expansion phase
			currToken += ch
		case `"`:
			if inSingleQuotes {
				currToken += ch
			} else {
				inDoubleQuotes = !inDoubleQuotes
			}
		case `\`:
			// need to check for $, `, ", ', \
			if inDoubleQuotes {
				if i < len(command)-1 {
					switch nextCh := command[i+1]; nextCh {
					case '$':
						currToken += string(nextCh)
					case '`':
						currToken += string(nextCh)
					case '"':
						currToken += string(nextCh)
					case '\'':
						currToken += string(nextCh)
					case '\\':
						currToken += string(nextCh)
					default:
						currToken += "\\" + string(nextCh)
					}
					i++
				} else {
					return nil, errors.New("Double quotes was not closed")
				}
			} else if inSingleQuotes {
				// can't escape anything in single quotes
				currToken += ch
			} else {
				if i < len(command)-1 && command[i+1] == '\n' {
					currToken += ch + string('\n')
					// Don't want to proces \n twice
				} else {
					currToken += string(command[i+1])
				}
				i++
			}
			/* Optional, mostly for actual scripting
			case `=`:
			case `#`:
			case `*`:
			case `?`:
			case `!`:
			case `%n`:
			*/
		default:
			currToken += ch
		}
	}
	if len(currToken) > 0 {
		tokens = append(tokens, New(TokenWord, strings.TrimSpace(currToken)))
	}

	return tokens, nil
}
