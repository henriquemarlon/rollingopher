import {
  Address,
  encodeFunctionData,
  encodePacked,
  Hex,
  hexToBytes,
  zeroAddress,
} from "viem";
import { inputsAbi, outputsAbi } from "./contracts";

type AdvanceInput = {
  chainId?: bigint;
  appContract?: Address;
  msgSender: Address;
  blockNumber?: bigint;
  blockTimestamp?: bigint;
  prevRandao?: bigint;
  index?: bigint;
  payload?: Hex;
};

type ERC20Deposit = {
  tokenAddress: Address;
  sender: Address;
  amount: bigint;
  execLayerData?: Hex;
};

type EtherDeposit = {
  sender: Address;
  amount: bigint;
  execLayerData?: Hex;
};

type Voucher = {
  destination: Address;
  value: bigint;
  payload: Hex;
};

type DelegateCallVoucher = {
  destination: Address;
  payload: Hex;
};

export const encodeAdvanceInput = (data: AdvanceInput): Buffer => {
  const {
    chainId = 0n,
    appContract = zeroAddress,
    msgSender = zeroAddress,
    blockNumber = 0n,
    blockTimestamp = 0n,
    prevRandao = 0n,
    index = 0n,
    payload = "0x",
  } = data;
  return Buffer.from(
    hexToBytes(
      encodeFunctionData({
        abi: inputsAbi,
        functionName: "EvmAdvance",
        args: [
          chainId,
          appContract,
          msgSender,
          blockNumber,
          blockTimestamp,
          prevRandao,
          index,
          payload,
        ],
      })
    )
  );
};

export const encodeVoucherOutput = (data: Voucher) => {
  return encodeFunctionData({
    abi: outputsAbi,
    functionName: "Voucher",
    args: [data.destination, data.value, data.payload],
  });
};

export const encodeDelegateCallVoucherOutput = (data: DelegateCallVoucher) => {
  return encodeFunctionData({
    abi: outputsAbi,
    functionName: "DelegateCallVoucher",
    args: [data.destination, data.payload],
  });
};

export const encodeNoticeOutput = (payload: Hex): Hex => {
  return encodeFunctionData({
    abi: outputsAbi,
    functionName: "Notice",
    args: [payload],
  });
};

export const encodeErc20Deposit = (data: ERC20Deposit) => {
  const { tokenAddress, sender, amount, execLayerData = "0x" } = data;
  return encodePacked(
    ["address", "address", "uint256", "bytes"],
    [tokenAddress, sender, amount, execLayerData]
  );
};


export const encodeEtherDeposit = (data: EtherDeposit): Hex => {
  const { sender, amount, execLayerData = "0x" } = data;
  return encodePacked(
    ["address", "uint256", "bytes"],
    [sender, amount, execLayerData]
  );
};
