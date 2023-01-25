package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var FS, inFile, outfile, datafile string

func main() {
	t1 := time.Now()

	flag.StringVar(&FS, "D", ",", "-D <char>")
	flag.StringVar(&inFile, "f", "", "-f <textfile>")
	flag.StringVar(&outfile, "o", "", "-o <awkfile>")
	flag.StringVar(&datafile, "A", "", "-A <datafile>")
	flag.Parse()

	var fileIN, fileOUT = os.Stdin, os.Stdout

	if inFile != "" {
		file, err := os.Open(inFile)
		if err != nil {
			log.Fatal("Cannot Open inputfile! Exiting")
		}
		fileIN = file
	}

	if outfile != "" {
		file, err := os.Create(outfile)
		if err != nil {
			log.Fatal("Cannot Open outputfile! Exiting")
		}
		fileOUT = file
	}

	if datafile != "" {
		f, err := ioutil.TempFile(os.TempDir(), "awk-")
		if err != nil {
			log.Fatal("could not create tmpfile: ", err)
		}
		defer os.Remove(f.Name())
		defer f.Close()
		fileData, err := os.Open(datafile)
		if err != nil {
			log.Fatal("Cannot Open datafile! Exiting", err)
		}
		convert(FS, fileIN, f)
		_, err = f.Seek(0, 0)
		if err != nil {
			log.Fatal("Could not reset script file to Start of file: ", err)
		}
		awk(f, fileData, fileOUT)

	} else {
		convert(FS, fileIN, fileOUT)
	}

	log.Println("Convert took ", time.Since(t1))
}
func convert(FS string, inFILE, outFILE *os.File) {
	osOutFile := os.Stdout
	osInFile := os.Stdin
	os.Stdout = outFILE
	data, err := io.ReadAll(inFILE)
	if err != nil {
		log.Fatal("STDIN is broken")
	}
	text := string(data)
	section := 1
	num := strings.Count(text, "<@>")
	if num > 1 {
		section = 0
	}

	text = strings.ReplaceAll(text, "<@>", "<@><->")
	text = strings.ReplaceAll(text, "\\<", "<(>")
	text = strings.ReplaceAll(text, "\\>", "<)>")
	text = strings.ReplaceAll(text, "\n", "\\n\n")
	text = strings.ReplaceAll(text, "<", "\n<")
	text = strings.ReplaceAll(text, ">", ">\n")
	lines := strings.Split(string(text), "\n")
	fmt.Println(`
###############################################################################
#                                                                             #
#         This is an automatic build awk-script                               #
#         It is constructed by a reimplementation of the Program MKAWK        #
#         Which was written by M.Wellner, Berlin Feb.1989                     #
#         Rewritten by P.Wellner, Berlin Jan.2023                             #
#                                                                             #
###############################################################################

BEGIN { FS="` + FS + `"
`)
	if section == 1 {
		fmt.Println("}\n{")
	}
	skip := false
	for _, line := range lines {
		if skip {
			if strings.Contains(line, "\\n") {
				skip = false
			}
			continue
		}

		if strings.Contains(line, "<") && strings.Contains(line, ">") {
			command, err := parseCommand(line, section)
			if err == nil {
				if command.Type() == "SKIP" {
					skip = true
				}
				if command.Type() == "NEXT_SECTION" {
					section++
				}
				fmt.Print(command)
				continue
			}
			log.Printf("ERR: Could not parse cmd %q : %v", line, err)
		}
		fmt.Printf(`printf "%s"`+"\n", line)
	}
	fmt.Println("}       # END OF AWK-SCRIPT")

	os.Stdout = osOutFile
	os.Stdin = osInFile
}

func parseCommand(line string, section int) (Command, error) {
	cmd := line[1 : len(line)-1]
	switch cmd[0] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '$':
		return parseNumberCommand(cmd)
	case '-':
		return SkipCommand{}, nil
	case '(':
		return TextCommand{Text: "<"}, nil
	case ')':
		return TextCommand{Text: ">"}, nil
	case '!':
		return AWKCommand{Text: cmd[1:]}, nil
	case '#':
		return CommentaryCommand{Text: cmd[1:]}, nil
	case '%':
		return parseFormatCommand(cmd)
	case '?':
		if strings.TrimSpace(cmd[1:]) == "" {
			return AWKCommand{Text: "}"}, nil
		}
		return AWKCommand{Text: fmt.Sprintf("if (%s) {", cmd[1:])}, nil
	case '@':
		return NextSectionCommand{Section: section}, nil
	}

	return nil, fmt.Errorf("no such command: %q", line)
}

func parseNumberCommand(cmd string) (Command, error) {
	if cmd[0] == '$' {
		cmd = cmd[1:]
	}
	number, err := strconv.Atoi(cmd)
	if err != nil {
		return nil, err
	}
	return NumberCommand{Number: number}, nil
}

func parseFormatCommand(cmd string) (Command, error) {
	before, after, found := strings.Cut(cmd, ",")
	if !found {
		return nil, errors.New("no data for format found")
	}

	return FormatCommand{Format: before, Data: after}, nil
}

func awk(fileScript, fileData, fileOut *os.File) {
	cmd := exec.Command("awk")
	cmd.Stdin = fileData
	cmd.Stdout = fileOut
	cmd.Args = append(cmd.Args, "-f", fileScript.Name())
	log.Println(cmd)
	err := cmd.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			log.Println(exitError)
			log.Println(exitError.ProcessState)
			log.Println(exitError.Stderr)
			log.Println(exitError.ExitCode())

		} else {
			log.Println("There was an unknown error executing Script: ", err)
		}
	}
}
