# Events

The resources module emits the following events:

## Msg's

### MsgInvestmint
| Type             | Attribute Key | Attribute Value        |
| ---------------- | ------------- | ---------------------- |
| message          | module        | energy                 |
| message          | action        | investmint             | 
| investmint       | agent         | {agentAddress}         |
| investmint       | amount        | {baseResourceAmount}   |
| investmint       | resource      | {resourceDenom}        |
| investmint       | length        | {seconds}              |
| investmint       | minted        | {mintedResourceAmount} |