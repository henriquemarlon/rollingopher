import { getAddress, toHex } from "viem";
import { afterAll, describe, expect, it } from "vitest";
import { encodeAdvanceInput, encodeNoticeOutput } from "../../encoder";
import { rollups, spawn } from "@tuler/node-cartesi-machine";

const ALICE = getAddress("0x0000000000000000000000000000000000000001");
const APP_CONTRACT = getAddress("0xab7528bb862fB57E8A2BCd567a2e929a0Be56a5e");

describe("Echo Tests", () => {
  const remoteMachine = spawn().load(".cartesi/image");
  const machine = rollups(remoteMachine, { noRollback: true });

  it("should echo payload as notice on advance", () => {
    const payload = toHex("Hello, Cartesi!");
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

    const expectedNotice = encodeNoticeOutput({ payload });
    expect(outputs[0].toString("hex")).toBe(expectedNotice.slice(2));
  });

  it("should echo payload as report on inspect", () => {
    const payload = "Hello, Inspect!";
    const reports = machine.inspect(Buffer.from(payload), { collect: true });
    expect(reports.length).toBe(1);
    expect(reports[0].toString()).toBe(payload);
  });

  afterAll(() => {
    machine.shutdown();
  });
});
