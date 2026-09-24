// Copyright © 2026 Ping Identity Corporation

//go:build !upgradeladder

package upgradeladder

// Enabled is false unless the build sets the 'upgradeladder' tag, so ladder
// tests inlined in the regular gen test files skip instantly (before any
// server access) in every normal build.
const Enabled = false
