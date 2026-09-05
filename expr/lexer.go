package expr

import (
	"fmt"
	"strconv"
	"strings"
)

// kind はトークンの種類。
type kind int

const (
	kindNum kind = iota
	kindPlus
	kindMinus
	kindStar
	kindSlash
	kindLParen
	kindRParen
	kindEOF
)

func (k kind) String() string {
	return [...]string{"数値", "+", "-", "*", "/", "(", ")", "式の終わり"}[k]
}

// token は字句解析の結果1つ分。num は kindNum のときだけ意味を持つ。
type token struct {
	kind kind
	num  float64
	pos  int // エラー表示用に、元の文字列での位置
}

// lex は入力文字列をトークン列に分解する（字句解析）。
//
//	"1+2*3"  →  [1] [+] [2] [*] [3] [EOF]
func lex(input string) ([]token, error) {
	var toks []token
	i := 0
	for i < len(input) {
		ch := input[i]

		// 空白は読み飛ばす
		if ch == ' ' || ch == '\t' {
			i++
			continue
		}

		// 数値：数字か '.' が続く限りまとめて取る
		if isDigit(ch) || ch == '.' {
			start := i
			for i < len(input) && (isDigit(input[i]) || input[i] == '.') {
				i++
			}
			text := input[start:i]
			v, err := strconv.ParseFloat(text, 64)
			if err != nil {
				return nil, fmt.Errorf("%d文字目: %q は数値として読めません", start+1, text)
			}
			toks = append(toks, token{kind: kindNum, num: v, pos: start})
			continue
		}

		// 記号1文字
		var k kind
		switch ch {
		case '+':
			k = kindPlus
		case '-':
			k = kindMinus
		case '*':
			k = kindStar
		case '/':
			k = kindSlash
		case '(':
			k = kindLParen
		case ')':
			k = kindRParen
		default:
			return nil, fmt.Errorf("%d文字目: 予期しない文字 %q", i+1, string(ch))
		}
		toks = append(toks, token{kind: k, pos: i})
		i++
	}

	toks = append(toks, token{kind: kindEOF, pos: len(input)})
	return toks, nil
}

func isDigit(ch byte) bool { return '0' <= ch && ch <= '9' }

// dump はデバッグ用にトークン列を文字列化する。
func dump(toks []token) string {
	var sb strings.Builder
	for i, t := range toks {
		if i > 0 {
			sb.WriteByte(' ')
		}
		if t.kind == kindNum {
			fmt.Fprintf(&sb, "[%v]", t.num)
		} else {
			fmt.Fprintf(&sb, "[%v]", t.kind)
		}
	}
	return sb.String()
}
