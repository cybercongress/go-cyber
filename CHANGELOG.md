# Change Log

## [v0.1.2](https://github.com/cybercongress/cyberd/tree/v0.1.2) (2019-02-23)
[Full Changelog](https://github.com/cybercongress/cyberd/compare/v0.1.1...v0.1.2)

**Implemented enhancements:**

- Update  to cosmos 0.32.0 [\#265](https://github.com/cybercongress/cyberd/issues/265)
- Download golang only once during docker build. [\#259](https://github.com/cybercongress/cyberd/issues/259)
- Do not calculate rank, if there is no new links on last window. [\#258](https://github.com/cybercongress/cyberd/issues/258)
- Http client only new links option for submit endpoint [\#254](https://github.com/cybercongress/cyberd/issues/254)
- Add fail at height debug flags [\#250](https://github.com/cybercongress/cyberd/issues/250)
- Extend isLinkExist endpoint with possibility not specify address param [\#242](https://github.com/cybercongress/cyberd/issues/242)
- Add staking/validators endpoint to node RPC [\#235](https://github.com/cybercongress/cyberd/issues/235)
- Add cyberd rest light-client to docker image [\#230](https://github.com/cybercongress/cyberd/issues/230)
- Improve building of links index for search queries [\#145](https://github.com/cybercongress/cyberd/issues/145)
- Webscoket API [\#76](https://github.com/cybercongress/cyberd/issues/76)

**Fixed bugs:**

- Rank calculation failed: out of memory [\#264](https://github.com/cybercongress/cyberd/issues/264)
- Different bandwidth price on same block height [\#247](https://github.com/cybercongress/cyberd/issues/247)
- Consensus failure after node restart [\#243](https://github.com/cybercongress/cyberd/issues/243)
- Search returns 500 error when previously unexisted "from" CID added [\#241](https://github.com/cybercongress/cyberd/issues/241)
- Rest server  blocks cli keystore file [\#233](https://github.com/cybercongress/cyberd/issues/233)
- Cannot interact with node [\#203](https://github.com/cybercongress/cyberd/issues/203)
- Update cyberd.ai [\#185](https://github.com/cybercongress/cyberd/issues/185)
- CI: Cyberd docker image building. Include genesis file [\#160](https://github.com/cybercongress/cyberd/issues/160)

**Closed issues:**

- Add more seeds into initial config [\#244](https://github.com/cybercongress/cyberd/issues/244)
- Create awesome gif for readme [\#239](https://github.com/cybercongress/cyberd/issues/239)
- Launch 0.1.1 testnet [\#205](https://github.com/cybercongress/cyberd/issues/205)
- Cyberd release scheme. [\#180](https://github.com/cybercongress/cyberd/issues/180)
- Genesis state generator [\#92](https://github.com/cybercongress/cyberd/issues/92)
- Distribution [\#27](https://github.com/cybercongress/cyberd/issues/27)
- Ranking research [\#24](https://github.com/cybercongress/cyberd/issues/24)
- Define distribution in paper [\#5](https://github.com/cybercongress/cyberd/issues/5)

**Merged pull requests:**

- 0.1.2 [\#269](https://github.com/cybercongress/cyberd/pull/269) ([hleb-albau](https://github.com/hleb-albau))
- Refactoring  of search index. Fix links export [\#268](https://github.com/cybercongress/cyberd/pull/268) ([arturalbov](https://github.com/arturalbov))
- Launch euler-2 testnet [\#267](https://github.com/cybercongress/cyberd/pull/267) ([hleb-albau](https://github.com/hleb-albau))
- \#265 Update to cosmos 0.32.0 [\#266](https://github.com/cybercongress/cyberd/pull/266) ([hleb-albau](https://github.com/hleb-albau))
- Consensus failure after node restart [\#263](https://github.com/cybercongress/cyberd/pull/263) ([arturalbov](https://github.com/arturalbov))
- Do not calculate rank, if there is no new links on last window. [\#261](https://github.com/cybercongress/cyberd/pull/261) ([arturalbov](https://github.com/arturalbov))
- Client submit only new links flag [\#260](https://github.com/cybercongress/cyberd/pull/260) ([hleb-albau](https://github.com/hleb-albau))
- Improve building of links index for search queries [\#255](https://github.com/cybercongress/cyberd/pull/255) ([arturalbov](https://github.com/arturalbov))
- \#244 Add moon as seed node [\#253](https://github.com/cybercongress/cyberd/pull/253) ([hleb-albau](https://github.com/hleb-albau))
- \#250 Add fail at height debug flag [\#251](https://github.com/cybercongress/cyberd/pull/251) ([hleb-albau](https://github.com/hleb-albau))
- \#247 Different bandwidth price on same block height [\#248](https://github.com/cybercongress/cyberd/pull/248) ([hleb-albau](https://github.com/hleb-albau))
- \#242 Extend isLinkExist endpoint with possibility not specify address [\#246](https://github.com/cybercongress/cyberd/pull/246) ([hleb-albau](https://github.com/hleb-albau))
- \#242 Extend isLinkExist endpoint with possibility not specify address [\#245](https://github.com/cybercongress/cyberd/pull/245) ([hleb-albau](https://github.com/hleb-albau))
- Add staking/validators endpoint to node RPC [\#236](https://github.com/cybercongress/cyberd/pull/236) ([arturalbov](https://github.com/arturalbov))
- Rest server blocks cli keystore file [\#234](https://github.com/cybercongress/cyberd/pull/234) ([arturalbov](https://github.com/arturalbov))
- UPD run validator md [\#232](https://github.com/cybercongress/cyberd/pull/232) ([SaveTheAles](https://github.com/SaveTheAles))
- Add cyberd rest light-client to docker image [\#231](https://github.com/cybercongress/cyberd/pull/231) ([arturalbov](https://github.com/arturalbov))
- Fix docs [\#228](https://github.com/cybercongress/cyberd/pull/228) ([arturalbov](https://github.com/arturalbov))
- upd run-validator.md [\#225](https://github.com/cybercongress/cyberd/pull/225) ([SaveTheAles](https://github.com/SaveTheAles))

## [v0.1.1](https://github.com/cybercongress/cyberd/tree/v0.1.1) (2019-01-31)
[Full Changelog](https://github.com/cybercongress/cyberd/compare/euler-1...v0.1.1)

**Fixed bugs:**

- Error on replay with wrong app hash calculation \(block \#914381\) [\#197](https://github.com/cybercongress/cyberd/issues/197)

**Closed issues:**

- Mass testing of cyberd-testnet 0.1.1 [\#219](https://github.com/cybercongress/cyberd/issues/219)
- Export 0.1.1 genesis from 0.1.0 testnet. [\#215](https://github.com/cybercongress/cyberd/issues/215)
- Setup 0.1.1 chain params for staking, distr, slashing, mint modules. [\#201](https://github.com/cybercongress/cyberd/issues/201)

## [euler-1](https://github.com/cybercongress/cyberd/tree/euler-1) (2019-01-31)
[Full Changelog](https://github.com/cybercongress/cyberd/compare/euler-1.1-dev1...euler-1)

**Merged pull requests:**

- Euler 0.1.1 [\#223](https://github.com/cybercongress/cyberd/pull/223) ([arturalbov](https://github.com/arturalbov))
- Fix import reader [\#222](https://github.com/cybercongress/cyberd/pull/222) ([arturalbov](https://github.com/arturalbov))
- Landing update [\#220](https://github.com/cybercongress/cyberd/pull/220) ([asadovka](https://github.com/asadovka))
- DISCUSS: Chain params for 0.1.1 network [\#202](https://github.com/cybercongress/cyberd/pull/202) ([hleb-albau](https://github.com/hleb-albau))

## [euler-1.1-dev1](https://github.com/cybercongress/cyberd/tree/euler-1.1-dev1) (2019-01-30)
[Full Changelog](https://github.com/cybercongress/cyberd/compare/euler-1.1-dev0...euler-1.1-dev1)

**Implemented enhancements:**

- Docker image  with possibility to select compute unit [\#211](https://github.com/cybercongress/cyberd/issues/211)

**Fixed bugs:**

- Link msg bandwidth cost doesn't take into account links count  [\#216](https://github.com/cybercongress/cyberd/issues/216)

**Merged pull requests:**

- \#215 Fix pou distribution uses share not percentage [\#221](https://github.com/cybercongress/cyberd/pull/221) ([hleb-albau](https://github.com/hleb-albau))
- \#215 Fix gentx commands. [\#218](https://github.com/cybercongress/cyberd/pull/218) ([hleb-albau](https://github.com/hleb-albau))

## [euler-1.1-dev0](https://github.com/cybercongress/cyberd/tree/euler-1.1-dev0) (2019-01-29)
[Full Changelog](https://github.com/cybercongress/cyberd/compare/v0.1.0...euler-1.1-dev0)

**Implemented enhancements:**

- Links Exporter [\#184](https://github.com/cybercongress/cyberd/issues/184)

**Fixed bugs:**

- Fix validator manual [\#189](https://github.com/cybercongress/cyberd/issues/189)
- Change mint module to default cosmos module [\#181](https://github.com/cybercongress/cyberd/issues/181)
- Current Bandwidth do not respect price [\#168](https://github.com/cybercongress/cyberd/issues/168)
- negative bandwidnt  [\#164](https://github.com/cybercongress/cyberd/issues/164)
- Wrong CYB allocation in Genesis [\#156](https://github.com/cybercongress/cyberd/issues/156)
- Update docker with 0.1.0 version [\#154](https://github.com/cybercongress/cyberd/issues/154)
- Move import\_private command to keys subcomand [\#152](https://github.com/cybercongress/cyberd/issues/152)
- Handle import of ethereum privkeys with 0x [\#150](https://github.com/cybercongress/cyberd/issues/150)

**Closed issues:**

- Fix typos [\#207](https://github.com/cybercongress/cyberd/issues/207)
- Update landing for 0.4 whitepaper [\#204](https://github.com/cybercongress/cyberd/issues/204)
- Update to tendermint v29.0 [\#187](https://github.com/cybercongress/cyberd/issues/187)
- Change Bandwidth Price to Average for 24h Sliding Window [\#179](https://github.com/cybercongress/cyberd/issues/179)
- High bandwidth cost for create validator msg [\#178](https://github.com/cybercongress/cyberd/issues/178)
- Bandwidth Specification Change [\#177](https://github.com/cybercongress/cyberd/issues/177)
- Add rpc endpoint to submit signed link and send messages [\#173](https://github.com/cybercongress/cyberd/issues/173)
- Tx size reduction: switch to more compact links tx encoding. [\#171](https://github.com/cybercongress/cyberd/issues/171)
- index entities count  endpoint [\#165](https://github.com/cybercongress/cyberd/issues/165)
- Add ipfs hashes to release [\#161](https://github.com/cybercongress/cyberd/issues/161)
- Change token name from CBD to CYB in the testnetwork [\#157](https://github.com/cybercongress/cyberd/issues/157)
- Add more issues templates. [\#153](https://github.com/cybercongress/cyberd/issues/153)
- Add trust-node by default in cli [\#151](https://github.com/cybercongress/cyberd/issues/151)
- Numerous small fixes in whitepaper [\#89](https://github.com/cybercongress/cyberd/issues/89)
- Cyberd landing make up [\#36](https://github.com/cybercongress/cyberd/issues/36)
- Whitepaper 0.4 [\#25](https://github.com/cybercongress/cyberd/issues/25)

**Merged pull requests:**

- CI: testnet ipfs files hashes   [\#217](https://github.com/cybercongress/cyberd/pull/217) ([arturalbov](https://github.com/arturalbov))
- WIP: CI Cyberd docker image building. Include genesis file [\#214](https://github.com/cybercongress/cyberd/pull/214) ([arturalbov](https://github.com/arturalbov))
- \#197 Error on replay with wrong app hash calculation [\#213](https://github.com/cybercongress/cyberd/pull/213) ([hleb-albau](https://github.com/hleb-albau))
- Links exporter [\#210](https://github.com/cybercongress/cyberd/pull/210) ([arturalbov](https://github.com/arturalbov))
- R4R \#181 Constant Block Reward [\#209](https://github.com/cybercongress/cyberd/pull/209) ([hleb-albau](https://github.com/hleb-albau))
- 157 change token protocol name [\#208](https://github.com/cybercongress/cyberd/pull/208) ([hleb-albau](https://github.com/hleb-albau))
- Fix typos get\_CYB.md [\#206](https://github.com/cybercongress/cyberd/pull/206) ([serejandmyself](https://github.com/serejandmyself))
- R4R: \#156 Add Missing euler tokens cmd [\#200](https://github.com/cybercongress/cyberd/pull/200) ([hleb-albau](https://github.com/hleb-albau))
- Update to cosmos-sdk v0.30.0 [\#198](https://github.com/cybercongress/cyberd/pull/198) ([hleb-albau](https://github.com/hleb-albau))
- Tx size reduction: switch to more compact links tx encoding [\#192](https://github.com/cybercongress/cyberd/pull/192) ([arturalbov](https://github.com/arturalbov))
- Fix rank calculation context data [\#191](https://github.com/cybercongress/cyberd/pull/191) ([hleb-albau](https://github.com/hleb-albau))
- Important fixes in validators manual [\#190](https://github.com/cybercongress/cyberd/pull/190) ([xhipster](https://github.com/xhipster))
- Change Bandwidth Price to Average for 24h Sliding Window [\#188](https://github.com/cybercongress/cyberd/pull/188) ([arturalbov](https://github.com/arturalbov))
- R4R Update cosmos to latest develop [\#186](https://github.com/cybercongress/cyberd/pull/186) ([hleb-albau](https://github.com/hleb-albau))
- \#153 Add more issues templates. [\#183](https://github.com/cybercongress/cyberd/pull/183) ([hleb-albau](https://github.com/hleb-albau))
- \#178 \#168 Bandwidth price + Msg cost for non link txes [\#182](https://github.com/cybercongress/cyberd/pull/182) ([hleb-albau](https://github.com/hleb-albau))
- \#173 Add rpc endpoint to submit signed link and send messages [\#174](https://github.com/cybercongress/cyberd/pull/174) ([hleb-albau](https://github.com/hleb-albau))
- Technical notes on euler release [\#169](https://github.com/cybercongress/cyberd/pull/169) ([xhipster](https://github.com/xhipster))
- Handle import of ethereum privkeys with 0x. Trust node by default  [\#167](https://github.com/cybercongress/cyberd/pull/167) ([arturalbov](https://github.com/arturalbov))
- \#164 fix rpc negative bw \#165 add index entities count endpoint [\#166](https://github.com/cybercongress/cyberd/pull/166) ([hleb-albau](https://github.com/hleb-albau))
- Move import\_private command to keys subcomand [\#159](https://github.com/cybercongress/cyberd/pull/159) ([arturalbov](https://github.com/arturalbov))

## [v0.1.0](https://github.com/cybercongress/cyberd/tree/v0.1.0) (2019-01-06)
[Full Changelog](https://github.com/cybercongress/cyberd/compare/v0.0.8...v0.1.0)

**Implemented enhancements:**

- CID rank merkle proofs [\#133](https://github.com/cybercongress/cyberd/issues/133)
- Simple go cyberd client [\#122](https://github.com/cybercongress/cyberd/issues/122)
- Add remainig bandwidth endpoint [\#109](https://github.com/cybercongress/cyberd/issues/109)
- Get rid of RPC proxy [\#94](https://github.com/cybercongress/cyberd/issues/94)

**Fixed bugs:**

- Cannot check balance using cli [\#149](https://github.com/cybercongress/cyberd/issues/149)
- Index out of range for cids with not calculated rank yet. [\#139](https://github.com/cybercongress/cyberd/issues/139)
- Account inmem balances should be updated by wrapping accountKeeper [\#99](https://github.com/cybercongress/cyberd/issues/99)

**Closed issues:**

- Cli: Add possibility to restore acc from priv keys [\#146](https://github.com/cybercongress/cyberd/issues/146)
- Go mod Ci dependency error  [\#128](https://github.com/cybercongress/cyberd/issues/128)
- Update to cosmos 29 [\#126](https://github.com/cybercongress/cyberd/issues/126)
- Write down bandwidth specification [\#114](https://github.com/cybercongress/cyberd/issues/114)
- Make docker container based on  nvidia-gpu image. [\#104](https://github.com/cybercongress/cyberd/issues/104)
- Index transactions by addresses. [\#103](https://github.com/cybercongress/cyberd/issues/103)
- Calculate rank in ||, post results each 600 blocks. [\#101](https://github.com/cybercongress/cyberd/issues/101)
- Linkchains support [\#91](https://github.com/cybercongress/cyberd/issues/91)
- Include addresses into knowledge graph [\#90](https://github.com/cybercongress/cyberd/issues/90)
- Graphics for whitepaper [\#88](https://github.com/cybercongress/cyberd/issues/88)
- Add bandwidth by stake [\#77](https://github.com/cybercongress/cyberd/issues/77)
- Launch testnet Euler [\#73](https://github.com/cybercongress/cyberd/issues/73)
- Test Ethereum Network Statistical Significance [\#52](https://github.com/cybercongress/cyberd/issues/52)
- Cyberd landing design [\#51](https://github.com/cybercongress/cyberd/issues/51)
- Simulation [\#31](https://github.com/cybercongress/cyberd/issues/31)
- Link Chain PoC | Zeronet [\#26](https://github.com/cybercongress/cyberd/issues/26)
- Write benchmark and test SpringRank with different amount of objects/edges [\#22](https://github.com/cybercongress/cyberd/issues/22)
- Parse Ethereum network and calculate SpringRank [\#21](https://github.com/cybercongress/cyberd/issues/21)
- Research basic chains fundamentals [\#15](https://github.com/cybercongress/cyberd/issues/15)
- Perfomance testing of Solana [\#10](https://github.com/cybercongress/cyberd/issues/10)
- Research on perfomance of consensus computers [\#9](https://github.com/cybercongress/cyberd/issues/9)
- Perfomance testing scenario for cyberd [\#8](https://github.com/cybercongress/cyberd/issues/8)

**Merged pull requests:**

- \[euler\] Launch euler testnet [\#148](https://github.com/cybercongress/cyberd/pull/148) ([hleb-albau](https://github.com/hleb-albau))
- Cli: Add possibility to restore acc from priv keys [\#147](https://github.com/cybercongress/cyberd/pull/147) ([arturalbov](https://github.com/arturalbov))
- Rank merkle proofs [\#144](https://github.com/cybercongress/cyberd/pull/144) ([arturalbov](https://github.com/arturalbov))
- Setup bw params [\#143](https://github.com/cybercongress/cyberd/pull/143) ([hleb-albau](https://github.com/hleb-albau))
- Merkle tree implementation [\#141](https://github.com/cybercongress/cyberd/pull/141) ([arturalbov](https://github.com/arturalbov))
- \#139 Index out of range for cids with not calculated rank yet. [\#140](https://github.com/cybercongress/cyberd/pull/140) ([hleb-albau](https://github.com/hleb-albau))
- Fix http client and rank logs [\#137](https://github.com/cybercongress/cyberd/pull/137) ([hleb-albau](https://github.com/hleb-albau))
- \[euler-dev3\] New testnet [\#136](https://github.com/cybercongress/cyberd/pull/136) ([hleb-albau](https://github.com/hleb-albau))
- \#101 Copy state before index creation [\#135](https://github.com/cybercongress/cyberd/pull/135) ([hleb-albau](https://github.com/hleb-albau))
- \#92 Populate state with random addresses [\#134](https://github.com/cybercongress/cyberd/pull/134) ([hleb-albau](https://github.com/hleb-albau))
- Small fixes [\#131](https://github.com/cybercongress/cyberd/pull/131) ([arturalbov](https://github.com/arturalbov))
- Parallel rank calculation [\#130](https://github.com/cybercongress/cyberd/pull/130) ([arturalbov](https://github.com/arturalbov))
- Update stake index every block [\#129](https://github.com/cybercongress/cyberd/pull/129) ([hleb-albau](https://github.com/hleb-albau))
- \#126 Update to cosmos 29 [\#127](https://github.com/cybercongress/cyberd/pull/127) ([hleb-albau](https://github.com/hleb-albau))
- Refactoring [\#125](https://github.com/cybercongress/cyberd/pull/125) ([hleb-albau](https://github.com/hleb-albau))
- Various Bug fixes [\#124](https://github.com/cybercongress/cyberd/pull/124) ([hleb-albau](https://github.com/hleb-albau))
- In-memory storages refactoring [\#123](https://github.com/cybercongress/cyberd/pull/123) ([arturalbov](https://github.com/arturalbov))
- \#114 Write bw specification [\#121](https://github.com/cybercongress/cyberd/pull/121) ([hleb-albau](https://github.com/hleb-albau))
- WIP \#114 Add cbdbank module [\#120](https://github.com/cybercongress/cyberd/pull/120) ([hleb-albau](https://github.com/hleb-albau))
- Update readme [\#119](https://github.com/cybercongress/cyberd/pull/119) ([hleb-albau](https://github.com/hleb-albau))
- Clean up folders [\#118](https://github.com/cybercongress/cyberd/pull/118) ([arturalbov](https://github.com/arturalbov))
- \#103 Index transactions by signers. [\#117](https://github.com/cybercongress/cyberd/pull/117) ([hleb-albau](https://github.com/hleb-albau))

## [v0.0.8](https://github.com/cybercongress/cyberd/tree/v0.0.8) (2018-12-11)
[Full Changelog](https://github.com/cybercongress/cyberd/compare/v0.0.7...v0.0.8)

**Implemented enhancements:**

- Cid validation [\#93](https://github.com/cybercongress/cyberd/issues/93)
- Update to cosmos-sdk version 0.26.1 [\#79](https://github.com/cybercongress/cyberd/issues/79)
- Build node releases with cleveldb [\#59](https://github.com/cybercongress/cyberd/issues/59)
- Remove 'cosmosaccaddr' prefix from cyberd address [\#39](https://github.com/cybercongress/cyberd/issues/39)

**Closed issues:**

- Make up cyberd landing [\#87](https://github.com/cybercongress/cyberd/issues/87)
- Add possibility to join for new validators. [\#75](https://github.com/cybercongress/cyberd/issues/75)
- Calculate rank using GPU [\#74](https://github.com/cybercongress/cyberd/issues/74)
- Create basic wiki cyberd indexer [\#71](https://github.com/cybercongress/cyberd/issues/71)
- Create cyberd PoC based on Cosmos SDK [\#37](https://github.com/cybercongress/cyberd/issues/37)
- Draw logo for cyberd [\#16](https://github.com/cybercongress/cyberd/issues/16)
- Build basic economic model [\#1](https://github.com/cybercongress/cyberd/issues/1)

**Merged pull requests:**

- \#76 Define Basic RPC specification [\#111](https://github.com/cybercongress/cyberd/pull/111) ([hleb-albau](https://github.com/hleb-albau))
- \#104 Make docker container based on nvidia-gpu image [\#110](https://github.com/cybercongress/cyberd/pull/110) ([hleb-albau](https://github.com/hleb-albau))
- \[DON'T MERGE\] Add bandwidth by stake. Part 2. [\#108](https://github.com/cybercongress/cyberd/pull/108) ([arturalbov](https://github.com/arturalbov))
- Add bandwidth by stake Part 1 [\#107](https://github.com/cybercongress/cyberd/pull/107) ([arturalbov](https://github.com/arturalbov))
- \#93 Cid validation [\#106](https://github.com/cybercongress/cyberd/pull/106) ([hleb-albau](https://github.com/hleb-albau))
- \#78 Add guide `How to join network as validator` [\#105](https://github.com/cybercongress/cyberd/pull/105) ([hleb-albau](https://github.com/hleb-albau))
-  \#1 Build basic economic model [\#102](https://github.com/cybercongress/cyberd/pull/102) ([hleb-albau](https://github.com/hleb-albau))
- Update cosmos to 0.27.0 [\#100](https://github.com/cybercongress/cyberd/pull/100) ([hleb-albau](https://github.com/hleb-albau))
- Small fixes for validators joining [\#98](https://github.com/cybercongress/cyberd/pull/98) ([arturalbov](https://github.com/arturalbov))
- Remove poc folder [\#97](https://github.com/cybercongress/cyberd/pull/97) ([hleb-albau](https://github.com/hleb-albau))
- Add possibility to join for new validators [\#96](https://github.com/cybercongress/cyberd/pull/96) ([arturalbov](https://github.com/arturalbov))
- Calculate eth network significance [\#85](https://github.com/cybercongress/cyberd/pull/85) ([hleb-albau](https://github.com/hleb-albau))
- 74 gpu rank calculation [\#83](https://github.com/cybercongress/cyberd/pull/83) ([hleb-albau](https://github.com/hleb-albau))
- Remove 'cosmosaccaddr' prefix from cyberd address [\#82](https://github.com/cybercongress/cyberd/pull/82) ([arturalbov](https://github.com/arturalbov))
- UPD docs\_upd job [\#81](https://github.com/cybercongress/cyberd/pull/81) ([SaveTheAles](https://github.com/SaveTheAles))
- Update cosmos-sdk version to 0.26.1 [\#80](https://github.com/cybercongress/cyberd/pull/80) ([arturalbov](https://github.com/arturalbov))

## [v0.0.7](https://github.com/cybercongress/cyberd/tree/v0.0.7) (2018-10-25)
[Full Changelog](https://github.com/cybercongress/cyberd/compare/v0.0.6...v0.0.7)

**Implemented enhancements:**

- Proxy service. Add search pagination [\#69](https://github.com/cybercongress/cyberd/issues/69)
- Add send tokens endpoint to proxy [\#62](https://github.com/cybercongress/cyberd/issues/62)

**Fixed bugs:**

- Proxy service. Search request with "spaces" fails [\#67](https://github.com/cybercongress/cyberd/issues/67)
- Non-deterministic rank calculation [\#66](https://github.com/cybercongress/cyberd/issues/66)

**Closed issues:**

- Claim service: increment tx sequence manually [\#64](https://github.com/cybercongress/cyberd/issues/64)

**Merged pull requests:**

- \#71 Create basic wiki cyberd indexer [\#72](https://github.com/cybercongress/cyberd/pull/72) ([hleb-albau](https://github.com/hleb-albau))
- Proxy service. Add search pagination [\#70](https://github.com/cybercongress/cyberd/pull/70) ([arturalbov](https://github.com/arturalbov))
- Proxy service. Search request with spaces fails [\#68](https://github.com/cybercongress/cyberd/pull/68) ([arturalbov](https://github.com/arturalbov))
- Claim service: increment tx sequence manually [\#65](https://github.com/cybercongress/cyberd/pull/65) ([arturalbov](https://github.com/arturalbov))

## [v0.0.6](https://github.com/cybercongress/cyberd/tree/v0.0.6) (2018-10-24)
[Full Changelog](https://github.com/cybercongress/cyberd/compare/v0.0.5...v0.0.6)

## [v0.0.5](https://github.com/cybercongress/cyberd/tree/v0.0.5) (2018-10-23)
[Full Changelog](https://github.com/cybercongress/cyberd/compare/v0.0.4...v0.0.5)

**Implemented enhancements:**

- \[RPC\] Fix small finding  [\#57](https://github.com/cybercongress/cyberd/issues/57)

## [v0.0.4](https://github.com/cybercongress/cyberd/tree/v0.0.4) (2018-10-23)
[Full Changelog](https://github.com/cybercongress/cyberd/compare/v0.0.3...v0.0.4)

**Closed issues:**

- Service to claim cyberd zeronet tokens [\#61](https://github.com/cybercongress/cyberd/issues/61)
- Update cosmos-sdk to latest dev branch version  [\#56](https://github.com/cybercongress/cyberd/issues/56)
- Perfomance Degradation: Heavy Disk Usage [\#50](https://github.com/cybercongress/cyberd/issues/50)

**Merged pull requests:**

- Add send tokens endpoint to proxy [\#63](https://github.com/cybercongress/cyberd/pull/63) ([arturalbov](https://github.com/arturalbov))
- Claim service [\#60](https://github.com/cybercongress/cyberd/pull/60) ([arturalbov](https://github.com/arturalbov))
- \#57 \[RPC\] Fix small finding [\#58](https://github.com/cybercongress/cyberd/pull/58) ([hleb-albau](https://github.com/hleb-albau))
- Update cosmos-sdk to latest dev branch version [\#53](https://github.com/cybercongress/cyberd/pull/53) ([hleb-albau](https://github.com/hleb-albau))

## [v0.0.3](https://github.com/cybercongress/cyberd/tree/v0.0.3) (2018-10-19)
[Full Changelog](https://github.com/cybercongress/cyberd/compare/v0.0.2...v0.0.3)

**Implemented enhancements:**

- CLI: Add "wait\_for\_confirmation" Flag [\#47](https://github.com/cybercongress/cyberd/issues/47)

**Fixed bugs:**

- RPC Client: /search on non existing cid return first added cid [\#48](https://github.com/cybercongress/cyberd/issues/48)

**Closed issues:**

- Write validation logic of IPFS hash for Losion Zeronet [\#18](https://github.com/cybercongress/cyberd/issues/18)
- Genesis Zeronet [\#17](https://github.com/cybercongress/cyberd/issues/17)
- Performance testing of Zeronet [\#4](https://github.com/cybercongress/cyberd/issues/4)

**Merged pull requests:**

- Add cyberdproxy process to docker container. Add status endpoint [\#55](https://github.com/cybercongress/cyberd/pull/55) ([arturalbov](https://github.com/arturalbov))
- Proxy rpc [\#54](https://github.com/cybercongress/cyberd/pull/54) ([arturalbov](https://github.com/arturalbov))

## [v0.0.2](https://github.com/cybercongress/cyberd/tree/v0.0.2) (2018-10-05)
[Full Changelog](https://github.com/cybercongress/cyberd/compare/v0.0.1...v0.0.2)

**Implemented enhancements:**

- Cosmos PoC: Integrate Rank Calculation [\#43](https://github.com/cybercongress/cyberd/issues/43)
- Cosmos PoC: Extenend Standart Tendermint RPC API  [\#42](https://github.com/cybercongress/cyberd/issues/42)

**Closed issues:**

- Implement persistent storage for links. [\#40](https://github.com/cybercongress/cyberd/issues/40)
- run extra node for cyberd [\#20](https://github.com/cybercongress/cyberd/issues/20)

**Merged pull requests:**

- Fix search on non existing cid [\#49](https://github.com/cybercongress/cyberd/pull/49) ([arturalbov](https://github.com/arturalbov))
- Add Circle CI build job [\#46](https://github.com/cybercongress/cyberd/pull/46) ([arturalbov](https://github.com/arturalbov))
- Extenend Standart Tendermint RPC API [\#45](https://github.com/cybercongress/cyberd/pull/45) ([arturalbov](https://github.com/arturalbov))
- \#43 Simplest Rank [\#44](https://github.com/cybercongress/cyberd/pull/44) ([hleb-albau](https://github.com/hleb-albau))
- \#40 introduce in-memory store [\#41](https://github.com/cybercongress/cyberd/pull/41) ([hleb-albau](https://github.com/hleb-albau))
- \#37 redesign db, app refactor [\#38](https://github.com/cybercongress/cyberd/pull/38) ([hleb-albau](https://github.com/hleb-albau))
- Cosmos POC: Clean up CLI [\#35](https://github.com/cybercongress/cyberd/pull/35) ([arturalbov](https://github.com/arturalbov))
- Update cyberd/cosmos README [\#34](https://github.com/cybercongress/cyberd/pull/34) ([arturalbov](https://github.com/arturalbov))

## [v0.0.1](https://github.com/cybercongress/cyberd/tree/v0.0.1) (2018-09-25)
**Closed issues:**

- Make cyberd docs to be included into common wiki. [\#32](https://github.com/cybercongress/cyberd/issues/32)
- Write LT/NLT logic to Losion Zeronet [\#19](https://github.com/cybercongress/cyberd/issues/19)
- Research basic technologies [\#14](https://github.com/cybercongress/cyberd/issues/14)
- Research basic papers [\#13](https://github.com/cybercongress/cyberd/issues/13)
- Perfomance testing of Ethermint [\#12](https://github.com/cybercongress/cyberd/issues/12)
- Perfomance testing of Plasma [\#11](https://github.com/cybercongress/cyberd/issues/11)
- Perfomance testing of PoA networks [\#7](https://github.com/cybercongress/cyberd/issues/7)
- Cleanup paper [\#6](https://github.com/cybercongress/cyberd/issues/6)
-  Perfomance testing of EOS [\#3](https://github.com/cybercongress/cyberd/issues/3)

**Merged pull requests:**

- Genesis zeronet: Cosmos SDK [\#33](https://github.com/cybercongress/cyberd/pull/33) ([arturalbov](https://github.com/arturalbov))
- \[WIP\] 21 calculate spring rank for ethereum [\#23](https://github.com/cybercongress/cyberd/pull/23) ([hleb-albau](https://github.com/hleb-albau))
- Fixing typo [\#2](https://github.com/cybercongress/cyberd/pull/2) ([trummax](https://github.com/trummax))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*