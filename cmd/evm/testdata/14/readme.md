## Difficulty calculation

This test shows how the `evm t8n` can be used to calculate the (ethash) difficulty, if none is provided by the caller. 

Calculating it (with an empty set of txs) using `London` rules (and no provided unclehash for the parent block):
```
[user@work evm]$ ./evm t8n --input.alloc=./testdata/14/alloc.json --input.txs=./testdata/14/txs.json --input.env=./testdata/14/env.json --output.result=stdout --state.fork=London
INFO [03-09|10:43:57.070] Trie dumping started                     root=6f0588..7f4bdc
INFO [03-09|10:43:57.070] Trie dumping complete                    accounts=2 elapsed="214.663µs"
INFO [03-09|10:43:57.071] Wrote file                               file=alloc.json
{
  "result": {
    "stateRoot": "0x6f058887ca01549716789c380ede95aecc510e6d1fdc4dbf67d053c7c07f4bdc",
    "txRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "logsHash": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "receipts": [],
    "currentDifficulty": "0x2000020000000",
    "gasUsed": "0x0",
    "currentBaseFee": "0x500"
  }
}
```
Same thing, but this time providing a non-empty (and non-`emptyKeccak`) unclehash, which leads to a slightly different result:
```
[user@work evm]$ ./evm t8n --input.alloc=./testdata/14/alloc.json --input.txs=./testdata/14/txs.json --input.env=./testdata/14/env.uncles.json --output.result=stdout --state.fork=London
INFO [03-09|10:44:20.511] Trie dumping started                     root=6f0588..7f4bdc
INFO [03-09|10:44:20.511] Trie dumping complete                    accounts=2 elapsed="184.319µs"
INFO [03-09|10:44:20.512] Wrote file                               file=alloc.json
{
  "result": {
    "stateRoot": "0x6f058887ca01549716789c380ede95aecc510e6d1fdc4dbf67d053c7c07f4bdc",
    "txRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "logsHash": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "receipts": [],
    "currentDifficulty": "0x1ff8020000000",
    "gasUsed": "0x0",
    "currentBaseFee": "0x500"
  }
}
```

