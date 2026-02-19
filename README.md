# Zipline

Secure P2P file transfer via 6-digit codes.

## Install

```bash
curl -sSL https://raw.githubusercontent.com/yourusername/zipline/main/install.sh | bash
```

Or build from source:

```bash
go build -o zipline
```

## Usage

**Send:**

```bash
zipline send file.pdf
```

**Receive:**

```bash
zipline get 123456
```

**Relay Server:**

```bash
zipline relay
```

## How it works

1. Sender gets 6-digit code
2. Receiver enters code and confirms
3. File encrypted with AES-GCM
4. Transferred in 64KB chunks via relay

## Security

- AES-GCM encryption
- Key derived from 6-digit code via SHA-256

*Like this tool? Give a ⭐️*
