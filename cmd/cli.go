package vault_go

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// 処理実行時の引数
type cmdOpts struct {
	Command      func(*cmdOpts) error
	Sub          string
	Path         string
	Password     string
	Editor       string
	RestoreSalt  bool
	PasswordFile string
	Label        string
}

// コマンド名と処理内容
var commands = map[string]func(*cmdOpts) error{
	"view":    view,
	"edit":    edit,
	"encrypt": encrypt,
	"decrypt": decrypt,
}

func initParse() (*cmdOpts, *flag.FlagSet) {
	opts := &cmdOpts{}

	fs := flag.NewFlagSet("vault_go", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: vault_go [edit|view|encrypt] [vaultfile.yml] [options]\n")
		fs.PrintDefaults()
	}

	fs.StringVar(&opts.Editor, "editor", "vim", "編集に使用するエディタを指定")
	fs.StringVar(&opts.PasswordFile, "password-file", "", "パスワードの書かれたファイルを指定")
	fs.StringVar(&opts.Label, "label", "", "暗号化時に付与するラベルを指定")
	fs.BoolVar(&opts.RestoreSalt, "restore-salt", false, "編集前後でソルトを変えない")

	return opts, fs
}

// 指定された引数を解析する
func Parse(args []string) (*cmdOpts, error) {
	opts, fs := initParse()

	if len(args) < 2 {
		fs.Usage()
		return nil, errors.New("")
	}

	cmd, ok := commands[args[0]]
	if !ok {
		fmt.Fprintln(os.Stderr, "サブコマンドが正しくありません。")
		fs.Usage()
		return nil, errors.New("")
	}

	opts.Sub = args[0]
	opts.Command = cmd
	opts.Path = args[1]
	fs.Parse(args[2:])

	return opts, nil
}

func Run(args []string) error {
	opts, err := Parse(args)
	if err != nil {
		return err
	}

	if err := SetPassword(opts); err != nil {
		return err
	}

	if err := opts.Command(opts); err != nil {
		return err
	}

	return nil
}
