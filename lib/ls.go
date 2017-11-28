package sman

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"sort"
	"text/tabwriter"
	"io"
	"strings"
)

var (
	porcelainFlag bool
)

func filterSnippets(p string, slice SnippetSlice) (matched SnippetSlice) {
	r, err := regexp.Compile(p)
	checkError(err, "Invalid search pattern")
	for _, s := range slice {
		if r.MatchString(s.Name) ||
			r.MatchString(s.Command) ||
			r.MatchString(s.Desc) {
			matched = append(matched, s)
		}
	}
	return matched
}

func doLs(pattern string) {
	c := getConfig()
	snippets := getSnippets(pattern, fileFlag, c.SnippetDir, tagFlag)
	snippets = filterSnippets(pattern, snippets)
	sort.Sort(snippets)
	if porcelainFlag {
		doLsPorcelain(snippets)
	} else {
		doLsSlice(snippets, os.Stdout)
	}
}

func doLsSlice(snippets SnippetSlice, output io.Writer) {
	c := getConfig()
	w := new(tabwriter.Writer)
	w.Init(output, 25, 2, 0, ' ', 0)
	var prevFile string
	for _, s := range snippets {
		if s.File != prevFile {
			fmt.Fprintln(w, c.LsFilesColor.SprintFunc()(s.File+":"))
			prevFile = s.File
		}
		line := fmt.Sprintf("   %v\t[%v]\t%v", s.Name, displaySlice(s.Tags),
			displayString(s.Desc))
		fmt.Fprintln(w, line)
	}
	err := w.Flush()
	checkError(err, "Flush error..")
}

func doLsPorcelain(snippets SnippetSlice) {
	for _, s := range snippets {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("%v\t%v\t%v\t%v", s.File, s.Name, strings.Join(s.Tags, ","), s.Desc))
	}
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:     "ls [-f FILE] [-t TAGS] [PATTERN]",
	Aliases: []string{"l"},
	Short:   "List and search pattern in all available snippets",
	Long: `
List and search pattern in all available snippets,

PATTERN is regexp matched against snippet name, description and command.

Examples:
s ls add
	- List all snippet matching pattern "add"
s ls -f docker
	- List all snippets in file 'docker'
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var p string
		if len(args) > 0 {
			p = args[0]
		}
		doLs(p)
	},
}

func init() {
	RootCmd.AddCommand(lsCmd)
	lsCmd.Flags().BoolVarP(&porcelainFlag, "porcelain", "", false, "produce machine-readable output")
}
