// SPDX-FileCopyrightText: 2025 The Kepler Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcryptHashing(t *testing.T) {
	password := "securePassword123"

	// Generate bcrypt hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate bcrypt hash: %v", err)
	}

	// Simulate a correct password check
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		t.Errorf("Expected password match, but got error: %v", err)
	}

	// Simulate an incorrect password check
	wrongPassword := "wrongPassword321"
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(wrongPassword))
	if err == nil {
		t.Errorf("Expected password mismatch, but it matched")
	}
}
