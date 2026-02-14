#  EchoStrike

> **High-Performance Syslog Attack Simulation & Traffic Generation**

##  Why EchoStrike?

**I built** a specialized tool to generate high-fidelity Syslog traffic because testing security tools shouldn't require waiting for real attacks.

**I use it** to verify that SIEMs (like Splunk, Elastic, Sentinel) and detection pipelines are working correctly. It ensures your alerts fire when they should, without the risk of running actual malware.

**It works** by combining an intelligent **Randomization Engine** with a low-level **Network Sender**. This allows you to generate millions of unique, realistic log lines (varying IPs, users, timestamps) and blast them over UDP, TCP or TLS to stress-test your infrastructure.

EchoStrike is a single-binary CLI tool designed for security professionals, Red Teamers, and Detection Engineers. It generates realistic syslog traffic, simulates attack patterns and stress-tests SIEM pipelines with high-volume log ingestion.

##  Features

### Core Capabilities (Implemented)
- **Multi-Protocol Support**: Send logs via UDP, TCP, or TLS.
- **RFC Compliance**: Full support for RFC 3164 (BSD) and RFC 5424 (IETF) message formats.
- **Attack Simulation**: Automated brute-force and port-scan log patterns (`simulate` command).
- **High-Volume Generation**: Rate-controlled traffic with `generate` command.
- **Replay Mode**: Replay existing log files with `replay` command.
- **Dry Run**: Preview logs formatting with `preview` command.
- **Zero-Dependency**: Static Go binary, runs anywhere.
- **Docker Support**: Includes Dockerfile for easy containerization.

### Planned Features (Coming Soon)
- **Jitter & Randomization**: More advanced variations in timestamps and user agents.

##  Installation

### From Source
Requires Go 1.21+:

```bash
git clone https://github.com/jomboi8/echostrike.git
cd echostrike
go install ./cmd/echostrike
```

### Run Directly
```bash
go run cmd/echostrike/main.go [command] [flags]
```

## 🛠 Usage

### Basic Sending
Send a single test message to a local syslog server over UDP:

```bash
echostrike send --host 127.0.0.1 --port 514 --message "User 'admin' failed login from 192.168.1.50"
```

### Advanced Protocol & Formatting
Send via TCP using the modern RFC 5424 format with a custom app tag:

```bash
echostrike send \
  --proto tcp \
  --format rfc5424 \
  --tag sshd \
  --message "Accepted publickey for root from 10.0.0.5 port 55412 ssh2"
```

### Secure Transmission (TLS)
Send over TLS (skips verify for self-signed certs by default):

```bash
echostrike send --proto tls --host syslog.corp.local --port 6514 --message "Secure audit event"
```

### Docker Usage
Build the container:
```bash
docker build -t echostrike .
```

Run a simulation via Docker:
```bash
docker run --rm echostrike simulate --type brute-force --host 192.168.1.50
```

## Architecture

EchoStrike is built with a modular architecture to support high throughput and extensibility:

- **`cmd/echostrike`**: The CLI entry point, built with `Cobra` for robust flag handling.
- **`internal/sender`**: Handles the network transport layer. It abstracts TCP/UDP/TLS connections and manages buffering for high-performance sending.
- **`pkg/syslog`**: A standalone RFC-compliant formatter. It ensures messages are strictly formatted according to syslog standards (RFC 3164/5424) before transmission.
- **`internal/generator`** *(Planned)*: The engine responsible for hydrating templates with randomized data (IPs, Usernames, Timestamps) to create realistic "fuzz" data.

## Contributing

Pull requests are welcome! I am currently looking for contributions in:

- **Template Packs**: Real-world log samples for various services (AWS, Cisco, Linux Auth).
- **Attack Scenarios**: Logic to generate multi-stage log sequences (e.g., failed login -> successful login -> sudo usage).

## License

MIT
