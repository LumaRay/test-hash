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

## 8'000'000 records

| Language | c | c++ | c# | go | rust |
|----------|---|-----|----|----|------|
|          |   |     |    |    |      |
|          |   |     |    |    |      |
|          |   |     |    |    |      |

## 16'000'000 records

| Language | c | c++ | c# | go | rust |
|----------|---|-----|----|----|------|
|          |   |     |    |    |      |
|          |   |     |    |    |      |
|          |   |     |    |    |      |

## 32'000'000 records

| Language | c | c++ | c# | go | rust |
|----------|---|-----|----|----|------|
|          |   |     |    |    |      |
|          |   |     |    |    |      |
|          |   |     |    |    |      |

## 64'000'000 records

| Language | c | c++ | c# | go | rust |
|----------|---|-----|----|----|------|
|          |   |     |    |    |      |
|          |   |     |    |    |      |
|          |   |     |    |    |      |
