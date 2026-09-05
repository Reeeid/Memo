package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"calc/expr"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	fmt.Println(`式を入力（例: (1+2)*3 ／ 先頭に ? でトークン表示 ／ 終了は q）`)

	for {
		fmt.Print("> ")
		if !sc.Scan() {
			return
		}
		line := strings.TrimSpace(sc.Text())
		switch {
		case line == "":
			continue
		case line == "q":
			return
		case strings.HasPrefix(line, "?"):
			s, err := expr.Tokens(strings.TrimPrefix(line, "?"))
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
			fmt.Println(s)
		default:
			v, err := expr.Eval(line)
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
			fmt.Println(v)
		}
	}
}
