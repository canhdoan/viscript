package cli

import (
	"fmt"
	"github.com/corpusc/viscript/msg"
	tm "github.com/corpusc/viscript/rpc/terminalmanager"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

func (c *CliManager) PrintHelp(_ []string) error {
	p := fmt.Printf
	p("\n<< [- HELP -] >>\n\n")

	p("> stp\t\tStart a new terminal with process.\n\n")

	p("> ltp\t\tList terminal Ids with Attached Process Ids.\n")
	p("> sett <tId>\tSet given terminal Id as default for all following commands.\n")
	p("> setp <pId>\tSet given process Id as default for all following commands.\n\n")

	p("> cft\t\tGet out channel info of terminal with default Id.\n\n")
	// p("> lpub\t\tList all publishers. --TODO\n\n")
	// p("> ld\t\tList all dbus objects. --TODO\n")

	p("> clear(c)\t\tClear the terminal.\n")
	p("> quit(q)\tQuit from cli.\n\n")

	return nil
}

func (c *CliManager) Quit(_ []string) error {
	println("See ya again! :>")
	c.SessionEnd = true
	return nil
}

func (c *CliManager) ClearTerminal(_ []string) error {

	runtimeOs := runtime.GOOS

	if runtimeOs == "linux" || runtimeOs == "darwin" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if runtimeOs == "windows" {
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		println("Your platform is unsupported! I can't clear terminal screen :(.")
	}

	return nil
}

func (c *CliManager) ListTermIDsWithAttachedProcesses(_ []string) error {
	termsWithProcessIDs, err := GetTermIDsWithProcessIDs(c.Client)

	if err != nil {
		return err
	}

	fmt.Printf("\nTerminals(%d total) defaults marked with [$]:\n\n", len(termsWithProcessIDs))
	fmt.Println("Index\tTerminalID\t\tAttached Process ID")
	fmt.Println()
	for index, term := range termsWithProcessIDs {
		fmt.Printf("%d\t", index)

		// mark selected default terminal id
		if term.TerminalId == c.ChosenTerminalId {
			fmt.Printf("[ %d ]\t", term.TerminalId)
		} else {
			fmt.Printf("  %d\t", term.TerminalId)
		}

		// mark selected default process id
		if term.AttachedProcessId == c.ChosenProcessId {
			fmt.Printf("[ %d ]\t", term.AttachedProcessId)
		} else {
			fmt.Printf("  %d\t", term.AttachedProcessId)
		}
		fmt.Printf("\n")
	}

	return nil
}

func (c *CliManager) SetDefaultTerminalId(args []string) error {
	if len(args) == 0 {
		fmt.Printf("\n\nPass the terminal Id as arguement please.")
		return nil
	}

	termId, err := strconv.Atoi(args[0])
	if err != nil || termId < 1 {
		fmt.Printf("\n\nArgument should be a number > 0, not %s\n\n", args[0])
		return nil
	}

	c.ChosenTerminalId = msg.TerminalId(termId)
	return nil
}

func (c *CliManager) SetDefaultProcessId(args []string) error {
	if len(args) == 0 {
		fmt.Printf("\n\nArgument should be a number > 0, not %s\n\n", args[0])
		return nil
	}

	processId, err := strconv.Atoi(args[0])
	if err != nil || processId < 1 {
		fmt.Printf("\n\nArgument should be a number > 0, not %s\n\n", args[0])
	}

	c.ChosenProcessId = msg.ProcessId(processId)
	return nil
}

func (c *CliManager) ShowChosenTermChannelInfo(_ []string) error {
	if c.ChosenTerminalId == 0 {
		fmt.Printf("\nDefault terminal Id is not set.\n\n")
		return nil
	}

	response, err := c.Client.SendToRPC("GetTermChannelInfo", []string{fmt.Sprintf("%d", c.ChosenTerminalId)})
	if err != nil {
		return err
	}

	var channelInfo msg.ChannelInfo
	err = msg.Deserialize(response, &channelInfo)
	if err != nil {
		return err
	}

	// TODO: print structured channel info
	// fmt.Printf("Terminal out channel info with subscribers (%d total):\n", len(channelInfo.Subscribers))
	// fmt.Printf("")
	fmt.Printf("Channel Info:\n%+v\n\n", channelInfo)

	return nil
}

func (c *CliManager) StartTerminalWithProcess(_ []string) error {
	fmt.Println("startTerminalWithProcess()")
	response, err := c.Client.SendToRPC("StartTerminalWithProcess", []string{})
	if err != nil {
		return err
	}

	var newID msg.TerminalId
	err = msg.Deserialize(response, &newID)
	if err != nil {
		return err
	}

	fmt.Println("New terminal was created with ID", newID)

	return nil
}

func GetTerminalIDs(client *tm.RPCClient) ([]msg.TerminalId, error) {
	response, err := client.SendToRPC("ListTerminalIDs", []string{})
	if err != nil {
		return []msg.TerminalId{}, err
	}

	var termIDs []msg.TerminalId
	err = msg.Deserialize(response, &termIDs)
	if err != nil {
		return []msg.TerminalId{}, err
	}
	return termIDs, nil
}

func GetTermIDsWithProcessIDs(client *tm.RPCClient) ([]msg.TermAndAttachedProcessID, error) {
	response, err := client.SendToRPC("ListTIDsWithProcessIDs", []string{})
	if err != nil {
		return []msg.TermAndAttachedProcessID{}, err
	}

	var termsAndAttachedProcesses []msg.TermAndAttachedProcessID
	err = msg.Deserialize(response, &termsAndAttachedProcesses)
	if err != nil {
		return []msg.TermAndAttachedProcessID{}, err
	}
	return termsAndAttachedProcesses, nil
}

func GetProcessIDs(client *tm.RPCClient) ([]msg.ProcessId, error) {
	response, err := client.SendToRPC("ListProcessIDs", []string{})
	if err != nil {
		return []msg.ProcessId{}, err
	}

	var processIDs []msg.ProcessId
	err = msg.Deserialize(response, &processIDs)
	if err != nil {
		return []msg.ProcessId{}, err
	}
	return processIDs, nil
}
