package gneslib

import "../core"
import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type debugger struct {
	emu *gnes.Emulator

	cmdFuncMap map[string]func(*debugger, string) error
	cmdHelpMap map[string]string

	prompt,
	invalidCmdString,
	helpCmd,
	helpHint string
}

// Creates a new debugger containing an emulator with
// the file at the path `path`
func newDebugger(path string) (*debugger, error) {
	dbg := &debugger{}
	dbg.prompt = "gnes > "
	dbg.helpCmd = "h"

	dbg.invalidCmdString = fmt.Sprintf("Invalid command. Enter '%s' for a list of valid commands.", dbg.helpCmd)

	err := dbg.createCmdMap()
	if err != nil {
		return nil, err
	}
	emu, err := gnes.NewEmulator(path)
	if err != nil {
		return nil, err
	}
	dbg.emu = emu
	return dbg, nil
}

func (dbg *debugger) createCmdMap() error {
	dbg.cmdFuncMap = make(map[string]func(*debugger, string) error)
	dbg.cmdHelpMap = make(map[string]string)

	dbg.cmdFuncMap[dbg.helpCmd] = cmdHelp
	dbg.cmdHelpMap[dbg.helpCmd] = "Display this help message"
	return nil
}

func (dbg *debugger) getValidCommands() []string {
	commands := make([]string, len(dbg.cmdFuncMap))
	for cmd := range dbg.cmdFuncMap {
		commands = append(commands, cmd)
	}
	return commands
}

func (dbg *debugger) commandIsValid(cmd string) bool {
	validCommands := dbg.getValidCommands()
	for _, validCmd := range validCommands {
		if validCmd == cmd {
			return true
		}
	}
	return false
}

func getCommandFromInput(input string) (string, error) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return "", errors.New("No command found in input")
	}
	return args[0], nil
}

/***********************************************/
/*         Debug dispatcher functions          */
/***********************************************/

func cmdHelp(dbg *debugger, input string) error {
	for helpCmd, helpText := range dbg.cmdHelpMap {
		fmt.Printf("%s: %s\n", helpCmd, helpText)
	}
	return nil
}

/***********************************************/
/*                                             */
/***********************************************/
// Displays the prompt and returns the user input
func (dbg *debugger) showPrompt() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(dbg.prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return input, nil
}

func (dbg *debugger) showPromptAndDispatch() error {
	input, err := dbg.showPrompt()
	if err != nil {
		return err
	}
	input = strings.TrimSpace(input)
	cmd, err := getCommandFromInput(input)
	if err != nil {
		fmt.Println(dbg.invalidCmdString)
		return nil
	}
	commandIsValid := dbg.commandIsValid(cmd)
	if commandIsValid {
		dbg.cmdFuncMap[cmd](dbg, input)
		return nil
	}
	fmt.Println(dbg.invalidCmdString)
	return nil
}

func RunCLIDebugger(path string) error {
	dbg, err := newDebugger(path)
	if err != nil {
		return err
	}
	for true {
		dbg.showPromptAndDispatch()
	}
	return nil
}
