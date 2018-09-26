> 1.需要对ring signature，zk-snark进行比较，ring signature decoy的数量在多少的时候更消耗时间空间，zk-sanrk在使用上占用多少空间，计算时间相比哪个更快

环签中的decoy越多，越消耗时间空间。
`TODO：zk-snark占用多少空间`

>2、将bitcoin，ethereum，monero，zcash，EOS 的交易、相关交易易属性、块大小以及填入多少交易写在report中

|project|交易字段|块大小|一个区块平均交易数|
|:--:|:--:|:--:|:--:|:--:|:--|
|Bitcoin|https://blog.brakmic.com/bitcoin-internals-part-2/|理论上1M，但实际会上下浮动|https://btc.com/stats/block-size|
|Ethereum|nonce、gasPrice、gasLimit、to、value、v,r,s、init\data|https://etherscan.io/chart/blocksize|与区块的gasLimit有关|
|Monero|https://monero.stackexchange.com/questions/2136/understanding-the-structure-of-a-monero-transaction#2150|大小与最后一百个区块的大小的中位数有关||
|Zcash|https://blog.z.cash/anatomy-of-zcash/|2M||
|EOS|https://developers.eos.io/eosio-cpp/reference#get_action-1|todo|10左右|
