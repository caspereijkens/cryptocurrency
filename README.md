## How to make a testnet address and receive some coins
1. Run 
```bash
bin/make-testnet-address
```

## How to create a testnet transaction
```bash
create-testnet-transaction \
-in 7df617a04d8e90d2786bdd74ba5d6c034b7ca72860019488da4b1aaecf55c6eb:1 \
-out 500000:mwJn1YPMq7y5F8J3LkC5Hxg9PHyZ5K4cFv \
-out 300000:mzdx3vTWBLQtG8robVqd5CADEY2LKyJvrK
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