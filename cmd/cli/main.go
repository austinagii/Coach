package main

import (
	"aisu.ai/api/v2/cmd/cli/api"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mitchellh/go-wordwrap"
)

var (
	userId string
	chatId string
)

type (
	errMsg error
)

type model struct {
	prompt    string
	userInput textinput.Model
	err       error
}

func main() {
	var err error
	userId, err = api.CreateUser("Kadeem")
	if err != nil {
		log.Printf("An error occurred while creating a new user: %s", err.Error())
		os.Exit(1)
	}

	var message string
	chatId, message, err = api.CreateChat(userId)
	if err != nil {
		log.Printf("An error occurred while creating a new chat: %s", err.Error())
		os.Exit(1)
	}

	p := tea.NewProgram(newModel(message))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func newModel(initialMsg string) model {
	ti := textinput.New()
	ti.Placeholder = "Type your response here"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	return model{
		prompt:    initialMsg,
		userInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			assistantResponse, err := api.Respond(userId, chatId, m.userInput.Value())
			if err != nil {
				m.err = err
				return m, nil
			}

			m.prompt = assistantResponse
			m.userInput.SetValue("")
			return m, textinput.Blink
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.userInput, cmd = m.userInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	wrappedPrompt := wordwrap.WrapString(m.prompt, uint(m.userInput.Width))
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		wrappedPrompt,
		m.userInput.View(),
		"(esc to quit)",
	) + "\n"
}
