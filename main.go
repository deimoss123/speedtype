package main

import (
	"fmt"
	"os"
	"speedtype/util"
	"time"

	// "github.com/charmbracelet/bubbles/cursor"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen int

const (
	screenSetup screen = iota
	screenTyping
	screenResults
)

type model struct {
	currentScreen screen

	words         []util.Word
	currentWord   int
	currentLetter int

	typingStarted bool

	typingStartTime time.Time
	typingEndTime   time.Time
	typingResult    util.TypingResult

	// terminal size
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		if key == "ctrl+c" || key == "ctrl+d" || key == "ctrl+z" {
			return m, tea.Quit
		}

		switch m.currentScreen {
		case screenSetup:
			switch key {
			case "enter":
				m.typingStarted = false
				m.words = util.GenerateWordList(50)
				m.currentScreen = screenTyping
			}
		case screenTyping:
			switch key {
			case " ":
				if m.currentLetter == len(m.words[m.currentWord].Letters) {
					m.currentWord++
					m.currentLetter = 0
				}
			case "backspace":
				if m.currentLetter == 0 {
					return m, nil
				}
				m.currentLetter--
			default:
				if len([]rune(key)) > 1 {
					return m, nil
				}

				if m.currentLetter >= len(m.words[m.currentWord].Letters) {
					return m, nil
				}

				if !m.typingStarted && m.currentWord == 0 && m.currentLetter == 0 {
					m.typingStarted = true
					m.typingStartTime = time.Now()
				}

				letter := m.words[m.currentWord].Letters[m.currentLetter]

				m.words[m.currentWord].Letters[m.currentLetter].Corrent = letter.Char == []rune(key)[0]

				if m.currentWord == len(m.words)-1 && m.currentLetter == len(m.words[m.currentWord].Letters)-1 {
					m.typingEndTime = time.Now()
					m.typingResult = util.CalcTypingResult(m.typingEndTime.Sub(m.typingStartTime).Milliseconds(), m.words)
					m.currentScreen = screenResults
				}

				m.currentLetter++
			}
		case screenResults:
			switch key {
			case "enter":
				m.currentScreen = screenTyping
				m.words = util.GenerateWordList(50)
				m.currentWord = 0
				m.currentLetter = 0
				m.typingStarted = false
			}
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	return m, nil
}

var (
	strGray      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	strWhite     = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	strRed       = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	strCursor    = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Background(lipgloss.Color("15"))
	strUnderline = lipgloss.NewStyle().Underline(true)
)

func (m model) View() string {
	switch m.currentScreen {
	case screenSetup:
		return "Press enter to start..."
	case screenTyping:
		var s string

		s += fmt.Sprintf("Words: %d/%d\n", m.currentWord+1, len(m.words))

		var textStr string

		for i := range m.words {
			wordStr := ""

			for j, letter := range m.words[i].Letters {
				if m.currentWord > i || (m.currentWord == i && m.currentLetter > j) {
					if letter.Corrent {
						wordStr += strWhite.Underline(m.currentWord == i).Render(string(letter.Char))
					} else {
						wordStr += strRed.Underline(m.currentWord == i).Render(string(letter.Char))
					}
					// } else if m.currentWord == i && m.currentLetter == j {
					// 	wordStr += strCursor(string(letter.Char))
				} else {
					wordStr += strGray.Underline(m.currentWord == i).Render(string(letter.Char))
				}
			}

			textStr += wordStr + " "
		}

		s += "\n\n" + lipgloss.NewStyle().Width(60).Align(lipgloss.Left).Render(textStr)

		return s
	case screenResults:
		return fmt.Sprintf("Results\n\nWPM: %.2f\n\nPress enter to type again...", m.typingResult.Wpm)
	}

	return ""
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

}
