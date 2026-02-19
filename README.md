# Zipline

Secure P2P file transfer via 6-digit codes.

## Install

```bash
go build -o zl
```

## Usage

**Send:**

```bash
zl send file.pdf
```

**Receive:**

```bash
zl get 123456
```

**Relay Server:**

```bash
zl relay
```

## How it works

1. Sender gets 6-digit code
2. Receiver enters code and confirms
3. File encrypted with AES-GCM
4. Transferred in 64KB chunks via relay

## Security

- AES-GCM encryption
- Key derived from 6-digit code via SHA-256
- ⚠️ Use longer codes for production

## License

MIT
