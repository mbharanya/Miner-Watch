package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
)

func main() {
	miningProcessPid := startMiningProcess()
	if miningProcessPid < 0 {
		log.Println("Could not start mining process")
		return
	}

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
			log.Println(processName + " was found, killing mining process")
			killMiningProcess(miningProcessPid)
		} else {
			// log.Println("Did not find "+processName + " to kill")
		}
	}
}

func startMiningProcess() int {
	// cmd := exec.Command("U:\\Work\\Ethereum\\PhoenixMiner_5.5c_Windows_AMD_NVIDIA.Password-phoenix\\start.bat")
	miningProcesses, err := readLines("mining-proc.txt")
	if err != nil {
		log.Println("Can't start mining process from mining-proc.txt: " + err.Error())
		return -1
	}

	miningProcess := miningProcesses[0]
	cmd := exec.Command(miningProcess)
	if err := cmd.Start(); err != nil {
		fmt.Println(">>0", err)
		return -1
	}
	log.Println("mining process pid: ", cmd.Process.Pid)
	return cmd.Process.Pid
}

func killMiningProcess(miningProcessPid int) {
	process, err := ps.FindProcess(miningProcessPid)
	if err != nil {
		fmt.Println(">>0", err)
		return
	}

	log.Println("Killing " + strconv.Itoa(process.Pid()) + "\t" + process.Executable())

	time.Sleep(1 * time.Second)

	var killCommand []string

	if runtime.GOOS == "windows" {
		killCommand = []string{"taskkill", "/T", "/F", "/PID"}
	} else {
		killCommand = []string{"kill", "-15"}
	}
	killCommand = append(killCommand, strconv.Itoa(process.Pid()))
	log.Println(killCommand)

	kill := exec.Command(killCommand[0], killCommand[1:]...)

	runErr := kill.Run()

	if runErr != nil {
		log.Fatalf("Error killing process " + runErr.Error())
	}

	log.Println("process was killed")
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
