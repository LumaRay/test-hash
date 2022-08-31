# Testing native vs custom hashmap algorithms

The idea is to benchmark performance over <byte[32], int32> hashmaps in different languages where byte[32] is itself SHA256 of some data.

Supposing RAM is not an issue, so some tests have consumed up to 16 Gb of RAM.

The custom algorithm can be way imperfect, so take these tests with a grain of salt and of course better alternatives are very welcome. 

The algorithm doesn't hash a key, instead it expects SHA256 as a key and stores each key and value pair in it's hash tree.

Everywhere release x64 configurations has been used.

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
|    1\*8 / 1\*8   | 3751/2230 |
|    2\*8 / 1\*8   | 3771/2265 |
|    2\*8 / 2\*8   |  mem over |
|    3\*8 / 1\*8   | 2911/1559 |
|    1\*8 / 2\*8   |  mem over |
|  3\*8+1 / 1\*8   | 2715/1429 |
|  3\*8+2 / 1\*8   | 2025/1108 |
|  3\*8+3 / 1\*8   | 2222/1215 |
|  3\*8+4 / 1\*8   | 2394/1246 |
|  3\*8+5 / 1\*8   | 3019/1363 |
|  3\*8+2 / 1\*8+1 | 2211/1127 |
| 3\*8+2 / 1\*8-1  | 1996/1136 |
| 3\*8+2 / 1\*8-2  | 1901/1060 |
| 3\*8+2 / 1\*8-3  | 1900/1049 |
| 3\*8+2 / 1\*8-4  | 1878/1070 |
| 3\*8+2 / 1\*8-5  | 1901/1045 |
| 3\*8+2 / 1\*8-6  | 1865/1056 |
| 3\*8+2 / 1\*8-7  | 1864/1079 |
| 3\*8+3 / 1\*8-6  | 2152/1179 |
| 3\*8+1 / 1\*8-6  | 1921/1285 |
| 3\*8-23 / 1\*8-7 | stack ovr |
| 3\*8-22 / 1\*8-6 | 8072/8726 |
| 3\*8-21 / 1\*8-5 | 6381/6659 |
| 3\*8-20 / 1\*8-4 | 4756/4725 |
| 3\*8-19 / 1\*8-3 | 4022/3695 |
| 3\*8-18 / 1\*8-2 | 3543/3102 |
| 3\*8-17 / 1\*8-1 | 3830/3060 |

## 8'000'000 records with test hash 2d array on a heap

| Tree size     | c         | c++       | c#          | go       | rust      |
|---------------|-----------|-----------|-------------|----------|-----------|
| 3\*8+2 / 1\*8-6 | 2341/1223 | 1831/1016 | 29558/38409 | 2472/830 | 2990/1919 |

## 16'000'000 records with test hash 2d array on a heap

| Tree size       | c++       | rust     |
|-----------------|-----------|----------|
| 3\*8+2 / 1\*8-6 | 4234/2700 |          |
| 3\*8+7 / 1\*8-4 |           | mem over |

## 32'000'000 records with test hash 2d array on a heap (native hash maps)

| Native     | c++         | c#         | go         | rust       |
|------------|-------------|------------|------------|------------|
| algorithms | 28591/12456 | 9861/15098 | 10925/5155 | 16700/9936 |

## 32'000'000 records with test hash 2d array on a heap

| Tree size       | c          | c++         | c#            | go         | rust       |
|-----------------|------------|-------------|---------------|------------|------------|
| 3\*8+2 / 1\*8-6 |            | 22717/22665 |               |            |            |
| 3\*8+1 / 1\*8-6 | 13557/8156 |             |               |            | 9353/11247 |
| 3\*8+2 / 1\*8-6 | __10470/6103__ | 22844/22823 |               |            | 9044/9615  |
| 3\*8+3 / 1\*8-7 | __10173/5845__ |             |               |            | 9429/8911  |
| 3\*8+3 / 1\*8-6 | __10204/5875__ | 19207/18580 |               | 12316/4920 | 9461/8662  |
| 3\*8+3 / 1\*8-5 | __10155/5862__ |             |               |            | 9507/8496  |
| 3\*8+3 / 1\*8-4 | __10181/5825__ |             |               |            | 10269/8460 |
| 3\*8+3 / 1\*8-3 |            |             |               |            | 11149/8435 |
| 3\*8+3 / 1\*8-0 | __10296/5841__ |             |               |            |            |
| 3\*8+4 / 1\*8-3 |            |             |               | __10988/4198__ | 12460/8008 |
| 3\*8+4 / 1\*8-4 |            |             |               | __10940/4020__ | __11895/8008__ |
| 3\*8+4 / 1\*8-5 |            |             |               | __10956/4049__ | __11597/8011__ |
| 3\*8+4 / 1\*8-6 | 10909/6640 | 18378/17201 | 155978/200232 | __11021/4064__ | 12074/8205 |
| 3\*8+4 / 1\*8-7 |            |             |               | 10920/4487 | 11438/8218 |
| 3\*8+5 / 1\*8-6 | 11487/6969 | 14940/13481 |               | 10359/4252 |            |
| 3\*8+6 / 1\*8-6 | 12371/7156 | 13441/9476  |               | 11398/4096 |            |
| 3\*8+7 / 1\*8-6 | 13706/7146 | 13887/8002  |               | 13486/3960 |            |
| 3\*8+7 / 1\*8-5 | 13774/7176 | __13435/7748__  |               |            |            |
| 3\*8+7 / 1\*8-4 | 13707/7138 | __13634/7650__  | mem over      | mem over   | mem over   |
| 3\*8+7 / 1\*8-3 | 13717/7140 | __13629/7660__  |               |            |            |
| 3\*8+7 / 1\*8-2 |            | mem over    |               |            |            |
| 3\*8+7 / 1\*8   |            | mem over    |               |            |            |

## 64'000'000 records with test hash 2d array on a heap

| Tree size       | c++        | go       |
|-----------------|------------|----------|
| 3\*8+2 / 1\*8-6 | stack over |          |
| 3\*8+2 / 1\*8   | mem over   |          |
| 3\*8+3 / 1\*8   | mem over   |          |
| 3\*8+4 / 1\*8   | mem over   |          |
| 3\*8+4 / 1\*8-4 |            | mem over |
| 3\*8+5 / 1\*8   | mem over   |          |
