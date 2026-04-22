# Deploying to Sepolia

Sepolia is a public Ethereum testnet, not required for local development. Use it only when testing against a shared network.

1. Create a wallet in MetaMask and fund it from a Sepolia faucet.
2. Create an API key on [Infura](https://infura.io) or [Alchemy](https://alchemy.com).
3. Set the environment variables and deploy the contract:

```bash
export ETHEREUM_RPC_URL=https://sepolia.infura.io/v3/<your-key>
export DEPLOYER_PRIVATE_KEY=0x<your-wallet-private-key>
mise run deploy:sepolia
```

The task writes the deployed contract address to `CONTRACT_ADDRESS` in `.env` automatically.
