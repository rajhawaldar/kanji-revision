package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xuri/excelize/v2"
)

func defaultStyle() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

func onyoumiStyle() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("26")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

func New(kanjis []JKanji) *model {
	styles := defaultStyle()
	onyoumiStyles := onyoumiStyle()
	kunyomiField := textinput.New()
	onyoumiField := textinput.New()
	kunyomiField.Focus()
	return &model{kanjis: kanjis, kunyomiField: kunyomiField, styles: styles, onyoumiField: onyoumiField, onyoumiStyles: onyoumiStyles}
}

func reviewAnswer(data []string, answers []string) (bool, []string) {
	var found bool
	var correct []string
	for _, k := range data {
		found = false
		for _, a := range answers {
			if k == a {
				found = true
			}
		}
		if found {
			correct = append(correct, k)
		}
		found = false
	}
	if len(correct) == 0 {
		return false, nil
	} else if len(correct) == len(data) {
		return true, correct
	} else {
		return false, correct
	}

}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getIndex(length int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(length)
}

func readInput(f *excelize.File) []JKanji {

	var kanjis []JKanji
	rows, err := f.Rows("漢字")
	if err != nil {
		log.Panic(err.Error())
	}
	rows.Next()
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			fmt.Println(err)
		}
		kanji := JKanji{
			Kanji:   row[2],
			Meaning: row[5],
		}
		kanji.Kunyomi = strings.Split(row[3], "、")
		// kanji.Kunyomi = strings.Split(row[3], ",")
		kanji.Onyoumi = strings.Split(row[4], "、")
		// kanji.Onyoumi = strings.Split(row[4], ",")
		kanjis = append(kanjis, kanji)
	}
	return kanjis
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.onyoumiField.Focus()
			m.onyoumiField, cmd = m.onyoumiField.Update(msg)
			m.kunyomiField.Blur()
			return m, cmd
		case "shift+tab":
			m.kunyomiField.Focus()
			m.kunyomiField, cmd = m.kunyomiField.Update(msg)
			m.onyoumiField.Blur()
			return m, cmd
		case "enter":
			m.attempt++
			if m.attempt == 3 {
				m.meaningField.SetValue(" (" + m.kanjis[m.index].Meaning + ") ")
			}
			m.kanjis[m.index].KunyomiAnswer = strings.Split(strings.Replace(m.kunyomiField.Value(), "。", "・", -1), "、")
			m.kanjis[m.index].OnyoumiAnswer = strings.Split(strings.Replace(m.onyoumiField.Value(), "。", "・", -1), "、")
			kunyoumiResult, correctK := reviewAnswer(m.kanjis[m.index].Kunyomi, m.kanjis[m.index].KunyomiAnswer)
			onyoumiResult, correctO := reviewAnswer(m.kanjis[m.index].Onyoumi, m.kanjis[m.index].OnyoumiAnswer)

			if kunyoumiResult && onyoumiResult {
				for {
					if len(m.alreadyCompleted) == len(m.kanjis) {
						return m, tea.Quit
					}
					index := getIndex(len(m.kanjis))
					if !contains(m.alreadyCompleted, index) {
						m.index = index
						break
					}
				}
				m.alreadyCompleted = append(m.alreadyCompleted, m.index)
				m.kunyomiField.Focus()
				m.onyoumiField.Blur()
				m.resultField.SetValue("")
				m.correctKField.SetValue("")
				m.correctOField.SetValue("")
				m.kunyomiField.SetValue("")
				m.onyoumiField.SetValue("")
				m.meaningField.SetValue("")
				m.attempt = 0
			} else {
				m.resultField.SetValue("Try again")
				if len(correctK) == 0 {
					m.correctKField.SetValue("All kunyoumi's are wrong")
				} else if len(correctK) == len(m.kanjis[m.index].Kunyomi) {
					m.correctKField.SetValue("All kunyoumi's are Correct!")
				} else {
					m.correctKField.SetValue("Correct kunyoumi: " + strings.Join(correctK[:], ",") + " (Remaining " + strconv.Itoa(len(m.kanjis[m.index].Kunyomi)-len(correctK)) + ")")
				}
				if len(correctO) == 0 {
					m.correctOField.SetValue("All onyoumi's are wrong")
				} else if len(correctO) == len(m.kanjis[m.index].Onyoumi) {
					m.correctOField.SetValue("All onyoumi's are Correct!")
				} else {
					m.correctOField.SetValue("Correct onyoumi: " + strings.Join(correctO[:], ",") + " (Remaining " + strconv.Itoa(len(m.kanjis[m.index].Onyoumi)-len(correctO)) + ")")
				}
			}
			return m, nil
		}
	}
	if m.kunyomiField.Focused() {
		m.kunyomiField, cmd = m.kunyomiField.Update(msg)
		return m, cmd
	}
	m.onyoumiField, cmd = m.onyoumiField.Update(msg)
	return m, cmd
}

func (m model) View() string {
	m.kunyomiField.Placeholder = "This kanji has " + strconv.Itoa(len(m.kanjis[m.index].Kunyomi)) + " Kunyomi"
	m.onyoumiField.Placeholder = "This kanji has " + strconv.Itoa(len(m.kanjis[m.index].Onyoumi)) + " Onyoumi"
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			"Please take a look at the given kanji "+m.kanjis[m.index].Kanji+m.meaningField.Value(),
			m.styles.InputField.Render(m.kunyomiField.View()),
			m.onyoumiStyles.InputField.Render(m.onyoumiField.View()),
			m.resultField.Value(),
			m.correctKField.Value(),
			m.correctOField.Value()),
	)
}

func main() {
	fileName := flag.String("file", "", "xlsx file to take input from")
	flag.Parse()
	f, err := excelize.OpenFile(*fileName)
	if err != nil {
		log.Panic(err.Error())
	}
	defer f.Close()
	if err != nil {
		log.Println(err.Error())
	}
	var kanjis = readInput(f)
	m := New(kanjis)

	logFile, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
