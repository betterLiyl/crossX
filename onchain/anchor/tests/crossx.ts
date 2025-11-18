import * as anchor from "@coral-xyz/anchor";
import { Program } from "@coral-xyz/anchor";
import { PublicKey, SystemProgram } from "@solana/web3.js";

describe("crossx-anchor", () => {
  const provider = anchor.AnchorProvider.env();
  anchor.setProvider(provider);
  const program = anchor.workspace.CrossxAnchor as Program<any>;

  it("initialize and update", async () => {
    const authority = provider.wallet.publicKey;
    const [state, bump] = PublicKey.findProgramAddressSync([
      Buffer.from("crossx"),
      authority.toBuffer(),
    ], program.programId);

    await program.methods.initialize(bump, new anchor.BN(42)).accounts({
      authority,
      state,
      systemProgram: SystemProgram.programId,
    }).rpc();

    const acc = await program.account.state.fetch(state);
    if (acc.value.toNumber() !== 42) throw new Error("init value mismatch");

    await program.methods.updateValue(new anchor.BN(100)).accounts({
      authority,
      state,
    }).rpc();

    const acc2 = await program.account.state.fetch(state);
    if (acc2.value.toNumber() !== 100) throw new Error("update value mismatch");
  });
});

