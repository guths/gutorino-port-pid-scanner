package port

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func GetPidByPort(protocol, port string) (string, error) {

	file, err := os.Open(fmt.Sprintf("/proc/net/%s", protocol))

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	defer file.Close()

	stats, err := file.Stat()

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	fmt.Printf("getting info of sockets in /proc/net/tcp: %s\n", stats.Name())

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "sl") {
			continue
		}

		fields := strings.Fields(line)

		if len(fields) >= 4 && fields[3] == "0A" {
			localAddrPort := strings.Split(fields[1], ":")
			localPort := localAddrPort[1]
			portNumber, err := strconv.ParseInt(localPort, 16, 64)

			if err != nil {
				continue
			}

			fmt.Printf("Port scanned: %v\n", portNumber)

			inode := fields[9]
			fmt.Printf("inode scanned: %v\n", inode)

			pid, err := findPIDByINode(inode)

			if err != nil {
				continue
			}

			fmt.Println("PID associated with local address 0.0.0.0:", pid)
		}
	}

	return "", nil
}

func findPIDByINode(inode string) (string, error) {
	procFiles, err := ioutil.ReadDir("/proc")

	if err != nil {
		return "", err
	}

	for _, file := range procFiles {
		if !file.IsDir() {
			continue
		}

		pid := file.Name()
		fdDir := fmt.Sprintf("/proc/%s/fd", pid)

		fdFiles, err := ioutil.ReadDir(fdDir)

		if err != nil {
			continue
		}

		for _, fdFile := range fdFiles {
			targetPath, err := os.Readlink(fmt.Sprintf("/proc/%s/fd/%s", pid, fdFile.Name()))

			if err != nil {
				fmt.Println("nao achou")
			}

			fmt.Printf("TARGET PATH: %s\n", targetPath)

			if strings.HasPrefix(fdFile.Name(), "socket:") && strings.HasSuffix(fdFile.Name(), inode) {
				return pid, nil
			}
		}
	}

	return "", nil
}
