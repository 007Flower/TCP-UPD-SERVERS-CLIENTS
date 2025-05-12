package main

import (
	"fmt"
	"math"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	numTCPClients = 50
	numUDPClients = 50
	tcpServerAddr = "localhost:3000"
	udpServerAddr = "localhost:3001"
	testDuration  = 20 * time.Second
)

var sampleMessages = []string{
	"Hello", "/time", "/date", "/echo test", "bye",
	"/help", "/nocknock", "How are you?", "/clients",
}

type Result struct {
	Protocol     string
	ResponseTime time.Duration
	Error        bool
}

var (
	results []Result
	mutex   sync.Mutex
)

func randomMessage() string {
	return sampleMessages[rand.Intn(len(sampleMessages))]
}

func recordResult(r Result) {
	mutex.Lock()
	defer mutex.Unlock()
	results = append(results, r)
}

func tcpClient(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	conn, err := net.Dial("tcp", tcpServerAddr)
	if err != nil {
		recordResult(Result{"TCP", 0, true})
		return
	}
	defer conn.Close()

	deadline := time.Now().Add(testDuration)
	for time.Now().Before(deadline) {
		msg := randomMessage()
		start := time.Now()
		_, err := fmt.Fprintf(conn, msg+"\n")
		if err != nil {
			recordResult(Result{"TCP", 0, true})
			break
		}
		buf := make([]byte, 2048)
		_, err = conn.Read(buf)
		latency := time.Since(start)
		if err != nil {
			recordResult(Result{"TCP", latency, true})
			break
		}
		recordResult(Result{"TCP", latency, false})
		time.Sleep(100 * time.Millisecond)
	}
}

func udpClient(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	raddr, err := net.ResolveUDPAddr("udp", udpServerAddr)
	if err != nil {
		recordResult(Result{"UDP", 0, true})
		return
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		recordResult(Result{"UDP", 0, true})
		return
	}
	defer conn.Close()

	buf := make([]byte, 2048)
	deadline := time.Now().Add(testDuration)
	for time.Now().Before(deadline) {
		msg := randomMessage()
		start := time.Now()
		_, err := conn.Write([]byte(msg))
		if err != nil {
			recordResult(Result{"UDP", 0, true})
			break
		}
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		_, _, err = conn.ReadFromUDP(buf)
		latency := time.Since(start)
		if err != nil {
			recordResult(Result{"UDP", latency, true})
			continue
		}
		recordResult(Result{"UDP", latency, false})
		time.Sleep(100 * time.Millisecond)
	}
}

func printSummary() {
	var tcpCount, udpCount int
	var tcpErrors, udpErrors int
	var tcpLatency, udpLatency time.Duration
	var tcpLatencies, udpLatencies []time.Duration

	minTCP, maxTCP := time.Hour, 0*time.Millisecond
	minUDP, maxUDP := time.Hour, 0*time.Millisecond

	for _, r := range results {
		if r.Protocol == "TCP" {
			tcpCount++
			if r.Error {
				tcpErrors++
				continue
			}
			tcpLatency += r.ResponseTime
			tcpLatencies = append(tcpLatencies, r.ResponseTime)
			if r.ResponseTime < minTCP {
				minTCP = r.ResponseTime
			}
			if r.ResponseTime > maxTCP {
				maxTCP = r.ResponseTime
			}
		} else if r.Protocol == "UDP" {
			udpCount++
			if r.Error {
				udpErrors++
				continue
			}
			udpLatency += r.ResponseTime
			udpLatencies = append(udpLatencies, r.ResponseTime)
			if r.ResponseTime < minUDP {
				minUDP = r.ResponseTime
			}
			if r.ResponseTime > maxUDP {
				maxUDP = r.ResponseTime
			}
		}
	}

	avgTCP := tcpLatency / time.Duration(math.Max(float64(len(tcpLatencies)), 1))
	avgUDP := udpLatency / time.Duration(math.Max(float64(len(udpLatencies)), 1))

	fmt.Println("\n========= STRESS TEST SUMMARY =========")
	fmt.Printf("TCP:\n")
	fmt.Printf("  Messages:     %d\n", tcpCount)
	fmt.Printf("  Errors:       %d\n", tcpErrors)
	fmt.Printf("  Avg Latency:  %v\n", avgTCP)
	fmt.Printf("  Min Latency:  %v\n", minTCP)
	fmt.Printf("  Max Latency:  %v\n", maxTCP)

	fmt.Printf("\nUDP:\n")
	fmt.Printf("  Messages:     %d\n", udpCount)
	fmt.Printf("  Errors:       %d\n", udpErrors)
	fmt.Printf("  Avg Latency:  %v\n", avgUDP)
	fmt.Printf("  Min Latency:  %v\n", minUDP)
	fmt.Printf("  Max Latency:  %v\n", maxUDP)

	fmt.Println("\nComparison:")
	if avgTCP < avgUDP {
		fmt.Println("  ✅ TCP has lower average latency.")
	} else if avgUDP < avgTCP {
		fmt.Println("  ✅ UDP has lower average latency.")
	} else {
		fmt.Println("  ⚖️  Both have similar latency.")
	}

	if tcpErrors < udpErrors {
		fmt.Println("  ✅ TCP experienced fewer errors.")
	} else if udpErrors < tcpErrors {
		fmt.Println("  ✅ UDP experienced fewer errors.")
	} else {
		fmt.Println("  ⚖️  Both have similar error rates.")
	}
	fmt.Println("========================================")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup

	for i := 0; i < numTCPClients; i++ {
		wg.Add(1)
		go tcpClient(&wg, i)
	}
	for i := 0; i < numUDPClients; i++ {
		wg.Add(1)
		go udpClient(&wg, i)
	}

	wg.Wait()
	printSummary()
}
