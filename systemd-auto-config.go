package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type service struct {
	name string
	unit struct {
		description string
		after       string
		requires    string
		wants       string
	}
	core struct {
		user             string
		group            string
		workingDirectory string
		environment      []string
		kind             string
		remainAfterExit  bool
		start            struct {
			pre   string
			start string
			post  string
		}
		reload string
		stop   struct {
			stop string
			post string
		}
	}
	install struct {
		alias     string
		autoStart bool
	}
}

var magenta = color.New(color.FgMagenta)
var cyan = color.New(color.FgCyan)

func main() {
	TestWritePermissions()
	color.Magenta("-------------------------------------------------------------------------------")
	fmt.Println("An interactive program to Automate systemd services creation by Sid Sun.")
	fmt.Println("Licensed under the MIT License.")
	fmt.Println("By using this program, you agree to abide by the MIT License.")
	fmt.Println("Copyright (c) 2019 Sidharth Soni (Sid Sun) <sid@sidsun.com>.")
	color.Magenta("-------------------------------------------------------------------------------")
	color.Blue("Let's  get started!\n\n")
	input := TakeInput()
	if input == 4 {
		return
	}
	var Service service
	switch input {
	case 1:
		Service.core.kind = "simple"
	case 2:
		Service.core.kind = "forking"
	case 3:
		Service.core.kind = "oneshot"
	}
	_, _ = cyan.Print("\nService Name: ")
	Service.name = GetString(false, true)
	//Unit section begins
	_, _ = cyan.Print("Service Description: ")
	Service.unit.description = GetString(false, false)
	fmt.Println("Please enter the services your service should start after (empty for none)")
	_, _ = cyan.Print("After: ")
	Service.unit.after = GetString(true, false)
	fmt.Println("Please enter the hard dependencies for your service (empty for none)")
	_, _ = cyan.Print("Requires: ")
	Service.unit.requires = GetString(true, false)
	fmt.Println("Please enter the soft dependencies for your service (empty for none)")
	_, _ = cyan.Print("Wants: ")
	Service.unit.wants = GetString(true, false)
	//Unit section ends
	//Service section begins
	fmt.Println("Please enter the user the service should as (empty to not specify default: root)")
	_, _ = cyan.Print("User: ")
	Service.core.user = GetString(true, true)
	fmt.Println("Please enter the group the service should as (empty to not specify)")
	_, _ = cyan.Print("Group: ")
	Service.core.group = GetString(true, true)
	fmt.Println("Please enter the working directory for the service (empty to not specify)")
	_, _ = cyan.Print("Working Directory: ")
	Service.core.workingDirectory = GetString(true, true)
	fmt.Println("Please specify the environment variables (empty to not specify)")
	fmt.Println("Example: \"TOKEN=XXX ID=YYY\" without double quotes and with space between variables only")
	_, _ = cyan.Print("Environment: ")
	Service.core.environment = strings.Fields(GetString(true, false))
	fmt.Println("Please enter the start command")
	_, _ = cyan.Print("ExecStart: ")
	Service.core.start.start = GetString(false, false)
	fmt.Println("Please enter the pre-start command (empty for none)")
	_, _ = cyan.Print("ExecStartPre: ")
	Service.core.start.pre = GetString(true, false)
	fmt.Println("Please enter the post-start command (empty for none)")
	_, _ = cyan.Print("ExecStartPost: ")
	Service.core.start.post = GetString(true, false)
	switch input {
	case 3:
		fmt.Println("Do you want the service to remain after the start command has finished executing?")
		_, _ = cyan.Print("Remain After Exit (yes/no): ")
		answer := GetString(false, true)
		if strings.ToLower(answer) == "yes" {
			Service.core.remainAfterExit = true
		}
	default:
		fmt.Println("Please enter the reload command (empty for none)")
		_, _ = cyan.Print("ExecReload: ")
		Service.core.reload = GetString(true, false)
		fmt.Println("Please enter the stop command (empty for none)")
		_, _ = cyan.Print("ExecStop: ")
		Service.core.stop.stop = GetString(true, false)
		fmt.Println("Please enter the post-stop command (empty for none)")
		_, _ = cyan.Print("ExecStopPost: ")
		Service.core.stop.post = GetString(true, false)
	}
	//Service section ends
	//Install section begins
	fmt.Println("Please enter any aliases for your service (empty for none)")
	_, _ = cyan.Print("Alias: ")
	Service.install.alias = GetString(true, false)
	fmt.Println("Do you want the service to be able to start at boot?")
	_, _ = cyan.Print("Auto Start (yes/no): ")
	answer := GetString(false, true)
	if strings.ToLower(answer) == "yes" {
		Service.install.autoStart = true
	}
	//Install section ends
	//We all love to write :D
	FileContents := PrepareServiceFileContents(Service)
	fmt.Print(FileContents)
	_, _ = cyan.Print("Is this correct? (yes/no): ")
	answer = GetString(false, true)
	if strings.ToLower(answer) == "yes" {
		FileName := CreateConfigFile(Service.name)
		WriteContentToFile(FileName, FileContents)
	}
}

func PrepareServiceFileContents(service service) string {
	//Unit Section Begins
	output := "[Unit]"
	output += "\nDescription=" + service.unit.description
	if service.unit.after != "" {
		output += "\nAfter=" + service.unit.after
	}
	if service.unit.requires != "" {
		output += "\nRequires=" + service.unit.requires
	}
	if service.unit.wants != "" {
		output += "\nWants=" + service.unit.wants
	}
	//Unit Section Ends
	//Service Section Begins
	output += "\n\n[Service]"
	output += "\nType=" + service.core.kind
	//Optional Service Start Executables
	if service.core.user != "" {
		output += "\nUser=" + service.core.user
	}
	if service.core.group != "" {
		output += "\nGroup=" + service.core.group
	}
	if service.core.workingDirectory != "" {
		output += "\nWorkingDirectory=" + service.core.workingDirectory
	}
	for _, env := range service.core.environment {
		output += "\nEnvironment=\"" + env + "\""
	}
	if service.core.start.pre != "" {
		output += "\nExecStartPre=" + service.core.start.pre
	}
	output += "\nExecStart=" + service.core.start.start
	if service.core.start.post != "" {
		output += "\nExecStartPost=" + service.core.start.post
	}
	if service.core.remainAfterExit {
		output += "\nRemainAfterExit=yes"
	}
	//Optional Service Reload Executables
	if service.core.reload != "" {
		output += "\nExecReload=" + service.core.reload
	}
	//Optional Service Stop Executables
	if service.core.stop.stop != "" {
		output += "\nExecStop=" + service.core.stop.stop
	}
	if service.core.stop.post != "" {
		output += "\nExecStopPost=" + service.core.stop.post
	}
	//Service Section Ends
	//Install Section Begins
	output += "\n\n[Install]"
	if service.install.alias != "" {
		output += "\nAlias=" + service.install.alias
	}
	if service.install.autoStart {
		output += "\nWantedBy=multi-user.target"
	}
	//Install Section Ends WantedBy=multi-user.target
	output += "\n" //Add Newline at EOF
	return output
}

func TakeInput() uint {
	var input int
	_, _ = magenta.Print("Options: ")
	fmt.Println("\n(1) Simple - A single command is executed and is considered to be the main process/service\n(2) Forking - A Service which expects the main process/service to fork and launch a child process\n(3) Oneshot - A Service which is expected to start and finish its job before other units are launched\n(4) Exit")
	_, _ = cyan.Print("Which type of service do you want to create: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	scannedText := scanner.Text()
	input, err := strconv.Atoi(scannedText)
	if err != nil {
		fmt.Println("Something went wrong, please try again.")
		return TakeInput()
	}
	if uint(input) > 4 || uint(input) == 0 {
		fmt.Println("Enter a valid number.")
		return TakeInput()
	}
	return uint(input)
}

func GetString(EmptyAllowed bool, SingleWord bool) string {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	if SingleWord {
		ans := strings.Fields(input.Text())
		if len(ans) == 0 {
			if EmptyAllowed {
				return ""
			}
			fmt.Println("This cannot be empty, please enter some text")
			return GetString(false, true)
		}
		if len(ans) > 1 {
			fmt.Println("More than one words entered, only first word will be used as name")
		}
		return ans[0]
	}
	if !EmptyAllowed && input.Text() == "" {
		fmt.Println("This cannot be empty, please enter some text")
		return GetString(false, false)
	}
	return input.Text()
}

func WriteContentToFile(FileName string, FileContents string) {
	TestWritePermissions() //Before Writing, test if it is possible
	err := ioutil.WriteFile(FileName, []byte(FileContents), 0644)
	if err != nil {
		fmt.Println("Something went wrong, please send the log below to Sid Sun at <sid@sidsun.com>.")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	color.Red("Please inspect the config for any issues before installing it, while this program tries to do everything possible to avoid issues, there might be issues.\n")
	fmt.Printf("Config written to %q, Enjoy!\n", FileName)
}

func CreateConfigFile(ServiceName string) string {
	ConfigFileName := ServiceName + ".service"
	ConfigFile, _ := os.Create(ConfigFileName)
	_ = ConfigFile.Close()
	return ConfigFileName
}

func TestWritePermissions() {
	newFile, err := os.Create("systemd-auto-config.test.txt")
	if err != nil {
		if os.IsPermission(err) {
			log.Println("Error: Write permission denied, please cd into a workable dir, exiting.")
			os.Exit(1)
		}
	} else {
		_ = os.Remove("systemd-auto-config.test.txt")
		_ = newFile.Close()
	}
}
