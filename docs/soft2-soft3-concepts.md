# soft2 and soft3

## content from presentation

## particle

A particle is a content identifier (CID) of the file in IPFS network. CIDs are based on the file’s cryptographic hash. That means:

- Any difference in the content will produce a different CID and
- The same content added to two different IPFS nodes using the same settings will produce the same CID.

A file can be retrieved from IPFS network using this hash. Particles are written into Bostrom blockchain in a form of cyberlinks.

## cyberlink

A cyberlink is a link between two particles registered in Bostrom blockchain by a particular neuron.

Cyberlinks are the edges of the knowledge graph, particles are the vertexes (aka nodes).

## neuron

Neuron is an agent that creates cyberlinks. A neuron can be:

- a private key holder;
- a cosm-wasm contract (autonomus program).

## knowledge graph

The knowledge graph is a directed weighted graph between particles.

The knowledge graph of Bostrom blockchain consists of pairs: each source particle is connected to a destination particle via cyberlink (with additional information of neuron's address and the height info). So each cyberlink recorded into blockchain in a form of:

`source_particle - destination_particle - neuron - height`

## content oracle
