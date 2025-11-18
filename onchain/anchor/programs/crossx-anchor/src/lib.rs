use anchor_lang::prelude::*;

declare_id!("11111111111111111111111111111111");

#[program]
pub mod crossx_anchor {
    use super::*;

    pub fn initialize(ctx: Context<Initialize>, bump: u8, init_value: u64) -> Result<()> {
        let state = &mut ctx.accounts.state;
        state.authority = ctx.accounts.authority.key();
        state.bump = bump;
        state.value = init_value;
        Ok(())
    }

    pub fn update_value(ctx: Context<UpdateValue>, new_value: u64) -> Result<()> {
        require!(new_value > 0, CrossxError::InvalidValue);
        let state = &mut ctx.accounts.state;
        require_keys_eq!(ctx.accounts.authority.key(), state.authority, CrossxError::Unauthorized);
        state.value = new_value;
        Ok(())
    }
}

#[derive(Accounts)]
pub struct Initialize<'info> {
    #[account(mut)]
    pub authority: Signer<'info>,
    #[account(
        init,
        payer = authority,
        seeds = [b"crossx", authority.key().as_ref()],
        bump,
        space = 8 + State::SIZE,
    )]
    pub state: Account<'info, State>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct UpdateValue<'info> {
    pub authority: Signer<'info>,
    #[account(
        mut,
        seeds = [b"crossx", authority.key().as_ref()],
        bump = state.bump,
        has_one = authority,
    )]
    pub state: Account<'info, State>,
}

#[account]
pub struct State {
    pub authority: Pubkey,
    pub bump: u8,
    pub value: u64,
}

impl State {
    pub const SIZE: usize = 32 + 1 + 8;
}

#[error_code]
pub enum CrossxError {
    #[msg("Unauthorized")]
    Unauthorized,
    #[msg("Invalid value")]
    InvalidValue,
}

