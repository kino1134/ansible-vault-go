package main

import (
	"fmt"
	"os"

	cli "github.com/kino1134/ansible-vault-go/cmd"
)

func main() {
	// スタックトレースを出したいので、コメントアウト
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Fprintln(os.Stderr, err)
	// 		os.Exit(1)
	// 	}
	// }()

	os.Exit(_main())
}

func _main() int {
	err := cli.Run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	return 0
}
