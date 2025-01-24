# Changelog

## [2.3.2](https://github.com/devopsarr/terraform-provider-radarr/compare/v2.3.1...v2.3.2) (2025-01-24)


### Bug Fixes

* **deps:** update module github.com/devopsarr/radarr-go to v1.2.0 ([a4bcd3e](https://github.com/devopsarr/terraform-provider-radarr/commit/a4bcd3ebec2cffeca6faceae8fe485aa2e3d7863))

## [2.3.1](https://github.com/devopsarr/terraform-provider-radarr/compare/v2.3.0...v2.3.1) (2024-12-10)


### Bug Fixes

* **#273:** allow smart colon replacement format ([c6dfb80](https://github.com/devopsarr/terraform-provider-radarr/commit/c6dfb80a9a75847f0bdf879d238e879b23cee65c))
* **deps:** update module github.com/devopsarr/radarr-go to v1.1.2 ([24d4176](https://github.com/devopsarr/terraform-provider-radarr/commit/24d41760cfc38f1c993c18c4a5c83eabe9fc6e73))
* **deps:** update module github.com/stretchr/testify to v1.10.0 ([9c510dd](https://github.com/devopsarr/terraform-provider-radarr/commit/9c510dd09628f1c7e1c3e8bf9759d267b75c1e81))

## [2.3.0](https://github.com/devopsarr/terraform-provider-radarr/compare/v2.2.0...v2.3.0) (2024-09-22)


### Features

* add host log_size_limit attribute ([899ff03](https://github.com/devopsarr/terraform-provider-radarr/commit/899ff0327089e85d281932d174a8a7b2db1cb856))
* add quality profile min_upgrade_format_score field ([b8ce1e7](https://github.com/devopsarr/terraform-provider-radarr/commit/b8ce1e7cac9e2c756063ee632ad2037d6f1095b8))
* update to go 1.23.1 ([d1996ab](https://github.com/devopsarr/terraform-provider-radarr/commit/d1996abb66663e33efca4956068dcd5f0f0865c0))


### Bug Fixes

* bump golangci-lint version ([cffe831](https://github.com/devopsarr/terraform-provider-radarr/commit/cffe83177e7470597eb4f37150641bdd2dbc21d0))
* correct goreleaser syntax ([42f3100](https://github.com/devopsarr/terraform-provider-radarr/commit/42f310022ea168ef5bc926c7a165c0f44d0138cf))
* **deps:** update terraform-framework ([c056e7b](https://github.com/devopsarr/terraform-provider-radarr/commit/c056e7b6b247ec374ea705eccae2925d8be1d929))
* **deps:** update terraform-framework ([5835bd0](https://github.com/devopsarr/terraform-provider-radarr/commit/5835bd0c4501a890e7b357ee97e51756c4b21089))

## [2.2.0](https://github.com/devopsarr/terraform-provider-radarr/compare/v2.1.0...v2.2.0) (2024-03-10)


### Features

* add fields to use client extra headers ([626afbf](https://github.com/devopsarr/terraform-provider-radarr/commit/626afbf2719a38c45355da6d0dd93ad881f9f89a))
* move to context based authentication ([1da831a](https://github.com/devopsarr/terraform-provider-radarr/commit/1da831a3d2fddccc33cc1cf829f5834479e2ca57))
* remove deprecated notification boxcar ([099a782](https://github.com/devopsarr/terraform-provider-radarr/commit/099a782d81d52046c86ab6652e4f51221384607f))
* update go version to 1.21 ([382bfad](https://github.com/devopsarr/terraform-provider-radarr/commit/382bfad4669cda1dc00fec15731e756a7f21bbdf))


### Bug Fixes

* email notification encryption field ([70c768c](https://github.com/devopsarr/terraform-provider-radarr/commit/70c768cb5799e9231783d22e494b809edd0f2f5c))
* update host authentication fields ([8d85262](https://github.com/devopsarr/terraform-provider-radarr/commit/8d852626aa28ce134a6cb4022ac85a0b38cd32a5))

## [2.1.0](https://github.com/devopsarr/terraform-provider-radarr/compare/v2.0.1...v2.1.0) (2023-10-14)


### Features

* **#222:** manage v5 sensitive values ([a799aa4](https://github.com/devopsarr/terraform-provider-radarr/commit/a799aa44d50ef2897d0a889fa16cbd3288fdf115))


### Bug Fixes

* quality profile use all formats and ordered quality groups ([4470f58](https://github.com/devopsarr/terraform-provider-radarr/commit/4470f587541d926af1f716a1489b6bc8aa621a14))

## [2.0.1](https://github.com/devopsarr/terraform-provider-radarr/compare/v2.0.0...v2.0.1) (2023-08-24)


### Bug Fixes

* manage encrypted host password ([fb8387b](https://github.com/devopsarr/terraform-provider-radarr/commit/fb8387b655ca8540e836c81292063ecccfc5313a))

## [2.0.0](https://github.com/devopsarr/terraform-provider-radarr/compare/v1.8.0...v2.0.0) (2023-08-18)


### ⚠ BREAKING CHANGES

* remove deprecated restriction
* remove deprecated import list sync_interval
* remove obsolete indexer omgwtfnzbs

### Features

* **#203:** add host data source ([be3a56a](https://github.com/devopsarr/terraform-provider-radarr/commit/be3a56acc1449d7ce43641a7248485becea877dc))
* **#203:** add host resource ([4d7aa59](https://github.com/devopsarr/terraform-provider-radarr/commit/4d7aa5910f7c0249860679911b48eb6dbef3df94))
* add apprise notification ([d84faec](https://github.com/devopsarr/terraform-provider-radarr/commit/d84faec7e8f07ce4292cfcec348fa5ec57b8848c))
* add auto tag conditions data source ([c58bb25](https://github.com/devopsarr/terraform-provider-radarr/commit/c58bb25c5f1a5b6a0f546a4364656e08bfff2d15))
* add auto tag data source ([6bd4726](https://github.com/devopsarr/terraform-provider-radarr/commit/6bd4726f13bba9b1178be300ef127784d4f003e1))
* add auto tag resource ([f689119](https://github.com/devopsarr/terraform-provider-radarr/commit/f68911907ae828c5897945a6e8abc13d063efb94))
* add auto tags data source ([d427bfe](https://github.com/devopsarr/terraform-provider-radarr/commit/d427bfef6e5fde654c9cae24b738eaed7d0ebfaf))
* add ireland certification country ([84785b6](https://github.com/devopsarr/terraform-provider-radarr/commit/84785b6b60b4a15a33e88e7157fc65c16326c93b))
* add support for new notifications flags ([777c6a6](https://github.com/devopsarr/terraform-provider-radarr/commit/777c6a6ecb2d5c84c5342b7e1444c191ae640d4f))
* add telegram topic id support ([842567f](https://github.com/devopsarr/terraform-provider-radarr/commit/842567f394fc5cdb3974577cdc6d81aa2471709d))
* improve diagnostics part 1 ([78a628f](https://github.com/devopsarr/terraform-provider-radarr/commit/78a628f1bc3e765f483a382ddb9197ec5d408cfa))
* improve diagnostics part 2 ([0e97afe](https://github.com/devopsarr/terraform-provider-radarr/commit/0e97afe2ac5889515db24226978f79d5e48db0ba))
* remove deprecated import list sync_interval ([17b4490](https://github.com/devopsarr/terraform-provider-radarr/commit/17b4490b33aa83c1638cf82c86245536cab8cf6e))
* remove deprecated restriction ([37b012c](https://github.com/devopsarr/terraform-provider-radarr/commit/37b012c911420de587db1c9ba626e43fb08be02e))
* remove obsolete indexer omgwtfnzbs ([3aa2cf9](https://github.com/devopsarr/terraform-provider-radarr/commit/3aa2cf9a008e7327684c628b4d26d99536f6ad7c))
* remove obsolete indexer rarbg ([e49da69](https://github.com/devopsarr/terraform-provider-radarr/commit/e49da6987bdd514b65ec01ac7032d5635fb3f943))
* use only ID for delete ([6cb897d](https://github.com/devopsarr/terraform-provider-radarr/commit/6cb897df0f1052ffa9511b001b4df316abbb16b5))


### Bug Fixes

* delete error message ([85ccd1c](https://github.com/devopsarr/terraform-provider-radarr/commit/85ccd1c606163222636fe8b1232f765d7e2df555))
* non empty state for nil fields ([05a969e](https://github.com/devopsarr/terraform-provider-radarr/commit/05a969eb376ca4f8825de39a77faeb8d359c9ab9))

## [1.8.0](https://github.com/devopsarr/terraform-provider-radarr/compare/v1.7.0...v1.8.0) (2023-02-01)


### Features

* add custom format condition data source ([a64393b](https://github.com/devopsarr/terraform-provider-radarr/commit/a64393b98eaaf22f840c9fa792acc2f506eef353))
* add custom format condition edition data source ([d15dfad](https://github.com/devopsarr/terraform-provider-radarr/commit/d15dfadcc7451fda01946554962bf98f934884b6))
* add custom format condition indexer flag data source ([7dc14df](https://github.com/devopsarr/terraform-provider-radarr/commit/7dc14dfdb5fb3452de79fce7d7b5fbd2faa6d0ac))
* add custom format condition language data source ([c48a11f](https://github.com/devopsarr/terraform-provider-radarr/commit/c48a11f5b084528cc6fd8c9da33754c2b075b7f7))
* add custom format condition quality modifier data source ([b863fb1](https://github.com/devopsarr/terraform-provider-radarr/commit/b863fb19ead54a32e14a850186259b317e18c746))
* add custom format condition release group data source ([3fecbd8](https://github.com/devopsarr/terraform-provider-radarr/commit/3fecbd8c9c22d298ad98173948d457a33e06d600))
* add custom format condition release title data source ([1a5eb59](https://github.com/devopsarr/terraform-provider-radarr/commit/1a5eb59786824849defa79ca0a0813cbcc4033e3))
* add custom format condition resolution data source ([8eddda3](https://github.com/devopsarr/terraform-provider-radarr/commit/8eddda383bf99fd4b27d9f249e66b0b39e1612c5))
* add custom format condition size data source ([7672dfd](https://github.com/devopsarr/terraform-provider-radarr/commit/7672dfdbf22aeced76141df71b4e197860f2d7fc))
* add custom format condition source data source ([5f33092](https://github.com/devopsarr/terraform-provider-radarr/commit/5f33092e153bc1739f0ac6ba36863638e1e4dc14))
* add import list config data source ([31134a7](https://github.com/devopsarr/terraform-provider-radarr/commit/31134a7c73c6c24c3d58b2509cf06b6033445bfe))
* add import list config resource ([fd4d71f](https://github.com/devopsarr/terraform-provider-radarr/commit/fd4d71f5a014ca20725ffa432996f72abcc51117))
* add import list couch potato resource ([f777653](https://github.com/devopsarr/terraform-provider-radarr/commit/f7776534d6303ed000897b29c7119c418d484e85))
* add import list custom resource ([27c4eb9](https://github.com/devopsarr/terraform-provider-radarr/commit/27c4eb997c1ccc11d77e67d8e467c257088018c9))
* add import list data source ([3cc3cf1](https://github.com/devopsarr/terraform-provider-radarr/commit/3cc3cf1a8d2390b0010d1c899f34367c5ff1241f))
* add import list exclusion data source ([09157b0](https://github.com/devopsarr/terraform-provider-radarr/commit/09157b0dc8254644517dc7a66231aa77d9f9ef25))
* add import list exclusion resource ([9b5c0c9](https://github.com/devopsarr/terraform-provider-radarr/commit/9b5c0c9224e9b8e918b273d47a7abc6b1009cc2b))
* add import list exclusions data source ([b525e09](https://github.com/devopsarr/terraform-provider-radarr/commit/b525e090901ff5539cf7544ac0215381653bc230))
* add import list imdb resource ([8ac9c1b](https://github.com/devopsarr/terraform-provider-radarr/commit/8ac9c1b6a59ffb952c18e47e088348362eceba73))
* add import list plex resource ([f53606e](https://github.com/devopsarr/terraform-provider-radarr/commit/f53606e947412587aefab5e24161788d5357bd5a))
* add import list radarr resource ([c138d40](https://github.com/devopsarr/terraform-provider-radarr/commit/c138d40f8ddbbd5f170d696feaf4d0a4ad534907))
* add import list resource ([a25d60e](https://github.com/devopsarr/terraform-provider-radarr/commit/a25d60e7acac6bfca35bf1b5b701c8b9878fe9b7))
* add import list rss resource ([4547336](https://github.com/devopsarr/terraform-provider-radarr/commit/45473362f3e6d48f22cb36f32d31b6579d5df557))
* add import list stevenlu resource ([36f0ad0](https://github.com/devopsarr/terraform-provider-radarr/commit/36f0ad087c3a2709d529c3d42d40e0024d688cb7))
* add import list stevenlu2 resource ([ec71afa](https://github.com/devopsarr/terraform-provider-radarr/commit/ec71afa8d699dc9939e488e1c8178d086b6ad43b))
* add import list tmdb company resource ([f2005df](https://github.com/devopsarr/terraform-provider-radarr/commit/f2005dfccb7ab29c2274d4c942883b405db6b094))
* add import list tmdb keyword resource ([0383fef](https://github.com/devopsarr/terraform-provider-radarr/commit/0383fef874b173383297e3672cdfe9c65f5acd69))
* add import list tmdb list resource ([2393836](https://github.com/devopsarr/terraform-provider-radarr/commit/2393836e18b5e633bf3b61b43ced6ee3f3ae4e06))
* add import list tmdb person resource ([c6b4cff](https://github.com/devopsarr/terraform-provider-radarr/commit/c6b4cfffa009a26c9a0cd009d680d862bdea0074))
* add import list tmdb popular resource ([8852e45](https://github.com/devopsarr/terraform-provider-radarr/commit/8852e45caacdbeb02f3fd5f2a0537ba5e9da6bc4))
* add import list tmdb user resource ([0b62c81](https://github.com/devopsarr/terraform-provider-radarr/commit/0b62c81e30e38100f2d50ea71e7df57b6c989bf5))
* add import list trakt list resource ([7effe92](https://github.com/devopsarr/terraform-provider-radarr/commit/7effe922b471801608d3a6386779751ad5dab647))
* add import list trakt popular resource ([a64e56b](https://github.com/devopsarr/terraform-provider-radarr/commit/a64e56b47c99e6a400a803c9e0157b13a0e017ae))
* add import list trakt user resource ([fba8af7](https://github.com/devopsarr/terraform-provider-radarr/commit/fba8af756f2e380df4e225b5104384df0f10aa6a))
* add import lists data source ([67ed715](https://github.com/devopsarr/terraform-provider-radarr/commit/67ed715972bf737fc487a6f4b0aa7733a8258af5))
* add metadata config data source ([3711ba3](https://github.com/devopsarr/terraform-provider-radarr/commit/3711ba32ff1421a19c8a8f1d8cd1f4d0ee27d59b))
* add metadata config resource ([ab8a9f8](https://github.com/devopsarr/terraform-provider-radarr/commit/ab8a9f8fe3bd0578f98ab59865681c6fb4e4088f))
* add metadata consumers datasource ([6605bc9](https://github.com/devopsarr/terraform-provider-radarr/commit/6605bc9510876fb38e4bf257d0e49028cde450e2))
* add metadata datasource ([d6d701b](https://github.com/devopsarr/terraform-provider-radarr/commit/d6d701bea1ac8fa8f3784a88d9c230eab764c534))
* add metadata emby resource ([c2f0927](https://github.com/devopsarr/terraform-provider-radarr/commit/c2f0927ba1d0b4d240bbde88ef1710a8554cf4b7))
* add metadata kodi resource ([3dc6ceb](https://github.com/devopsarr/terraform-provider-radarr/commit/3dc6cebadc62107e77caa82fcc27bf7ac34df493))
* add metadata resource ([146745d](https://github.com/devopsarr/terraform-provider-radarr/commit/146745d8340b1df0da4a85eebc9e48a1adc17b2e))
* add metadata roksbox resource ([b0c06a0](https://github.com/devopsarr/terraform-provider-radarr/commit/b0c06a03a770d2a1901d384d534a90c51d8e8079))
* add metadata wdtv resource ([a83c1a7](https://github.com/devopsarr/terraform-provider-radarr/commit/a83c1a769f5488aabc5e36869536c00ea9ee18b6))
* add quality data source ([766c53d](https://github.com/devopsarr/terraform-provider-radarr/commit/766c53d9ba4a321216910880056bb4014129c8dc))
* add quality definition data source ([5a31c83](https://github.com/devopsarr/terraform-provider-radarr/commit/5a31c838320f3830506c0c6066a02f61cc5a2308))
* add quality definition resource ([5094681](https://github.com/devopsarr/terraform-provider-radarr/commit/5094681c88905ea7d7dbcdd78a678a499efa3bb4))
* add quality definitions data source ([846ef00](https://github.com/devopsarr/terraform-provider-radarr/commit/846ef0070543d3dc7ac71ec3c626b71d4593a06b))
* improve slice init ([4a20f86](https://github.com/devopsarr/terraform-provider-radarr/commit/4a20f867976a6849b053724ab4822facff5791d9))


### Bug Fixes

* indexer seedcriteria fields not written ([d74a6f0](https://github.com/devopsarr/terraform-provider-radarr/commit/d74a6f0690f073f3f29e9877660aff4297168085))
* quality profile with single item group ([a2a1137](https://github.com/devopsarr/terraform-provider-radarr/commit/a2a1137b5419a82f144de7fefe70600baed3eee4))
* read data source from request ([fc36939](https://github.com/devopsarr/terraform-provider-radarr/commit/fc36939d46bcae2f5e0781c30087a20a93a75796))
* use get function for sdk fields ([478d6c8](https://github.com/devopsarr/terraform-provider-radarr/commit/478d6c8d466f6bf5c21b39b4d9051a2f60c3dbb3))

## [1.7.0](https://github.com/devopsarr/terraform-provider-radarr/compare/v1.6.0...v1.7.0) (2023-01-23)


### Features

* add freebox download client ([9bd8cc4](https://github.com/devopsarr/terraform-provider-radarr/commit/9bd8cc47790ebcf51c1c6ef3da1554ff610965c4))
* add language data source ([a9b75d6](https://github.com/devopsarr/terraform-provider-radarr/commit/a9b75d6383be2d1c0273c837427be47bff158c02))
* add languages data source ([94da3d4](https://github.com/devopsarr/terraform-provider-radarr/commit/94da3d413dee96195c11c855e4baddd7971e6670))
* add movie datasource ([381c551](https://github.com/devopsarr/terraform-provider-radarr/commit/381c551e82d88ed9c41f0af5d373221f69406924))
* add movie resource ([04b1f95](https://github.com/devopsarr/terraform-provider-radarr/commit/04b1f957237f9d131a2e6b7073cf3b6d05849dcf))
* add movies datasource ([e36ead9](https://github.com/devopsarr/terraform-provider-radarr/commit/e36ead9fb45e0efcc8033ca658f5c35813ccece8))
* make notification flags optional ([419b9c6](https://github.com/devopsarr/terraform-provider-radarr/commit/419b9c6638aa217d5c888a1cbda254cef0c14906))


### Bug Fixes

* quality profile custom format write ([fbf263c](https://github.com/devopsarr/terraform-provider-radarr/commit/fbf263c9fdbf99f0524acb8b2ddc652a9565a8a6))
* remove unused parameter form notifiarr ([2db4824](https://github.com/devopsarr/terraform-provider-radarr/commit/2db482417ebf8514065a28857bf30a72482255dd))
* update sdk method naming ([0f48fc2](https://github.com/devopsarr/terraform-provider-radarr/commit/0f48fc2c2209cead9309c5eb799eb77a7f3588cd))

## [1.6.0](https://github.com/devopsarr/terraform-provider-radarr/compare/v1.5.0...v1.6.0) (2022-12-09)


### Features

* add download client aria2 resource ([7cb733c](https://github.com/devopsarr/terraform-provider-radarr/commit/7cb733c10a53d9f586a5407607b683d8069eeaff))
* add download client blackhole resource ([bb19984](https://github.com/devopsarr/terraform-provider-radarr/commit/bb199842784f4419999eb7ffb05d2153a8450a6e))
* add download client deluge resource ([200d7ab](https://github.com/devopsarr/terraform-provider-radarr/commit/200d7abe024e62562a26eacab73e4654790111f8))
* add download client flood resource ([593bad5](https://github.com/devopsarr/terraform-provider-radarr/commit/593bad555f48cc4c2a07c25a990bafdb15275cda))
* add download client hadouken resource ([5a45185](https://github.com/devopsarr/terraform-provider-radarr/commit/5a451858788d4aafabb4466ba41817a3b26f6c9f))
* add download client nzbget resource ([5175c1a](https://github.com/devopsarr/terraform-provider-radarr/commit/5175c1a4beb0a07c26fa85ee6b612d9e6a96f646))
* add download client nzbvortex resource ([0cb8e03](https://github.com/devopsarr/terraform-provider-radarr/commit/0cb8e031d128b1bf932cf217d6546b3dff6dbdad))
* add download client pneumatic resource ([0cdca61](https://github.com/devopsarr/terraform-provider-radarr/commit/0cdca617200e556d2b92e5a1e97b0b21f900e9ea))
* add download client qbittorrent resource ([926f368](https://github.com/devopsarr/terraform-provider-radarr/commit/926f368ecef5feb907b74f6b1a6f2ff2c3205696))
* add download client rtorrent resource ([750a46f](https://github.com/devopsarr/terraform-provider-radarr/commit/750a46f43f04a9f8d18e1a75b51d1e5d122790a8))
* add download client sabnzbd resource ([e18d1b0](https://github.com/devopsarr/terraform-provider-radarr/commit/e18d1b0ecd04fc56de2544715691aef9105b6560))
* add download client torrent download station resource ([05d8a00](https://github.com/devopsarr/terraform-provider-radarr/commit/05d8a00352aa163ddbad09d2cb35f545003874f6))
* add download client usenet blackhole resource ([0692347](https://github.com/devopsarr/terraform-provider-radarr/commit/0692347cfd6a48274e9420b361ac01a4f92edd57))
* add download client usenet download station resource ([73176e3](https://github.com/devopsarr/terraform-provider-radarr/commit/73176e3d461b14c38063e3271038fac9eae66cab))
* add download client utorrent resource ([be1a82c](https://github.com/devopsarr/terraform-provider-radarr/commit/be1a82c9d05658df337b42ea41255c0dd6d5756e))
* add download client vuze resource ([84ca9dd](https://github.com/devopsarr/terraform-provider-radarr/commit/84ca9ddc7c5f5f3c09ce36c380f1c78fcd36f1c5))
* add indexer filelist resource ([0dccd77](https://github.com/devopsarr/terraform-provider-radarr/commit/0dccd7793b90e4139d3542cab2b607b45aab41e7))
* add indexer hdbits resource ([8337ec2](https://github.com/devopsarr/terraform-provider-radarr/commit/8337ec23cdc97065d9b262fadfd820c4c4a51465))
* add indexer iptorrents resource ([05cc40a](https://github.com/devopsarr/terraform-provider-radarr/commit/05cc40a7942e26e8452d449e182c1b80aaa7a870))
* add indexer nyaa resource ([4c83f9b](https://github.com/devopsarr/terraform-provider-radarr/commit/4c83f9b3f25fe87af8e984ec715f14ad4df96359))
* add indexer omgwtfnzbs resource ([fc7dff4](https://github.com/devopsarr/terraform-provider-radarr/commit/fc7dff4d0820742c88c5a3789eec623ccd868389))
* add indexer pass the popcorn resource ([89d2d0c](https://github.com/devopsarr/terraform-provider-radarr/commit/89d2d0ca4178a6aa3fad0763597f115ff37c76ee))
* add indexer torrent potato resource ([07dafb8](https://github.com/devopsarr/terraform-provider-radarr/commit/07dafb810bae427a1bbcdbe3947f4544d8616571))
* add indexer torrent rss resource ([a948496](https://github.com/devopsarr/terraform-provider-radarr/commit/a948496c624fa67d27a123cb80df64b4a021d770))
* add indexer torznab resource ([cf886c6](https://github.com/devopsarr/terraform-provider-radarr/commit/cf886c6039f26c38ffb1722e8314e91b24517a01))
* add notification boxcar resource ([0cd4b10](https://github.com/devopsarr/terraform-provider-radarr/commit/0cd4b10efc0f965229cb890430411809db1bf741))
* add notification discord resource ([767afa5](https://github.com/devopsarr/terraform-provider-radarr/commit/767afa5afcb14d79674b7d9fe670f5ae28773927))
* add notification email resource ([174877c](https://github.com/devopsarr/terraform-provider-radarr/commit/174877caad7506a1d18d0e2dd479096b21a468f5))
* add notification emby resource ([cec8aae](https://github.com/devopsarr/terraform-provider-radarr/commit/cec8aaeea483528120b0f338d740322c6406b432))
* add notification gotify resource ([30cd267](https://github.com/devopsarr/terraform-provider-radarr/commit/30cd267ee0ea23872ec14f1d1fef8d6b1331cb3d))
* add notification join resource ([56062d3](https://github.com/devopsarr/terraform-provider-radarr/commit/56062d3a198d09126c9983c345c76edcaaad105b))
* add notification kodi resource ([c7fcef3](https://github.com/devopsarr/terraform-provider-radarr/commit/c7fcef3b8b171eef2a956169ec1e8c6124c7483d))
* add notification mailgun resource ([62090f7](https://github.com/devopsarr/terraform-provider-radarr/commit/62090f7a3b7c916ca11c7b65f9da16c6ca872625))
* add notification notifiarr resource ([ed4365a](https://github.com/devopsarr/terraform-provider-radarr/commit/ed4365add560657d025c95f81fcbfef4a3184d78))
* add notification ntfy resource ([7cf6f31](https://github.com/devopsarr/terraform-provider-radarr/commit/7cf6f3100116c84d6e2f10f0b02fd485cf84b46d))
* add notification plex resource ([ded48a7](https://github.com/devopsarr/terraform-provider-radarr/commit/ded48a7fd44c58c386ef3e4016407578f226973e))
* add notification prowl resource ([769393c](https://github.com/devopsarr/terraform-provider-radarr/commit/769393c7883502dd3412190a06b26e44a4658274))
* add notification pushbullet resource ([13cbe21](https://github.com/devopsarr/terraform-provider-radarr/commit/13cbe213e498074427e3a0c84580b62a73f1b41d))
* add notification pushover resource ([384922b](https://github.com/devopsarr/terraform-provider-radarr/commit/384922b2b11ec23efc7b7fc30f91e81881a7f0d6))
* add notification sendgrid resource ([2d9885d](https://github.com/devopsarr/terraform-provider-radarr/commit/2d9885d0db0bcdfc714dd445e12d967b17b2107d))
* add notification simplepush resource ([fc2ca35](https://github.com/devopsarr/terraform-provider-radarr/commit/fc2ca35d5844bd2c7daa2c467011403df93f7cd4))
* add notification slack resource ([fc95196](https://github.com/devopsarr/terraform-provider-radarr/commit/fc95196d3c03271501a1ca2c71f83b3860cb9aa0))
* add notification synology resource ([c811240](https://github.com/devopsarr/terraform-provider-radarr/commit/c8112407e0c994873a347a488e12e70a81ad846f))
* add notification telegram resource ([617caef](https://github.com/devopsarr/terraform-provider-radarr/commit/617caef193c86c81ae4c0514bfeb8c5bbf121e7d))
* add notification trakt resource ([40a7417](https://github.com/devopsarr/terraform-provider-radarr/commit/40a7417c0fbbde388de3e31ec94d7a72286651b3))
* add notification twitter resource ([fc972ac](https://github.com/devopsarr/terraform-provider-radarr/commit/fc972ac05c8c82ec894542c7064f33def6cdaec4))


### Bug Fixes

* download client fields were sonarr related ([d12be45](https://github.com/devopsarr/terraform-provider-radarr/commit/d12be45b1f555d8db73f3d770bf044ac7d1a3821))
* post_im_tags and watch_floder fields ([b56db33](https://github.com/devopsarr/terraform-provider-radarr/commit/b56db335a6b26c40bae3e5a7cc08a1bb68eb7140))

## [1.5.0](https://github.com/devopsarr/terraform-provider-radarr/compare/v1.4.0...v1.5.0) (2022-11-23)


### Features

* add custom format datasource ([b969064](https://github.com/devopsarr/terraform-provider-radarr/commit/b96906463f52eb54b0af7b44f875d947d776c4c9))
* add custom format resource ([f4b7fa0](https://github.com/devopsarr/terraform-provider-radarr/commit/f4b7fa032300e90a5f4be58592d1c6f017310869))
* add custom formats datasource ([908bc1b](https://github.com/devopsarr/terraform-provider-radarr/commit/908bc1b63fe07f62bea71878e9a95aa2f77e1b47))
* add quality profile datasource ([c0db28c](https://github.com/devopsarr/terraform-provider-radarr/commit/c0db28c2574b49a47d989bbde70e3423c66df298))
* add quality profile resource ([8a7e8bb](https://github.com/devopsarr/terraform-provider-radarr/commit/8a7e8bba2f820507cf37dd4ddbd3508aa207b4bb))
* add quality profiles datasource ([c176978](https://github.com/devopsarr/terraform-provider-radarr/commit/c176978c63fbae6fe3b5009c1e0ce60f21a2f313))

## [1.4.0](https://github.com/devopsarr/terraform-provider-radarr/compare/v1.3.0...v1.4.0) (2022-11-16)


### Features

* add download client config datasource ([72f88de](https://github.com/devopsarr/terraform-provider-radarr/commit/72f88de68d91d128e5fca365161f17dcf675317f))
* add download client config resource ([891e8d9](https://github.com/devopsarr/terraform-provider-radarr/commit/891e8d9b01f9b313f5666bbe350c7e006cedcadb))
* add download client datasource ([e3e6b4d](https://github.com/devopsarr/terraform-provider-radarr/commit/e3e6b4d03c416b0570e3b356bd4c375af97ad788))
* add download client resource ([07b177f](https://github.com/devopsarr/terraform-provider-radarr/commit/07b177fc4af796e3dd563cfb93e734de04ea6f8e))
* add download client transmission resource ([8952e56](https://github.com/devopsarr/terraform-provider-radarr/commit/8952e569cd2b6b51749ef595d7cc7d06f53167b4))
* add download clients datasource ([9205c6c](https://github.com/devopsarr/terraform-provider-radarr/commit/9205c6ce449f47e9646ee64c201a39f02a2b061c))
* add indexer config datasource ([3f78e87](https://github.com/devopsarr/terraform-provider-radarr/commit/3f78e878f6db144562f42ff144c411f973fc09bb))
* add indexer config resource ([e55a128](https://github.com/devopsarr/terraform-provider-radarr/commit/e55a1282fc11da8e8f304cdaa024ab1584585491))
* add indexer datasource ([b6db72e](https://github.com/devopsarr/terraform-provider-radarr/commit/b6db72ef4776d460345c1781b175a0d6c30097fe))
* add indexer newznab resource ([b5e20e6](https://github.com/devopsarr/terraform-provider-radarr/commit/b5e20e6d0415dcc87ece3a0a395d154b0e5173ce))
* add indexer rargb resource ([44f27df](https://github.com/devopsarr/terraform-provider-radarr/commit/44f27df6cbe1f29714dc76c840432d1aa0059977))
* add indexer resource ([3d260df](https://github.com/devopsarr/terraform-provider-radarr/commit/3d260dff37bd2e03b588cbbba7308340e60ea26a))
* add indexers datasource ([998f2dd](https://github.com/devopsarr/terraform-provider-radarr/commit/998f2dd523f801f0e42911bba98c6f1b0e49369c))
* add remote path mapping datasource ([194b59a](https://github.com/devopsarr/terraform-provider-radarr/commit/194b59a653d9474f0787947a412497e3d9513ca8))
* add remote path mapping resource ([07ccc74](https://github.com/devopsarr/terraform-provider-radarr/commit/07ccc7438c3d3bb2ac97fbe5759711d8951b57de))
* add remote path mappings datasource ([19d2d59](https://github.com/devopsarr/terraform-provider-radarr/commit/19d2d590fb2fd651b70fc626350f38ac7f13c787))
* add restriction data source ([5969d6f](https://github.com/devopsarr/terraform-provider-radarr/commit/5969d6fb1fb55459114a00d441b24097083ed3b8))
* add restriction resource ([d5ee234](https://github.com/devopsarr/terraform-provider-radarr/commit/d5ee2341895cce7fe34fcc2ae66a7819dd24e7ba))
* add restrictions datasource ([89623f3](https://github.com/devopsarr/terraform-provider-radarr/commit/89623f305071f704c94e977cd83c98bbb6db24a8))

## [1.3.0](https://github.com/devopsarr/terraform-provider-radarr/compare/v1.2.0...v1.3.0) (2022-11-14)


### Features

* add delay profile datasource ([9fc9669](https://github.com/devopsarr/terraform-provider-radarr/commit/9fc9669914b6cc9b2451627e8fe403e6a958d0d5))
* add delay profile resource ([a839652](https://github.com/devopsarr/terraform-provider-radarr/commit/a839652566bc02e1a7682685cd267539e7d96166))
* add delay profiles datasource ([4bddbd1](https://github.com/devopsarr/terraform-provider-radarr/commit/4bddbd199fcd7d11ca8e765c1bb7a8c3292e249e))
* add notification custom script resource ([8f1ca4c](https://github.com/devopsarr/terraform-provider-radarr/commit/8f1ca4cfc7179fa6a4dbe2a198fa3f36a2cf30ae))
* add notification datasource ([5ddf2d3](https://github.com/devopsarr/terraform-provider-radarr/commit/5ddf2d3f92f1a48121d6d190f6dc3cb2cefc8b4b))
* add notification resource ([f59290c](https://github.com/devopsarr/terraform-provider-radarr/commit/f59290c69d9a9483eded499a5c94c3f80212abf5))
* add notification webhook resource ([808c3b2](https://github.com/devopsarr/terraform-provider-radarr/commit/808c3b239e7fe04f8603a9b528a547f3342ce72d))
* add notifications data source ([3b39c57](https://github.com/devopsarr/terraform-provider-radarr/commit/3b39c5732f8a187266d9efa78b474fa2d7262624))

## [1.2.0](https://github.com/devopsarr/terraform-provider-radarr/compare/v1.1.0...v1.2.0) (2022-11-04)


### Features

* add media_management datasource ([299ac85](https://github.com/devopsarr/terraform-provider-radarr/commit/299ac85241b519fb7a00553608bf23edc2087965))
* add media_management resource ([fa458ad](https://github.com/devopsarr/terraform-provider-radarr/commit/fa458adc50a0cfd0a3206e4aa25c641f79d1e50c))
* add naming datasource ([7c2aa8a](https://github.com/devopsarr/terraform-provider-radarr/commit/7c2aa8a130d4cb877f890b5d2e06a243668f3f47))
* add naming resource ([62183c6](https://github.com/devopsarr/terraform-provider-radarr/commit/62183c6e06b76de87380d2538b174bd7bfc835e6))
* add root_folder datasource ([322ccc0](https://github.com/devopsarr/terraform-provider-radarr/commit/322ccc0a2a2bb426c58249d71f5b50d6193e93c1))
* add root_folder resource ([374f87f](https://github.com/devopsarr/terraform-provider-radarr/commit/374f87f4738fa17fd1feac2630f1fbb1db6b038d))
* add root_folders datasource ([0f1709a](https://github.com/devopsarr/terraform-provider-radarr/commit/0f1709ae59ce11ec99df4ea6a4078d63e2a548ec))
* add system status datasource ([05fa5aa](https://github.com/devopsarr/terraform-provider-radarr/commit/05fa5aa8f1128ca9178c693a7375ad6ce7ed2f1f))
* add tag datasource ([defab86](https://github.com/devopsarr/terraform-provider-radarr/commit/defab8632c25f6bd997e22aa69461247f92566ea))

## [1.1.0](https://github.com/devopsarr/terraform-provider-radarr/compare/v1.0.0...v1.1.0) (2022-08-29)


### Features

* add validators ([34e6668](https://github.com/devopsarr/terraform-provider-radarr/commit/34e666859c874de5b6d3a34788824126ff8519f9))


### Bug Fixes

* remove set parameter for framework 0.9.0 ([0fbe0d5](https://github.com/devopsarr/terraform-provider-radarr/commit/0fbe0d5145a03c97b016d3fb6e6383249fcb44de))

## 1.0.0 (2022-03-15)


### Features

* first configuration ([3d3b720](https://github.com/devopsarr/terraform-provider-radarr/commit/3d3b720ca0d43b640611940f831a724d3f0f7027))
