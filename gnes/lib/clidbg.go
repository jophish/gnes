package gneslib

import "../core"
import (
	"errors"
	"fmt"
	"github.com/chzyer/readline"
	"strconv"
	"strings"
)

type debugger struct {
	emu *gnes.Emulator

	cmdFuncMap map[string]func(*debugger, string) error
	cmdHelpMap map[string]string

	prompt,
	invalidCmdString,
	helpCmd,
	quitCmd,
	helpHint string

	runN uint64

	continueFlag,
	runFlag bool

	breakpoints []uint16
}

// Creates a new debugger containing an emulator with
// the file at the path `path`
func newDebugger(path string) (*debugger, error) {
	dbg := &debugger{}
	dbg.prompt = "[%#04x] gnes > "
	dbg.helpCmd = "h"
	dbg.quitCmd = "q"

	dbg.runFlag = true
	dbg.continueFlag = false

	dbg.runN = 0

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

	dbg.cmdFuncMap[dbg.quitCmd] = cmdQuit
	dbg.cmdHelpMap[dbg.quitCmd] = "Quit the debugger"

	dbg.cmdFuncMap["s"] = cmdStepInstruction
	dbg.cmdHelpMap["s"] = "Step a single instruction"

	dbg.cmdFuncMap["sn"] = cmdStepInstructions
	dbg.cmdHelpMap["sn"] = "Step n instructions (sn n)"

	dbg.cmdFuncMap["rs"] = cmdReadAddress
	dbg.cmdHelpMap["rs"] = "Read single memory address 'addr' (rs addr)"

	dbg.cmdFuncMap["rn"] = cmdReadN
	dbg.cmdHelpMap["rn"] = "Read n memory addresses starting from 'addr' (rn addr n)"

	dbg.cmdFuncMap["ris"] = cmdReadInstructionAddress
	dbg.cmdHelpMap["ris"] = "Read single memory address 'addr' as instruction (rs addr)"

	dbg.cmdFuncMap["rin"] = cmdReadNInstructions
	dbg.cmdHelpMap["rin"] = "Read n memory addresses as instructions starting from 'addr' (rn addr n)"

	dbg.cmdFuncMap["r"] = cmdShowRegisters
	dbg.cmdHelpMap["r"] = "Show contents of internal CPU regsiters"

	dbg.cmdFuncMap["bp"] = setBreakpoint
	dbg.cmdHelpMap["bp"] = "Add breakpoint at address 'addr' (bp addr)"

	dbg.cmdFuncMap["sbp"] = showBreakpoints
	dbg.cmdHelpMap["sbp"] = "Show all breakpoints"

	dbg.cmdFuncMap["rbp"] = removeBreakpoint
	dbg.cmdHelpMap["rbp"] = "Remove breakpoint at address 'addr' (rbp addr)"

	dbg.cmdFuncMap["c"] = continueExecution
	dbg.cmdHelpMap["c"] = "Continue execution until reaching a breakpoint"

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

func splitString(input string) []string {
	return strings.Split(input, " ")
}

func getArgs(input string) []string {
	split := splitString(input)
	if len(split) > 0 {
		return split[1:]
	}
	return split
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

func cmdQuit(dbg *debugger, input string) error {
	dbg.runFlag = false
	return nil
}

func continueExecution(dbg *debugger, input string) error {
	dbg.continueFlag = true
	return nil
}

func setBreakpoint(dbg *debugger, input string) error {
	args := getArgs(input)
	if len(args) != 1 {
		return errors.New("Command requires single hexadecimal argument")
	}
	addr, err := strconv.ParseUint(args[0], 16, 16)
	if err != nil {
		return errors.New("Command requires single hexadecimal argument")
	}

	addr16 := uint16(addr)

	for _, bp := range dbg.breakpoints {
		if bp == addr16 {
			return fmt.Errorf("Breakpoint already set for %#04x!", addr16)
		}
	}

	dbg.breakpoints = append(dbg.breakpoints, addr16)
	fmt.Printf("Breakpoint set for %#04x\n", addr16)
	return nil
}

func showBreakpoints(dbg *debugger, input string) error {
	if len(dbg.breakpoints) == 0 {
		fmt.Println("No breakpoints set!")
	} else {
		fmt.Println("Breakpoints set for:")
		for _, bp := range dbg.breakpoints {
			fmt.Printf("    %#04x\n", bp)
		}
	}
	return nil
}

func removeBreakpoint(dbg *debugger, input string) error {
	args := getArgs(input)
	if len(args) != 1 {
		return errors.New("Command requires single hexadecimal argument")
	}
	addr, err := strconv.ParseUint(args[0], 16, 16)
	if err != nil {
		return errors.New("Command requires single hexadecimal argument")
	}

	addr16 := uint16(addr)

	for i, bp := range dbg.breakpoints {
		if bp == addr16 {
			copy(dbg.breakpoints[i:], dbg.breakpoints[i+1:])
			dbg.breakpoints[len(dbg.breakpoints)-1] = 0
			dbg.breakpoints = dbg.breakpoints[:len(dbg.breakpoints)-1]
			fmt.Printf("Removed breakpoint at %#04x\n", addr16)
			return nil
		}
	}
	return fmt.Errorf("No breakpoint was set at %#04x", addr16)
}

func cmdStepInstruction(dbg *debugger, input string) error {
	dbg.runN = 1
	return nil
}

func cmdStepInstructions(dbg *debugger, input string) error {
	args := getArgs(input)
	if len(args) != 1 {
		return errors.New("Command requires single integer argument")
	}
	numIns, err := strconv.ParseUint(args[0], 10, 16)
	if err != nil {
		return errors.New("Argument must be integer")
	}

	dbg.runN = numIns
	return nil
}

func cmdReadInstructionAddress(dbg *debugger, input string) error {
	args := getArgs(input)
	if len(args) != 1 {
		return errors.New("Command requires single hexadecimal argument")
	}
	addr, err := strconv.ParseUint(args[0], 16, 16)
	if err != nil {
		return errors.New("Command requires single hexadecimal argument")
	}
	val, err := dbg.emu.GetOpMnemonic(uint16(addr))
	if err != nil {
		return err
	}
	fmt.Printf("    [%#04x]: %s\n", addr, val)
	return nil
}

func cmdReadNInstructions(dbg *debugger, input string) error {
	args := getArgs(input)
	if len(args) != 2 {
		return errors.New("Command requires hex, int args")
	}
	addr, err := strconv.ParseUint(args[0], 16, 16)
	if err != nil {
		return errors.New("First argument must be hexadecimal")
	}

	numAddr, err := strconv.ParseUint(args[1], 10, 16)
	if err != nil {
		return errors.New("Second argument must be integer")
	}

	var i uint64
	for i = 0; i < numAddr; i++ {
		val, err := dbg.emu.GetOpMnemonic(uint16(addr))
		if err != nil {
			return err
		}
		fmt.Printf("    [%#04x]: %s\n", addr, val)
		opLength, err := dbg.emu.GetOpLength(uint16(addr))
		if err != nil {
			return err
		}
		addr += uint64(opLength)
	}
	return nil
}

func cmdReadAddress(dbg *debugger, input string) error {
	args := getArgs(input)
	if len(args) != 1 {
		return errors.New("Command requires single hexadecimal argument")
	}
	addr, err := strconv.ParseUint(args[0], 16, 16)
	if err != nil {
		return errors.New("Command requires single hexadecimal argument")
	}
	val, err := dbg.emu.ReadAddr(uint16(addr))
	if err != nil {
		return err
	}
	fmt.Printf("    [%#04x]: (%#02x)\n", addr, val)
	return nil
}

func cmdReadN(dbg *debugger, input string) error {
	args := getArgs(input)
	if len(args) != 2 {
		return errors.New("Command requires hex, int args")
	}
	addr, err := strconv.ParseUint(args[0], 16, 16)
	if err != nil {
		return errors.New("First argument must be hexadecimal")
	}

	numAddr, err := strconv.ParseUint(args[1], 10, 16)
	if err != nil {
		return errors.New("Second argument must be integer")
	}

	var i uint64
	for i = 0; i < numAddr; i++ {
		val, err := dbg.emu.ReadAddr(uint16(addr + i))
		if err != nil {
			return err
		}
		fmt.Printf("    [%#04x]: (%#02x)\n", addr+i, val)
	}
	return nil
}

func cmdShowRegisters(dbg *debugger, input string) error {
	regs := dbg.emu.GetCPUState()
	fmt.Printf("    PC: %#04x\n", regs.PC)
	fmt.Printf("    SP: %#02x\n", regs.SP)
	fmt.Printf("    A: %#02x, X: %#02x, Y: %#02x\n", regs.A, regs.X, regs.Y)
	fmt.Printf("    Flag Register:\n")
	fmt.Printf("        N: %t, V: %t, B: %t, D: %t\n",
		regs.N, regs.V, regs.B, regs.D)
	fmt.Printf("        I: %t, Z: %t, C: %t\n",
		regs.I, regs.Z, regs.C)

	return nil
}

/***********************************************/
/*                                             */
/***********************************************/

// Displays the prompt and returns the user input
func (dbg *debugger) showPrompt() (string, error) {
	input, err := readline.Line(fmt.Sprintf(dbg.prompt, dbg.emu.GetPC()))
	if err != nil {
		return "", err
	}
	return input, nil
}

// checkBreakpoint checks if the current execution has reached
// a breakpoint. If so, a message is printed and any further
// execution is halted.
func (dbg *debugger) checkBreakpoint() {
	pc := dbg.emu.GetPC()
	for _, bp := range dbg.breakpoints {
		if pc == bp {
			dbg.runN = 0
			dbg.continueFlag = false
			fmt.Printf("Halted for breakpoint at %#04x\n", pc)
		}
	}
}

func (dbg *debugger) stepExecutionAndHaltOnError() error {
	err := dbg.emu.Step()
	if err != nil {
		dbg.runN = 0
		dbg.continueFlag = false
		fmt.Println(err)
		return err
	}
	return nil
}

func (dbg *debugger) showPromptAndDispatch() error {

	if dbg.runN > 0 {
		err := dbg.stepExecutionAndHaltOnError()
		if err != nil {
			return err
		}
		dbg.runN -= 1
		dbg.checkBreakpoint()
		return nil
	}

	if dbg.continueFlag {
		err := dbg.stepExecutionAndHaltOnError()
		if err != nil {
			return err
		}
		dbg.checkBreakpoint()
		return nil
	}

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
		err = dbg.cmdFuncMap[cmd](dbg, input)
		if err != nil {
			fmt.Println(err)
		}
		return err
	}
	fmt.Println(dbg.invalidCmdString)
	return nil
}

func RunCLIDebugger(path string) error {
	dbg, err := newDebugger(path)
	if err != nil {
		return err
	}
	for dbg.runFlag {
		dbg.showPromptAndDispatch()
	}
	return nil
}
