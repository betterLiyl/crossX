# crossx-anchor (Solana/Anchor)

This subproject provides a minimal on-chain MVP using Solana + Anchor.

## Structure
- `Anchor.toml`: project configuration
- `Cargo.toml`: workspace
- `programs/crossx-anchor`: Rust program implementing PDA state
- `tests/crossx.ts`: mocha+ts tests

## Program Logic
- PDA `state` derived as seeds `[b"crossx", authority]`
- `initialize(bump, init_value)`: creates `state` account, sets `authority`, `value`
- `update_value(new_value)`: only `authority` can update, requires `new_value > 0`

## Build & Test
```bash
anchor build
anchor test
```

## Deploy
```bash
anchor deploy
```

## Notes
- Uses Anchor `0.30.x`
- Follows best practices: PDA auth binding, input validation, error handling
- Requires localnet: `solana-test-validator` and Anchor default wallet
