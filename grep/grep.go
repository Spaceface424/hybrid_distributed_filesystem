// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// grep searches file contents using regular Expressions.
//
// Synopsis:
//
//	grep [-clFivnhqre] [FILE]...
//
// Options:
//
//  -c, --Count                Just show Counts
//  -l, --files-with-matches   list only files
//  -F, --Fixed-strings        Match using Fixed strings
//  -i, --ignore-case          case-insensitive matching
//  -v, --Invert-match         Print only non-matching lines
//  -n, --line-Number          Show line Numbers
//  -h, --no-filename          Suppress file name prefixes on output
//  -q, --Quiet                Don't print matches; exit on first match
//  -r, --Recursive            Recursive
//  -e, --regexp string        Pattern to match

package grep

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

var errQuiet = fmt.Errorf("not found")

type Params struct {
	Expr, Logfile, Hosts string
	Headers, Invert, Recursive, CaseInsensitive, Fixed,
	NoShowMatch, Quiet, Count, Number, Swim bool
}

type grepCommand struct {
	rc   io.ReadCloser
	name string
}

func StructToProto(p Params) *GRPCParams {
	return &GRPCParams{
		Expr:            p.Expr,
		Logfile:         p.Logfile,
		Hosts:           p.Hosts,
		Headers:         p.Headers,
		Invert:          p.Invert,
		Recursive:       p.Recursive,
		CaseInsensitive: p.CaseInsensitive,
		Fixed:           p.Fixed,
		NoShowMatch:     p.NoShowMatch,
		Quiet:           p.Quiet,
		Count:           p.Count,
		Number:          p.Number,
	}
}

func ProtoToStruct(p *GRPCParams) Params {
	return Params{
		Expr:            p.GetExpr(),
		Logfile:         p.GetLogfile(),
		Hosts:           p.GetHosts(),
		Headers:         p.GetHeaders(),
		Invert:          p.GetInvert(),
		Recursive:       p.GetRecursive(),
		CaseInsensitive: p.GetCaseInsensitive(),
		Fixed:           p.GetFixed(),
		NoShowMatch:     p.GetNoShowMatch(),
		Quiet:           p.GetQuiet(),
		Count:           p.GetCount(),
		Number:          p.GetNumber(),
	}
}

func ParseParams() Params {
	p := Params{}
	flag.StringVarP(&p.Hosts, "hosts_file", "o", "", "path to file of host addresses")
	flag.StringVarP(&p.Expr, "regexp", "e", "", "Pattern to match")
	flag.BoolVarP(&p.Headers, "no-filename", "h", false, "Suppress file name prefixes on output")
	flag.BoolVarP(&p.Invert, "Invert-match", "v", false, "Print only non-matching lines")
	flag.BoolVarP(&p.Recursive, "Recursive", "r", false, "Recursive")
	flag.BoolVarP(&p.NoShowMatch, "files-with-matches", "l", false, "list only files")
	flag.BoolVarP(&p.Count, "Count", "c", false, "Just show Counts")
	flag.BoolVarP(&p.CaseInsensitive, "ignore-case", "i", false, "case-insensitive matching")
	flag.BoolVarP(&p.Number, "line-Number", "n", false, "Show line Numbers")
	flag.BoolVarP(&p.Fixed, "Fixed-strings", "F", false, "Match using Fixed strings")
	flag.BoolVarP(&p.Quiet, "Quiet", "q", false, "Don't print matches; exit on first match")
	flag.BoolVarP(&p.Quiet, "silent", "s", false, "Don't print matches; exit on first match")
	flag.BoolVarP(&p.Swim, "SWIM", "m", false, "Set True to execute swim cmd")
	flag.Parse()

	return p
}

func main() {
	if err := Command(os.Stdin, os.Stdout, os.Stderr, ParseParams(), flag.Args()).Run(); err != nil {
		if err == errQuiet {
			os.Exit(1)
		}
		log.Fatal(err)
	}
}

// cmd contains the actually business logic of grep
type cmd struct {
	stdin  io.ReadCloser
	stdout *bufio.Writer
	stderr io.Writer
	args   []string
	Params
	matchCount int
	showName   bool
}

func Command(stdin io.ReadCloser, stdout io.Writer, stderr io.Writer, p Params, args []string) *cmd {
	return &cmd{
		stdin:  stdin,
		stdout: bufio.NewWriter(stdout),
		stderr: stderr,
		Params: p,
		args:   args,
	}
}

// grep reads data from the os.File embedded in grepCommand.
// It matches each line against the re and prints the matching result
// If we are only looking for a match, we exit as soon as the condition is met.
// "match" means result of re.Match == match flag.
func (c *cmd) grep(f *grepCommand, re *regexp.Regexp) (ok bool) {
	r := bufio.NewScanner(f.rc)
	defer f.rc.Close()
	var lineNum int
	for r.Scan() {
		line := r.Text()
		var m bool
		switch {
		case c.Fixed && c.CaseInsensitive:
			m = strings.Contains(strings.ToLower(line), strings.ToLower(c.Expr))
		case c.Fixed && !c.CaseInsensitive:
			m = strings.Contains(line, c.Expr)
		default:
			m = re.MatchString(line)
		}
		if m != c.Invert {
			// in Quiet mode, exit before the first match
			if c.Quiet {
				return false
			}
			c.printMatch(f, line, lineNum+1, m)
			if c.NoShowMatch {
				break
			}
		}
		lineNum++
	}
	c.stdout.Flush()
	return true
}

func (c *cmd) printMatch(cmd *grepCommand, line string, lineNum int, match bool) {
	if match == !c.Invert {
		c.matchCount++
	}
	if c.Count {
		return
	}
	// at this point, we have committed to writing a line
	defer func() {
		c.stdout.WriteByte('\n')
	}()
	// if showName, write name to stdout
	if c.showName {
		c.stdout.WriteString(cmd.name)
	}
	// if dont show match, then newline and return, we are done
	if c.NoShowMatch {
		return
	}
	if match == !c.Invert {
		// if showName, need a :
		if c.showName {
			c.stdout.WriteByte(':')
		}
		// if showing line Number, print the line Number then a :
		if c.Number {
			c.stdout.Write(strconv.AppendUint(nil, uint64(lineNum), 10))
			c.stdout.WriteByte(':')
		}
		// now write the line to stdout
		c.stdout.WriteString(line)
	}
}

func (c *cmd) Run() error {
	defer c.stdout.Flush()
	// parse the Expression into valid regex
	if c.Expr != "" {
		c.args = append([]string{c.Expr}, c.args...)
	}

	r := ".*"
	if len(c.args) > 0 {
		r = c.args[0]
	}
	if c.CaseInsensitive && !bytes.HasPrefix([]byte(r), []byte("(?i)")) && !c.Fixed {
		r = "(?i)" + r
	}
	var re *regexp.Regexp
	if !c.Fixed {
		re = regexp.MustCompile(r)
	} else if c.Expr == "" {
		c.Expr = c.args[0]
	}

	// if len(c.args) < 2, then we read from stdin
	if len(c.args) < 2 {
		if !c.grep(&grepCommand{c.stdin, "<stdin>"}, re) {
			return nil
		}
	} else {
		c.showName = (len(c.args[1:]) >= 1 || c.Recursive || c.NoShowMatch) && !c.Headers
		var ok bool
		for _, v := range c.args[1:] {
			err := filepath.Walk(v, func(name string, fi os.FileInfo, err error) error {
				if err != nil {
					fmt.Fprintf(c.stderr, "grep: %v: %v\n", name, err)
					return nil
				}
				if fi.IsDir() && !c.Recursive {
					fmt.Fprintf(c.stderr, "grep: %v: Is a directory\n", name)
					return filepath.SkipDir
				}
				fp, err := os.Open(name)
				if err != nil {
					fmt.Fprintf(c.stderr, "can't open %s: %v\n", name, err)
					return nil
				}
				defer fp.Close()
				if !c.grep(&grepCommand{fp, name}, re) {
					ok = true
					return nil
				}
				return nil
			})
			if ok {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
	if c.Quiet {
		return errQuiet
	}
	if c.Count {
		c.stdout.WriteString(c.Logfile + ":")
		c.stdout.Write(strconv.AppendUint(nil, uint64(c.matchCount), 10))
		c.stdout.WriteByte('\n')
	}
	return nil
}
