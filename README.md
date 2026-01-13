<br>
<p align="center">
    <img src="https://github.com/user-attachments/assets/080bb0be-060c-4813-85b4-6d9bf25af01f" align="center" width="20%">
</p>
<br>
<div align="center">
	<i>Cartesi Rollups Go High-Level Framework</i>
</div>
<div align="center">
	<b>Any Code. Ethereum’s Security.</b>
</div>
<br>
<p align="center">
	<img src="https://img.shields.io/github/license/henriquemarlon/rollingopher?style=default&logo=opensourceinitiative&logoColor=white&color=008DA5" alt="license">
	<img src="https://img.shields.io/github/last-commit/henriquemarlon/rollingopher?style=default&logo=git&logoColor=white&color=000000" alt="last-commit">
</p>

## Table of Contents

- [Overview](#overview)
- [Packages](#packages)
- [Examples](#examples)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Testing](#testing)

## Overview

**Rollingopher** provides Go bindings for building Cartesi Rollups applications. It wraps the low-level C libraries (`libcmt` and `libcma`) to provide idiomatic Go APIs for:

- **Rollup operations**: Reading inputs, emitting vouchers, notices, and reports
- **Asset management**: Handling deposits, withdrawals, and transfers for Ether, ERC20, ERC721, and ERC1155 tokens
- **Ledger management**: Tracking account balances and asset supplies

## Packages

| Package      | Description                                                                              |
| ------------ | ---------------------------------------------------------------------------------------- |
| `pkg/rollup` | CGO bindings for `libcmt` - handles rollup state machine operations                      |
| `pkg/ledger` | CGO bindings for `libcma` - manages asset ledger and account balances                    |
| `pkg/parser` | Go implementation for decoding inputs                                                    |
| `pkg/router` | TBD                                                                                      |
| `pkg/tester` | TBD                                                                                      |

## Examples

| Example           | Description                                                                 |
| ----------------- | --------------------------------------------------------------------------- |
| `echo`            | Simple example using only `rollup` - echoes inputs                          |
| `handling-assets` | Full asset management example using `rollup`, `ledger`, and `parser`        |

## Getting Started

### Prerequisites

- [Go 1.24+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) with RISC-V support (for building Cartesi images)
- [Cartesi CLI](https://docs.cartesi.io/cartesi-rollups/1.5/development/installation/)
- [Node.js](https://nodejs.org/) and [pnpm](https://pnpm.io/) (for running tests)

### Testing

1. Install dependencies:

   ```sh
   pnpm install
   ```

2. Build the example you want to test:

   ```sh
   make build-echo

   # or
   
   make build-handling-assets
   ```

3. Generate contracts ABIs:

   ```sh
   pnpm codegen
   ```

4. Run tests for a specific example:

   ```sh
   make test-echo

   # or
   
   make test-handling-assets
   ```
