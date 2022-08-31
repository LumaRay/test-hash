# Testing native vs custom hashmap algorithms

The idea is to benchmark performance over <byte[32], int32> hashmaps in different languages where byte[32] is itself SHA256 of some data.

Supposing RAM is not an issue, so some tests have consumed up to 16 Gb of RAM.

The custom algorithm can be way imperfect, so take these tests with a grain of salt and of course better alternatives are very welcome. 

The algorithm doesn't hash a key, instead it expects SHA256 as a key and stores each key and value pair in it's hash tree.

__"Tree size"__ means size of pointer arrays in the top lovel of hash map and in the consequent layers.

So __3\*8+7 / 1\*8-4__ means 3 bytes plus 7 bits of key data will be used for addressing the top pointer array of the hash tree, while 1 byte minus 4 bits of subsequent key data will be used for addressing pointer arrays of the lower levels of the hash tree.

Basically a pointer array size in bytes of any hash tree level on x64 system can be calculated as __(1 << bits) \* 8__ and these are the most RAM consuming parts of the algorithms.

The results are in milliseconds and consist of two tests in a form of __T1/T2__ where:
- T1: adding pairs
- T2: getting and validating values

The testing has been performed on Windows 7x64.

- CPU: Intel Core i7-4790 3.6 GHz
- RAM: 32 Gb

---

## 8'000'000 records with test hash 2d array on a stack

|    Tree size   |    c++    |
|:--------------:|:---------:|
|      1 / 1     | 3751/2230 |
|      2 / 1     | 3771/2265 |
|      2 / 2     |  mem over |
|      3 / 1     | 2911/1559 |
|      1 / 2     |  mem over |
|    3*8+1 / 1   | 2715/1429 |
|    3*8+2 / 1   | 2025/1108 |
|    3*8+3 / 1   | 2222/1215 |
|    3*8+4 / 1   | 2394/1246 |
|    3*8+5 / 1   | 3019/1363 |
|  3*8+2 / 1*8+1 | 2211/1127 |
| 3*8+2 / 1*8-1  | 1996/1136 |
| 3*8+2 / 1*8-2  | 1901/1060 |
| 3*8+2 / 1*8-3  | 1900/1049 |
| 3*8+2 / 1*8-4  | 1878/1070 |
| 3*8+2 / 1*8-5  | 1901/1045 |
| 3*8+2 / 1*8-6  | 1865/1056 |
| 3*8+2 / 1*8-7  | 1864/1079 |
| 3*8+3 / 1*8-6  | 2152/1179 |
| 3*8+1 / 1*8-6  | 1921/1285 |
| 3*8-23 / 1*8-7 | stack ovr |
| 3*8-22 / 1*8-6 | 8072/8726 |
| 3*8-21 / 1*8-5 | 6381/6659 |
| 3*8-20 / 1*8-4 | 4756/4725 |
| 3*8-19 / 1*8-3 | 4022/3695 |
| 3*8-18 / 1*8-2 | 3543/3102 |
| 3*8-17 / 1*8-1 | 3830/3060 |

## 8'000'000 records with test hash 2d array on a heap

| Language | c | c++ | c# | go | rust |
|----------|---|-----|----|----|------|
|          |   |     |    |    |      |
|          |   |     |    |    |      |
|          |   |     |    |    |      |

## 16'000'000 records with test hash 2d array on a heap

| Language | c | c++ | c# | go | rust |
|----------|---|-----|----|----|------|
|          |   |     |    |    |      |
|          |   |     |    |    |      |
|          |   |     |    |    |      |

## 32'000'000 records with test hash 2d array on a heap

| Language | c | c++ | c# | go | rust |
|----------|---|-----|----|----|------|
|          |   |     |    |    |      |
|          |   |     |    |    |      |
|          |   |     |    |    |      |

## 64'000'000 records with test hash 2d array on a heap

| Language | c | c++ | c# | go | rust |
|----------|---|-----|----|----|------|
|          |   |     |    |    |      |
|          |   |     |    |    |      |
|          |   |     |    |    |      |
