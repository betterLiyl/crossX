# On-chain Communication (Go ↔ Anchor/Solana)

## Overview
- Gateway and matching engine are in Go; on-chain program is in `onchain/anchor` (Solana + Anchor).
- Go can call Anchor program instructions via Solana RPC using a Go SDK.

## Recommended Setup
- RPC URL: `SOLANA_RPC_URL` (e.g., `http://localhost:8899` for localnet)
- Wallet: `SOLANA_WALLET_PATH` pointing to a keypair JSON (Anchor default `~/.config/solana/id.json`).
- Program ID: set the actual deployed program id after `anchor deploy`.

## Instruction Calls
- `initialize(bump, init_value)` accounts: `authority (signer)`, `state (PDA)`, `system_program`.
- `update_value(new_value)` accounts: `authority (signer)`, `state (PDA)`.
- PDA seeds: `[b"crossx", authority_pubkey]` with `bump`.

## Flow
- Derive PDA → Build instruction (borsh-encoded per IDL) → Sign and send transaction → Await confirmation → Read account state.

## Go SDK Options
- Use `github.com/portto/solana-go-sdk` or similar to:
  - Load wallet, derive PDA, build transactions, send via RPC.
  - Encode instruction data according to Anchor IDL or hardcoded layout.

## Environment
- `ANCHOR_PROVIDER_URL` and `ANCHOR_WALLET` for tests.
- `SOLANA_RPC_URL`, `SOLANA_WALLET_PATH`, `ANCHOR_PROGRAM_ID` for Go integration.

## Notes
- Keep keys out of logs; use env or KMS.
- Use `confirmed` or `finalized` commitment per business needs.
- Validate inputs and check transaction errors.