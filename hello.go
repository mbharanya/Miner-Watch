package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/mitchellh/go-ps"
)

func main() {
	processList, err := ps.Processes()
	if err != nil {
		log.Println("ps.Processes() Failed")
		return
	}

	lines, err := readLines("stoplist.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	// map ages
	for x := range processList {
		var process ps.Process
		process = processList[x]
		// log.Printf("%d\t%s\n", process.Pid(), process.Executable())
		processName := process.Executable()
		if arrayContainsString(strings.TrimSpace(processName), lines) {
			log.Println(processName)
			killMiningProcess()
		}else{
			// log.Println("Did not find "+processName + " to kill")
		}
	}
}

func killMiningProcess() {
	mining_proc_list, err := readLines("mining-proc.txt")

	processList, err := ps.Processes()
	if err != nil {
		log.Println("ps.Processes() Failed")
		return
	}
	for x := range processList {
		var process ps.Process
		process = processList[x]
		// log.Printf("%d\t%s\n", process.Pid(), process.Executable())
		processName := process.Executable()
		if arrayContainsString(processName, mining_proc_list) {
			log.Println("Killing " + processName)

			if runtime.GOOS == "windows" {
				kill := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(process.Pid()))
				err := kill.Run()
				if err != nil {
					log.Fatalf("Error killing process " + err.Error())
				}

			} else {
				kill := exec.Command("kill", "-15", strconv.Itoa(process.Pid()))
				err := kill.Run()
				if err != nil {
					log.Fatalf("Error killing process " + err.Error())
				}

			}
			log.Println("process was killed")
		}
		// do os.* stuff on the pid
	}

}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func arrayContainsString(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
