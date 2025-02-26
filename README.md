# PASETO CLI

A simple command-line utility for working with PASETO v3.public tokens using the go-paseto library.

## Installation

```bash
go build -o paseto-cli
```

## Usage

### Generate Keys

The utility can generate key pairs directly:

```bash
./paseto-cli generate
```

This will output a base64-encoded private and public key pair that you can use with the other commands.

### Sign a message

```bash
./paseto-cli sign -message "Hello, PASETO!" -key "YOUR_BASE64_PRIVATE_KEY" [-expiration "5m"]
```

Parameters:
- `-message`: The message you want to include in the token (required)
- `-key`: Base64-encoded private key (required)
- `-expiration`: Optional duration string specifying token expiration (default: "5m")
    - Examples: "1h" (1 hour), "30m" (30 minutes), "24h" (24 hours)

This will output a signed PASETO v3.public token.

### Verify a token

```bash
./paseto-cli verify -token "v3.public.eyJleHAiOiIyMD..." -key "YOUR_BASE64_PUBLIC_KEY"
```

Parameters:
- `-token`: The PASETO token to verify (required)
- `-key`: Base64-encoded public key (required)

If the token is valid, the utility will display the message and all claims contained in the token.

## Notes

- Tokens include standard claims: issued at (iat), not before (nbf), and expiration (exp)
- By default, tokens expire after 5 minutes if not specified otherwise
- All times are expressed in UTC

## Learn more

- [PASETO spec](https://github.com/paseto-standard/paseto-spec/blob/master/README.md)
    - [Version 3: NIST Modern](https://github.com/paseto-standard/paseto-spec/blob/master/docs/01-Protocol-Versions/README.md#version-3-nist-modern)
- [go-paseto](https://github.com/aidantwoods/go-paseto)
