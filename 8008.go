package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func usage(exitcode int) {
	message := fmt.Sprintf(`usage: 8008 [-h] [-s] [-o {dst}] {src}`)

	if exitcode == 0 {
		fmt.Println(message)
	} else {
		fmt.Fprintln(os.Stderr, message)
	}

	os.Exit(exitcode)
}

type Option struct {
	help bool
	show bool
	src  string
	dst  string
}

func newOption() *Option {
	opt := new(Option)

	flag.BoolVar(&opt.help, "h", false, "Show usage")
	flag.BoolVar(&opt.show, "s", false, "Interpret sources in dry (show only)")
	flag.StringVar(&opt.dst, "o", "", "Specify target")

	flag.Parse()

	opt.src = flag.Arg(0)

	return opt
}

func main() {
	opt := newOption()
	if opt.help {
		usage(0)
	} else if opt.src == "" {
		usage(1)
	}

	src_file, err := os.Open(opt.src)
	if err != nil {
		panic(err)
	}
	defer src_file.Close()

	if opt.show {
		reader := bufio.NewReader(src_file)
		show(reader)
		return
	} else {
		panic("Only dry mode is supported currently")
	}

	// var dst_file *os.File
	// if opt.dst == "" {
	// 	dst_file = os.Stdout
	// } else {
	// 	dst_file, _ = os.Create(opt.dst)
	// }
	// defer dst_file.Close()
	//
	// cpu := NewCPU()
	// cpu.LoadROM(src_file)
	// cpu.Run()
	// cpu.SaveRAM(dst_file)
}
