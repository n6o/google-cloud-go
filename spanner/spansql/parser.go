package spansql

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"cloud.google.com/go/civil"
	"log"
)

const debug = false

func gologoo__debugf_531485d04a38cff0d09d3c55685c367d(format string, args ...interface {
}) {
	if !debug {
		return
	}
	fmt.Fprintf(os.Stderr, "spansql debug: "+format+"\n", args...)
}
func gologoo__ParseDDL_531485d04a38cff0d09d3c55685c367d(filename, s string) (*DDL, error) {
	p := newParser(filename, s)
	ddl := &DDL{Filename: filename}
	for {
		p.skipSpace()
		if p.done {
			break
		}
		stmt, err := p.parseDDLStmt()
		if err != nil {
			return nil, err
		}
		ddl.List = append(ddl.List, stmt)
		tok := p.next()
		if tok.err == eof {
			break
		} else if tok.err != nil {
			return nil, tok.err
		}
		if tok.value == ";" {
			continue
		} else {
			return nil, p.errorf("unexpected token %q", tok.value)
		}
	}
	if p.Rem() != "" {
		return nil, fmt.Errorf("unexpected trailing contents %q", p.Rem())
	}
	for _, com := range p.comments {
		c := &Comment{Marker: com.marker, Isolated: com.isolated, Start: com.start, End: com.end, Text: com.text}
		var prefix string
		for i, line := range c.Text {
			line = strings.TrimRight(line, " \b\t")
			c.Text[i] = line
			trim := len(line) - len(strings.TrimLeft(line, " \b\t"))
			if i == 0 {
				prefix = line[:trim]
			} else {
				for !strings.HasPrefix(line, prefix) {
					prefix = prefix[:len(prefix)-1]
				}
			}
			if prefix == "" {
				break
			}
		}
		if prefix != "" {
			for i, line := range c.Text {
				c.Text[i] = strings.TrimPrefix(line, prefix)
			}
		}
		ddl.Comments = append(ddl.Comments, c)
	}
	return ddl, nil
}
func gologoo__ParseDDLStmt_531485d04a38cff0d09d3c55685c367d(s string) (DDLStmt, error) {
	p := newParser("-", s)
	stmt, err := p.parseDDLStmt()
	if err != nil {
		return nil, err
	}
	if p.Rem() != "" {
		return nil, fmt.Errorf("unexpected trailing contents %q", p.Rem())
	}
	return stmt, nil
}
func gologoo__ParseDMLStmt_531485d04a38cff0d09d3c55685c367d(s string) (DMLStmt, error) {
	p := newParser("-", s)
	stmt, err := p.parseDMLStmt()
	if err != nil {
		return nil, err
	}
	if p.Rem() != "" {
		return nil, fmt.Errorf("unexpected trailing contents %q", p.Rem())
	}
	return stmt, nil
}
func gologoo__ParseQuery_531485d04a38cff0d09d3c55685c367d(s string) (Query, error) {
	p := newParser("-", s)
	q, err := p.parseQuery()
	if err != nil {
		return Query{}, err
	}
	if p.Rem() != "" {
		return Query{}, fmt.Errorf("unexpected trailing query contents %q", p.Rem())
	}
	return q, nil
}

type token struct {
	value        string
	err          *parseError
	line, offset int
	typ          tokenType
	float64      float64
	string       string
	int64Base    int
}
type tokenType int

const (
	unknownToken tokenType = iota
	int64Token
	float64Token
	stringToken
	bytesToken
	unquotedID
	quotedID
)

func (t *token) gologoo__String_531485d04a38cff0d09d3c55685c367d() string {
	if t.err != nil {
		return fmt.Sprintf("parse error: %v", t.err)
	}
	return strconv.Quote(t.value)
}

type parseError struct {
	message  string
	filename string
	line     int
	offset   int
}

func (pe *parseError) gologoo__Error_531485d04a38cff0d09d3c55685c367d() string {
	if pe == nil {
		return "<nil>"
	}
	if pe.line == 1 {
		return fmt.Sprintf("%s:1.%d: %v", pe.filename, pe.offset, pe.message)
	}
	return fmt.Sprintf("%s:%d: %v", pe.filename, pe.line, pe.message)
}

var eof = &parseError{message: "EOF"}

type parser struct {
	s            string
	done         bool
	backed       bool
	cur          token
	filename     string
	line, offset int
	comments     []comment
}
type comment struct {
	marker     string
	isolated   bool
	start, end Position
	text       []string
}

func (p *parser) gologoo__Pos_531485d04a38cff0d09d3c55685c367d() Position {
	return Position{Line: p.cur.line, Offset: p.cur.offset}
}
func gologoo__newParser_531485d04a38cff0d09d3c55685c367d(filename, s string) *parser {
	return &parser{s: s, cur: token{line: 1}, filename: filename, line: 1}
}
func (p *parser) gologoo__Rem_531485d04a38cff0d09d3c55685c367d() string {
	rem := p.s
	if p.backed {
		rem = p.cur.value + rem
	}
	i := 0
	for ; i < len(rem); i++ {
		if !isSpace(rem[i]) {
			break
		}
	}
	return rem[i:]
}
func (p *parser) gologoo__String_531485d04a38cff0d09d3c55685c367d() string {
	if p.backed {
		return fmt.Sprintf("next tok: %s (rem: %q)", &p.cur, p.s)
	}
	return fmt.Sprintf("rem: %q", p.s)
}
func (p *parser) gologoo__errorf_531485d04a38cff0d09d3c55685c367d(format string, args ...interface {
}) *parseError {
	pe := &parseError{message: fmt.Sprintf(format, args...), filename: p.filename, line: p.cur.line, offset: p.cur.offset}
	p.cur.err = pe
	p.done = true
	return pe
}
func gologoo__isInitialIdentifierChar_531485d04a38cff0d09d3c55685c367d(c byte) bool {
	switch {
	case 'A' <= c && c <= 'Z':
		return true
	case 'a' <= c && c <= 'z':
		return true
	case c == '_':
		return true
	}
	return false
}
func gologoo__isIdentifierChar_531485d04a38cff0d09d3c55685c367d(c byte) bool {
	switch {
	case 'A' <= c && c <= 'Z':
		return true
	case 'a' <= c && c <= 'z':
		return true
	case '0' <= c && c <= '9':
		return true
	case c == '_':
		return true
	}
	return false
}
func gologoo__isHexDigit_531485d04a38cff0d09d3c55685c367d(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}
func gologoo__isOctalDigit_531485d04a38cff0d09d3c55685c367d(c byte) bool {
	return '0' <= c && c <= '7'
}
func (p *parser) gologoo__consumeNumber_531485d04a38cff0d09d3c55685c367d() {
	i, neg, base := 0, false, 10
	float, e, dot := false, false, false
	if p.s[i] == '-' {
		neg = true
		i++
	} else if p.s[i] == '+' {
		i++
	}
	if strings.HasPrefix(p.s[i:], "0x") || strings.HasPrefix(p.s[i:], "0X") {
		base = 16
		i += 2
	}
	d0 := i
digitLoop:
	for i < len(p.s) {
		switch c := p.s[i]; {
		case '0' <= c && c <= '9':
			i++
		case base == 16 && 'A' <= c && c <= 'F':
			i++
		case base == 16 && 'a' <= c && c <= 'f':
			i++
		case base == 10 && (c == 'e' || c == 'E'):
			if e {
				p.errorf("bad token %q", p.s[:i])
				return
			}
			float, e = true, true
			i++
			if i < len(p.s) && (p.s[i] == '+' || p.s[i] == '-') {
				i++
			}
		case base == 10 && c == '.':
			if dot || e {
				p.errorf("bad token %q", p.s[:i])
				return
			}
			float, dot = true, true
			i++
		default:
			break digitLoop
		}
	}
	if d0 == i {
		p.errorf("no digits in numeric literal")
		return
	}
	sign := ""
	if neg {
		sign = "-"
	}
	p.cur.value, p.s = p.s[:i], p.s[i:]
	p.offset += i
	var err error
	if float {
		p.cur.typ = float64Token
		p.cur.float64, err = strconv.ParseFloat(sign+p.cur.value[d0:], 64)
	} else {
		p.cur.typ = int64Token
		p.cur.value = sign + p.cur.value[d0:]
		p.cur.int64Base = base
	}
	if err != nil {
		p.errorf("bad numeric literal %q: %v", p.cur.value, err)
	}
}
func (p *parser) gologoo__consumeString_531485d04a38cff0d09d3c55685c367d() {
	delim := p.stringDelimiter()
	if p.cur.err != nil {
		return
	}
	p.cur.string, p.cur.err = p.consumeStringContent(delim, false, true, "string literal")
	p.cur.typ = stringToken
}
func (p *parser) gologoo__consumeRawString_531485d04a38cff0d09d3c55685c367d() {
	p.s = p.s[1:]
	delim := p.stringDelimiter()
	if p.cur.err != nil {
		return
	}
	p.cur.string, p.cur.err = p.consumeStringContent(delim, true, true, "raw string literal")
	p.cur.typ = stringToken
}
func (p *parser) gologoo__consumeBytes_531485d04a38cff0d09d3c55685c367d() {
	p.s = p.s[1:]
	delim := p.stringDelimiter()
	if p.cur.err != nil {
		return
	}
	p.cur.string, p.cur.err = p.consumeStringContent(delim, false, false, "bytes literal")
	p.cur.typ = bytesToken
}
func (p *parser) gologoo__consumeRawBytes_531485d04a38cff0d09d3c55685c367d() {
	p.s = p.s[2:]
	delim := p.stringDelimiter()
	if p.cur.err != nil {
		return
	}
	p.cur.string, p.cur.err = p.consumeStringContent(delim, true, false, "raw bytes literal")
	p.cur.typ = bytesToken
}
func (p *parser) gologoo__stringDelimiter_531485d04a38cff0d09d3c55685c367d() string {
	c := p.s[0]
	if c != '"' && c != '\'' {
		p.errorf("invalid string literal")
		return ""
	}
	if len(p.s) >= 3 && p.s[1] == c && p.s[2] == c {
		return p.s[:3]
	}
	return p.s[:1]
}
func (p *parser) gologoo__consumeStringContent_531485d04a38cff0d09d3c55685c367d(delim string, raw, unicode bool, name string) (string, *parseError) {
	if len(delim) == 3 {
		name = "triple-quoted " + name
	}
	i := len(delim)
	var content []byte
	for i < len(p.s) {
		if strings.HasPrefix(p.s[i:], delim) {
			i += len(delim)
			p.s = p.s[i:]
			p.offset += i
			return string(content), nil
		}
		if p.s[i] == '\\' {
			i++
			if i >= len(p.s) {
				return "", p.errorf("unclosed %s", name)
			}
			if raw {
				content = append(content, '\\', p.s[i])
				i++
				continue
			}
			switch p.s[i] {
			case 'a':
				i++
				content = append(content, '\a')
			case 'b':
				i++
				content = append(content, '\b')
			case 'f':
				i++
				content = append(content, '\f')
			case 'n':
				i++
				content = append(content, '\n')
			case 'r':
				i++
				content = append(content, '\r')
			case 't':
				i++
				content = append(content, '\t')
			case 'v':
				i++
				content = append(content, '\v')
			case '\\':
				i++
				content = append(content, '\\')
			case '?':
				i++
				content = append(content, '?')
			case '"':
				i++
				content = append(content, '"')
			case '\'':
				i++
				content = append(content, '\'')
			case '`':
				i++
				content = append(content, '`')
			case 'x', 'X':
				i++
				if !(i+1 < len(p.s) && isHexDigit(p.s[i]) && isHexDigit(p.s[i+1])) {
					return "", p.errorf("illegal escape sequence: hex escape sequence must be followed by 2 hex digits")
				}
				c, err := strconv.ParseUint(p.s[i:i+2], 16, 64)
				if err != nil {
					return "", p.errorf("illegal escape sequence: invalid hex digits: %q: %v", p.s[i:i+2], err)
				}
				content = append(content, byte(c))
				i += 2
			case 'u', 'U':
				t := p.s[i]
				if !unicode {
					return "", p.errorf("illegal escape sequence: \\%c", t)
				}
				i++
				size := 4
				if t == 'U' {
					size = 8
				}
				if i+size-1 >= len(p.s) {
					return "", p.errorf("illegal escape sequence: \\%c escape sequence must be followed by %d hex digits", t, size)
				}
				for j := 0; j < size; j++ {
					if !isHexDigit(p.s[i+j]) {
						return "", p.errorf("illegal escape sequence: \\%c escape sequence must be followed by %d hex digits", t, size)
					}
				}
				c, err := strconv.ParseUint(p.s[i:i+size], 16, 64)
				if err != nil {
					return "", p.errorf("illegal escape sequence: invalid \\%c digits: %q: %v", t, p.s[i:i+size], err)
				}
				if 0xD800 <= c && c <= 0xDFFF || 0x10FFFF < c {
					return "", p.errorf("illegal escape sequence: invalid codepoint: %x", c)
				}
				var buf [utf8.UTFMax]byte
				n := utf8.EncodeRune(buf[:], rune(c))
				content = append(content, buf[:n]...)
				i += size
			case '0', '1', '2', '3', '4', '5', '6', '7':
				if !(i+2 < len(p.s) && isOctalDigit(p.s[i+1]) && isOctalDigit(p.s[i+2])) {
					return "", p.errorf("illegal escape sequence: octal escape sequence must be followed by 3 octal digits")
				}
				c, err := strconv.ParseUint(p.s[i:i+3], 8, 64)
				if err != nil {
					return "", p.errorf("illegal escape sequence: invalid octal digits: %q: %v", p.s[i:i+3], err)
				}
				if c >= 256 {
					return "", p.errorf("illegal escape sequence: octal digits overflow: %q (%d)", p.s[i:i+3], c)
				}
				content = append(content, byte(c))
				i += 3
			default:
				return "", p.errorf("illegal escape sequence: \\%c", p.s[i])
			}
			continue
		}
		if p.s[i] == '\n' {
			if len(delim) != 3 {
				return "", p.errorf("newline forbidden in %s", name)
			}
			p.line++
		}
		content = append(content, p.s[i])
		i++
	}
	return "", p.errorf("unclosed %s", name)
}

var operators = map[string]bool{"-": true, "~": true, "*": true, "/": true, "||": true, "+": true, "<<": true, ">>": true, "&": true, "^": true, "|": true, "<": true, "<=": true, ">": true, ">=": true, "=": true, "!=": true, "<>": true}

func gologoo__isSpace_531485d04a38cff0d09d3c55685c367d(c byte) bool {
	switch c {
	case ' ', '\b', '\t', '\n':
		return true
	}
	return false
}
func (p *parser) gologoo__skipSpace_531485d04a38cff0d09d3c55685c367d() bool {
	initLine := p.line
	var com *comment
	i := 0
	for i < len(p.s) {
		if isSpace(p.s[i]) {
			if p.s[i] == '\n' {
				p.line++
			}
			i++
			continue
		}
		marker, term := "", ""
		if p.s[i] == '#' {
			marker, term = "#", "\n"
		} else if i+1 < len(p.s) && p.s[i] == '-' && p.s[i+1] == '-' {
			marker, term = "--", "\n"
		} else if i+1 < len(p.s) && p.s[i] == '/' && p.s[i+1] == '*' {
			marker, term = "/*", "*/"
		}
		if term == "" {
			break
		}
		ti := strings.Index(p.s[i+len(marker):], term)
		if ti < 0 {
			p.errorf("unterminated comment")
			return false
		}
		ti += len(marker)
		if com != nil && (com.end.Line+1 < p.line || com.marker != marker) {
			com = nil
		}
		if com == nil {
			p.comments = append(p.comments, comment{marker: marker, isolated: (p.line != initLine) || p.line == 1, start: Position{Line: p.line, Offset: p.offset + i}})
			com = &p.comments[len(p.comments)-1]
		}
		textLines := strings.Split(p.s[i+len(marker):i+ti], "\n")
		com.text = append(com.text, textLines...)
		com.end = Position{Line: p.line + len(textLines) - 1, Offset: p.offset + i + ti}
		p.line = com.end.Line
		if term == "\n" {
			p.line++
		}
		i += ti + len(term)
		if !com.isolated {
			com = nil
		}
	}
	p.s = p.s[i:]
	p.offset += i
	if p.s == "" {
		p.done = true
	}
	return i > 0
}
func (p *parser) gologoo__advance_531485d04a38cff0d09d3c55685c367d() {
	prevID := p.cur.typ == quotedID || p.cur.typ == unquotedID
	p.skipSpace()
	if p.done {
		return
	}
	if prevID && p.s[0] == '.' {
		p.cur.err = nil
		p.cur.line, p.cur.offset = p.line, p.offset
		p.cur.typ = unknownToken
		p.cur.value, p.s = p.s[:1], p.s[1:]
		p.offset++
		return
	}
	p.cur.err = nil
	p.cur.line, p.cur.offset = p.line, p.offset
	p.cur.typ = unknownToken
	switch p.s[0] {
	case ',', ';', '(', ')', '{', '}', '[', ']', '*', '+', '-':
		p.cur.value, p.s = p.s[:1], p.s[1:]
		p.offset++
		return
	case 'B', 'b', 'R', 'r', '"', '\'':
		raw, bytes := false, false
		for i := 0; i < 4 && i < len(p.s); i++ {
			switch {
			case !raw && (p.s[i] == 'R' || p.s[i] == 'r'):
				raw = true
				continue
			case !bytes && (p.s[i] == 'B' || p.s[i] == 'b'):
				bytes = true
				continue
			case p.s[i] == '"' || p.s[i] == '\'':
				switch {
				case raw && bytes:
					p.consumeRawBytes()
				case raw:
					p.consumeRawString()
				case bytes:
					p.consumeBytes()
				default:
					p.consumeString()
				}
				return
			}
			break
		}
	case '`':
		p.cur.string, p.cur.err = p.consumeStringContent("`", false, true, "quoted identifier")
		p.cur.typ = quotedID
		return
	}
	if p.s[0] == '@' || isInitialIdentifierChar(p.s[0]) {
		i := 1
		for i < len(p.s) && isIdentifierChar(p.s[i]) {
			i++
		}
		p.cur.value, p.s = p.s[:i], p.s[i:]
		p.cur.typ = unquotedID
		p.offset += i
		return
	}
	if len(p.s) >= 2 && p.s[0] == '.' && ('0' <= p.s[1] && p.s[1] <= '9') {
		p.consumeNumber()
		return
	}
	if '0' <= p.s[0] && p.s[0] <= '9' {
		p.consumeNumber()
		return
	}
	for i := 2; i >= 1; i-- {
		if i <= len(p.s) && operators[p.s[:i]] {
			p.cur.value, p.s = p.s[:i], p.s[i:]
			p.offset += i
			return
		}
	}
	p.errorf("unexpected byte %#x", p.s[0])
}
func (p *parser) gologoo__back_531485d04a38cff0d09d3c55685c367d() {
	if p.backed {
		panic("parser backed up twice")
	}
	p.done = false
	p.backed = true
	if p.cur.err != eof {
		p.cur.err = nil
	}
}
func (p *parser) gologoo__next_531485d04a38cff0d09d3c55685c367d() *token {
	if p.backed || p.done {
		p.backed = false
		return &p.cur
	}
	p.advance()
	if p.done && p.cur.err == nil {
		p.cur.value = ""
		p.cur.err = eof
	}
	debugf("parser·next(): returning [%v] [err: %v] @l%d,o%d", p.cur.value, p.cur.err, p.cur.line, p.cur.offset)
	return &p.cur
}
func (t *token) gologoo__caseEqual_531485d04a38cff0d09d3c55685c367d(x string) bool {
	return t.err == nil && t.typ != quotedID && strings.EqualFold(t.value, x)
}
func (p *parser) gologoo__sniff_531485d04a38cff0d09d3c55685c367d(want ...string) bool {
	orig := *p
	defer func() {
		*p = orig
	}()
	for _, w := range want {
		if !p.next().caseEqual(w) {
			return false
		}
	}
	return true
}
func (p *parser) gologoo__sniffTokenType_531485d04a38cff0d09d3c55685c367d(want tokenType) bool {
	orig := *p
	defer func() {
		*p = orig
	}()
	if p.next().typ == want {
		return true
	}
	return false
}
func (p *parser) gologoo__eat_531485d04a38cff0d09d3c55685c367d(want ...string) bool {
	orig := *p
	for _, w := range want {
		if !p.next().caseEqual(w) {
			*p = orig
			return false
		}
	}
	return true
}
func (p *parser) gologoo__expect_531485d04a38cff0d09d3c55685c367d(want ...string) *parseError {
	for _, w := range want {
		tok := p.next()
		if tok.err != nil {
			return tok.err
		}
		if !tok.caseEqual(w) {
			return p.errorf("got %q while expecting %q", tok.value, w)
		}
	}
	return nil
}
func (p *parser) gologoo__parseDDLStmt_531485d04a38cff0d09d3c55685c367d() (DDLStmt, *parseError) {
	debugf("parseDDLStmt: %v", p)
	if p.sniff("CREATE", "TABLE") {
		ct, err := p.parseCreateTable()
		return ct, err
	} else if p.sniff("CREATE", "INDEX") || p.sniff("CREATE", "UNIQUE", "INDEX") || p.sniff("CREATE", "NULL_FILTERED", "INDEX") || p.sniff("CREATE", "UNIQUE", "NULL_FILTERED", "INDEX") {
		ci, err := p.parseCreateIndex()
		return ci, err
	} else if p.sniff("CREATE", "VIEW") || p.sniff("CREATE", "OR", "REPLACE", "VIEW") {
		cv, err := p.parseCreateView()
		return cv, err
	} else if p.sniff("ALTER", "TABLE") {
		a, err := p.parseAlterTable()
		return a, err
	} else if p.eat("DROP") {
		pos := p.Pos()
		tok := p.next()
		if tok.err != nil {
			return nil, tok.err
		}
		switch {
		default:
			return nil, p.errorf("got %q, want TABLE, VIEW or INDEX", tok.value)
		case tok.caseEqual("TABLE"):
			name, err := p.parseTableOrIndexOrColumnName()
			if err != nil {
				return nil, err
			}
			return &DropTable{Name: name, Position: pos}, nil
		case tok.caseEqual("INDEX"):
			name, err := p.parseTableOrIndexOrColumnName()
			if err != nil {
				return nil, err
			}
			return &DropIndex{Name: name, Position: pos}, nil
		case tok.caseEqual("VIEW"):
			name, err := p.parseTableOrIndexOrColumnName()
			if err != nil {
				return nil, err
			}
			return &DropView{Name: name, Position: pos}, nil
		}
	} else if p.sniff("ALTER", "DATABASE") {
		a, err := p.parseAlterDatabase()
		return a, err
	}
	return nil, p.errorf("unknown DDL statement")
}
func (p *parser) gologoo__parseCreateTable_531485d04a38cff0d09d3c55685c367d() (*CreateTable, *parseError) {
	debugf("parseCreateTable: %v", p)
	if err := p.expect("CREATE"); err != nil {
		return nil, err
	}
	pos := p.Pos()
	if err := p.expect("TABLE"); err != nil {
		return nil, err
	}
	tname, err := p.parseTableOrIndexOrColumnName()
	if err != nil {
		return nil, err
	}
	ct := &CreateTable{Name: tname, Position: pos}
	err = p.parseCommaList("(", ")", func(p *parser) *parseError {
		if p.sniffTableConstraint() {
			tc, err := p.parseTableConstraint()
			if err != nil {
				return err
			}
			ct.Constraints = append(ct.Constraints, tc)
			return nil
		}
		cd, err := p.parseColumnDef()
		if err != nil {
			return err
		}
		ct.Columns = append(ct.Columns, cd)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := p.expect("PRIMARY"); err != nil {
		return nil, err
	}
	if err := p.expect("KEY"); err != nil {
		return nil, err
	}
	ct.PrimaryKey, err = p.parseKeyPartList()
	if err != nil {
		return nil, err
	}
	if p.eat(",", "INTERLEAVE") {
		if err := p.expect("IN"); err != nil {
			return nil, err
		}
		if err := p.expect("PARENT"); err != nil {
			return nil, err
		}
		pname, err := p.parseTableOrIndexOrColumnName()
		if err != nil {
			return nil, err
		}
		ct.Interleave = &Interleave{Parent: pname, OnDelete: NoActionOnDelete}
		if p.eat("ON", "DELETE") {
			od, err := p.parseOnDelete()
			if err != nil {
				return nil, err
			}
			ct.Interleave.OnDelete = od
		}
	}
	if p.eat(",", "ROW", "DELETION", "POLICY") {
		rdp, err := p.parseRowDeletionPolicy()
		if err != nil {
			return nil, err
		}
		ct.RowDeletionPolicy = &rdp
	}
	return ct, nil
}
func (p *parser) gologoo__sniffTableConstraint_531485d04a38cff0d09d3c55685c367d() bool {
	if p.sniff("FOREIGN", "KEY") || p.sniff("CHECK") {
		return true
	}
	orig := *p
	defer func() {
		*p = orig
	}()
	if !p.eat("CONSTRAINT") {
		return false
	}
	if _, err := p.parseTableOrIndexOrColumnName(); err != nil {
		return false
	}
	return p.sniff("FOREIGN") || p.sniff("CHECK")
}
func (p *parser) gologoo__parseCreateIndex_531485d04a38cff0d09d3c55685c367d() (*CreateIndex, *parseError) {
	debugf("parseCreateIndex: %v", p)
	var unique, nullFiltered bool
	if err := p.expect("CREATE"); err != nil {
		return nil, err
	}
	pos := p.Pos()
	if p.eat("UNIQUE") {
		unique = true
	}
	if p.eat("NULL_FILTERED") {
		nullFiltered = true
	}
	if err := p.expect("INDEX"); err != nil {
		return nil, err
	}
	iname, err := p.parseTableOrIndexOrColumnName()
	if err != nil {
		return nil, err
	}
	if err := p.expect("ON"); err != nil {
		return nil, err
	}
	tname, err := p.parseTableOrIndexOrColumnName()
	if err != nil {
		return nil, err
	}
	ci := &CreateIndex{Name: iname, Table: tname, Unique: unique, NullFiltered: nullFiltered, Position: pos}
	ci.Columns, err = p.parseKeyPartList()
	if err != nil {
		return nil, err
	}
	if p.eat("STORING") {
		ci.Storing, err = p.parseColumnNameList()
		if err != nil {
			return nil, err
		}
	}
	if p.eat(",", "INTERLEAVE", "IN") {
		ci.Interleave, err = p.parseTableOrIndexOrColumnName()
		if err != nil {
			return nil, err
		}
	}
	return ci, nil
}
func (p *parser) gologoo__parseCreateView_531485d04a38cff0d09d3c55685c367d() (*CreateView, *parseError) {
	debugf("parseCreateView: %v", p)
	var orReplace bool
	if err := p.expect("CREATE"); err != nil {
		return nil, err
	}
	pos := p.Pos()
	if p.eat("OR", "REPLACE") {
		orReplace = true
	}
	if err := p.expect("VIEW"); err != nil {
		return nil, err
	}
	vname, err := p.parseTableOrIndexOrColumnName()
	if err := p.expect("SQL", "SECURITY", "INVOKER", "AS"); err != nil {
		return nil, err
	}
	query, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	return &CreateView{Name: vname, OrReplace: orReplace, Query: query, Position: pos}, nil
}
func (p *parser) gologoo__parseAlterTable_531485d04a38cff0d09d3c55685c367d() (*AlterTable, *parseError) {
	debugf("parseAlterTable: %v", p)
	if err := p.expect("ALTER"); err != nil {
		return nil, err
	}
	pos := p.Pos()
	if err := p.expect("TABLE"); err != nil {
		return nil, err
	}
	tname, err := p.parseTableOrIndexOrColumnName()
	if err != nil {
		return nil, err
	}
	a := &AlterTable{Name: tname, Position: pos}
	tok := p.next()
	if tok.err != nil {
		return nil, tok.err
	}
	switch {
	default:
		return nil, p.errorf("got %q, expected ADD or DROP or SET or ALTER", tok.value)
	case tok.caseEqual("ADD"):
		if p.sniff("CONSTRAINT") || p.sniff("FOREIGN") || p.sniff("CHECK") {
			tc, err := p.parseTableConstraint()
			if err != nil {
				return nil, err
			}
			a.Alteration = AddConstraint{Constraint: tc}
			return a, nil
		}
		if p.eat("ROW", "DELETION", "POLICY") {
			rdp, err := p.parseRowDeletionPolicy()
			if err != nil {
				return nil, err
			}
			a.Alteration = AddRowDeletionPolicy{RowDeletionPolicy: rdp}
			return a, nil
		}
		if err := p.expect("COLUMN"); err != nil {
			return nil, err
		}
		cd, err := p.parseColumnDef()
		if err != nil {
			return nil, err
		}
		a.Alteration = AddColumn{Def: cd}
		return a, nil
	case tok.caseEqual("DROP"):
		if p.eat("CONSTRAINT") {
			name, err := p.parseTableOrIndexOrColumnName()
			if err != nil {
				return nil, err
			}
			a.Alteration = DropConstraint{Name: name}
			return a, nil
		}
		if p.eat("ROW", "DELETION", "POLICY") {
			a.Alteration = DropRowDeletionPolicy{}
			return a, nil
		}
		if err := p.expect("COLUMN"); err != nil {
			return nil, err
		}
		name, err := p.parseTableOrIndexOrColumnName()
		if err != nil {
			return nil, err
		}
		a.Alteration = DropColumn{Name: name}
		return a, nil
	case tok.caseEqual("SET"):
		if err := p.expect("ON"); err != nil {
			return nil, err
		}
		if err := p.expect("DELETE"); err != nil {
			return nil, err
		}
		od, err := p.parseOnDelete()
		if err != nil {
			return nil, err
		}
		a.Alteration = SetOnDelete{Action: od}
		return a, nil
	case tok.caseEqual("ALTER"):
		if err := p.expect("COLUMN"); err != nil {
			return nil, err
		}
		name, err := p.parseTableOrIndexOrColumnName()
		if err != nil {
			return nil, err
		}
		ca, err := p.parseColumnAlteration()
		if err != nil {
			return nil, err
		}
		a.Alteration = AlterColumn{Name: name, Alteration: ca}
		return a, nil
	case tok.caseEqual("REPLACE"):
		if p.eat("ROW", "DELETION", "POLICY") {
			rdp, err := p.parseRowDeletionPolicy()
			if err != nil {
				return nil, err
			}
			a.Alteration = ReplaceRowDeletionPolicy{RowDeletionPolicy: rdp}
			return a, nil
		}
	}
	return a, nil
}
func (p *parser) gologoo__parseAlterDatabase_531485d04a38cff0d09d3c55685c367d() (*AlterDatabase, *parseError) {
	debugf("parseAlterDatabase: %v", p)
	if err := p.expect("ALTER"); err != nil {
		return nil, err
	}
	pos := p.Pos()
	if err := p.expect("DATABASE"); err != nil {
		return nil, err
	}
	dbname, err := p.parseTableOrIndexOrColumnName()
	if err != nil {
		return nil, err
	}
	a := &AlterDatabase{Name: dbname, Position: pos}
	tok := p.next()
	if tok.err != nil {
		return nil, tok.err
	}
	switch {
	default:
		return nil, p.errorf("got %q, expected SET", tok.value)
	case tok.caseEqual("SET"):
		options, err := p.parseDatabaseOptions()
		if err != nil {
			return nil, err
		}
		a.Alteration = SetDatabaseOptions{Options: options}
		return a, nil
	}
}
func (p *parser) gologoo__parseDMLStmt_531485d04a38cff0d09d3c55685c367d() (DMLStmt, *parseError) {
	debugf("parseDMLStmt: %v", p)
	if p.eat("DELETE") {
		p.eat("FROM")
		tname, err := p.parseTableOrIndexOrColumnName()
		if err != nil {
			return nil, err
		}
		if err := p.expect("WHERE"); err != nil {
			return nil, err
		}
		where, err := p.parseBoolExpr()
		if err != nil {
			return nil, err
		}
		return &Delete{Table: tname, Where: where}, nil
	}
	if p.eat("UPDATE") {
		tname, err := p.parseTableOrIndexOrColumnName()
		if err != nil {
			return nil, err
		}
		u := &Update{Table: tname}
		if err := p.expect("SET"); err != nil {
			return nil, err
		}
		for {
			ui, err := p.parseUpdateItem()
			if err != nil {
				return nil, err
			}
			u.Items = append(u.Items, ui)
			if p.eat(",") {
				continue
			}
			break
		}
		if err := p.expect("WHERE"); err != nil {
			return nil, err
		}
		where, err := p.parseBoolExpr()
		if err != nil {
			return nil, err
		}
		u.Where = where
		return u, nil
	}
	return nil, p.errorf("unknown DML statement")
}
func (p *parser) gologoo__parseUpdateItem_531485d04a38cff0d09d3c55685c367d() (UpdateItem, *parseError) {
	col, err := p.parseTableOrIndexOrColumnName()
	if err != nil {
		return UpdateItem{}, err
	}
	ui := UpdateItem{Column: col}
	if err := p.expect("="); err != nil {
		return UpdateItem{}, err
	}
	if p.eat("DEFAULT") {
		return ui, nil
	}
	ui.Value, err = p.parseExpr()
	if err != nil {
		return UpdateItem{}, err
	}
	return ui, nil
}
func (p *parser) gologoo__parseColumnDef_531485d04a38cff0d09d3c55685c367d() (ColumnDef, *parseError) {
	debugf("parseColumnDef: %v", p)
	name, err := p.parseTableOrIndexOrColumnName()
	if err != nil {
		return ColumnDef{}, err
	}
	cd := ColumnDef{Name: name, Position: p.Pos()}
	cd.Type, err = p.parseType()
	if err != nil {
		return ColumnDef{}, err
	}
	if p.eat("NOT", "NULL") {
		cd.NotNull = true
	}
	if p.eat("AS", "(") {
		cd.Generated, err = p.parseExpr()
		if err != nil {
			return ColumnDef{}, err
		}
		if err := p.expect(")"); err != nil {
			return ColumnDef{}, err
		}
		if err := p.expect("STORED"); err != nil {
			return ColumnDef{}, err
		}
	}
	if p.sniff("OPTIONS") {
		cd.Options, err = p.parseColumnOptions()
		if err != nil {
			return ColumnDef{}, err
		}
	}
	return cd, nil
}
func (p *parser) gologoo__parseColumnAlteration_531485d04a38cff0d09d3c55685c367d() (ColumnAlteration, *parseError) {
	debugf("parseColumnAlteration: %v", p)
	if p.eat("SET") {
		co, err := p.parseColumnOptions()
		if err != nil {
			return nil, err
		}
		return SetColumnOptions{Options: co}, nil
	}
	typ, err := p.parseType()
	if err != nil {
		return nil, err
	}
	sct := SetColumnType{Type: typ}
	if p.eat("NOT", "NULL") {
		sct.NotNull = true
	}
	return sct, nil
}
func (p *parser) gologoo__parseColumnOptions_531485d04a38cff0d09d3c55685c367d() (ColumnOptions, *parseError) {
	debugf("parseColumnOptions: %v", p)
	if err := p.expect("OPTIONS"); err != nil {
		return ColumnOptions{}, err
	}
	if err := p.expect("("); err != nil {
		return ColumnOptions{}, err
	}
	var co ColumnOptions
	if p.eat("allow_commit_timestamp", "=") {
		tok := p.next()
		if tok.err != nil {
			return ColumnOptions{}, tok.err
		}
		allowCommitTimestamp := new(bool)
		switch tok.value {
		case "true":
			*allowCommitTimestamp = true
		case "null":
			*allowCommitTimestamp = false
		default:
			return ColumnOptions{}, p.errorf("got %q, want true or null", tok.value)
		}
		co.AllowCommitTimestamp = allowCommitTimestamp
	}
	if err := p.expect(")"); err != nil {
		return ColumnOptions{}, err
	}
	return co, nil
}
func (p *parser) gologoo__parseDatabaseOptions_531485d04a38cff0d09d3c55685c367d() (DatabaseOptions, *parseError) {
	debugf("parseDatabaseOptions: %v", p)
	if err := p.expect("OPTIONS"); err != nil {
		return DatabaseOptions{}, err
	}
	if err := p.expect("("); err != nil {
		return DatabaseOptions{}, err
	}
	var opts DatabaseOptions
	for {
		if p.eat("enable_key_visualizer", "=") {
			tok := p.next()
			if tok.err != nil {
				return DatabaseOptions{}, tok.err
			}
			enableKeyVisualizer := new(bool)
			switch tok.value {
			case "true":
				*enableKeyVisualizer = true
			case "null":
				*enableKeyVisualizer = false
			default:
				return DatabaseOptions{}, p.errorf("invalid enable_key_visualizer_value: %v", tok.value)
			}
			opts.EnableKeyVisualizer = enableKeyVisualizer
		} else if p.eat("optimizer_version", "=") {
			tok := p.next()
			if tok.err != nil {
				return DatabaseOptions{}, tok.err
			}
			optimizerVersion := new(int)
			if tok.value == "null" {
				*optimizerVersion = 0
			} else {
				if tok.typ != int64Token {
					return DatabaseOptions{}, p.errorf("invalid optimizer_version value: %v", tok.value)
				}
				version, err := strconv.Atoi(tok.value)
				if err != nil {
					return DatabaseOptions{}, p.errorf("invalid optimizer_version value: %v", tok.value)
				}
				optimizerVersion = &version
			}
			opts.OptimizerVersion = optimizerVersion
		} else if p.eat("version_retention_period", "=") {
			tok := p.next()
			if tok.err != nil {
				return DatabaseOptions{}, tok.err
			}
			retentionPeriod := new(string)
			if tok.value == "null" {
				*retentionPeriod = ""
			} else {
				if tok.typ != stringToken {
					return DatabaseOptions{}, p.errorf("invalid version_retention_period: %v", tok.value)
				}
				retentionPeriod = &tok.string
			}
			opts.VersionRetentionPeriod = retentionPeriod
		} else {
			tok := p.next()
			return DatabaseOptions{}, p.errorf("unknown database option: %v", tok.value)
		}
		if p.sniff(")") {
			break
		}
		if !p.eat(",") {
			return DatabaseOptions{}, p.errorf("missing ',' in options list")
		}
	}
	if err := p.expect(")"); err != nil {
		return DatabaseOptions{}, err
	}
	return opts, nil
}
func (p *parser) gologoo__parseKeyPartList_531485d04a38cff0d09d3c55685c367d() ([]KeyPart, *parseError) {
	var list []KeyPart
	err := p.parseCommaList("(", ")", func(p *parser) *parseError {
		kp, err := p.parseKeyPart()
		if err != nil {
			return err
		}
		list = append(list, kp)
		return nil
	})
	return list, err
}
func (p *parser) gologoo__parseKeyPart_531485d04a38cff0d09d3c55685c367d() (KeyPart, *parseError) {
	debugf("parseKeyPart: %v", p)
	name, err := p.parseTableOrIndexOrColumnName()
	if err != nil {
		return KeyPart{}, err
	}
	kp := KeyPart{Column: name}
	if p.eat("ASC") {
	} else if p.eat("DESC") {
		kp.Desc = true
	}
	return kp, nil
}
func (p *parser) gologoo__parseTableConstraint_531485d04a38cff0d09d3c55685c367d() (TableConstraint, *parseError) {
	debugf("parseTableConstraint: %v", p)
	if p.eat("CONSTRAINT") {
		pos := p.Pos()
		cname, err := p.parseTableOrIndexOrColumnName()
		if err != nil {
			return TableConstraint{}, err
		}
		c, err := p.parseConstraint()
		if err != nil {
			return TableConstraint{}, err
		}
		return TableConstraint{Name: cname, Constraint: c, Position: pos}, nil
	}
	c, err := p.parseConstraint()
	if err != nil {
		return TableConstraint{}, err
	}
	return TableConstraint{Constraint: c, Position: c.Pos()}, nil
}
func (p *parser) gologoo__parseConstraint_531485d04a38cff0d09d3c55685c367d() (Constraint, *parseError) {
	if p.sniff("FOREIGN") {
		fk, err := p.parseForeignKey()
		return fk, err
	}
	c, err := p.parseCheck()
	return c, err
}
func (p *parser) gologoo__parseForeignKey_531485d04a38cff0d09d3c55685c367d() (ForeignKey, *parseError) {
	debugf("parseForeignKey: %v", p)
	if err := p.expect("FOREIGN"); err != nil {
		return ForeignKey{}, err
	}
	fk := ForeignKey{Position: p.Pos()}
	if err := p.expect("KEY"); err != nil {
		return ForeignKey{}, err
	}
	var err *parseError
	fk.Columns, err = p.parseColumnNameList()
	if err != nil {
		return ForeignKey{}, err
	}
	if err := p.expect("REFERENCES"); err != nil {
		return ForeignKey{}, err
	}
	fk.RefTable, err = p.parseTableOrIndexOrColumnName()
	if err != nil {
		return ForeignKey{}, err
	}
	fk.RefColumns, err = p.parseColumnNameList()
	if err != nil {
		return ForeignKey{}, err
	}
	return fk, nil
}
func (p *parser) gologoo__parseCheck_531485d04a38cff0d09d3c55685c367d() (Check, *parseError) {
	debugf("parseCheck: %v", p)
	if err := p.expect("CHECK"); err != nil {
		return Check{}, err
	}
	c := Check{Position: p.Pos()}
	if err := p.expect("("); err != nil {
		return Check{}, err
	}
	var err *parseError
	c.Expr, err = p.parseBoolExpr()
	if err != nil {
		return Check{}, err
	}
	if err := p.expect(")"); err != nil {
		return Check{}, err
	}
	return c, nil
}
func (p *parser) gologoo__parseColumnNameList_531485d04a38cff0d09d3c55685c367d() ([]ID, *parseError) {
	var list []ID
	err := p.parseCommaList("(", ")", func(p *parser) *parseError {
		n, err := p.parseTableOrIndexOrColumnName()
		if err != nil {
			return err
		}
		list = append(list, n)
		return nil
	})
	return list, err
}

var baseTypes = map[string]TypeBase{"BOOL": Bool, "INT64": Int64, "FLOAT64": Float64, "NUMERIC": Numeric, "STRING": String, "BYTES": Bytes, "DATE": Date, "TIMESTAMP": Timestamp, "JSON": JSON}

func (p *parser) gologoo__parseBaseType_531485d04a38cff0d09d3c55685c367d() (Type, *parseError) {
	return p.parseBaseOrParameterizedType(false)
}
func (p *parser) gologoo__parseType_531485d04a38cff0d09d3c55685c367d() (Type, *parseError) {
	return p.parseBaseOrParameterizedType(true)
}

var extractPartTypes = map[string]TypeBase{"DAY": Int64, "MONTH": Int64, "YEAR": Int64, "DATE": Date}

func (p *parser) gologoo__parseExtractType_531485d04a38cff0d09d3c55685c367d() (Type, string, *parseError) {
	var t Type
	tok := p.next()
	if tok.err != nil {
		return Type{}, "", tok.err
	}
	base, ok := extractPartTypes[strings.ToUpper(tok.value)]
	if !ok {
		return Type{}, "", p.errorf("got %q, want valid EXTRACT types", tok.value)
	}
	t.Base = base
	return t, strings.ToUpper(tok.value), nil
}
func (p *parser) gologoo__parseBaseOrParameterizedType_531485d04a38cff0d09d3c55685c367d(withParam bool) (Type, *parseError) {
	debugf("parseBaseOrParameterizedType: %v", p)
	var t Type
	tok := p.next()
	if tok.err != nil {
		return Type{}, tok.err
	}
	if tok.caseEqual("ARRAY") {
		t.Array = true
		if err := p.expect("<"); err != nil {
			return Type{}, err
		}
		tok = p.next()
		if tok.err != nil {
			return Type{}, tok.err
		}
	}
	base, ok := baseTypes[strings.ToUpper(tok.value)]
	if !ok {
		return Type{}, p.errorf("got %q, want scalar type", tok.value)
	}
	t.Base = base
	if withParam && (t.Base == String || t.Base == Bytes) {
		if err := p.expect("("); err != nil {
			return Type{}, err
		}
		tok = p.next()
		if tok.err != nil {
			return Type{}, tok.err
		}
		if tok.caseEqual("MAX") {
			t.Len = MaxLen
		} else if tok.typ == int64Token {
			n, err := strconv.ParseInt(tok.value, tok.int64Base, 64)
			if err != nil {
				return Type{}, p.errorf("%v", err)
			}
			t.Len = n
		} else {
			return Type{}, p.errorf("got %q, want MAX or int64", tok.value)
		}
		if err := p.expect(")"); err != nil {
			return Type{}, err
		}
	}
	if t.Array {
		if err := p.expect(">"); err != nil {
			return Type{}, err
		}
	}
	return t, nil
}
func (p *parser) gologoo__parseQuery_531485d04a38cff0d09d3c55685c367d() (Query, *parseError) {
	debugf("parseQuery: %v", p)
	if err := p.expect("SELECT"); err != nil {
		return Query{}, err
	}
	p.back()
	sel, err := p.parseSelect()
	if err != nil {
		return Query{}, err
	}
	q := Query{Select: sel}
	if p.eat("ORDER", "BY") {
		for {
			o, err := p.parseOrder()
			if err != nil {
				return Query{}, err
			}
			q.Order = append(q.Order, o)
			if !p.eat(",") {
				break
			}
		}
	}
	if p.eat("LIMIT") {
		lim, err := p.parseLiteralOrParam()
		if err != nil {
			return Query{}, err
		}
		q.Limit = lim
		if p.eat("OFFSET") {
			off, err := p.parseLiteralOrParam()
			if err != nil {
				return Query{}, err
			}
			q.Offset = off
		}
	}
	return q, nil
}
func (p *parser) gologoo__parseSelect_531485d04a38cff0d09d3c55685c367d() (Select, *parseError) {
	debugf("parseSelect: %v", p)
	if err := p.expect("SELECT"); err != nil {
		return Select{}, err
	}
	var sel Select
	if p.eat("ALL") {
	} else if p.eat("DISTINCT") {
		sel.Distinct = true
	}
	list, aliases, err := p.parseSelectList()
	if err != nil {
		return Select{}, err
	}
	sel.List, sel.ListAliases = list, aliases
	if p.eat("FROM") {
		padTS := func() {
			for len(sel.TableSamples) < len(sel.From) {
				sel.TableSamples = append(sel.TableSamples, nil)
			}
		}
		for {
			from, err := p.parseSelectFrom()
			if err != nil {
				return Select{}, err
			}
			sel.From = append(sel.From, from)
			if p.sniff("TABLESAMPLE") {
				ts, err := p.parseTableSample()
				if err != nil {
					return Select{}, err
				}
				padTS()
				sel.TableSamples[len(sel.TableSamples)-1] = &ts
			}
			if p.eat(",") {
				continue
			}
			break
		}
		if sel.TableSamples != nil {
			padTS()
		}
	}
	if p.eat("WHERE") {
		where, err := p.parseBoolExpr()
		if err != nil {
			return Select{}, err
		}
		sel.Where = where
	}
	if p.eat("GROUP", "BY") {
		list, err := p.parseExprList()
		if err != nil {
			return Select{}, err
		}
		sel.GroupBy = list
	}
	return sel, nil
}
func (p *parser) gologoo__parseSelectList_531485d04a38cff0d09d3c55685c367d() ([]Expr, []ID, *parseError) {
	var list []Expr
	var aliases []ID
	padAliases := func() {
		for len(aliases) < len(list) {
			aliases = append(aliases, "")
		}
	}
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, nil, err
		}
		list = append(list, expr)
		if p.eat("AS") {
			alias, err := p.parseAlias()
			if err != nil {
				return nil, nil, err
			}
			padAliases()
			aliases[len(aliases)-1] = alias
		}
		if p.eat(",") {
			continue
		}
		break
	}
	if aliases != nil {
		padAliases()
	}
	return list, aliases, nil
}
func (p *parser) gologoo__parseSelectFromTable_531485d04a38cff0d09d3c55685c367d() (SelectFrom, *parseError) {
	if p.eat("UNNEST") {
		if err := p.expect("("); err != nil {
			return nil, err
		}
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if err := p.expect(")"); err != nil {
			return nil, err
		}
		sfu := SelectFromUnnest{Expr: e}
		if p.eat("AS") {
			alias, err := p.parseAlias()
			if err != nil {
				return nil, err
			}
			sfu.Alias = alias
		}
		return sfu, nil
	}
	tname, err := p.parseTableOrIndexOrColumnName()
	if err != nil {
		return nil, err
	}
	sf := SelectFromTable{Table: tname}
	if p.eat("@") {
		hints, err := p.parseHints(map[string]string{})
		if err != nil {
			return nil, err
		}
		sf.Hints = hints
	}
	if p.eat("AS") {
		alias, err := p.parseAlias()
		if err != nil {
			return nil, err
		}
		sf.Alias = alias
	}
	return sf, nil
}
func (p *parser) gologoo__parseSelectFromJoin_531485d04a38cff0d09d3c55685c367d(lhs SelectFrom) (SelectFrom, *parseError) {
	tok := p.next()
	if tok.err != nil {
		p.back()
		return nil, nil
	}
	var hashJoin bool
	if tok.caseEqual("HASH") {
		hashJoin = true
		tok = p.next()
		if tok.err != nil {
			return nil, tok.err
		}
	}
	var jt JoinType
	if tok.caseEqual("JOIN") {
		jt = InnerJoin
	} else if j, ok := joinKeywords[tok.value]; ok {
		jt = j
		switch jt {
		case FullJoin, LeftJoin, RightJoin:
			p.eat("OUTER")
		}
		if err := p.expect("JOIN"); err != nil {
			return nil, err
		}
	} else {
		p.back()
		return nil, nil
	}
	sfj := SelectFromJoin{Type: jt, LHS: lhs}
	var hints map[string]string
	if hashJoin {
		hints = map[string]string{}
		hints["JOIN_METHOD"] = "HASH_JOIN"
	}
	if p.eat("@") {
		h, err := p.parseHints(hints)
		if err != nil {
			return nil, err
		}
		hints = h
	}
	sfj.Hints = hints
	rhs, err := p.parseSelectFromTable()
	if err != nil {
		return nil, err
	}
	sfj.RHS = rhs
	if p.eat("ON") {
		sfj.On, err = p.parseBoolExpr()
		if err != nil {
			return nil, err
		}
	}
	if p.eat("USING") {
		if sfj.On != nil {
			return nil, p.errorf("join may not have both ON and USING clauses")
		}
		sfj.Using, err = p.parseColumnNameList()
		if err != nil {
			return nil, err
		}
	}
	return sfj, nil
}
func (p *parser) gologoo__parseSelectFrom_531485d04a38cff0d09d3c55685c367d() (SelectFrom, *parseError) {
	debugf("parseSelectFrom: %v", p)
	leftHandSide, err := p.parseSelectFromTable()
	if err != nil {
		return nil, err
	}
	for {
		sfj, err := p.parseSelectFromJoin(leftHandSide)
		if err != nil {
			return nil, err
		}
		if sfj == nil {
			break
		}
		leftHandSide = sfj
	}
	return leftHandSide, nil
}

var joinKeywords = map[string]JoinType{"INNER": InnerJoin, "CROSS": CrossJoin, "FULL": FullJoin, "LEFT": LeftJoin, "RIGHT": RightJoin}

func (p *parser) gologoo__parseTableSample_531485d04a38cff0d09d3c55685c367d() (TableSample, *parseError) {
	var ts TableSample
	if err := p.expect("TABLESAMPLE"); err != nil {
		return ts, err
	}
	tok := p.next()
	switch {
	case tok.err != nil:
		return ts, tok.err
	case tok.caseEqual("BERNOULLI"):
		ts.Method = Bernoulli
	case tok.caseEqual("RESERVOIR"):
		ts.Method = Reservoir
	default:
		return ts, p.errorf("got %q, want BERNOULLI or RESERVOIR", tok.value)
	}
	if err := p.expect("("); err != nil {
		return ts, err
	}
	size, err := p.parseExpr()
	if err != nil {
		return ts, err
	}
	ts.Size = size
	tok = p.next()
	switch {
	case tok.err != nil:
		return ts, tok.err
	case tok.caseEqual("PERCENT"):
		ts.SizeType = PercentTableSample
	case tok.caseEqual("ROWS"):
		ts.SizeType = RowsTableSample
	default:
		return ts, p.errorf("got %q, want PERCENT or ROWS", tok.value)
	}
	if err := p.expect(")"); err != nil {
		return ts, err
	}
	return ts, nil
}
func (p *parser) gologoo__parseOrder_531485d04a38cff0d09d3c55685c367d() (Order, *parseError) {
	expr, err := p.parseExpr()
	if err != nil {
		return Order{}, err
	}
	o := Order{Expr: expr}
	if p.eat("ASC") {
	} else if p.eat("DESC") {
		o.Desc = true
	}
	return o, nil
}
func (p *parser) gologoo__parseLiteralOrParam_531485d04a38cff0d09d3c55685c367d() (LiteralOrParam, *parseError) {
	tok := p.next()
	if tok.err != nil {
		return nil, tok.err
	}
	if tok.typ == int64Token {
		n, err := strconv.ParseInt(tok.value, tok.int64Base, 64)
		if err != nil {
			return nil, p.errorf("%v", err)
		}
		return IntegerLiteral(n), nil
	}
	if strings.HasPrefix(tok.value, "@") {
		return Param(tok.value[1:]), nil
	}
	return nil, p.errorf("got %q, want literal or parameter", tok.value)
}
func (p *parser) gologoo__parseExprList_531485d04a38cff0d09d3c55685c367d() ([]Expr, *parseError) {
	var list []Expr
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		list = append(list, expr)
		if p.eat(",") {
			continue
		}
		break
	}
	return list, nil
}
func (p *parser) gologoo__parseParenExprList_531485d04a38cff0d09d3c55685c367d() ([]Expr, *parseError) {
	return p.parseParenExprListWithParseFunc(func(p *parser) (Expr, *parseError) {
		return p.parseExpr()
	})
}
func (p *parser) gologoo__parseParenExprListWithParseFunc_531485d04a38cff0d09d3c55685c367d(f func(*parser) (Expr, *parseError)) ([]Expr, *parseError) {
	var list []Expr
	err := p.parseCommaList("(", ")", func(p *parser) *parseError {
		e, err := f(p)
		if err != nil {
			return err
		}
		list = append(list, e)
		return nil
	})
	return list, err
}

var typedArgParser = func(p *parser) (Expr, *parseError) {
	e, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if err := p.expect("AS"); err != nil {
		return nil, err
	}
	toType, err := p.parseBaseType()
	if err != nil {
		return nil, err
	}
	return TypedExpr{Expr: e, Type: toType}, nil
}
var extractArgParser = func(p *parser) (Expr, *parseError) {
	partType, part, err := p.parseExtractType()
	if err != nil {
		return nil, err
	}
	if err := p.expect("FROM"); err != nil {
		return nil, err
	}
	e, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.eat("AT", "TIME", "ZONE") {
		tok := p.next()
		if tok.err != nil {
			return nil, err
		}
		return ExtractExpr{Part: part, Type: partType, Expr: AtTimeZoneExpr{Expr: e, Zone: tok.string, Type: Type{Base: Timestamp}}}, nil
	}
	return ExtractExpr{Part: part, Expr: e, Type: partType}, nil
}

func (p *parser) gologoo__parseExpr_531485d04a38cff0d09d3c55685c367d() (Expr, *parseError) {
	debugf("parseExpr: %v", p)
	return orParser.parse(p)
}

type binOpParser struct {
	LHS, RHS func(*parser) (Expr, *parseError)
	Op       string
	ArgCheck func(Expr) error
	Combiner func(lhs, rhs Expr) Expr
}

func (bin binOpParser) gologoo__parse_531485d04a38cff0d09d3c55685c367d(p *parser) (Expr, *parseError) {
	expr, err := bin.LHS(p)
	if err != nil {
		return nil, err
	}
	for {
		if !p.eat(bin.Op) {
			break
		}
		rhs, err := bin.RHS(p)
		if err != nil {
			return nil, err
		}
		if bin.ArgCheck != nil {
			if err := bin.ArgCheck(expr); err != nil {
				return nil, p.errorf("%v", err)
			}
			if err := bin.ArgCheck(rhs); err != nil {
				return nil, p.errorf("%v", err)
			}
		}
		expr = bin.Combiner(expr, rhs)
	}
	return expr, nil
}
func gologoo__init_531485d04a38cff0d09d3c55685c367d() {
	orParser = orParserShim
}

var (
	boolExprCheck = func(expr Expr) error {
		if _, ok := expr.(BoolExpr); !ok {
			return fmt.Errorf("got %T, want a boolean expression", expr)
		}
		return nil
	}
	orParser     binOpParser
	orParserShim = binOpParser{LHS: andParser.parse, RHS: andParser.parse, Op: "OR", ArgCheck: boolExprCheck, Combiner: func(lhs, rhs Expr) Expr {
		return LogicalOp{LHS: lhs.(BoolExpr), Op: Or, RHS: rhs.(BoolExpr)}
	}}
	andParser = binOpParser{LHS: (*parser).parseLogicalNot, RHS: (*parser).parseLogicalNot, Op: "AND", ArgCheck: boolExprCheck, Combiner: func(lhs, rhs Expr) Expr {
		return LogicalOp{LHS: lhs.(BoolExpr), Op: And, RHS: rhs.(BoolExpr)}
	}}
	bitOrParser  = newBinArithParser("|", BitOr, bitXorParser.parse)
	bitXorParser = newBinArithParser("^", BitXor, bitAndParser.parse)
	bitAndParser = newBinArithParser("&", BitAnd, bitShrParser.parse)
	bitShrParser = newBinArithParser(">>", BitShr, bitShlParser.parse)
	bitShlParser = newBinArithParser("<<", BitShl, subParser.parse)
	subParser    = newBinArithParser("-", Sub, addParser.parse)
	addParser    = newBinArithParser("+", Add, concatParser.parse)
	concatParser = newBinArithParser("||", Concat, divParser.parse)
	divParser    = newBinArithParser("/", Div, mulParser.parse)
	mulParser    = newBinArithParser("*", Mul, (*parser).parseUnaryArithOp)
)

func gologoo__newBinArithParser_531485d04a38cff0d09d3c55685c367d(opStr string, op ArithOperator, nextPrec func(*parser) (Expr, *parseError)) binOpParser {
	return binOpParser{LHS: nextPrec, RHS: nextPrec, Op: opStr, Combiner: func(lhs, rhs Expr) Expr {
		return ArithOp{LHS: lhs, Op: op, RHS: rhs}
	}}
}
func (p *parser) gologoo__parseLogicalNot_531485d04a38cff0d09d3c55685c367d() (Expr, *parseError) {
	if !p.eat("NOT") {
		return p.parseIsOp()
	}
	be, err := p.parseBoolExpr()
	if err != nil {
		return nil, err
	}
	return LogicalOp{Op: Not, RHS: be}, nil
}
func (p *parser) gologoo__parseIsOp_531485d04a38cff0d09d3c55685c367d() (Expr, *parseError) {
	debugf("parseIsOp: %v", p)
	expr, err := p.parseInOp()
	if err != nil {
		return nil, err
	}
	if !p.eat("IS") {
		return expr, nil
	}
	isOp := IsOp{LHS: expr}
	if p.eat("NOT") {
		isOp.Neg = true
	}
	tok := p.next()
	if tok.err != nil {
		return nil, tok.err
	}
	switch {
	case tok.caseEqual("NULL"):
		isOp.RHS = Null
	case tok.caseEqual("TRUE"):
		isOp.RHS = True
	case tok.caseEqual("FALSE"):
		isOp.RHS = False
	default:
		return nil, p.errorf("got %q, want NULL or TRUE or FALSE", tok.value)
	}
	return isOp, nil
}
func (p *parser) gologoo__parseInOp_531485d04a38cff0d09d3c55685c367d() (Expr, *parseError) {
	debugf("parseInOp: %v", p)
	expr, err := p.parseComparisonOp()
	if err != nil {
		return nil, err
	}
	inOp := InOp{LHS: expr}
	if p.eat("NOT", "IN") {
		inOp.Neg = true
	} else if p.eat("IN") {
	} else {
		return expr, nil
	}
	if p.eat("UNNEST") {
		inOp.Unnest = true
	}
	inOp.RHS, err = p.parseParenExprList()
	if err != nil {
		return nil, err
	}
	return inOp, nil
}

var symbolicOperators = map[string]ComparisonOperator{"<": Lt, "<=": Le, ">": Gt, ">=": Ge, "=": Eq, "!=": Ne, "<>": Ne}

func (p *parser) gologoo__parseComparisonOp_531485d04a38cff0d09d3c55685c367d() (Expr, *parseError) {
	debugf("parseComparisonOp: %v", p)
	expr, err := p.parseArithOp()
	if err != nil {
		return nil, err
	}
	for {
		var op ComparisonOperator
		var rhs2 bool
		if p.eat("NOT", "LIKE") {
			op = NotLike
		} else if p.eat("NOT", "BETWEEN") {
			op, rhs2 = NotBetween, true
		} else if p.eat("LIKE") {
			op = Like
		} else if p.eat("BETWEEN") {
			op, rhs2 = Between, true
		} else {
			tok := p.next()
			if tok.err != nil {
				p.back()
				break
			}
			var ok bool
			op, ok = symbolicOperators[tok.value]
			if !ok {
				p.back()
				break
			}
		}
		rhs, err := p.parseArithOp()
		if err != nil {
			return nil, err
		}
		co := ComparisonOp{LHS: expr, Op: op, RHS: rhs}
		if rhs2 {
			if err := p.expect("AND"); err != nil {
				return nil, err
			}
			rhs2, err := p.parseArithOp()
			if err != nil {
				return nil, err
			}
			co.RHS2 = rhs2
		}
		expr = co
	}
	return expr, nil
}
func (p *parser) gologoo__parseArithOp_531485d04a38cff0d09d3c55685c367d() (Expr, *parseError) {
	return bitOrParser.parse(p)
}

var unaryArithOperators = map[string]ArithOperator{"-": Neg, "~": BitNot, "+": Plus}

func (p *parser) gologoo__parseUnaryArithOp_531485d04a38cff0d09d3c55685c367d() (Expr, *parseError) {
	tok := p.next()
	if tok.err != nil {
		return nil, tok.err
	}
	op := tok.value
	if op == "-" || op == "+" {
		ntok := p.next()
		if ntok.err == nil {
			switch ntok.typ {
			case int64Token:
				comb := op + ntok.value
				n, err := strconv.ParseInt(comb, ntok.int64Base, 64)
				if err != nil {
					return nil, p.errorf("%v", err)
				}
				return IntegerLiteral(n), nil
			case float64Token:
				f := ntok.float64
				if op == "-" {
					f = -f
				}
				return FloatLiteral(f), nil
			}
		}
		p.back()
	}
	if op, ok := unaryArithOperators[op]; ok {
		e, err := p.parseLit()
		if err != nil {
			return nil, err
		}
		return ArithOp{Op: op, RHS: e}, nil
	}
	p.back()
	return p.parseLit()
}
func (p *parser) gologoo__parseLit_531485d04a38cff0d09d3c55685c367d() (Expr, *parseError) {
	tok := p.next()
	if tok.err != nil {
		return nil, tok.err
	}
	switch tok.typ {
	case int64Token:
		n, err := strconv.ParseInt(tok.value, tok.int64Base, 64)
		if err != nil {
			return nil, p.errorf("%v", err)
		}
		return IntegerLiteral(n), nil
	case float64Token:
		return FloatLiteral(tok.float64), nil
	case stringToken:
		return StringLiteral(tok.string), nil
	case bytesToken:
		return BytesLiteral(tok.string), nil
	}
	if tok.value == "(" {
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if err := p.expect(")"); err != nil {
			return nil, err
		}
		return Paren{Expr: e}, nil
	}
	if name := strings.ToUpper(tok.value); funcs[name] && p.sniff("(") {
		var list []Expr
		var err *parseError
		if f, ok := funcArgParsers[name]; ok {
			list, err = p.parseParenExprListWithParseFunc(f)
		} else {
			list, err = p.parseParenExprList()
		}
		if err != nil {
			return nil, err
		}
		return Func{Name: name, Args: list}, nil
	}
	switch {
	case tok.caseEqual("TRUE"):
		return True, nil
	case tok.caseEqual("FALSE"):
		return False, nil
	case tok.caseEqual("NULL"):
		return Null, nil
	case tok.value == "*":
		return Star, nil
	default:
	}
	switch {
	case tok.caseEqual("ARRAY") || tok.value == "[":
		p.back()
		return p.parseArrayLit()
	case tok.caseEqual("DATE"):
		if p.sniffTokenType(stringToken) {
			p.back()
			return p.parseDateLit()
		}
	case tok.caseEqual("TIMESTAMP"):
		if p.sniffTokenType(stringToken) {
			p.back()
			return p.parseTimestampLit()
		}
	}
	if strings.HasPrefix(tok.value, "@") {
		return Param(tok.value[1:]), nil
	}
	p.back()
	pe, err := p.parsePathExp()
	if err != nil {
		return nil, err
	}
	if len(pe) == 1 {
		return pe[0], nil
	}
	return pe, nil
}
func (p *parser) gologoo__parseArrayLit_531485d04a38cff0d09d3c55685c367d() (Array, *parseError) {
	p.eat("ARRAY")
	var arr Array
	err := p.parseCommaList("[", "]", func(p *parser) *parseError {
		e, err := p.parseLit()
		if err != nil {
			return err
		}
		arr = append(arr, e)
		return nil
	})
	return arr, err
}
func (p *parser) gologoo__parseDateLit_531485d04a38cff0d09d3c55685c367d() (DateLiteral, *parseError) {
	if err := p.expect("DATE"); err != nil {
		return DateLiteral{}, err
	}
	s, err := p.parseStringLit()
	if err != nil {
		return DateLiteral{}, err
	}
	d, perr := civil.ParseDate(string(s))
	if perr != nil {
		return DateLiteral{}, p.errorf("bad date literal %q: %v", s, perr)
	}
	return DateLiteral(d), nil
}

var timestampFormats = []string{"2006-01-02", "2006-01-02 15:04:05", "2006-01-02 15:04:05.000000", "2006-01-02 15:04:05 -07:00", "2006-01-02 15:04:05.000000 -07:00"}
var defaultLocation = func() *time.Location {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err == nil {
		return loc
	}
	return time.UTC
}()

func (p *parser) gologoo__parseTimestampLit_531485d04a38cff0d09d3c55685c367d() (TimestampLiteral, *parseError) {
	if err := p.expect("TIMESTAMP"); err != nil {
		return TimestampLiteral{}, err
	}
	s, err := p.parseStringLit()
	if err != nil {
		return TimestampLiteral{}, err
	}
	for _, format := range timestampFormats {
		t, err := time.ParseInLocation(format, string(s), defaultLocation)
		if err == nil {
			return TimestampLiteral(t), nil
		}
	}
	return TimestampLiteral{}, p.errorf("invalid timestamp literal %q", s)
}
func (p *parser) gologoo__parseStringLit_531485d04a38cff0d09d3c55685c367d() (StringLiteral, *parseError) {
	tok := p.next()
	if tok.err != nil {
		return "", tok.err
	}
	if tok.typ != stringToken {
		return "", p.errorf("got %q, want string literal", tok.value)
	}
	return StringLiteral(tok.string), nil
}
func (p *parser) gologoo__parsePathExp_531485d04a38cff0d09d3c55685c367d() (PathExp, *parseError) {
	var pe PathExp
	for {
		tok := p.next()
		if tok.err != nil {
			return nil, tok.err
		}
		switch tok.typ {
		case quotedID:
			pe = append(pe, ID(tok.string))
		case unquotedID:
			pe = append(pe, ID(tok.value))
		default:
			return nil, p.errorf("expected identifer")
		}
		if !p.eat(".") {
			break
		}
	}
	return pe, nil
}
func (p *parser) gologoo__parseBoolExpr_531485d04a38cff0d09d3c55685c367d() (BoolExpr, *parseError) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	be, ok := expr.(BoolExpr)
	if !ok {
		return nil, p.errorf("got non-bool expression %T", expr)
	}
	return be, nil
}
func (p *parser) gologoo__parseAlias_531485d04a38cff0d09d3c55685c367d() (ID, *parseError) {
	return p.parseTableOrIndexOrColumnName()
}
func (p *parser) gologoo__parseHints_531485d04a38cff0d09d3c55685c367d(hints map[string]string) (map[string]string, *parseError) {
	if hints == nil {
		hints = map[string]string{}
	}
	if err := p.expect("{"); err != nil {
		return nil, err
	}
	for {
		if p.sniff("}") {
			break
		}
		tok := p.next()
		if tok.err != nil {
			return nil, tok.err
		}
		k := tok.value
		if err := p.expect("="); err != nil {
			return nil, err
		}
		tok = p.next()
		if tok.err != nil {
			return nil, tok.err
		}
		v := tok.value
		hints[k] = v
		if !p.eat(",") {
			break
		}
	}
	if err := p.expect("}"); err != nil {
		return nil, err
	}
	return hints, nil
}
func (p *parser) gologoo__parseTableOrIndexOrColumnName_531485d04a38cff0d09d3c55685c367d() (ID, *parseError) {
	tok := p.next()
	if tok.err != nil {
		return "", tok.err
	}
	switch tok.typ {
	case quotedID:
		return ID(tok.string), nil
	case unquotedID:
		return ID(tok.value), nil
	default:
		return "", p.errorf("expected identifier")
	}
}
func (p *parser) gologoo__parseOnDelete_531485d04a38cff0d09d3c55685c367d() (OnDelete, *parseError) {
	tok := p.next()
	if tok.err != nil {
		return 0, tok.err
	}
	if tok.caseEqual("CASCADE") {
		return CascadeOnDelete, nil
	}
	if !tok.caseEqual("NO") {
		return 0, p.errorf("got %q, want NO or CASCADE", tok.value)
	}
	if err := p.expect("ACTION"); err != nil {
		return 0, err
	}
	return NoActionOnDelete, nil
}
func (p *parser) gologoo__parseRowDeletionPolicy_531485d04a38cff0d09d3c55685c367d() (RowDeletionPolicy, *parseError) {
	if err := p.expect("(", "OLDER_THAN", "("); err != nil {
		return RowDeletionPolicy{}, err
	}
	cname, err := p.parseTableOrIndexOrColumnName()
	if err != nil {
		return RowDeletionPolicy{}, err
	}
	if err := p.expect(",", "INTERVAL"); err != nil {
		return RowDeletionPolicy{}, err
	}
	tok := p.next()
	if tok.err != nil {
		return RowDeletionPolicy{}, tok.err
	}
	if tok.typ != int64Token {
		return RowDeletionPolicy{}, p.errorf("got %q, expected int64 token", tok.value)
	}
	n, serr := strconv.ParseInt(tok.value, tok.int64Base, 64)
	if serr != nil {
		return RowDeletionPolicy{}, p.errorf("%v", serr)
	}
	if err := p.expect("DAY", ")", ")"); err != nil {
		return RowDeletionPolicy{}, err
	}
	return RowDeletionPolicy{Column: cname, NumDays: n}, nil
}
func (p *parser) gologoo__parseCommaList_531485d04a38cff0d09d3c55685c367d(bra, ket string, f func(*parser) *parseError) *parseError {
	if err := p.expect(bra); err != nil {
		return err
	}
	for {
		if p.eat(ket) {
			return nil
		}
		err := f(p)
		if err != nil {
			return err
		}
		tok := p.next()
		if tok.err != nil {
			return err
		}
		if tok.value == ket {
			return nil
		} else if tok.value == "," {
			continue
		} else {
			return p.errorf(`got %q, want %q or ","`, tok.value, ket)
		}
	}
}
func debugf(format string, args ...interface {
}) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__debugf_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v %v\n", format, args)
	gologoo__debugf_531485d04a38cff0d09d3c55685c367d(format, args...)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func ParseDDL(filename, s string) (*DDL, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ParseDDL_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v %v\n", filename, s)
	r0, r1 := gologoo__ParseDDL_531485d04a38cff0d09d3c55685c367d(filename, s)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func ParseDDLStmt(s string) (DDLStmt, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ParseDDLStmt_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", s)
	r0, r1 := gologoo__ParseDDLStmt_531485d04a38cff0d09d3c55685c367d(s)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func ParseDMLStmt(s string) (DMLStmt, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ParseDMLStmt_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", s)
	r0, r1 := gologoo__ParseDMLStmt_531485d04a38cff0d09d3c55685c367d(s)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func ParseQuery(s string) (Query, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ParseQuery_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", s)
	r0, r1 := gologoo__ParseQuery_531485d04a38cff0d09d3c55685c367d(s)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *token) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0 := t.gologoo__String_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (pe *parseError) Error() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Error_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0 := pe.gologoo__Error_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *parser) Pos() Position {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Pos_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0 := p.gologoo__Pos_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func newParser(filename, s string) *parser {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__newParser_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v %v\n", filename, s)
	r0 := gologoo__newParser_531485d04a38cff0d09d3c55685c367d(filename, s)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *parser) Rem() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Rem_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0 := p.gologoo__Rem_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *parser) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0 := p.gologoo__String_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *parser) errorf(format string, args ...interface {
}) *parseError {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errorf_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v %v\n", format, args)
	r0 := p.gologoo__errorf_531485d04a38cff0d09d3c55685c367d(format, args...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func isInitialIdentifierChar(c byte) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isInitialIdentifierChar_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", c)
	r0 := gologoo__isInitialIdentifierChar_531485d04a38cff0d09d3c55685c367d(c)
	log.Printf("Output: %v\n", r0)
	return r0
}
func isIdentifierChar(c byte) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isIdentifierChar_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", c)
	r0 := gologoo__isIdentifierChar_531485d04a38cff0d09d3c55685c367d(c)
	log.Printf("Output: %v\n", r0)
	return r0
}
func isHexDigit(c byte) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isHexDigit_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", c)
	r0 := gologoo__isHexDigit_531485d04a38cff0d09d3c55685c367d(c)
	log.Printf("Output: %v\n", r0)
	return r0
}
func isOctalDigit(c byte) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isOctalDigit_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", c)
	r0 := gologoo__isOctalDigit_531485d04a38cff0d09d3c55685c367d(c)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *parser) consumeNumber() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__consumeNumber_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	p.gologoo__consumeNumber_531485d04a38cff0d09d3c55685c367d()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *parser) consumeString() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__consumeString_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	p.gologoo__consumeString_531485d04a38cff0d09d3c55685c367d()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *parser) consumeRawString() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__consumeRawString_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	p.gologoo__consumeRawString_531485d04a38cff0d09d3c55685c367d()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *parser) consumeBytes() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__consumeBytes_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	p.gologoo__consumeBytes_531485d04a38cff0d09d3c55685c367d()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *parser) consumeRawBytes() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__consumeRawBytes_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	p.gologoo__consumeRawBytes_531485d04a38cff0d09d3c55685c367d()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *parser) stringDelimiter() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__stringDelimiter_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0 := p.gologoo__stringDelimiter_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *parser) consumeStringContent(delim string, raw, unicode bool, name string) (string, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__consumeStringContent_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v %v %v %v\n", delim, raw, unicode, name)
	r0, r1 := p.gologoo__consumeStringContent_531485d04a38cff0d09d3c55685c367d(delim, raw, unicode, name)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func isSpace(c byte) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isSpace_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", c)
	r0 := gologoo__isSpace_531485d04a38cff0d09d3c55685c367d(c)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *parser) skipSpace() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__skipSpace_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0 := p.gologoo__skipSpace_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *parser) advance() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__advance_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	p.gologoo__advance_531485d04a38cff0d09d3c55685c367d()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *parser) back() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__back_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	p.gologoo__back_531485d04a38cff0d09d3c55685c367d()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *parser) next() *token {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__next_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0 := p.gologoo__next_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *token) caseEqual(x string) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__caseEqual_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", x)
	r0 := t.gologoo__caseEqual_531485d04a38cff0d09d3c55685c367d(x)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *parser) sniff(want ...string) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__sniff_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", want)
	r0 := p.gologoo__sniff_531485d04a38cff0d09d3c55685c367d(want...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *parser) sniffTokenType(want tokenType) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__sniffTokenType_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", want)
	r0 := p.gologoo__sniffTokenType_531485d04a38cff0d09d3c55685c367d(want)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *parser) eat(want ...string) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__eat_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", want)
	r0 := p.gologoo__eat_531485d04a38cff0d09d3c55685c367d(want...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *parser) expect(want ...string) *parseError {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__expect_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", want)
	r0 := p.gologoo__expect_531485d04a38cff0d09d3c55685c367d(want...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *parser) parseDDLStmt() (DDLStmt, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseDDLStmt_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseDDLStmt_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseCreateTable() (*CreateTable, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseCreateTable_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseCreateTable_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) sniffTableConstraint() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__sniffTableConstraint_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0 := p.gologoo__sniffTableConstraint_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *parser) parseCreateIndex() (*CreateIndex, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseCreateIndex_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseCreateIndex_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseCreateView() (*CreateView, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseCreateView_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseCreateView_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseAlterTable() (*AlterTable, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseAlterTable_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseAlterTable_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseAlterDatabase() (*AlterDatabase, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseAlterDatabase_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseAlterDatabase_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseDMLStmt() (DMLStmt, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseDMLStmt_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseDMLStmt_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseUpdateItem() (UpdateItem, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseUpdateItem_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseUpdateItem_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseColumnDef() (ColumnDef, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseColumnDef_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseColumnDef_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseColumnAlteration() (ColumnAlteration, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseColumnAlteration_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseColumnAlteration_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseColumnOptions() (ColumnOptions, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseColumnOptions_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseColumnOptions_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseDatabaseOptions() (DatabaseOptions, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseDatabaseOptions_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseDatabaseOptions_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseKeyPartList() ([]KeyPart, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseKeyPartList_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseKeyPartList_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseKeyPart() (KeyPart, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseKeyPart_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseKeyPart_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseTableConstraint() (TableConstraint, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseTableConstraint_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseTableConstraint_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseConstraint() (Constraint, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseConstraint_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseConstraint_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseForeignKey() (ForeignKey, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseForeignKey_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseForeignKey_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseCheck() (Check, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseCheck_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseCheck_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseColumnNameList() ([]ID, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseColumnNameList_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseColumnNameList_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseBaseType() (Type, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseBaseType_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseBaseType_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseType() (Type, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseType_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseType_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseExtractType() (Type, string, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseExtractType_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1, r2 := p.gologoo__parseExtractType_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func (p *parser) parseBaseOrParameterizedType(withParam bool) (Type, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseBaseOrParameterizedType_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", withParam)
	r0, r1 := p.gologoo__parseBaseOrParameterizedType_531485d04a38cff0d09d3c55685c367d(withParam)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseQuery() (Query, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseQuery_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseQuery_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseSelect() (Select, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseSelect_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseSelect_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseSelectList() ([]Expr, []ID, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseSelectList_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1, r2 := p.gologoo__parseSelectList_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func (p *parser) parseSelectFromTable() (SelectFrom, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseSelectFromTable_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseSelectFromTable_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseSelectFromJoin(lhs SelectFrom) (SelectFrom, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseSelectFromJoin_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", lhs)
	r0, r1 := p.gologoo__parseSelectFromJoin_531485d04a38cff0d09d3c55685c367d(lhs)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseSelectFrom() (SelectFrom, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseSelectFrom_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseSelectFrom_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseTableSample() (TableSample, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseTableSample_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseTableSample_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseOrder() (Order, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseOrder_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseOrder_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseLiteralOrParam() (LiteralOrParam, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseLiteralOrParam_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseLiteralOrParam_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseExprList() ([]Expr, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseExprList_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseExprList_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseParenExprList() ([]Expr, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseParenExprList_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseParenExprList_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseParenExprListWithParseFunc(f func(*parser) (Expr, *parseError)) ([]Expr, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseParenExprListWithParseFunc_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", f)
	r0, r1 := p.gologoo__parseParenExprListWithParseFunc_531485d04a38cff0d09d3c55685c367d(f)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseExpr() (Expr, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseExpr_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseExpr_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (bin binOpParser) parse(p *parser) (Expr, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parse_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", p)
	r0, r1 := bin.gologoo__parse_531485d04a38cff0d09d3c55685c367d(p)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func init() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__init_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	gologoo__init_531485d04a38cff0d09d3c55685c367d()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func newBinArithParser(opStr string, op ArithOperator, nextPrec func(*parser) (Expr, *parseError)) binOpParser {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__newBinArithParser_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v %v %v\n", opStr, op, nextPrec)
	r0 := gologoo__newBinArithParser_531485d04a38cff0d09d3c55685c367d(opStr, op, nextPrec)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *parser) parseLogicalNot() (Expr, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseLogicalNot_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseLogicalNot_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseIsOp() (Expr, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseIsOp_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseIsOp_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseInOp() (Expr, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseInOp_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseInOp_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseComparisonOp() (Expr, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseComparisonOp_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseComparisonOp_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseArithOp() (Expr, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseArithOp_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseArithOp_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseUnaryArithOp() (Expr, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseUnaryArithOp_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseUnaryArithOp_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseLit() (Expr, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseLit_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseLit_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseArrayLit() (Array, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseArrayLit_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseArrayLit_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseDateLit() (DateLiteral, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseDateLit_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseDateLit_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseTimestampLit() (TimestampLiteral, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseTimestampLit_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseTimestampLit_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseStringLit() (StringLiteral, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseStringLit_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseStringLit_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parsePathExp() (PathExp, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parsePathExp_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parsePathExp_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseBoolExpr() (BoolExpr, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseBoolExpr_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseBoolExpr_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseAlias() (ID, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseAlias_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseAlias_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseHints(hints map[string]string) (map[string]string, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseHints_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v\n", hints)
	r0, r1 := p.gologoo__parseHints_531485d04a38cff0d09d3c55685c367d(hints)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseTableOrIndexOrColumnName() (ID, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseTableOrIndexOrColumnName_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseTableOrIndexOrColumnName_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseOnDelete() (OnDelete, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseOnDelete_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseOnDelete_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseRowDeletionPolicy() (RowDeletionPolicy, *parseError) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseRowDeletionPolicy_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : (none)\n")
	r0, r1 := p.gologoo__parseRowDeletionPolicy_531485d04a38cff0d09d3c55685c367d()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *parser) parseCommaList(bra, ket string, f func(*parser) *parseError) *parseError {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseCommaList_531485d04a38cff0d09d3c55685c367d")
	log.Printf("Input : %v %v %v\n", bra, ket, f)
	r0 := p.gologoo__parseCommaList_531485d04a38cff0d09d3c55685c367d(bra, ket, f)
	log.Printf("Output: %v\n", r0)
	return r0
}
