package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
    "io"
	"bytes"
	
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

	err := http.ListenAndServe(":80", mux)
	if err != nil {
		panic(err)
	}
}

func startHTTPS() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println()
		fmt.Println("Remote:")
		fmt.Println(r.RemoteAddr)
		fmt.Println("Method:")
		fmt.Println(r.Method)
		fmt.Println("Host:")
		fmt.Println(r.Host)
		fmt.Println("URL:")
		fmt.Println(r.URL.String())
		fmt.Println("Protocol:")
		fmt.Println(r.Proto)
		fmt.Println()
		fmt.Println("Headers:")

		for k, v := range r.Header {
			fmt.Println(k+":", v)
		}

		// Read body
		if r.Body != nil {
			body, err := io.ReadAll(r.Body)

			if err == nil && len(body) > 0 {
				fmt.Println()
				fmt.Println("Body:")

				//print readable body
				fmt.Println(string(body))
				// restore body for handlers later
				r.Body = io.NopCloser(bytes.NewReader(body),)
			}
		}

		fmt.Println()
		fmt.Println("TLS Info:")

		if r.TLS != nil {
			fmt.Println("Version:", tlsVersionName(r.TLS.Version),)
			fmt.Println("Cipher:", tls.CipherSuiteName(r.TLS.CipherSuite),)
			fmt.Println("SNI:", r.TLS.ServerName,)

		} else {
			fmt.Println("No TLS")
		}

		fmt.Println()
		w.WriteHeader(200)
		w.Write([]byte("OK"))

	})

	cert, err := tls.LoadX509KeyPair("crt/certs/gosredirectoreacom-fullchain.pem", "crt/certs/gosredirectoreacom-priv.pem",)

	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{
			cert,
		},

		MinVersion: tls.VersionTLS10,
		MaxVersion: tls.VersionTLS11,

		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},

		PreferServerCipherSuites: true,
		SessionTicketsDisabled: true,
	}

	tlsConfig.GetConfigForClient = func(hello *tls.ClientHelloInfo,) (*tls.Config, error) {

		fmt.Println()

		fmt.Println("Remote:")
		fmt.Println(hello.Conn.RemoteAddr())

		fmt.Println("SNI:")
		if hello.ServerName == "" {
			fmt.Println("<none>")
		} else {
			fmt.Println(hello.ServerName)
		}

		fmt.Println("TLS Versions:")
		for _, v := range hello.SupportedVersions {
			fmt.Printf("  0x%x\n", v)
		}

		fmt.Println("Cipher Suites:")
		for _, c := range hello.CipherSuites {
			fmt.Printf("  0x%x\n", c)
		}

		fmt.Println("Signature Algorithms:")
		for _, s := range hello.SignatureSchemes {
			fmt.Printf("  0x%x\n", s)
		}

		fmt.Println()
		return nil, nil
	}

	rawListener, err := net.Listen("tcp", ":443",)
	if err != nil {
		panic(err)
	}

	fmt.Println("[+] HTTPS listening :443")

	listener := &TLSListener{
		Listener: rawListener,
		Config: tlsConfig,
	}

	server := &http.Server{
		Handler: mux,
	}

	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}

type TLSListener struct {
	net.Listener
	Config *tls.Config
}


func (l *TLSListener) Accept() (net.Conn,error){
	for {
		conn,err := l.Listener.Accept()
		if err != nil {
			continue
		}

		tlsConn := tls.Server(conn, l.Config,)

		err = tlsConn.Handshake()
		if err != nil {
			fmt.Println()
			fmt.Println("Remote:",conn.RemoteAddr())
			fmt.Println("Error:",err)
			fmt.Println()
			conn.Close()
			continue
		}

		state := tlsConn.ConnectionState()
		fmt.Println("[+] TLS CONNECTED")
		fmt.Println("Remote:",conn.RemoteAddr())
		fmt.Println("Version:",tlsVersionName(state.Version))
		fmt.Println("Cipher:",tls.CipherSuiteName(state.CipherSuite))
		return tlsConn,nil
	}
}

func tlsVersionName(v uint16) string {
	switch v {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	default:
		return fmt.Sprintf("0x%x",v)
	}
}

func handleBlaze(conn net.Conn){
	defer conn.Close()
	fmt.Println("[BLAZE] Client:",conn.RemoteAddr())
	buf:=make([]byte,8192)

	for {
		n,err:=conn.Read(buf)
		if err != nil {
			fmt.Println("[BLAZE] Disconnected")
			return
		}

		data:=buf[:n]
		packet:=blaze.Parse(data)
		packet.Dump()
		tdfs:=blaze.ReadTDF(packet.Payload)

		for _,t:=range tdfs {
			fmt.Println("TDF:", t.Tag, t.Value,)
		}

		reply:=components.HandlePacket(data)
		if reply != nil {
			conn.Write(reply)
		}
	}
}

func startTCP(port string){
	listener,err:=net.Listen("tcp", ":"+port,)
	if err != nil {
		panic(err)
	}

	fmt.Println("[+] TCP listening:",port)

	for {
		conn,err:=listener.Accept()
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
	select{}
}
