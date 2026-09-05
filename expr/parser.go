package expr

import "fmt"

// 文法（BNF風）。この3行がそのまま下の3つの関数になる。
//
//	expression := term   ( ("+" | "-") term   )*
//	term       := factor ( ("*" | "/") factor )*
//	factor     := NUMBER | "(" expression ")" | "-" factor
//
// expression が term を呼び、term が factor を呼ぶ、という「深い方が先に計算される」
// 入れ子構造そのものが演算子の優先順位になっている。

type parser struct {
	toks []token
	pos  int
}

func (p *parser) peek() token { return p.toks[p.pos] }

func (p *parser) next() token {
	t := p.toks[p.pos]
	if t.kind != kindEOF {
		p.pos++
	}
	return t
}

// accept は次が k なら1つ進めて true。違えば何もせず false。
func (p *parser) accept(k kind) bool {
	if p.peek().kind == k {
		p.pos++
		return true
	}
	return false
}

// expression := term ( ("+"|"-") term )*
func (p *parser) expression() (float64, error) {
	left, err := p.term()
	if err != nil {
		return 0, err
	}
	for {
		switch {
		case p.accept(kindPlus):
			right, err := p.term()
			if err != nil {
				return 0, err
			}
			left += right
		case p.accept(kindMinus):
			right, err := p.term()
			if err != nil {
				return 0, err
			}
			left -= right
		default:
			return left, nil
		}
	}
}

// term := factor ( ("*"|"/") factor )*
func (p *parser) term() (float64, error) {
	left, err := p.factor()
	if err != nil {
		return 0, err
	}
	for {
		switch {
		case p.accept(kindStar):
			right, err := p.factor()
			if err != nil {
				return 0, err
			}
			left *= right
		case p.accept(kindSlash):
			right, err := p.factor()
			if err != nil {
				return 0, err
			}
			if right == 0 {
				return 0, fmt.Errorf("0 では割れません")
			}
			left /= right
		default:
			return left, nil
		}
	}
}

// factor := NUMBER | "(" expression ")" | "-" factor
func (p *parser) factor() (float64, error) {
	t := p.peek()

	switch {
	case t.kind == kindNum:
		p.next()
		return t.num, nil

	case t.kind == kindMinus: // 単項マイナス： -3, -(1+2)
		p.next()
		v, err := p.factor()
		if err != nil {
			return 0, err
		}
		return -v, nil

	case t.kind == kindLParen:
		p.next()
		v, err := p.expression() // ← 再帰。カッコの中は独立した式
		if err != nil {
			return 0, err
		}
		if !p.accept(kindRParen) {
			return 0, fmt.Errorf("%d文字目: ')' が閉じていません", t.pos+1)
		}
		return v, nil
	}

	return 0, fmt.Errorf("%d文字目: %v が来るべきではありません", t.pos+1, t.kind)
}

// Eval は中置記法の式を評価する。
//
//	Eval("1+2*3")      → 7
//	Eval("(1+2)*3")    → 9
func Eval(input string) (float64, error) {
	toks, err := lex(input)
	if err != nil {
		return 0, err
	}

	p := &parser{toks: toks}
	v, err := p.expression()
	if err != nil {
		return 0, err
	}

	// 式を読み終えた時点で EOF に到達していないなら、余りがある（例: "1+2)"）
	if t := p.peek(); t.kind != kindEOF {
		return 0, fmt.Errorf("%d文字目: %v が余っています", t.pos+1, t.kind)
	}
	return v, nil
}

// Tokens はデバッグ用に、字句解析の結果を文字列で返す。
func Tokens(input string) (string, error) {
	toks, err := lex(input)
	if err != nil {
		return "", err
	}
	return dump(toks), nil
}
