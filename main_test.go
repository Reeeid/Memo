package main

import (
	"strings"
	"testing"
)

func TestCalculator(t *testing.T) {
	tests := []struct {
		name string
		keys string
		want string
	}{
		{"足し算", "1 + 2 =", "3"},
		{"複数桁", "12 + 34 =", "46"},
		{"1桁ずつ", "1 2 + 3 4 =", "46"},
		{"引き算", "10 - 3 =", "7"},
		{"掛け算", "6 * 7 =", "42"},
		{"割り算", "9 / 2 =", "4.5"},
		{"連続演算", "1 + 2 + 3 =", "6"},
		{"小数点", "1 . 5 + 2 =", "3.5"},
		{"小数点は1つだけ", "1 . 2 . 3 =", "1.23"},
		{"AC", "1 + 2 AC", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Calculator{}
			for _, key := range strings.Fields(tt.keys) {
				if err := c.Input(key); err != nil {
					t.Fatalf("Input(%q) failed: %v", key, err)
				}
			}
			if got := c.Display(); got != tt.want {
				t.Errorf("keys %q: got %q, want %q", tt.keys, got, tt.want)
			}
		})
	}
}

func TestDivideByZero(t *testing.T) {
	if _, err := apply(1, 0, "/"); err == nil {
		t.Error("1/0 はエラーになるべき")
	}
}

// --- ここから次の課題 ---

func TestChainFromResult(t *testing.T) {
	// 結果に続けて演算する。実機の電卓は当然できる。
	tests := []struct {
		keys string
		want string
	}{
		{"1 + 2 = + 3 =", "6"},
		{"10 / 4 = * 2 =", "5"},
	}
	for _, tt := range tests {
		c := &Calculator{}
		for _, key := range strings.Fields(tt.keys) {
			if err := c.Input(key); err != nil {
				t.Fatalf("keys %q: Input(%q) failed: %v", tt.keys, key, err)
			}
		}
		if got := c.Display(); got != tt.want {
			t.Errorf("keys %q: got %q, want %q", tt.keys, got, tt.want)
		}
	}
}

func TestRepeatEquals(t *testing.T) {
	// = 連打。直前の「演算子と右辺」を繰り返し適用する実機の挙動。
	tests := []struct {
		keys string
		want string
	}{
		{"1 + 2 =", "3"},
		{"1 + 2 = =", "5"},
		{"1 + 2 = = =", "7"},
		{"10 - 1 = =", "8"},
	}
	for _, tt := range tests {
		c := &Calculator{}
		for _, key := range strings.Fields(tt.keys) {
			if err := c.Input(key); err != nil {
				t.Fatalf("keys %q: Input(%q) failed: %v", tt.keys, key, err)
			}
		}
		if got := c.Display(); got != tt.want {
			t.Errorf("keys %q: got %q, want %q", tt.keys, got, tt.want)
		}
	}
}
