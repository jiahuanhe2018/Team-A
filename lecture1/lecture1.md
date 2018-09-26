##  Question 1

> 需要对ring signature，zk-snark进行比较，ring signature decoy的数量在多少的时候更消耗时
> 间空间，zk-snark在使用上占用多少空间，计算时间相比哪个更快



| decoy number | time(ms) | memory(KB) |
| ------------ | -------- | ---------- |
| 5            | 30       | 1500       |
| 20           | 131      | 2930       |
| 100          | 722      | 1820       |
| 500          | 3510     | 3204       |

ring-signature执行速度与decoy数量有关，基本呈线性增长关系。ring-signature和zk-SNARKs比较，后者的运行更有效。



## Question 2

> 将bitcoin，ethereum，monero，zcash，EOS 的交易、相关交易易属性、块大小以及填入多
> 少交易写在report



| Token    | 出块间隔 | 交易属性                 | 块大小                         | 交易数       |
| -------- | -------- | ------------------------ | ------------------------------ | ------------ |
| bitcoin  | 10min    | PoW最终一致性            | 1M                             | 2000左右     |
| ethereum | 15s      | PoW,有图灵完备的智能合约 | 1M                             | 70           |
| monero   | 600s     | ring-signature实现匿名   | 最大为前100个块大小中位数的2倍 | 1000左右     |
| zcash    | 2.5min   | zk-SNARKs实现匿名交易    | 2M                             | 3000左右     |
| eos      | 0.5s     | DPoS+PBFT，拥有权限体系  | 不固定                         | 现在1500左右 |

