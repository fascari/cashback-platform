import { ethers } from "hardhat";

async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("Deploying with account:", deployer.address);

  const CashbackToken = await ethers.getContractFactory("CashbackToken");
  const token = await CashbackToken.deploy();
  await token.waitForDeployment();

  const address = await token.getAddress();
  console.log("CashbackToken deployed to:", address);
  console.log("Set CONTRACT_ADDRESS=" + address + " in your .env");
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
