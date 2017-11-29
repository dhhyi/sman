package sman

import (
	"github.com/spf13/viper"
	"github.com/fatih/color"
)

// Config file struct
type Config struct {
	SnippetDir                 string
	ExecConfirm, AppendHistory bool
	LsFilesColor               *color.Color
	MinMatchPercentage         float64
}

func init() {
	viper.SetDefault("snippet_dir", "~/snippets")
	viper.SetDefault("append_history", "true")
	viper.SetDefault("exec_confirm", "true")
	viper.SetDefault("ls_color_files", "34")
	viper.SetDefault("min_match_percentage", "0")
}

//getConfig reads config and returns struct
func getConfig() (c Config) {
	c.SnippetDir = expandPath(viper.GetString("snippet_dir"))
	c.AppendHistory = viper.GetBool("append_history")
	c.ExecConfirm = viper.GetBool("exec_confirm")
	c.LsFilesColor = parseColor(viper.GetString("ls_color_files"))
	c.MinMatchPercentage = viper.GetFloat64("min_match_percentage")
	if c.MinMatchPercentage < 0 || c.MinMatchPercentage > 1 {
		panic("only values between 0 and 1 are allowed for min_match_percentage")
	}
	return c
}
