package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Calculator は記事の3つの状態変数をそのまま持つ。
//
//	buffer   … 入力中の数字（文字列で持つ）
//	stack    … 確定した数値
//	operator … 直前に押された演算子（未入力なら ""）
type Calculator struct {
	buffer   string
	stack    []float64
	operator string
}

// Display は画面に出す文字列を返す。
// 入力中なら buffer、なければスタックの一番上、それも無ければ "0"。
func (c *Calculator) Display() string {
	if c.buffer != "" {
		return c.buffer
	}
	if len(c.stack) > 0 {
		return formatFloat(c.stack[len(c.stack)-1])
	}
	return "0"
}

// Input はキーを1つ受け取って状態を進める。
func (c *Calculator) Input(key string) error {
	switch key {
	case "+", "-", "*", "/":
		if err := c.push(); err != nil {
			return err
		}
		if c.operator != "" {
			if err := c.reduce(); err != nil {
				return err
			}
		}
		c.operator = key
		// buffer を数値化して push し、operator を key に更新する。
		//       すでに operator が入っていれば、その時点で一度計算する（連続演算）。
		return nil

	case "=":
		// buffer を push して operator を適用し、結果だけを stack に残す。
		if c.operator == "" {
			return nil
		}
		if err := c.reduce(); err != nil {
			return err
		}
		c.operator = ""
		return nil

	case ".":
		if strings.Contains(c.buffer, ".") {
			return nil
		}
		if c.buffer == "" {
			c.buffer = "0."
		} else {
			c.buffer += "."
		}
		//  buffer に既に "." があるなら無視。無ければ追加。
		//       buffer が空なら "0." にすると自然。
		return nil

	case "C":
		c.buffer = ""
		return nil

	case "AC":
		c.buffer = ""
		c.stack = nil
		c.operator = ""
		return nil

	default:
		c.buffer += key
		// 数字キーとして buffer に足す。
		return nil
	}
}

// push は buffer を数値化して stack に積み、buffer を空にする。
func (c *Calculator) push() error {
	// strconv.ParseFloat を使う。buffer が空のときの扱いも決めること。
	str, err := strconv.ParseFloat(c.buffer, 64)
	if err != nil {
		return fmt.Errorf("push: %v", err)
	}
	c.stack = append(c.stack, str)
	c.buffer = ""
	return nil
}

func (c *Calculator) reduce() error {
	b, err := c.pop()
	if err != nil {
		return err
	}
	a, err := c.pop()
	if err != nil {
		return err
	}
	result, err := apply(a, b, c.operator)
	if err != nil {
		return err
	}
	c.stack = append(c.stack, result)
	return nil
}

// pop は stack の末尾を取り出す。
func (c *Calculator) pop() (float64, error) {
	if len(c.stack) == 0 {
		return 0, fmt.Errorf("pop: stack is empty")
	}
	v := c.stack[len(c.stack)-1]
	c.stack = c.stack[:len(c.stack)-1]
	return v, nil
	// len(c.stack) == 0 のときは error を返す。
}

// apply は二項演算を行う。
func apply(a, b float64, op string) (float64, error) {
	// switch op で四則演算。"/" は b == 0 を弾く。
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("apply: division by zero")
		}
		return a / b, nil
	}
	return 0, fmt.Errorf("apply: unknown operator")
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}

func main() {
	c := &Calculator{}
	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanWords)

	fmt.Println(`数字 / + - * / = / C / AC を入力（終了は q）`)
	fmt.Printf("> %s\n", c.Display())

	for sc.Scan() {
		key := sc.Text()
		if key == "q" {
			return
		}
		if err := c.Input(key); err != nil {
			fmt.Println("error:", err)
			c.buffer, c.stack, c.operator = "", nil, ""
		}
		fmt.Printf("> %s\n", c.Display())
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "read error:", err)
	}
}
