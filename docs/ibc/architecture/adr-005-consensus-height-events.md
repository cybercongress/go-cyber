# ADR 005: UpdateClient Events - ClientState Consensus Heights

## Changelog
* 25/04/2022: initial draft

## Status

Accepted 

## Context

The `ibc-go` implementation leverages the [Cosmos-SDK's EventManager](https://github.com/cosmos/cosmos-sdk/blob/v0.45.4/docs/core/events.md#EventManager) to provide subscribers a method of reacting to application specific events.
Some IBC relayers depend on the [`consensus_height`](https://github.com/cosmos/ibc-go/blob/v3.0.0/modules/core/02-client/keeper/events.go#L33) attribute emitted as part of `UpdateClient` events in order to run `07-tendermint` misbehaviour detection by cross-checking the details of the *Header* emitted at a given consensus height against those of the *Header* from the originating chain. This includes such details as:

- The `SignedHeader` containing the commitment root.
- The `ValidatorSet` that signed the *Header*.
- The `TrustedHeight` seen by the client at less than or equal to the height of *Header*.
- The last `TrustedValidatorSet` at the trusted height.

Following the refactor of the `02-client` submodule and associated `ClientState` interfaces, it will now be possible for
light client implementations to perform such actions as batch updates, inserting `N` number of `ConsensusState`s into the application state tree with a single `UpdateClient` message. This flexibility is provided in `ibc-go` by the usage of the [Protobuf `Any`](https://developers.google.com/protocol-buffers/docs/proto3#any) field contained within the [`UpdateClient`](https://github.com/cosmos/ibc-go/blob/v3.0.0/proto/ibc/core/client/v1/tx.proto#L44) message.
For example, a batched client update message serialized as a Protobuf `Any` type for the `07-tendermint` lightclient implementation could be defined as follows:

```protobuf
message BatchedHeaders {
    repeated Header headers = 1;
}
```

To complement this flexibility, the `UpdateClient` handler will now support the submission of [client misbehaviour](https://github.com/cosmos/ibc/tree/master/spec/core/ics-002-client-semantics#misbehaviour) by consolidating the `Header` and `Misbehaviour` interfaces into a single `ClientMessage` interface type:

```go
// ClientMessage is an interface used to update an IBC client.
// The update may be done by a single header, a batch of headers, misbehaviour, or any type which when verified produces
// a change to state of the IBC client
type ClientMessage interface {
	proto.Message

	ClientType() string
	ValidateBasic() error
}
```

To support this functionality the `GetHeight()` method has been omitted from the new `ClientMessage` interface.
Emission of standardised events from the `02-client` submodule now becomes problematic and is two-fold:

1. The `02-client` submodule previously depended upon the `GetHeight()` method of `Header` types in order to [retrieve the updated consensus height](https://github.com/cosmos/ibc-go/blob/v3.0.0/modules/core/02-client/keeper/client.go#L90).
2. Emitting a single `consensus_height` event attribute is not sufficient in the case of a batched client update containing multiple *Headers*.

## Decision

The following decisions have been made in order to provide flexibility to consumers of `UpdateClient` events in a non-breaking fashion:

1. Return a list of updated consensus heights `[]exported.Height` from the new `UpdateState` method of the `ClientState` interface.

```go
// UpdateState updates and stores as necessary any associated information for an IBC client, such as the ClientState and corresponding ConsensusState.
// Upon successful update, a list of consensus heights is returned. It assumes the ClientMessage has already been verified.
UpdateState(sdk.Context, codec.BinaryCodec, sdk.KVStore, ClientMessage) []Height
```

2. Maintain the `consensus_height` event attribute emitted from the `02-client` update handler, but mark as deprecated for future removal. For example, with tendermint lightclients this will simply be `consensusHeights[0]` following a successful update using a single *Header*.

3. Add an additional `consensus_heights` event attribute, containing a comma separated list of updated heights. This provides flexibility for emitting a single consensus height or multiple consensus heights in the example use-case of batched header updates.

## Consequences

### Positive

- Subscribers of IBC core events can act upon `UpdateClient` events containing one or more consensus heights.
- Deprecation of the existing `consensus_height` attribute allows consumers to continue to process `UpdateClient` events as normal, with a path to upgrade to using the `consensus_heights` attribute moving forward.

### Negative

- Consumers of IBC core `UpdateClient` events are forced to make future code changes.

### Neutral

## References

Discussions:
- [#1208](https://github.com/cosmos/ibc-go/pull/1208#discussion_r839691927)

Issues:
- [#594](https://github.com/cosmos/ibc-go/issues/594)

PRs:
- [#1285](https://github.com/cosmos/ibc-go/pull/1285)
