package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

type JKanji struct {
	Kanji         string
	Kunyomi       []string
	Onyoumi       []string
	KunyomiAnswer []string
	OnyoumiAnswer []string
	Meaning       string
}
type model struct {
	kanjis           []JKanji
	kunyomiField     textinput.Model
	onyoumiField     textinput.Model
	styles           *Styles
	onyoumiStyles    *Styles
	width            int
	height           int
	index            int
	resultField      textinput.Model
	correctKField    textinput.Model
	correctOField    textinput.Model
	attempt          int
	meaningField     textinput.Model
	alreadyCompleted []int
}
