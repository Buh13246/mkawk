package main

import (
	"fmt"
)

type Command interface {
	Type() string
	String() string
}

type NumberCommand struct {
	Number int
}

func (n NumberCommand) String() string {
	return fmt.Sprintf(`printf "%%s",$%d  #User-Only-Variable`+"\n", n.Number)
}

func (n NumberCommand) Type() string {
	return "NUMBER"
}

type SkipCommand struct{}

func (cmd SkipCommand) String() string {
	return ""
}

func (cmd SkipCommand) Type() string {
	return "SKIP"
}

type TextCommand struct {
	Text string
}

func (cmd TextCommand) String() string {
	return fmt.Sprintf(`printf "%s"`+"\n", cmd.Text)
}

func (cmd TextCommand) Type() string {
	return "TEXT"
}

type AWKCommand struct {
	Text string
}

func (cmd AWKCommand) String() string {
	return fmt.Sprintf(`%s`+"\n", cmd.Text)
}

func (cmd AWKCommand) Type() string {
	return "AWK"
}

type CommentaryCommand struct {
	Text string
}

func (cmd CommentaryCommand) String() string {
	return fmt.Sprintf(`# %s`+"\n", cmd.Text)
}

func (cmd CommentaryCommand) Type() string {
	return "COMMENTARY"
}

type FormatCommand struct {
	Format string
	Data   string
}

func (cmd FormatCommand) String() string {
	return fmt.Sprintf(`printf "%s", %s`+"\n", cmd.Format, cmd.Data)
}

func (cmd FormatCommand) Type() string {
	return "FORMAT"
}

type NextSectionCommand struct {
	Section int
}

func (cmd NextSectionCommand) String() string {
	switch cmd.Section {
	case 1:
		return fmt.Sprintln("\n}\nEND {")
	}
	return fmt.Sprintln("\n}\n{")
}

func (cmd NextSectionCommand) Type() string {
	return "NEXT_SECTION"
}
