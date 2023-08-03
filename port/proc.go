package port

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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

		fdDir := "/proc/"

		err := filepath.Walk(fdDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				regexPattern := `^/proc/\d+/fd$`

				regex := regexp.MustCompile(regexPattern)

				if regex.MatchString(path) {
					p := path

					err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {

						if !info.IsDir() {
							targetInode, _ := os.Readlink(fmt.Sprintf("%s/%s", p, info.Name()))

							if strings.Contains(targetInode, inode) {
								fmt.Printf("INODE: %s\nPID: %s", inode, path)
							}
						}

						return nil
					})

					if err != nil {
						return err
					}
				}

			}

			return nil
		})

		return "", err

	}

	return "", nil
}
