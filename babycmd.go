package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/go-ini/ini"
)

var babyrc string = "/Users/Kios/.babyrc"
var babysh string = "/Users/Kios/baby.sh"

func init() {}

func env_check() {
	// check sshpass command is valid
	// command: brew install https://raw.githubusercontent.com/kadwanev/bigboybrew/master/Library/Formula/sshpass.rb
	chk_cmd := "sshpass"
	install_sshpass := "brew install https://raw.githubusercontent.com/kadwanev/bigboybrew/master/Library/Formula/sshpass.rb"
	cmd := exec.Command(chk_cmd)
	log.Printf("Running chk command and waiting for it to finish...")
	err := cmd.Run()
	if (err != nil) {
		fmt.Println("PREPARE TO INSTALL SSHPASS VIA HOMEBREW")
		install_cmd := exec.Command(install_sshpass)
		err := install_cmd.Run()

		if err != nil {
			fmt.Println("INSTALL COMPLETE!")
		}
	}
	log.Printf("Command finished with error: %v", err)
}

func banner() {
	fmt.Printf("%c[1;40;32m##################################################%c[0m\n", 0x1B, 0x1B)
	fmt.Printf("%c[1;40;32m#\tBABY SERVER MANAGER                      #%c[0m\n", 0x1B, 0x1B)
	fmt.Printf("%c[1;40;32m#\tChoose server to connect!                #%c[0m\n", 0x1B, 0x1B)
	fmt.Printf("%c[1;40;32m#\t2019 Jeremy Ting. All Rights Reserved    #%c[0m\n", 0x1B, 0x1B)
	fmt.Printf("%c[1;40;32m##################################################%c[0m\n", 0x1B, 0x1B)
}

func clearTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func baby_list() {
	cfg, err := ini.Load(babyrc)
	if err != nil {
		fmt.Printf("Fail to read baby config file: %v", err)
		os.Exit(1)
	}
	names := cfg.SectionStrings()

	banner()
	for index, name := range(names) {
		if name != "DEFAULT" {
			ip := cfg.Section(name).Key("ip")
			fmt.Printf("%c[1;40;32m[%d]: %s%c[0m\t%c[1;40;31m%s%c[0m\n",0x1B, index, name, 0x1B, 0x1B, ip,0x1B)
		}
	}
}

func connect_via_ssh(index int64) {
	cfg, err := ini.Load(babyrc)
	if err != nil {
		fmt.Printf("Fail to read baby config file: %v", err)
		os.Exit(1)
	}
	names := cfg.SectionStrings()
	serverName := names[index+1]
	ip := cfg.Section(serverName).Key("ip")
	port := cfg.Section(serverName).Key("port")
	username := cfg.Section(serverName).Key("username")
	password := cfg.Section(serverName).Key("password")

	cmd_str := fmt.Sprintf("sshpass -p %s ssh %s@%s -p %s", password, username, ip, port)
	// 需要将这个命令写入 sh 文件中 然后用 iTerm2 someshell.sh 打开才行 不能直接使用命令 会出错 （iTerm2 会报错）

	f, err := os.Create(babysh)
	check(err)

	f.WriteString(cmd_str)
	f.Sync()
	f.Close()

	err = os.Chmod(babysh, 0777)

	if err != nil {
		log.Println(err)
	}

	if err != nil {
		fmt.Printf("%c[1;40;31m[!] CAN NOT WRITE SSH COMMAND TO SHELL FILE %c[0m\n", 0x1B, 0x1B)
		return
	}

	// fmt.Printf("%c[1;40;32m# Command: %s%c[0m\t\n",0x1B, cmd_str, 0x1B)
	cmd := exec.Command("/Applications/iTerm.app/Contents/MacOS/iTerm2", babysh)
	err = cmd.Run()

	if err != nil {
		fmt.Printf("%c[1;40;31m[!] CAN NOT OPEN SHELL VIA iTerm2 %c[0m\n", 0x1B, 0x1B)
		return
	}

	remove_sh(babysh)
}


func check(e error) {
	if e != nil {
		panic(e)
	}
}

func remove_sh(filename string) {
	err := os.Remove(filename)
	if err != nil {
		log.Println(err)
	}
}

func file_test() {
	f, err := os.Create(babysh)
	check(err)

	var cmd string = "echo 999"
	f.WriteString(cmd)
	f.Sync()
	f.Close()

	err = os.Chmod(babysh, 0777)

	if err != nil {
		log.Println(err)
	}

	remove_sh(babysh)
}

func list_detail(index int64) {
	cfg, err := ini.Load(babyrc)
	if err != nil {
		fmt.Printf("Fail to read baby config file: %v", err)
		os.Exit(1)
	}
	names := cfg.SectionStrings()
	serverName := names[index+1]
	ip := cfg.Section(serverName).Key("ip")
	port := cfg.Section(serverName).Key("port")
	username := cfg.Section(serverName).Key("username")
	password := cfg.Section(serverName).Key("password")
	fmt.Printf("%c[1;40;32m##################################################%c[0m\n", 0x1B, 0x1B)
	fmt.Printf("%c[1;40;32m#\tSERVER [%s] DETAILS%c[0m\t\n",0x1B, serverName, 0x1B)
	fmt.Printf("%c[1;40;32m#\tSERVER IP:      \t %s %c[0m\t\n",0x1B, ip, 0x1B)
	fmt.Printf("%c[1;40;32m#\tSERVER Port:    \t %s %c[0m\t\n",0x1B, port, 0x1B)
	fmt.Printf("%c[1;40;32m#\tSERVER USERNAME:\t %s %c[0m\t\n",0x1B, username, 0x1B)
	fmt.Printf("%c[1;40;32m#\tSERVER PASSWORD:\t %s %c[0m\t\n",0x1B, password, 0x1B)
	fmt.Printf("%c[1;40;32m##################################################%c[0m\n", 0x1B, 0x1B)
}

func logic() {
	clearTerminal()
	if (len(os.Args) < 2) {
		banner()
		fmt.Printf("%c[1;40;31m[!] TOO LESS ARGUMENTS%c[0m\n", 0x1B, 0x1B)
		return
	} else if (len(os.Args) == 2){
		command := os.Args[1]
		if command == "ls" {
			baby_list()
		} else if(command == "env_check") {
			env_check()
		} else if(command == "test") {
			file_test()
		} else {
			banner()
			fmt.Printf("%c[1;40;31m[!] COMMAND NOT SUPPORTED%c[0m\n", 0x1B, 0x1B)
		}
	} else if(len(os.Args) > 2 && len(os.Args) <= 3) {
		cmd_main := os.Args[1]
		cmd_param, err := strconv.ParseInt(os.Args[2], 10, 64)

		if err != nil {
			fmt.Printf("%c[1;40;31m[!] YOUR CHOICE IS INVALID%c[0m\n", 0x1B, 0x1B)
		}

		if(cmd_main == "ls") {
			list_detail(cmd_param-1)
		} else if(cmd_main == "con") {
			connect_via_ssh(cmd_param-1)
		}

	}
}

func main() {
	logic()
}
