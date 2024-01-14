## How to send testnet coins to yourself
1. Run 
```bash
bin/make-testnet-address
```

## How to lookup a transaction
1. Run
```bash
bin/fetch-tx <transaction-id>
```

Example:
```
bin/fetch-tx 7dff938918f07619abd38e4510890396b1cef4fbeca154fb7aafba8843295ea2
```
(You just looked up the first btc transaction)

TODOs
1. Implement ProtoBuf for serializing messages.