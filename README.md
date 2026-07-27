# BF4 Blaze Emulator

![Status](https://img.shields.io/badge/status-WIP-orange)
![Language](https://img.shields.io/badge/language-Go-blue)
![Platform](https://img.shields.io/badge/platform-PS3-lightgrey)

# ⚠️ Work In Progress

This project is a **work in progress** implementation of a Battlefield 4 Blaze backend replacement.

The goal is to recreate the required online services for legacy clients by implementing the server-side protocols and services while keeping the original game client untouched.

---

# Overview

BF4 originally relied on EA's backend infrastructure for:

- Authentication
- Server discovery
- Redirect services
- Blaze RPC communication
- Multiplayer sessions
- Presence
- Matchmaking
- Game reporting

This project provides a custom backend implementation that recreates these services locally.

The client is redirected using **DNS control**, allowing the original game executable to communicate with the replacement server without modifying the game binary.

---

# Features

## Implemented

✅ Custom Blaze TCP server  
✅ Blaze packet parsing  
✅ TDF structure parsing  
✅ Component based request handling  
✅ UDP service listeners  
✅ HTTP service endpoint  
✅ HTTPS/TLS service support  
✅ Custom certificate handling  
✅ DNS based service redirection  
✅ Multi-service EA endpoint emulation  

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
        |
        v
 gosredirector.ea.com
 bf4.gos.ea.com
        |
        |
        v
BF4 Blaze Emulator
        |
        |
 +------+-------+
 |              |
 v              v
Blaze TCP     UDP Services
42130         Game Services
42131

```

---

# DNS Redirection

The client normally connects to EA infrastructure:

```
gosredirector.ea.com
bf4.gos.ea.com
```

Instead of changing the game executable, DNS responses redirect these domains to the emulator server.

Example:

```
gosredirector.ea.com  ->  151.xxx.xxx.xxx
bf4.gos.ea.com        ->  151.xxx.xxx.xxx
```

The client believes it is communicating with the original backend while requests are handled by this server.

---

# No Client Binary Modification

One of the main goals of this project is:

- No modified executable
- No patched game files
- No altered network code

The original client is used as-is.

All compatibility work happens on the server side through:

- Protocol recreation
- Packet handling
- DNS routing
- TLS compatibility

---

# Blaze Protocol

The server implements parts of EA's Blaze networking protocol.

Current flow:

```
Client
 |
 | TCP connection
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

Example packet processing:

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

Handles:

- Blaze RPC requests
- Authentication flow
- Component communication
- Session handling


---

## UDP Services

Ports:

```
25100
25101
25102
25103
```

Used for:

- Game traffic
- Session communication
- Real-time services

---

# TLS / ProtoSSL Compatibility

Older EA clients use a legacy TLS implementation.

The emulator provides TLS endpoints compatible with older clients.

Certificate handling includes:

```
Fake EA Root CA

        |
        |

gosredirector.ea.com
bf4.gos.ea.com

```

A single certificate can contain multiple EA service domains using SAN entries:

```
Subject:
CN=gosredirector.ea.com

SAN:
DNS:gosredirector.ea.com
DNS:bf4.gos.ea.com
```

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
│   
│   
│
├── servers/
│   ├── udp.go
│   
│
├── crt/
│   ├── fullchain.pem
│   └── privkey.pem
│
└── main.go

```

---

# Requirements

- Go 1.20+
- DNS server capable of custom records
- Server/VPS running TCP and UDP ports
- Battlefield 4 client

---

## Not implemented yet

- Full matchmaking replacement
- Dedicated server integration
- Complete statistics backend
- All Blaze components

---

# Purpose

This project is for:

- Game preservation research
- Network protocol research
- Learning about legacy online architectures
- Understanding Blaze server communication

---

# Disclaimer

This project is an independent implementation and is not affiliated with Electronic Arts.

Battlefield and related trademarks belong to their respective owners.

---

## Project Status

This project is currently a **Work In Progress (WIP)** and is not yet fully functional.

The current version is provided **as-is** and should be considered experimental. At this stage, the emulator/server will not work correctly out of the box, as additional research and development is required.

Some components are incomplete, unstable, or missing. Current development includes investigating network protocols, server behavior, and generating the required TLS certificates and configurations needed for compatibility with the original clients without modifying their binaries.

Features may change, break, or be redesigned as development continues.

This project is intended for research, preservation, and educational purposes.

---
