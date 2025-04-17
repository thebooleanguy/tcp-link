# tcp-link 🧵

A minimalist multiplayer TCP chat server — connect via `telnet` or `nc`, choose a username, join rooms, and chat in real-time over the terminal.

---

### 🌐 Overview

**tcp-link** is a CLI-based TCP chat server written in **Go**, designed for learning and experimenting with socket programming, concurrency, and real-time messaging.

- ⚙️ Go-powered: `net`, goroutines, channels
- 👥 Multi-client support
- 🧑 Usernames + 🔁 Chat Rooms
- 📝 Simple architecture, no external dependencies

---

### 🚀 Features

- ✅ Echo Server (Phase 1)  
- ✅ Multi-Client Support (Phase 2)  
- ✅ Usernames & Chat Rooms (Phase 3)  
- 🛠️ [Next] Chat history, TLS, Admin tools  

---

### 📦 Run It

Compile & run the server:

```bash
go run chat-server.go
# or
go build chat-server.go && ./chat-server
```

Then connect in another terminal:

```bash
telnet localhost 1234
# or
nc localhost 1234
```

Example session:

```
> Enter your username:
komez_jk

> /join general
> hello from Go chat 🌊
```

---

### 🧠 How It Works

- Each client is handled via a separate **goroutine**
- Server manages:
  - `map[string]net.Conn` → usernames
  - `map[string][]Client` → chat rooms
- Command parsing supports `/join`, `/users`, etc.

---

### 🛠️ Tech Stack

- **Language**: Go (Golang)
- **Core Libs**: `bufio`, `fmt`, `io`, `log`, `net`, `slices`, `strings`, `time`
- **Testing tools**: `telnet`, `netcat (nc)`

---

### 🎯 Why?

**tcp-link** was built to:

- Explore low-level networking with Go
- Practice goroutines, maps, and concurrent I/O
- Create something fun and retro-feeling in the terminal

---

### 📅 Roadmap

- [ ] Persistent message logging
- [ ] Load recent history for new users
- [ ] TLS encryption (`crypto/tls`)
- [ ] Admin commands (`/kick`, `/mute`)
- [ ] File transfer / binary support
- [ ] A minimalistic retro themed UI

---

### 🌀 License

MIT – free to learn from, remix, and use.

---

### ☕ Credits & Inspiration

- Classic tools like `nc`, `telnet`, and BBS chatrooms
