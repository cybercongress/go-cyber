
import sys

from web3.auto import w3

"""
Collects ethereum transactions as adjacency list.

Output - txt file with next structure:
FROM TO VALUE
str str int
...

Usage:
1) Start local Parity full node
2) Run Script with two params: start block  and end blocknumbers ex:
python ethereum-chain-to-adjacency-list.py 6235000 6235009
"""

firstBlockToDownload = int(sys.argv[1])
lastBlockToDownload = int(sys.argv[2])
fileName = "{}-{}_blocks_data.txt".format(firstBlockToDownload, lastBlockToDownload)

print("")
print("-----------------------------------------------")
print("About to download {}-{} blocks".format(firstBlockToDownload, lastBlockToDownload))
print("")
resultDataFile = open(fileName, "w")


for blockNumber in range(firstBlockToDownload, lastBlockToDownload + 1):

    print("Downloading {} block".format(blockNumber))
    block = w3.eth.getBlock(blockNumber, True)

    for tx in block.transactions:
        if tx.to is not None:
            resultDataFile.write("{} {} {}\r\n".format(tx.to, getattr(tx, 'from'), tx.value))
        else:
            resultDataFile.write("{} {} {}\r\n".format(tx.creates, getattr(tx, 'from'), tx.value))

resultDataFile.close()
print("Finished to download data into {}".format(fileName))
print("-----------------------------------------------")
print("")
