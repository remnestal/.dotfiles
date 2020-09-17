package csgods

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.com/abios/motattack/pkg/csgods/events"
)

const (
	DATETIME_FORMAT = "01/02/2006 - 15:04:05.999"
)

type parser struct {
	l       *lexer
	current token
	err     error
}

func NewParser(input []byte) *parser {
	l := lex(input)
	return &parser{l: l, current: l.Next()}
}

func (p *parser) More() bool {
	return p.current.typ != tokenEOF && p.err == nil
}

func (p *parser) Next() (events.Event, error) {
	return p.parse(), p.err
}

func (p *parser) eat(expected tokenType) {
	if p.current.typ == expected {
		p.current = p.l.Next()
	} else {
		p.panik(fmt.Sprintf("expected %v, have %v", expected, p.current))
	}
}

func (p *parser) panik(msg string, a ...interface{}) {
	slice_start := p.l.pos

	slice_end := slice_start + 20
	if slice_end > len(p.l.input) {
		slice_end = len(p.l.input)
	}

	p.err = fmt.Errorf("current input: %q\n\tcurrent token: %v:\n\t%w", p.l.input[slice_start:slice_end], p.current, fmt.Errorf(msg, a...))
}

func (p *parser) discardWhitespace() {
	for p.current.typ == tokenSpace || p.current.typ == tokenTab {
		p.eat(p.current.typ)
	}
}

func (p *parser) discardWord(expected string) {
	if string(p.current.val) == expected {
		p.eat(tokenWord)
	} else {
		p.panik(fmt.Sprintf("expected word %v, have %v", expected, string(p.current.val)))
	}
}

func (p *parser) discardWordAndWhitespace(expected string) {
	p.discardWhitespace()
	p.discardWord(expected)
	p.discardWhitespace()
}

func (p *parser) eatUntil(stop tokenType) string {
	result := strings.Builder{}
	for p.current.typ != tokenEOF && p.current.typ != stop {
		result.Write(p.current.val)
		p.eat(p.current.typ)
	}
	return result.String()
}

func (p *parser) eatUntilFollowedBy(stop, after tokenType) string {
	result := strings.Builder{}
	for {
		result.WriteString(p.eatUntil(stop))
		peek := p.l.Peek()
		if p.current.typ == tokenEOF || peek.typ == tokenEOF || peek.typ == after {
			break
		}
		result.Write(p.current.val)
		p.eat(stop)
	}
	return result.String()
}

func (p *parser) eatWord() string {
	word := string(p.current.val)
	p.eat(tokenWord)
	return word
}

func (p *parser) eatNum() float64 {
	word := p.eatWord()

	f, err := strconv.ParseFloat(word, 64)
	if err != nil {
		p.panik("parsing string as float: %w", err)
	}
	return f
}

func (p *parser) parse() events.Event {
	/*
		Each event starts with either
			* tokenHeader -> tokenStrange -> tokenSeparator -> tokenDatetime
			* tokenHeader -> tokenStrange -> tokenUniqueToken ->  tokenSeparator -> tokenDatetime
			* tokenDatetime
			* tokenWord -> tokenspace -> tokenQuote -> tokenDatetime -> tokenQuote
	*/
	var token []byte
	if p.current.typ == tokenHeader {
		p.eat(tokenHeader)
		p.eat(tokenStrange)
		if p.current.typ == tokenUniqueToken {
			token = p.current.val
			p.eat(tokenUniqueToken)
		}
		p.eat(tokenSpace)
		p.eat(tokenSeparator)
	}

	dt := p.current
	p.eat(tokenDatetime)
	datetime, err := time.Parse(DATETIME_FORMAT, string(dt.val))
	if err != nil {
		p.panik("parsing datetime format: %w", err)
	}

	e := p.event()
	if e != nil {
		e.SetTimestamp(datetime)
		if token != nil {
			e.SetUniqueToken(token)
		}
	}

	// Clear any potential line feeds
	if p.current.typ == tokenNewline {
		p.eat(tokenNewline)
	}

	return e
}

func (p *parser) event() events.Event {
	switch p.current.typ {
	case tokenBell:
		/*
			E.g

			\a[FACEIT^] 15 minutes left to !ready up.
			\a(ADMIN)\x01: ? ???? ??? ??? ????, ????? ?????? ??????
		*/
		p.eat(tokenBell)
		switch p.current.typ {
		case tokenLeftBracket:
			p.eat(tokenLeftBracket)
			p.discardWord("FACEIT^")
			p.eat(tokenRightBracket)
			return p.faceit()
		case tokenLeftParentheses:
			p.eat(tokenLeftParentheses)
			p.discardWord("ADMIN")
			p.eat(tokenRightParentheses)
			return p.admin()
		default:
			p.panik("unrecognized token following bell: %q", p.current)
		}
	case tokenLeftBracket:
		/* E.g

		[META] Loaded 0 plugins (1 already loaded)
		[FACEIT] Uploading GOTV demo
		[FACEIT MATCH] The heartbeat failed. No response from the server.

		*/
		p.eat(tokenLeftBracket)
		word := p.eatUntil(tokenRightBracket)
		p.eat(tokenRightBracket)
		switch word {
		case "META":
			return p.meta()
		case "FACEIT", "FACEIT^", "FACEIT MATCH":
			return p.faceit()
		default:
			p.panik("unrecognized word inside brackets: %q", word)
		}
	case tokenQuote:
		// Read the entire quoted text and do something with it. It's either a cvar or "<player>"
		peek := p.l.Peek()
		if _, ok := CVARS[string(peek.val)]; ok {
			return p.cvar()
		}

		// It's a player
		return p.PlayerEventData()
	case tokenWord:
		first_word := p.eatWord()
		switch first_word {
		case "Team":
			return p.team()
		case "Molotov":
			return p.molotov()
		case "World":
			return p.world()
		case "Log":
			return p.log()
		case "Loading":
			return p.loading()
		case "server":
			return p.server()
		case "server_cvar:":
			return p.server_cvar()
		case "server_message:":
			return p.server_message()
		case "Started":
			return p.started()
		case "Starting":
			return p.starting()
		case "Game":
			return p.game()
		case "Your":
			return p.server_restart()
		case "Player":
			return p.encryptionkey()
		case "Vote":
			return p.vote()
		case "Match":
			return p.match()
		case "ACCOLADE":
			return p.accolade()
		case "Error:":
			return p.errormessage()
		case "rcon":
			return p.rcon()
		case "eBot":
			return p.ebot()
		case "get5_event:":
			return p.get5_event()
		case "STEAMAUTH:":
			return p.steamauth()
		default:
			p.panik("unrecognized first word of event: %q", first_word)
		}
	}
	p.panik("unrecognized start of event %v", p.current)
	return nil
}
