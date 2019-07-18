# An Overview of Cyber Governance

The Cyberd as the Cosmos Hub has launched and one of its core features is the ability for CYB holders to collectively govern the blockchain. CYB holders can submit proposals and signal their approval or disapproval to proposals submitted to the network by signing a special type of transaction. The following article will cover the governance procedure on the Cyber and introduce our governance tool that allows any CYB holder to participate in the on-chain governance mechanism.

## The Governance Procedure

In the current Cyber governance implementation anyone can submit a text proposal to the system. A minimum deposit is required for a proposal to enter the voting period during which CYB holders will be able to vote on whether the proposal is accepted or not. The numbers used in the following are based on the parameters implemented on the Cyberd at the time of writing (18th July 2019).

![governance2](https://ipfs.io/ipfs/QmRR3jeAdW8mpkxFoinrMRzMRmARgmz1hKCtHEvLSTiJ2u)

<p align="center">An Illustration of the Governance Process on the Cyber</p>

## Phase 1: The Deposit Period

For a proposal to be considered for voting, a minimum deposit of 50 MCYBs needs to be deposited within 2 weeks from when the proposal was submitted. Any CYB holder can contribute to this deposit to support proposals, meaning that the party submitting the proposal doesn’t necessarily need to provide the deposit itself. The deposit is required as spam protection, CYB holders that contributed CYBs to a proposal will be able to collect their deposit when the proposal was accepted or when it did not reach the minimum threshold after 2 weeks.

## Phase 2: The Voting Period

When the minimum deposit for a particular proposal is reached the 2 week (336h) long voting period begins. During this period, CYB holders are able to cast their vote on that proposal. There are 4 voting options ("Yes", "No", "No with Veto", "Abstain").

Important details in the governance implementation of Cyber are:

- Only staked (bonded) tokens can participate in governance.
- Voting power is measured in terms of stake. The number of CYBs you stake determine your influence on the decision (coin voting).
- Delegators inherit the vote of the validators they are delegated to unless they cast their own vote, which will overwrite validator decisions.

## Phase 3: Tallying Results
If a proposal is accepted depends on the result of the coin voting by CYB holders. The following requirements need to be satisfied for a proposal to be considered accepted:

- **Quorum:** More than 40% of the total staked tokens at the end of the voting period need to have participated in the vote.
- **Threshold:** More than 50% of the tokens that participated in the vote (after excluding “Abstain” votes) need to have voted in favor of the proposal (“Yes”).
- **Veto:** Less than 33.4% of the tokens that participated in the vote (after excluding “Abstain” votes) need to have vetoed the decision (“No with Veto”).

If one of these requirements is not met at the end of the voting period, e.g. because the quorum was not met, the proposal is denied. As a result, the deposit associated with the denied proposal will not be refunded and instead be awarded to the community pool.

## Phase 4: Implementing the Proposal

An accepted proposal will need to be implemented as part of the software that is run by the networks’ validators.

For an exemplary governance vote, check out the [first governance Cosmos Hub proposal](https://hubble.figment.network/cosmos/chains/cosmoshub-1/governance/proposals/1) from validator team B-Harvest about adjusting the block time rate used to calculate network inflation to mirror actual conditions in the network.

## How to Participate in Governance

At this time only CLI tool implemented for submitting and voting. Running docker cyberd container, synced node and staked CYBs.

You can get the list of proposals by following command:

```bash
docker exec cyberd cyberdcli q gov proposals
```

And vote by command above:

```bash
docker exec cyberd cyberdcli tx gov vote [proposal-id] [option] --from [keyname] --chain-id [chain id]
```

The **[proposal-id]** you can find at the list of proposals by command below.

The **[option]** should be ***yes***, ***no***, ***abstain*** and ***no_with_veto***.

## Conclusion

We believe the governance implementation of Cyber based Cosmos with delegators inheriting validator votes, but delegators being able to overwrite validator choices and make their own decision is a great solution to the low governance turnouts that can be observed in other blockchain networks. We think this approach does good in striking a balance between representative democracy (“proxy voting”) for less interested holders, while still allowing for direct influence for token holders that want to directly engage in network governance.
