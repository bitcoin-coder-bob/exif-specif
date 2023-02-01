package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	cmd := exec.Command("open", "pepeaccountant.jpg")
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("original pid: %d\n", cmd.Process.Pid)

	// pgrep -f pepeaccount.jpg
	cmd = exec.Command("pidof", "eog")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
	pidToKill := out.String()[0 : len(out.String())-1]

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter tags: ")
	keywords, _ := reader.ReadString('\n')
	fmt.Println(keywords)

	keywordRoot := "-keywords+="
	cmd = exec.Command("exiftool", keywordRoot+keywords, "pepeaccountant.jpg")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	pid, err := strconv.Atoi(pidToKill)
	if err != nil {
		log.Fatalf("conversion error: %s", err.Error())
	}

	hoverflyProcess := os.Process{Pid: pid}
	err = hoverflyProcess.Kill()
	if err != nil {
		log.Fatal(err)
	}

}
