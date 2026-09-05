package expr

import (
	"math"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		in   string
		want float64
	}{
		{"1+2", 3},
		{"1+2*3", 7},    // 優先順位：* が先
		{"(1+2)*3", 9},  // カッコで上書き
		{"10-3-2", 5},   // 左結合：(10-3)-2 であって 10-(3-2) ではない
		{"100/10/2", 5}, // 左結合
		{"2*3+4*5", 26},
		{"-3+5", 2}, // 単項マイナス
		{"-(2+3)*2", -10},
		{"1.5*4", 6},
		{"((1+2)*(3+4))", 21},
		{"  1 +  2 * 3 ", 7}, // 空白は無視
	}
	for _, tt := range tests {
		got, err := Eval(tt.in)
		if err != nil {
			t.Errorf("Eval(%q) failed: %v", tt.in, err)
			continue
		}
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("Eval(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestEvalErrors(t *testing.T) {
	bad := []string{
		"1+",   // 右辺がない
		"(1+2", // 閉じ括弧がない
		"1+2)", // 余りがある
		"1/0",  // ゼロ除算
		"1$2",  // 未知の文字
		"",     // 空
		"1..2", // 数値として不正
	}
	for _, in := range bad {
		if v, err := Eval(in); err == nil {
			t.Errorf("Eval(%q) = %v, エラーになるべき", in, v)
		}
	}
}
