import {
  Address,
  concat,
  encodeAbiParameters,
  encodeFunctionData,
  Hex,
  hexToBytes,
  numberToHex,
  pad,
  parseAbiParameters,
  zeroAddress,
} from "viem";
import { inputsAbi, outputsAbi } from "./contracts";

// Ledger selectors
const SELECTOR_WITHDRAW_ETHER = "0x8cf70f0b";
const SELECTOR_WITHDRAW_ERC20 = "0x4f94d342";
const SELECTOR_WITHDRAW_ERC721 = "0x33acf293";
const SELECTOR_WITHDRAW_ERC1155_SINGLE = "0x8bb0a811";
const SELECTOR_WITHDRAW_ERC1155_BATCH = "0x50c80019";
const SELECTOR_TRANSFER_ETHER = "0xff67c903";
const SELECTOR_TRANSFER_ERC20 = "0x03d61dcd";

// Types
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

type EtherDeposit = {
  sender: Address;
  amount: bigint;
  execLayerData?: Hex;
};

type ERC20Deposit = {
  tokenAddress: Address;
  sender: Address;
  amount: bigint;
  execLayerData?: Hex;
};

type ERC721Deposit = {
  tokenAddress: Address;
  sender: Address;
  tokenId: bigint;
  execLayerData?: Hex;
};

type ERC1155SingleDeposit = {
  tokenAddress: Address;
  sender: Address;
  tokenId: bigint;
  amount: bigint;
  execLayerData?: Hex;
};

type ERC1155BatchDeposit = {
  tokenAddress: Address;
  sender: Address;
  tokenIds: bigint[];
  amounts: bigint[];
  execLayerData?: Hex;
};

type EtherWithdrawal = {
  amount: bigint;
  execLayerData?: Hex;
};

type ERC20Withdrawal = {
  token: Address;
  amount: bigint;
  execLayerData?: Hex;
};

type ERC721Withdrawal = {
  token: Address;
  tokenId: bigint;
  execLayerData?: Hex;
};

type ERC1155SingleWithdrawal = {
  token: Address;
  tokenId: bigint;
  amount: bigint;
  execLayerData?: Hex;
};

type ERC1155BatchWithdrawal = {
  token: Address;
  tokenIds: bigint[];
  amounts: bigint[];
  execLayerData?: Hex;
};

type EtherTransfer = {
  receiver: Hex;
  amount: bigint;
  execLayerData?: Hex;
};

type ERC20Transfer = {
  token: Address;
  receiver: Hex;
  amount: bigint;
  execLayerData?: Hex;
};

type Notice = {
  payload: Hex;
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
      }),
    ),
  );
};

// Outputs

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

export const encodeNoticeOutput = (data: Notice): Hex => {
  return encodeFunctionData({
    abi: outputsAbi,
    functionName: "Notice",
    args: [data.payload],
  });
};

// Ledger

export const encodeEtherWithdrawal = (data: EtherWithdrawal): Hex => {
  const { amount, execLayerData = "0x" } = data;
  const encoded = encodeAbiParameters(
    parseAbiParameters("uint256 amount, bytes execLayerData"),
    [amount, execLayerData],
  );
  return concat([SELECTOR_WITHDRAW_ETHER as Hex, encoded]);
};

export const encodeERC20Withdrawal = (data: ERC20Withdrawal): Hex => {
  const { token, amount, execLayerData = "0x" } = data;
  const encoded = encodeAbiParameters(
    parseAbiParameters("address token, uint256 amount, bytes execLayerData"),
    [token, amount, execLayerData],
  );
  return concat([SELECTOR_WITHDRAW_ERC20 as Hex, encoded]);
};

export const encodeERC721Withdrawal = (data: ERC721Withdrawal): Hex => {
  const { token, tokenId, execLayerData = "0x" } = data;
  const encoded = encodeAbiParameters(
    parseAbiParameters("address token, uint256 tokenId, bytes execLayerData"),
    [token, tokenId, execLayerData],
  );
  return concat([SELECTOR_WITHDRAW_ERC721 as Hex, encoded]);
};

export const encodeERC1155SingleWithdrawal = (
  data: ERC1155SingleWithdrawal,
): Hex => {
  const { token, tokenId, amount, execLayerData = "0x" } = data;
  const encoded = encodeAbiParameters(
    parseAbiParameters(
      "address token, uint256 tokenId, uint256 amount, bytes execLayerData",
    ),
    [token, tokenId, amount, execLayerData],
  );
  return concat([SELECTOR_WITHDRAW_ERC1155_SINGLE as Hex, encoded]);
};

export const encodeERC1155BatchWithdrawal = (
  data: ERC1155BatchWithdrawal,
): Hex => {
  const { token, tokenIds, amounts, execLayerData = "0x" } = data;
  const encoded = encodeAbiParameters(
    parseAbiParameters(
      "address token, uint256[] tokenIds, uint256[] amounts, bytes execLayerData",
    ),
    [token, tokenIds, amounts, execLayerData],
  );
  return concat([SELECTOR_WITHDRAW_ERC1155_BATCH as Hex, encoded]);
};

export const encodeEtherTransfer = (data: EtherTransfer): Hex => {
  const { receiver, amount, execLayerData = "0x" } = data;
  const encoded = encodeAbiParameters(
    parseAbiParameters("bytes32 receiver, uint256 amount, bytes execLayerData"),
    [receiver, amount, execLayerData],
  );
  return concat([SELECTOR_TRANSFER_ETHER as Hex, encoded]);
};

export const encodeERC20Transfer = (data: ERC20Transfer): Hex => {
  const { token, receiver, amount, execLayerData = "0x" } = data;
  const encoded = encodeAbiParameters(
    parseAbiParameters(
      "address token, bytes32 receiver, uint256 amount, bytes execLayerData",
    ),
    [token, receiver, amount, execLayerData],
  );
  return concat([SELECTOR_TRANSFER_ERC20 as Hex, encoded]);
};

// Portals

export const encodeEtherDeposit = (data: EtherDeposit): Hex => {
  const { sender, amount, execLayerData = "0x" } = data;
  return concat([
    sender,
    pad(numberToHex(amount), { size: 32 }),
    execLayerData,
  ]);
};

export const encodeERC20Deposit = (data: ERC20Deposit): Hex => {
  const { tokenAddress, sender, amount, execLayerData = "0x" } = data;
  return concat([
    tokenAddress,
    sender,
    pad(numberToHex(amount), { size: 32 }),
    execLayerData,
  ]);
};

export const encodeERC721Deposit = (data: ERC721Deposit): Hex => {
  const { tokenAddress, sender, tokenId, execLayerData = "0x" } = data;
  return concat([
    tokenAddress,
    sender,
    pad(numberToHex(tokenId), { size: 32 }),
    execLayerData,
  ]);
};

export const encodeERC1155SingleDeposit = (data: ERC1155SingleDeposit): Hex => {
  const { tokenAddress, sender, tokenId, amount, execLayerData = "0x" } = data;
  return concat([
    tokenAddress,
    sender,
    pad(numberToHex(tokenId), { size: 32 }),
    pad(numberToHex(amount), { size: 32 }),
    execLayerData,
  ]);
};

export const encodeERC1155BatchDeposit = (data: ERC1155BatchDeposit): Hex => {
  const {
    tokenAddress,
    sender,
    tokenIds,
    amounts,
    execLayerData = "0x",
  } = data;
  const encodedData = encodeAbiParameters(
    parseAbiParameters(
      "uint256[] tokenIds, uint256[] amounts, bytes baseLayerData, bytes execLayerData",
    ),
    [tokenIds, amounts, "0x", execLayerData],
  );
  return concat([tokenAddress, sender, encodedData]);
};
