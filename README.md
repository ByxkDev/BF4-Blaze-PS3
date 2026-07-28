# BF4 Blaze Emulator

![Status](https://img.shields.io/badge/status-WIP-orange)
![Language](https://img.shields.io/badge/language-Go-blue)
![Platform](https://img.shields.io/badge/platform-PS3-lightgrey)

---

# ⚠️ Work In Progress

This project is an experimental implementation of a replacement backend for **Battlefield 4 (PlayStation 3)**.

The primary goal is to recreate EA's original online infrastructure so that the original game client can communicate with a custom server **without modifying the game executable**.

Development is currently focused on one major milestone:

> **Successfully completing the TLS handshake.**

Once the PS3 client establishes a secure HTTPS connection, development can continue with implementing the Blaze backend, authentication, matchmaking, and the remaining online services.

At the moment, the TLS handshake is the primary blocker.

---

# Current TLS Status

After generating the certificates using the included `generate.bat` script, the server currently receives the following output:

```text
[TLS] ClientHello
SNI:
Ciphers:
0x5 TLS_RSA_WITH_RC4_128_SHA
0x4 0x0004
0x2f TLS_RSA_WITH_AES_128_CBC_SHA
0x35 TLS_RSA_WITH_AES_256_CBC_SHA
2026/07/28 19:59:26 http: TLS handshake error EOF
```

This confirms that:

- The PS3 successfully connects to the server.
- The TLS ClientHello is received.
- The client advertises legacy RSA cipher suites.
- The client immediately aborts the handshake before completion.

At the moment, the exact reason for the disconnect has not yet been identified.

Research currently includes:

- Legacy EA ProtoSSL behavior
- Certificate chain compatibility
- Signature algorithm compatibility
- Legacy RSA cipher suites
- TLS version compatibility
- Differences between OpenSSL and EA's ProtoSSL implementation

The repository also contains tooling to generate certificates that resemble EA's certificate chain.

Even after applying the known **Bug_OldProtoSSL** patch (replacing the second occurrence of `2a864886f70d010104`), Battlefield 4 PS3 v1.20 still terminates the handshake with the same EOF error.

Until the TLS handshake completes successfully, no HTTPS requests are exchanged, preventing further implementation of the Blaze backend.

---

# Overview

Battlefield 4 originally relied on numerous EA backend services, including:

- Authentication
- Redirect services
- Blaze RPC communication
- Server discovery
- Matchmaking
- Presence
- Multiplayer sessions
- Game reporting

This project aims to recreate those services locally while preserving compatibility with the original PlayStation 3 client.

The client is redirected through **DNS control**, allowing all traffic to be handled by the emulator instead of EA's original infrastructure.

---

# Features

## Currently Implemented

- ✅ Custom Blaze TCP server
- ✅ Blaze packet parsing
- ✅ TDF structure parsing
- ✅ Component-based request handling
- ✅ UDP service listeners
- ✅ HTTP service endpoint
- ✅ HTTPS listener
- ✅ Custom certificate generation
- ✅ DNS-based service redirection
- ✅ Multi-service EA endpoint emulation
- ✅ TLS ClientHello inspection

## Currently In Progress

- ⏳ Legacy ProtoSSL compatibility
- ⏳ TLS handshake completion
- ⏳ HTTPS request analysis

---

# Architecture

```
                Battlefield 4 Client
                       |
                       |
                       v
              Custom DNS Resolver
                       |
                       |
        +--------------+--------------+
        |                             |
        v                             v
 gosredirector.ea.com          bf4.gos.ea.com
        |                             |
        +--------------+--------------+
                       |
                       v
              BF4 Blaze Emulator
                       |
          +------------+------------+
          |                         |
          v                         v
     Blaze TCP                 UDP Services
     42130                     25100
     42131                     25101
                               25102
                               25103
```

---

# DNS Redirection

The original game normally connects to EA servers:

```
gosredirector.ea.com
bf4.gos.ea.com
```

DNS responses redirect those domains to the emulator server instead.

Example:

```
gosredirector.ea.com  ->  151.xxx.xxx.xxx
bf4.gos.ea.com        ->  151.xxx.xxx.xxx
```

The client believes it is communicating with EA while all requests are handled locally.

---

# No Client Binary Modification

One of the primary goals of this project is preserving compatibility with the original game.

The emulator does **not** require:

- Modified EBOOT files
- Patched executables
- Altered game assets
- Custom network libraries

Instead, compatibility is achieved through:

- DNS redirection
- Protocol recreation
- Packet parsing
- Blaze implementation
- TLS compatibility research

---

# Blaze Protocol

The emulator implements portions of EA's Blaze networking protocol.

Current packet flow:

```
Client
   |
TCP Connection
   |
   v
Blaze Parser
   |
   v
Packet Decoder
   |
   v
Component Handler
   |
   v
Response Builder
   |
   v
Client
```

Example:

```go
packet := blaze.Parse(data)

tdfs := blaze.ReadTDF(packet.Payload)

reply := components.HandlePacket(data)
```

---

# Services

## Blaze TCP

Ports:

```
42130
42131
```

Responsibilities:

- Blaze RPC
- Authentication
- Component communication
- Session management

---

## UDP Services

Ports:

```
25100
25101
25102
25103
```

Responsibilities:

- Session traffic
- Game traffic
- Real-time communication

---

# TLS / ProtoSSL Research

Older EA titles use EA's proprietary **ProtoSSL** implementation.

Before any Blaze communication begins, Battlefield 4 establishes an HTTPS connection.

The emulator includes tooling to generate a custom certificate chain attempting to match the behavior expected by the client.

Certificate layout:

```
Fake EA Root CA

        |
        |

gosredirector.ea.com
bf4.gos.ea.com
```

The generated server certificate supports multiple EA domains using SAN entries:

```
Subject:
CN=gosredirector.ea.com

Subject Alternative Names:

DNS:gosredirector.ea.com
DNS:bf4.gos.ea.com
```

Current TLS compatibility remains under active investigation.

---

# Project Structure

```
BlazeEmu/

├── blaze/
│   ├── parser.go
│   ├── packet.go
│   └── tdf.go
│
├── components/
│   ├── authentication.go
│   └── ...
│
├── servers/
│   ├── udp.go
│   └── ...
│
├── crt/
│   ├── fullchain.pem
│   └── privkey.pem
│
├── generate.bat
│
└── main.go
```

---

# Requirements

- Go 1.20+
- Battlefield 4 PS3 v1.20
- DNS server capable of custom records
- Server or VPS with TCP/UDP support

---

# Current Roadmap

1. Complete TLS handshake.
2. Capture the first HTTPS requests.
3. Implement GOS Redirector.
4. Implement Blaze authentication.
5. Implement session management.
6. Implement matchmaking.
7. Implement presence.
8. Implement game reporting.
9. Continue protocol reverse engineering.

---

# Not Yet Implemented

- Full matchmaking
- Presence system
- Dedicated server backend
- Statistics backend
- Leaderboards
- Friends
- Clubs
- Complete Blaze components
- Full redirector implementation

---

# Purpose

This project exists for:

- Game preservation
- Reverse engineering research
- Network protocol research
- Learning about legacy online architectures
- Educational purposes

Its goal is to better understand EA's original Battlefield 4 backend while developing a compatible replacement implementation.

---

# Client Analysis

Development relies heavily on information extracted from the original Battlefield 4 PlayStation 3 client.

A `strings.txt` file is generated by dumping readable strings from the original `EBOOT.ELF`.

These strings help identify:

- Domains
- URLs
- Blaze components
- Service names
- Redirect endpoints
- Internal protocol identifiers

Combined with packet captures, TLS debugging, and server logs, this information helps reconstruct the behavior expected by the original client.

---

# Project Status

| Component | Status |
|------------|--------|
| DNS Redirection | ✅ Working |
| HTTP Server | ✅ Working |
| HTTPS Listener | ✅ Working |
| Blaze TCP Listener | ✅ Working |
| UDP Services | ✅ Working |
| Blaze Packet Parser | ✅ Working |
| TDF Parser | ✅ Working |
| TLS ClientHello Reception | ✅ Working |
| TLS Handshake Completion | ❌ Not Yet |
| HTTPS Request Processing | ❌ Not Yet |
| Blaze Authentication | ⏳ Planned |
| Matchmaking | ⏳ Planned |
| Presence | ⏳ Planned |
| Game Reporting | ⏳ Planned |

---

# Disclaimer

This project is an independent implementation created for research and preservation purposes.

It is **not affiliated with, endorsed by, or associated with Electronic Arts or DICE**.

Battlefield and all related trademarks are the property of their respective owners.

---

# Contributing

Contributions are always welcome.

If you are interested in helping with:

- Blaze protocol research
- Packet analysis
- TLS compatibility
- Reverse engineering
- Development
- Documentation
- Testing

feel free to reach out on Discord:

**@Byxk**

Any assistance is greatly appreciated.

---

# License

This project is intended solely for:

- Research
- Education
- Preservation

No original EA server software is included.

Only original source code written for this emulator is distributed.

---

**Current Development Focus**

> Successfully completing the TLS handshake so the original Battlefield 4 PS3 client can begin exchanging HTTPS requests with the custom backend.
