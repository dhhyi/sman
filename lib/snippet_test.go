package sman

import (
	"reflect"
	"testing"
)

func TestSnippetReplacePlaceholders(t *testing.T) {
	tests := []struct {
		name         string
		command      string
		placeholders []Placeholder
		wantCommand  string
	}{
		{"no placeholder", "hello world", []Placeholder(nil), "hello world"},
		{"single patterns", "hello <<name>>",
			[]Placeholder{
				Placeholder{
					Name:     "name",
					Patterns: []string{"<<name>>"},
					Input:    "test",
				},
			},
			"hello test",
		},
		{"multiple patterns", "hello <<name#desc>> sup <<name>>",
			[]Placeholder{
				Placeholder{
					Name:     "name",
					Patterns: []string{"<<name#desc>>", "<<name>>"},
					Input:    "test",
				},
			},
			"hello test sup test",
		},
	}
	for _, tt := range tests {
		s := &Snippet{
			Command:      tt.command,
			Placeholders: tt.placeholders,
		}
		s.ReplacePlaceholders()
		if s.Command != tt.wantCommand {
			t.Errorf("%q. ReplacePlaceholders() = %#v, want %#v", tt.name, s.Command, tt.wantCommand)
		}
	}
}

func TestSnippetParseCommand(t *testing.T) {
	tests := []struct {
		name             string
		command          string
		wantPlaceholders []Placeholder
	}{
		{"no placeholder", "hello world", []Placeholder(nil)},
		{"full placeholder", "hello <<name(one,two)#desc>>",
			[]Placeholder{
				Placeholder{
					Name:     "name",
					Desc:     "#desc",
					Options:  []string{"one", "two"},
					Patterns: []string{"<<name(one,two)#desc>>"},
				},
			},
		},
		{"multiple patterns", `hello <<name#desc>> <<name>>`,
			[]Placeholder{
				Placeholder{
					Name:     "name",
					Desc:     "#desc",
					Patterns: []string{"<<name#desc>>", "<<name>>"},
				},
			},
		},
		{"multiple placeholders", `hello <<name>> <<last>>`,
			[]Placeholder{
				Placeholder{
					Name:     "name",
					Patterns: []string{"<<name>>"},
				},
				Placeholder{
					Name:     "last",
					Patterns: []string{"<<last>>"},
				},
			},
		},
	}
	for _, tt := range tests {
		s := &Snippet{
			Command: tt.command,
		}
		s.ParseCommand()
		if !reflect.DeepEqual(s.Placeholders, tt.wantPlaceholders) {
			t.Errorf("%q. ParseCommand() = %#v, want %#v", tt.name, s.Placeholders, tt.wantPlaceholders)
		}
	}
}

func TestInitSnippets(t *testing.T) {
	tests := []struct {
		name         string
		snippetMap   map[string]Snippet
		file         string
		dir          string
		wantSnippets SnippetSlice
	}{
		{"t",
			map[string]Snippet{"echo": Snippet{Command: "hello world"}},
			"file", "",
			SnippetSlice{
				Snippet{
					Name:    "echo",
					Command: "hello world",
					File:    "file",
				},
			},
		},
		{"ext_command",
			map[string]Snippet{"ext_command": Snippet{}},
			"examples", testPath,
			SnippetSlice{
				Snippet{
					Name:    "ext_command",
					Command: "test command",
					File:    "examples",
				},
			},
		},
		{"invalid snippet",
			map[string]Snippet{"invalid_snippet": Snippet{}},
			"examples", testPath,
			SnippetSlice(nil),
		},
	}
	makeTestFiles(t)
	for _, tt := range tests {
		if gotSnippets := initSnippets(tt.snippetMap, tt.file, tt.dir); !reflect.DeepEqual(gotSnippets, tt.wantSnippets) {
			t.Errorf("%q. initSnippets() = %#v, want %#v", tt.name, gotSnippets, tt.wantSnippets)
		}
	}
	defer cleanTestFiles(t)
}

func TestFilterByTag(t *testing.T) {
	var snippet = Snippet{Name:"s"}
	var snippet1 = Snippet{Name:"s1", Tags:[]string{"tag1"}}
	var snippet12 = Snippet{Name:"s12", Tags:[]string{"tag1", "tag2"}}
	var snippet2 = Snippet{Name:"s2", Tags:[]string{"tag2"}}
	var snippet3 = Snippet{Name:"s3", Tags:[]string{"tag3"}}
	var all = SnippetSlice{snippet, snippet1, snippet12, snippet2, snippet3}

	tests := []struct {
		name        string
		snippets    SnippetSlice
		tag         string
		wantMatched SnippetSlice
	}{
		{"no tag filter",
			all,
			"",
			SnippetSlice{},
		},
		{"single tag filter tag1",
			all,
			"tag1",
			SnippetSlice{snippet1, snippet12},
		},
		{"single tag filter tag3",
			all,
			"tag3",
			SnippetSlice{snippet3},
		},
		{"multiple tag filter tag1 and tag2",
			all,
			"tag1+tag2",
			SnippetSlice{snippet12},
		},
		{"multiple tag filter tag1 and tag3",
			all,
			"tag1+tag3",
			SnippetSlice{},
		},
		{"multiple tag filter tag1 or tag2",
			all,
			"tag1,tag2",
			SnippetSlice{snippet1, snippet12, snippet2},
		},
		{"multiple tag filter tag1 or tag3",
			all,
			"tag1,tag3",
			SnippetSlice{snippet1, snippet12, snippet3},
		},
		{"multiple tag filter tag1 and tag2 or tag3",
			all,
			"tag1+tag2,tag3",
			SnippetSlice{snippet12, snippet3},
		},
		{"multiple tag filter tag1 and tag2 or tag1 and tag3",
			all,
			"tag1+tag2,tag1+tag3",
			SnippetSlice{snippet12},
		},
	}
	for _, tt := range tests {
		if gotMatched := filterByTag(tt.snippets, tt.tag);
			!((gotMatched.Len() == 0 && tt.wantMatched.Len() == 0) || reflect.DeepEqual(gotMatched, tt.wantMatched)) {
			t.Errorf("%q. filterByTag() = %v, want %v", tt.name, gotMatched, tt.wantMatched)
		}
	}
}
