package main

import (
	"fmt"
	"net"
	"net/http"
    "io"
	"crypto/tls"
	"bf4/blaze"
    "bf4/components"
	"bf4/servers"
)

func startHTTP() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[HTTP]")
		fmt.Println("Host:", r.Host)
		fmt.Println("Path:", r.URL.Path)
		fmt.Println("Method:", r.Method)
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	fmt.Println("[+] HTTP listening :80")
	err := http.ListenAndServe(":80", mux,)
	if err != nil {
		panic(err)
	}
}

func startHTTPS() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("==============================")
		fmt.Println("[HTTPS REQUEST]")
		fmt.Println("Host:", r.Host)
		fmt.Println("Path:", r.URL.Path)
		fmt.Println("Method:", r.Method)
		fmt.Println("User-Agent:", r.UserAgent())

		if r.Body != nil {
			defer r.Body.Close()
			body, err := io.ReadAll(r.Body)
			if err == nil && len(body) > 0 {
				fmt.Printf("Body (%d bytes): % X\n", len(body), body,)
			}
		}

		w.Header().Set("Content-Type", "application/octet-stream",)
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	cert, err := tls.LoadX509KeyPair("crt/fullchain.pem", "crt/privkey.pem")
	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS10,
		MaxVersion: tls.VersionTLS12,

		Certificates: []tls.Certificate{
			cert,
		},

		NextProtos: []string{
			"http/1.1",
		},

		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	server := &http.Server{
		Addr: ":443",
		Handler: mux,
		TLSConfig: tlsConfig,
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
	}

	fmt.Println("[+] HTTPS listening :443")

	listener, err := tls.Listen("tcp", ":443", tlsConfig,)
	if err != nil {
		panic(err)
	}

	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}

func handleBlaze(conn net.Conn) {
	defer conn.Close()
	fmt.Println("[BLAZE] Client:", conn.RemoteAddr())
	buf := make([]byte,8192)

	for {
		n,err := conn.Read(buf)
		if err != nil {
			fmt.Println("Disconnected")
			return
		}

		data := buf[:n]
		// Parse Blaze packet
		packet := blaze.Parse(data)
		// Print header info
		packet.Dump()
		// Parse TDF data
		tdfs := blaze.ReadTDF(packet.Payload)

		for _, t := range tdfs {
			fmt.Println("TDF:", t.Tag, t.Value,)
		}

		reply := components.HandlePacket(data)
		if reply != nil {
			conn.Write(reply)
		}
	}
}

func startTCP(port string) {
	listener,err := net.Listen("tcp", ":"+port,)
	if err != nil {
		panic(err)
	}

	fmt.Println("[+] TCP listening:", port,)

	for {
		conn,err := listener.Accept()
		if err != nil {
			continue
		}

		go handleBlaze(conn)
	}
}

func main(){
	fmt.Println("BF4 Blaze Emulator")
	go startTCP("42130")
	go startTCP("42131")
	go servers.StartUDP(25100)
	go servers.StartUDP(25101)
	go servers.StartUDP(25102)
	go servers.StartUDP(25103)
	go startHTTP()
	go startHTTPS()
	select {}
}
