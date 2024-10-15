<!--
order: 3
-->

# CosmWasm Integration

This `x/clock` module does not require any custom bindings. Rather, you must add a Sudo message to your contract. If your contract does not implement this Sudo message, it will be jailed. Please review the following sections below for more information.

> Note: you can find a basic [cw-clock contract here](https://github.com/Reecepbcups/cw-clock-example).

## Implementation

To satisfy the module's requirements, add the following message and entry point to your contract:

```rust
// msg.rs
#[cw_serde]
pub enum SudoMsg {    
    BeginBlock { },
    EndBlock { },
}

// contract.rs
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn sudo(deps: DepsMut, _env: Env, msg: SudoMsg) -> Result<Response, ContractError> {
    match msg {        
        SudoMsg::BeginBlock { } => {}
        SudoMsg::EndBlock { } => {
            
            // TODO: PERFORM LOGIC HERE

            Ok(Response::new())
        }
    }
}
```

At the end of every block, registered contracts will execute the `ClockEndBlock` Sudo message. This is where all of the contract's custom end block logic can be performed. Please keep in mind that contracts which exceed the gas limit specified in the params will be jailed.

## Examples

In the example below, at the end of every block the `val` Config variable will increase by 1. This is a simple example, but one can extrapolate upon this idea and perform actions such as cleanup, auto compounding, etc.

```rust
// msg.rs
#[cw_serde]
pub enum SudoMsg {   
    BeginBlock { }, 
    EndBlock { },
}

// contract.rs
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn sudo(deps: DepsMut, _env: Env, msg: SudoMsg) -> Result<Response, ContractError> {
    match msg {        
        SudoMsg::BeginBlock { } => {}
        SudoMsg::EndBlock { } => {
            let mut config = CONFIG.load(deps.storage)?;
            config.val += 1;
            CONFIG.save(deps.storage, &config)?;

            Ok(Response::new())
        }
    }
}
```

To perform an action occasionally rather than every block, use the `env` variable in the Sudo message to check the block height and then perform logic accordingly. The contract below will only increase the `val` Config variable by 1 if the block height is divisible by 10.

```rust
// msg.rs
#[cw_serde]
pub enum SudoMsg {
    BeginBlock { },    
    EndBlock { },
}

// contract.rs
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn sudo(deps: DepsMut, env: Env, msg: SudoMsg) -> Result<Response, ContractError> {
    match msg {        
        SudoMsg::BeginBlock { } => {}
        SudoMsg::EndBlock { } => {    
            // If the block is not divisible by ten, do nothing.      
            if env.block.height % 10 != 0 {
                return Ok(Response::new());
            }

            let mut config = CONFIG.load(deps.storage)?;
            config.val += 1;
            CONFIG.save(deps.storage, &config)?;

            Ok(Response::new())
        }
    }
}
```
