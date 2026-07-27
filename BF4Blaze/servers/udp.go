package servers

import (
	"fmt"
	"net"
)

func StartUDP(port int){
	addr :=
		net.UDPAddr{
			Port:port,
			IP:net.ParseIP("0.0.0.0"),
		}

	conn,err := net.ListenUDP("udp", &addr,)

	if err != nil {
		panic(err)
	}

	fmt.Println("[UDP] Listening:", port,)
	buf := make([]byte,8192)

	for {
		n,client,err := conn.ReadFromUDP(buf)

		if err != nil {
			continue
		}

		fmt.Printf("[UDP %d] %s %d bytes\n", port, client, n,)
		fmt.Printf("% X\n", buf[:n],)
	}
}