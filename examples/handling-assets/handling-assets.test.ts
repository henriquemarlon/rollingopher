import { encodeFunctionData, getAddress, padHex } from "viem";
import { afterAll, describe, expect, it } from "vitest";
import {
  encodeAdvanceInput,
  encodeEtherDeposit,
  encodeERC20Deposit,
  encodeERC721Deposit,
  encodeERC1155SingleDeposit,
  encodeERC1155BatchDeposit,
  encodeEtherWithdrawal,
  encodeERC20Withdrawal,
  encodeERC721Withdrawal,
  encodeERC1155SingleWithdrawal,
  encodeERC1155BatchWithdrawal,
  encodeEtherTransfer,
  encodeERC20Transfer,
  encodeVoucherOutput,
} from "../../encoder";
import { rollups, spawn } from "@tuler/node-cartesi-machine";

const BOB = getAddress("0x0000000000000000000000000000000000000002");
const TOKEN = getAddress("0x0000000000000000000000000000000000000ABC");
const ALICE = getAddress("0x0000000000000000000000000000000000000001");
const APP_CONTRACT = getAddress("0xab7528bb862fB57E8A2BCd567a2e929a0Be56a5e");
const ETHER_PORTAL = getAddress("0xFfdbe43d4c855BF7e0f105c400A50857f53AB044");
const ERC20_PORTAL = getAddress("0x9C21AEb2093C32DDbC53eEF24B873BDCd1aDa1DB");
const ERC721_PORTAL = getAddress("0x237F8DD094C0e47f4236f12b4Fa01d6Dae89fb87");
const ERC1155_SINGLE_PORTAL = getAddress(
  "0x7CFB0193Ca87eB6e48056885E026552c3A941FC4",
);
const ERC1155_BATCH_PORTAL = getAddress(
  "0xedB53860A6B52bbb7561Ad596416ee9965B055Aa",
);

describe("Handling Assets Tests", () => {
  const remoteMachine = spawn().load(".cartesi/image");
  const machine = rollups(remoteMachine, { noRollback: true });

  it("should handle ether deposit", () => {
    const payload = encodeEtherDeposit({
      sender: ALICE,
      amount: 1000n,
    });
    const { outputs } = machine.advance(
      encodeAdvanceInput({
        appContract: APP_CONTRACT,
        msgSender: ETHER_PORTAL,
        blockTimestamp: BigInt(Math.floor(Date.now() / 1000)),
        payload,
      }),
      { collect: true },
    );
    expect(outputs.length).toBe(0);
  });

  it("should handle ether withdrawal", () => {
    const payload = encodeEtherWithdrawal({
      amount: 100n,
    });
    const { outputs } = machine.advance(
      encodeAdvanceInput({
        appContract: APP_CONTRACT,
        msgSender: ALICE,
        blockTimestamp: BigInt(Math.floor(Date.now() / 1000)),
        payload,
      }),
      { collect: true },
    );
    expect(outputs.length).toBe(1);

    const expectedVoucher = encodeVoucherOutput({
      destination: ALICE,
      value: 100n,
      payload: "0x",
    });
    expect(outputs[0].toString("hex")).toBe(expectedVoucher.slice(2));
  });

  it("should handle ether transfer", () => {
    const payload = encodeEtherTransfer({
      receiver: padHex(BOB, { size: 32 }),
      amount: 100n,
    });
    const { outputs } = machine.advance(
      encodeAdvanceInput({
        appContract: APP_CONTRACT,
        msgSender: ALICE,
        blockTimestamp: BigInt(Math.floor(Date.now() / 1000)),
        payload,
      }),
      { collect: true },
    );
    expect(outputs.length).toBe(0);
  });

  it("should handle ERC20 deposit", () => {
    const payload = encodeERC20Deposit({
      tokenAddress: TOKEN,
      sender: ALICE,
      amount: 500n,
    });
    const { outputs } = machine.advance(
      encodeAdvanceInput({
        appContract: APP_CONTRACT,
        msgSender: ERC20_PORTAL,
        blockTimestamp: BigInt(Math.floor(Date.now() / 1000)),
        payload,
      }),
      { collect: true },
    );
    expect(outputs.length).toBe(0);
  });

  it("should handle ERC20 withdrawal", () => {
    const payload = encodeERC20Withdrawal({
      token: TOKEN,
      amount: 100n,
    });
    const { outputs } = machine.advance(
      encodeAdvanceInput({
        appContract: APP_CONTRACT,
        msgSender: ALICE,
        blockTimestamp: BigInt(Math.floor(Date.now() / 1000)),
        payload,
      }),
      { collect: true },
    );
    expect(outputs.length).toBe(1);

    const erc20TransferPayload = encodeFunctionData({
      abi: [
        {
          name: "transfer",
          type: "function",
          inputs: [
            { name: "to", type: "address" },
            { name: "amount", type: "uint256" },
          ],
        },
      ],
      functionName: "transfer",
      args: [ALICE, 100n],
    });
    const expectedVoucher = encodeVoucherOutput({
      destination: TOKEN,
      value: 0n,
      payload: erc20TransferPayload,
    });
    expect(outputs[0].toString("hex")).toBe(expectedVoucher.slice(2));
  });

  it("should handle ERC20 transfer", () => {
    const payload = encodeERC20Transfer({
      token: TOKEN,
      receiver: padHex(BOB, { size: 32 }),
      amount: 100n,
    });
    const { outputs } = machine.advance(
      encodeAdvanceInput({
        appContract: APP_CONTRACT,
        msgSender: ALICE,
        blockTimestamp: BigInt(Math.floor(Date.now() / 1000)),
        payload,
      }),
      { collect: true },
    );
    expect(outputs.length).toBe(0);
  });

  it("should handle ERC721 deposit", () => {
    const payload = encodeERC721Deposit({
      tokenAddress: TOKEN,
      sender: ALICE,
      tokenId: 1n,
    });
    const { outputs } = machine.advance(
      encodeAdvanceInput({
        appContract: APP_CONTRACT,
        msgSender: ERC721_PORTAL,
        blockTimestamp: BigInt(Math.floor(Date.now() / 1000)),
        payload,
      }),
      { collect: true },
    );
    expect(outputs.length).toBe(0);
  });

  it("should handle ERC721 withdrawal", () => {
    const payload = encodeERC721Withdrawal({
      token: TOKEN,
      tokenId: 1n,
    });
    const { outputs } = machine.advance(
      encodeAdvanceInput({
        appContract: APP_CONTRACT,
        msgSender: ALICE,
        blockTimestamp: BigInt(Math.floor(Date.now() / 1000)),
        payload,
      }),
      { collect: true },
    );
    expect(outputs.length).toBe(1);

    const erc721SafeTransferPayload = encodeFunctionData({
      abi: [
        {
          name: "safeTransferFrom",
          type: "function",
          inputs: [
            { name: "from", type: "address" },
            { name: "to", type: "address" },
            { name: "tokenId", type: "uint256" },
          ],
        },
      ],
      functionName: "safeTransferFrom",
      args: [APP_CONTRACT, ALICE, 1n],
    });
    const expectedVoucher = encodeVoucherOutput({
      destination: TOKEN,
      value: 0n,
      payload: erc721SafeTransferPayload,
    });
    expect(outputs[0].toString("hex")).toBe(expectedVoucher.slice(2));
  });

  it("should handle ERC1155 single deposit", () => {
    const payload = encodeERC1155SingleDeposit({
      tokenAddress: TOKEN,
      sender: ALICE,
      tokenId: 1n,
      amount: 100n,
    });
    const { outputs } = machine.advance(
      encodeAdvanceInput({
        appContract: APP_CONTRACT,
        msgSender: ERC1155_SINGLE_PORTAL,
        blockTimestamp: BigInt(Math.floor(Date.now() / 1000)),
        payload,
      }),
      { collect: true },
    );
    expect(outputs.length).toBe(0);
  });

  it("should handle ERC1155 single withdrawal", () => {
    const payload = encodeERC1155SingleWithdrawal({
      token: TOKEN,
      tokenId: 1n,
      amount: 50n,
    });
    const { outputs } = machine.advance(
      encodeAdvanceInput({
        appContract: APP_CONTRACT,
        msgSender: ALICE,
        blockTimestamp: BigInt(Math.floor(Date.now() / 1000)),
        payload,
      }),
      { collect: true },
    );
    expect(outputs.length).toBe(1);

    const erc1155SafeTransferPayload = encodeFunctionData({
      abi: [
        {
          name: "safeTransferFrom",
          type: "function",
          inputs: [
            { name: "from", type: "address" },
            { name: "to", type: "address" },
            { name: "id", type: "uint256" },
            { name: "value", type: "uint256" },
            { name: "data", type: "bytes" },
          ],
        },
      ],
      functionName: "safeTransferFrom",
      args: [APP_CONTRACT, ALICE, 1n, 50n, "0x"],
    });
    const expectedVoucher = encodeVoucherOutput({
      destination: TOKEN,
      value: 0n,
      payload: erc1155SafeTransferPayload,
    });
    expect(outputs[0].toString("hex")).toBe(expectedVoucher.slice(2));
  });

  it("should handle ERC1155 batch deposit", () => {
    const payload = encodeERC1155BatchDeposit({
      tokenAddress: TOKEN,
      sender: ALICE,
      tokenIds: [1n, 2n, 3n],
      amounts: [10n, 20n, 30n],
    });
    const { outputs } = machine.advance(
      encodeAdvanceInput({
        appContract: APP_CONTRACT,
        msgSender: ERC1155_BATCH_PORTAL,
        blockTimestamp: BigInt(Math.floor(Date.now() / 1000)),
        payload,
      }),
      { collect: true },
    );
    expect(outputs.length).toBe(0);
  });

  it("should handle ERC1155 batch withdrawal", () => {
    const payload = encodeERC1155BatchWithdrawal({
      token: TOKEN,
      tokenIds: [2n, 3n],
      amounts: [5n, 10n],
    });
    const { outputs } = machine.advance(
      encodeAdvanceInput({
        appContract: APP_CONTRACT,
        msgSender: ALICE,
        blockTimestamp: BigInt(Math.floor(Date.now() / 1000)),
        payload,
      }),
      { collect: true },
    );
    expect(outputs.length).toBe(1);

    const erc1155SafeBatchTransferPayload = encodeFunctionData({
      abi: [
        {
          name: "safeBatchTransferFrom",
          type: "function",
          inputs: [
            { name: "from", type: "address" },
            { name: "to", type: "address" },
            { name: "ids", type: "uint256[]" },
            { name: "values", type: "uint256[]" },
            { name: "data", type: "bytes" },
          ],
        },
      ],
      functionName: "safeBatchTransferFrom",
      args: [APP_CONTRACT, ALICE, [2n, 3n], [5n, 10n], "0x"],
    });
    const expectedVoucher = encodeVoucherOutput({
      destination: TOKEN,
      value: 0n,
      payload: erc1155SafeBatchTransferPayload,
    });
    expect(outputs[0].toString("hex")).toBe(expectedVoucher.slice(2));
  });

  it("should handle ether balance query", () => {
    const payload = JSON.stringify({
      method: "ledger_getBalance",
      params: [padHex(ALICE, { size: 32 })],
    });
    const reports = machine.inspect(Buffer.from(payload), { collect: true });
    expect(reports.length).toBe(1);
    const balance = BigInt("0x" + reports[0].toString("hex"));
    expect(balance).toBe(800n);
  });

  it("should handle ether supply query", () => {
    const payload = JSON.stringify({
      method: "ledger_getTotalSupply",
      params: [],
    });
    const reports = machine.inspect(Buffer.from(payload), { collect: true });
    expect(reports.length).toBe(1);
    const supply = BigInt("0x" + reports[0].toString("hex"));
    expect(supply).toBe(900n);
  });

  it("should handle ERC20 balance query", () => {
    const payload = JSON.stringify({
      method: "ledger_getBalance",
      params: [padHex(ALICE, { size: 32 }), TOKEN],
    });
    const reports = machine.inspect(Buffer.from(payload), { collect: true });
    expect(reports.length).toBe(1);
    const balance = BigInt("0x" + reports[0].toString("hex"));
    expect(balance).toBe(300n);
  });

  it("should handle ERC20 supply query", () => {
    const payload = JSON.stringify({
      method: "ledger_getTotalSupply",
      params: [TOKEN],
    });
    const reports = machine.inspect(Buffer.from(payload), { collect: true });
    expect(reports.length).toBe(1);
    const supply = BigInt("0x" + reports[0].toString("hex"));
    expect(supply).toBe(400n);
  });

  it("should handle tokenId 1 balance query", () => {
    const payload = JSON.stringify({
      method: "ledger_getBalance",
      params: [padHex(ALICE, { size: 32 }), TOKEN, "1"],
    });
    const reports = machine.inspect(Buffer.from(payload), { collect: true });
    expect(reports.length).toBe(1);
    const balance = BigInt("0x" + reports[0].toString("hex"));
    expect(balance).toBe(60n);
  });

  it("should handle tokenId 1 supply query", () => {
    const payload = JSON.stringify({
      method: "ledger_getTotalSupply",
      params: [TOKEN, "1"],
    });
    const reports = machine.inspect(Buffer.from(payload), { collect: true });
    expect(reports.length).toBe(1);
    const supply = BigInt("0x" + reports[0].toString("hex"));
    expect(supply).toBe(60n);
  });

  it("should handle tokenId 2 balance query", () => {
    const payload = JSON.stringify({
      method: "ledger_getBalance",
      params: [padHex(ALICE, { size: 32 }), TOKEN, "2"],
    });
    const reports = machine.inspect(Buffer.from(payload), { collect: true });
    expect(reports.length).toBe(1);
    const balance = BigInt("0x" + reports[0].toString("hex"));
    expect(balance).toBe(15n);
  });

  it("should handle tokenId 3 supply query", () => {
    const payload = JSON.stringify({
      method: "ledger_getTotalSupply",
      params: [TOKEN, "3"],
    });
    const reports = machine.inspect(Buffer.from(payload), { collect: true });
    expect(reports.length).toBe(1);
    const supply = BigInt("0x" + reports[0].toString("hex"));
    expect(supply).toBe(20n);
  });

  afterAll(() => {
    machine.shutdown();
  });
});
