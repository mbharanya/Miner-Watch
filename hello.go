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

var miningProcessPid int

func main() {
	miningProcessPid = startMiningProcess()
	if miningProcessPid < 0 {
		log.Println("Could not start mining process")
		return
	}
	processStopped := ""
	minerIsRunning := true

	for ok := true; ok; ok = true {
		processList, err := ps.Processes()
		if err != nil {
			log.Println("ps.Processes() Failed")
			return
		}

		lines, err := readLines("stoplist.txt")
		if err != nil {
			log.Fatalf("readLines: %s", err)
		}

		var runningProcessNames []string

		// map ages
		for x := range processList {
			var process ps.Process
			process = processList[x]
			// log.Printf("%d\t%s\n", process.Pid(), process.Executable())
			processName := strings.TrimSpace(process.Executable())
			runningProcessNames = append(runningProcessNames, processName)
		}

		if minerIsRunning {
			restartMinerIfCrashed()
			for _, stopListedProcess := range lines {
				if arrayContainsString(stopListedProcess, runningProcessNames) {
					log.Println(stopListedProcess + " was found, killing mining process")
					killMiningProcess(miningProcessPid)
					log.Println("Killed it, watching for changes...")
					processStopped = stopListedProcess
					minerIsRunning = false
					break
				}
			}
		} else {
			if !arrayContainsString(processStopped, runningProcessNames) {
				log.Println(processStopped + " is not running anymore, start mining")
				miningProcessPid = startMiningProcess()
				processStopped = ""
				minerIsRunning = true
			}
		}

		time.Sleep(1 * time.Second)
	}

}

func restartMinerIfCrashed() {
	process, err := ps.FindProcess(miningProcessPid)
	if err != nil || process == nil {
		log.Println("Miner seems to have crashed, restarting it")
		miningProcessPid = startMiningProcess()
	}
}

func startMiningProcess() int {
	miningProcesses, err := readLines("mining-proc2.txt")
	if err != nil {
		log.Println("Can't start mining process from mining-proc.txt: " + err.Error())
		return -1
	}

	miningProcess := miningProcesses[0]
	cmd := exec.Command("cmd.exe", "/C", miningProcess)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
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
		log.Println("Could not find running PID " + strconv.Itoa(miningProcessPid) + err.Error())
		return
	} else if process == nil {
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
