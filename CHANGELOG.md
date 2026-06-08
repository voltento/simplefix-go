## [1.18.1](https://gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/compare/v1.18.0...v1.18.1) (2026-06-08)


### Bug Fixes

* accept zero-count repeating groups in unmarshaler (B2CT-18604) ([3690472](https://gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/commit/3690472d3794d6a4ff421608be1d02a908f875db))


### Documentation

* fix default branch to master in CLAUDE.md (B2CT-18604) ([18b6f70](https://gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/commit/18b6f701430c641509959467fd697f91e4695551))

## [1.18.0](https://gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/compare/v1.17.3...v1.18.0) (2026-05-15)


### Features

* **B2CT-19469:** rename module path to GitLab canonical path ([db959eb](https://gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/commit/db959ebf76887f758062b4729b044e51b0cb31b9))


### Bug Fixes

* **ci:** restore Go 1.24 coverage image pin ([6a3449b](https://gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/commit/6a3449ba09c86b07ce6238e7e71bcda8956e27e8))
* relax golangci rules for legacy simplefix-go code (B2CT-19468) ([893a960](https://gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/commit/893a9607047021c5055c33c48c3e8d3e46b7506d))
* **review:** add fixgen import placeholder to CustomLogon example ([c1834a0](https://gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/commit/c1834a0c8ca27218ebf3cffcae6e726ecf16d958))
* **review:** address R2 findings — state machine, layer order, CI description, RunNewInitiator sig, Disconnect, infinite recursion, TOC anchors ([9212f0d](https://gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/commit/9212f0dc0bfe70d6316bc1ed9039ac611e3037c4))
* **review:** address R3 findings — CustomLogon type-assert, SendingMessage interface, state machine diagram, Disconnect edge, Remove() race, packages table, link labels ([2e05c0f](https://gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/commit/2e05c0f7ebe373feda9dad7f999ffeb8ca095da2))
* **review:** address R4 findings — MO-01 SetCounterPartyID, P-02 initiator state, F-01 sessionOpts wiring ([47b8523](https://gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/commit/47b8523f6e891c09ee51a32b874010b991e1cd8c))
* **review:** address R5 findings — Value types list, messages import, typo ([2e9c4ab](https://gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/commit/2e9c4ab0e45c8c7152643ad0ef4c9f434e87a5f2))


### Documentation

* **B2CT-19469:** add .claude/ AI guidance files from PR [#53](https://gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/issues/53) ([ef5ab29](https://gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/commit/ef5ab29428f4c208c87f44e68c495803b0129e12))
* fix .claude guidance and README badge/install errors (ralph-loop R1+CV) ([517ecc6](https://gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/commit/517ecc67cdeab5070b7af2f4c9e57e1eee8610a7))
