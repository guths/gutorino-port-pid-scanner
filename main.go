package main

import "gihub.com/guths/gutorino-port-scanner/port"

func main() {
	// fmt.Println("Port Scanning")
	// widescanresults := port.WideScan("localhost")
	// fmt.Println(widescanresults)

	port.GetPidByPort("tcp", "80")
}
