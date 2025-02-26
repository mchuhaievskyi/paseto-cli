package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
)

func main() {
	// Define command line flags
	signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
	signCmdMessage := signCmd.String("message", "", "Message to sign (required)")
	signCmdKey := signCmd.String("key", "", "Base64-encoded private key (required)")
	signCmdExpiration := signCmd.String("expiration", "5m", "Expiration time (optional)")

	verifyCmd := flag.NewFlagSet("verify", flag.ExitOnError)
	verifyCmdToken := verifyCmd.String("token", "", "PASETO token to verify (required)")
	verifyCmdKey := verifyCmd.String("key", "", "Base64-encoded public key (required)")

	// Check if we have enough arguments
	if len(os.Args) < 2 {
		fmt.Println("Expected 'sign', 'verify', or 'generate' subcommands")
		os.Exit(1)
	}

	// Parse the appropriate subcommand
	switch os.Args[1] {
	case "sign":
		signCmd.Parse(os.Args[2:])
		if *signCmdMessage == "" || *signCmdKey == "" {
			signCmd.PrintDefaults()
			os.Exit(1)
		}
		expirationDuration, err := time.ParseDuration(*signCmdExpiration)
		if err != nil {
			fmt.Printf("Error parsing expiration: %v\n", err)
			signCmd.PrintDefaults()
			os.Exit(1)
		}
		signToken(*signCmdMessage, *signCmdKey, expirationDuration)

	case "verify":
		verifyCmd.Parse(os.Args[2:])
		if *verifyCmdToken == "" || *verifyCmdKey == "" {
			verifyCmd.PrintDefaults()
			os.Exit(1)
		}
		verifyToken(*verifyCmdToken, *verifyCmdKey)

	case "generate":
		generateKeys()

	default:
		fmt.Printf("Unknown subcommand: %s\n", os.Args[1])
		fmt.Println("Expected 'sign' or 'verify'")
		os.Exit(1)
	}
}

func generateKeys() {
	secretKey := paseto.NewV3AsymmetricSecretKey()

	privateKeyBase64 := base64.StdEncoding.EncodeToString(secretKey.ExportBytes())
	publicKeyBase64 := base64.StdEncoding.EncodeToString(secretKey.Public().ExportBytes())

	fmt.Println("private: '" + privateKeyBase64 + "'")
	fmt.Println("public: '" + publicKeyBase64 + "'")
}

func signToken(message, keyBase64 string, expiration time.Duration) {
	// Decode the private key from base64
	keyBytes, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		fmt.Println("Error decoding private key:", err)
		os.Exit(1)
	}

	// Create a V3 asymmetric secret key from the provided key bytes
	secretKey, err := paseto.NewV3AsymmetricSecretKeyFromBytes(keyBytes)
	if err != nil {
		fmt.Println("Error creating secret key:", err)
		os.Exit(1)
	}

	// Create a new token
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(expiration))
	token.SetString("message", message)

	// Sign the token using V3
	signed := token.V3Sign(secretKey, nil)

	fmt.Println("Signed PASETO v3.public token:")
	fmt.Println(signed)
}

func verifyToken(tokenString, keyBase64 string) {
	// Decode the public key from base64
	keyBytes, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		fmt.Println("Error decoding public key:", err)
		os.Exit(1)
	}

	// Create a V3 asymmetric public key from the provided key bytes
	publicKey, err := paseto.NewV3AsymmetricPublicKeyFromBytes(keyBytes)
	if err != nil {
		fmt.Println("Error creating public key:", err)
		os.Exit(1)
	}

	// Create a parser with expiry check
	parser := paseto.NewParser()

	// Parse and verify the token
	token, err := parser.ParseV3Public(publicKey, tokenString, nil)
	if err != nil {
		fmt.Printf("Token verification failed: %v\n", err)
		os.Exit(1)
	}

	// Extract the message from the verified token
	message, err := token.GetString("message")
	if err != nil {
		fmt.Printf("Failed to get message from token: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Token verification successful!")
	fmt.Println("Message:", message)

	// Print all claims
	fmt.Println("\nAll claims:")
	fmt.Println(string(token.ClaimsJSON()))
}
