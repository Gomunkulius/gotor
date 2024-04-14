package ui

import (
	"github.com/anacrolix/torrent"
	tea "github.com/charmbracelet/bubbletea"
	torrent2 "gotor/internal/torrent"
	"log"
	"time"
)

type MainModel struct {
	state ProgramState
	conn  *torrent.Client
	// Models
	inputModel   *InputModel
	programModel *ProgramModel
	chooseModel  *ChooseModel
}

type ProgramState int

const (
	Main ProgramState = iota
	Input
	Choose
)

type TickMsg time.Time

type ChangeStateMsg ProgramState

func NewModel(table *TorrentTable, conn *torrent.Client, storage torrent2.Storage) MainModel {
	inpModel := NewInputModel(0, 0, table, conn, storage)
	prgModel := NewProgramModel(0, 0, storage, table, keys)
	chooseModel := NewChooseModel(0, 0, table, storage, conn)
	return MainModel{
		state:        Main,
		inputModel:   inpModel,
		programModel: prgModel,
		chooseModel:  chooseModel,
		conn:         conn,
	}
}

// Send a message every second.
func tickEvery() tea.Cmd {
	return tea.Every(100*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type TorrentInfoUpdate string

func (m MainModel) Init() tea.Cmd {
	m.inputModel.Init()
	m.programModel.Init()
	m.chooseModel.Init()
	m.state = Choose
	return tickEvery()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case ChangeStateMsg:
		prev := m.state
		m.state = ProgramState(msg)
		switch ProgramState(msg) {
		case Main:
			log.Println("Change state to Main")
			switch prev {
			case Input:
				m.programModel.width = m.inputModel.width
				m.programModel.height = m.inputModel.height
			case Choose:
				m.programModel.width = m.chooseModel.width
				m.programModel.height = m.chooseModel.height
			}
			m.inputModel.inputField.Blur()
			m.programModel.table = m.inputModel.table
		case Input:
			log.Println("Change state to Input")
			m.inputModel.inputField.Focus()
			m.inputModel.width = m.programModel.width
			m.inputModel.height = m.programModel.height
		case Choose:
			log.Println("Change state to Choose")
			m.chooseModel.width = m.programModel.width
			m.chooseModel.height = m.programModel.height
		}

		return m, tickEvery()
	}
	switch m.state {
	case Input:
		inpMod, cmd := m.inputModel.Update(msg)
		inpModel, ok := inpMod.(InputModel)
		if !ok {
			panic("wrong type")
		}
		m.inputModel = &inpModel
		cmds = append(cmds, cmd)

	case Main:
		programMod, cmd := m.programModel.Update(msg)
		programModel, ok := programMod.(ProgramModel)
		if !ok {
			panic("wrong type")
		}
		m.programModel = &programModel
		cmds = append(cmds, cmd)
	case Choose:
		chooseMod, cmd := m.chooseModel.Update(msg)
		chooseModel, ok := chooseMod.(ChooseModel)
		if !ok {
			panic("wrong type")
		}
		m.chooseModel = &chooseModel
		cmds = append(cmds, cmd)
	default:
		return m, tickEvery()
	}
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	//textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	switch m.state {
	case Input:
		return m.inputModel.View()
	case Main:
		return m.programModel.View()
	case Choose:
		return m.chooseModel.View()
	}
	return m.programModel.View()
}
