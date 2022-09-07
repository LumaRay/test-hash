// ConsoleApplication1.cpp : Этот файл содержит функцию "main". Здесь начинается и заканчивается выполнение программы.
//

#include <iostream>
//#include <memory>
//#include <cstring>
#include <cassert>
#include "HashTree.h"

#include <map>
#include <unordered_map>

#include <chrono>

#include <iomanip>
#include <sstream>
#include <string>

//template class HashTree<int>;

using namespace std;

#define MAX_HASHES	(32 * 1024 * 1024)

//byte hashes[MAX_HASHES][HASH_SIZE] = { 0 };
//string hash_strings[MAX_HASHES];

byte (*hashes)[HASH_SIZE];
string* hash_strings;

/*string hexStr(const uint8_t* data, int len)
{
	stringstream ss;
	ss << hex;

	for (int i(0); i < len; ++i)
		ss << setw(2) << setfill('0') << (int)data[i];

	return ss.str();
}*/

void GenerateRandomHashes() {
	hashes = new byte[MAX_HASHES][HASH_SIZE];
	hash_strings = new string[MAX_HASHES];

	for (int idxHash = 0; idxHash < MAX_HASHES; idxHash++) {
		srand(static_cast<unsigned int>(time(0)) + idxHash);
//		for (int idxByte = 0; idxByte < HASH_SIZE; idxByte++)
//			//hashes[idxHash][idxByte] = rand() % 256;
//			hashes[idxHash][idxByte] = 1 + rand() % 255;
		for (int idxByte = 0; idxByte < HASH_SIZE; idxByte += 2)
			*(uint16_t*)&hashes[idxHash][idxByte] = 0x0101 | (uint16_t)rand();
		////hash_strings[idxHash] = hexStr(hashes[idxHash], HASH_SIZE);
		hashes[idxHash][(int)((float)HASH_SIZE / 2)] = 0;
		hash_strings[idxHash] = move(string((char*)hashes[idxHash]));
	}

	//cout << hexStr(hashes[0], HASH_SIZE);
}

void MakeHashTree(HashTree<int>& hashTree, int hashes_count) {
	for (int idxHash = 0; idxHash < hashes_count; idxHash++) {
		byte* hash = hashes[idxHash];
		hashTree.Add(hash, hash[0]);
	}
}

void TestHashTree(HashTree<int>& hashTree, int hashes_count) {
	for (int idxHash = 0; idxHash < hashes_count; idxHash++) {
		byte* hash = hashes[idxHash];
		byte b = hashTree[hash];
		byte hashByte = hash[0];
		if (hashTree[hash] != hash[0]) cout << "TestHashTree!" << endl;
	}
}

void MakeStdMap(map<string, int>& hashMap, int hashes_count) {
	for (int idxHash = 0; idxHash < hashes_count; idxHash++) {
		byte* hash = hashes[idxHash];
		hashMap[hash_strings[idxHash]] = hash[0];
	}
}

void TestStdMap(map<string, int>& hashMap, int hashes_count) {
	long l = 0;
	for (int idxHash = 0; idxHash < hashes_count; idxHash++) {
		byte* hash = hashes[idxHash];
		l += hash[0];
		if (hashMap[hash_strings[idxHash]] != hash[0]) cout << "TestStdMap!" << endl;
	}
	//cout << l << endl;
}

void MakeStdUMap(unordered_map<string, int>& hashMap, int hashes_count) {
	for (int idxHash = 0; idxHash < hashes_count; idxHash++) {
		byte* hash = hashes[idxHash];
		hashMap[hash_strings[idxHash]] = hash[0];
	}
}

void TestStdUMap(unordered_map<string, int>& hashMap, int hashes_count) {
	long l = 0;
	for (int idxHash = 0; idxHash < hashes_count; idxHash++) {
		byte* hash = hashes[idxHash];
		l += hash[0];
		if (hashMap[hash_strings[idxHash]] != hash[0]) cout << "TestStdUMap!" << endl;
	}
	//cout << l << endl;
}

void TestStdUMapWrongGet(unordered_map<string, int>& hashMap, int hashes_count) {
	long l = 0;
	byte hash[HASH_SIZE] = { 0 };
	srand((uint32_t)time(0));
	for (int idxHash = 0; idxHash < hashes_count; idxHash++) {
		l += hash[0];
		for (int idxByte = 0; idxByte < HASH_SIZE; idxByte++)
			hash[idxByte] = rand() % 256;
		hash[(int)((float)HASH_SIZE / 2)] = 0;
		//if (hashMap.find(hexStr(hash, HASH_SIZE)) != hashMap.end()) cout << "TestStdUMapWrongGet!" << endl;
		if (hashMap.find((char*)hash) != hashMap.end()) cout << "TestStdUMapWrongGet!" << endl;
	}
	//cout << l << endl;
}

int main(int argc, char* argv[]) {

	if (argc < 3) return -1;

	cout << "Generating hashes...";

	GenerateRandomHashes();

	cout << "Done!" << endl;

	int hashes_count = atoi(argv[1]);

	if (hashes_count > MAX_HASHES) return -1;

	for (int i = 2; i < argc; i++) switch (argv[i][0]) {
	case 'h': {
		auto start = chrono::high_resolution_clock::now();
		HashTree<int> hashTree;
		MakeHashTree(hashTree, hashes_count);
		cout << "MakeHashTree "
			<< chrono::duration_cast<chrono::milliseconds>(chrono::high_resolution_clock::now() - start).count()
			<< endl;

		start = chrono::high_resolution_clock::now();
		TestHashTree(hashTree, hashes_count);
		cout << "TestHashTree "
			<< chrono::duration_cast<chrono::milliseconds>(chrono::high_resolution_clock::now() - start).count()
			<< endl;
		break;
	}
	case 'm': {
		map<string, int> hashMap;

		auto start = chrono::high_resolution_clock::now();
		MakeStdMap(hashMap, hashes_count);
		cout << "MakeStdMap "
			<< chrono::duration_cast<chrono::milliseconds>(chrono::high_resolution_clock::now() - start).count()
			<< endl;

		start = chrono::high_resolution_clock::now();
		TestStdMap(hashMap, hashes_count);
		cout << "TestStdMap "
			<< chrono::duration_cast<chrono::milliseconds>(chrono::high_resolution_clock::now() - start).count()
			<< endl;
		break;
	}
	case 'u': {
		unordered_map<string, int> hashMap;

		auto start = chrono::high_resolution_clock::now();
		MakeStdUMap(hashMap, hashes_count);
		cout << "MakeStdUMap "
			<< chrono::duration_cast<chrono::milliseconds>(chrono::high_resolution_clock::now() - start).count()
			<< endl;

		start = chrono::high_resolution_clock::now();
		TestStdUMap(hashMap, hashes_count);
		cout << "TestStdUMap "
			<< chrono::duration_cast<chrono::milliseconds>(chrono::high_resolution_clock::now() - start).count()
			<< endl;

		TestStdUMapWrongGet(hashMap, hashes_count);
		break;
	}
	}

	delete [] hashes;
	//+++delete [] hash_strings;
	return 0;
}

//32'768

//Golang
//1*8+3 1*8-2
//11
//2

//Golang+Postgres
//15380
//9861

//Golang+Redis
//3389
//3398

//Golang+Memcached
//3180
//3328

//Golang+MongoDB-no-index
//9893
//245567

//Golang+MongoDB-with-index
//10611
//11382

//500'000

//1756
//93

//1816
//91

//1746
//89

//1785
//88

//1328  child nodes unique_ptr -> *
//92

//1341  all unique_ptr -> *
//90

//1072
//112

//355
//102


//8'000'000

//3751
//2230

//1 1
//3835
//2367

//2 1
//3771
//2265

//2 2 mem overflow

//3 1
//2911
//1559

//1 2 mem overflow

//3*8+1 1
//2715
//1429

//3*8+2 1
//2025
//1108

//3*8+3 1
//2222
//1215

//3*8+4 1
//2394
//1246

//3*8+5 1
//3019
//1363

//3*8+5 1
//4051
//1431

//3*8+2 1*8+1
//2211
//1127

//3*8+2 1*8-1
//1996
//1136

//3*8+2 1*8-2 +
//1902 1901 1922
//1076 1060 1068

//3*8+2 1*8-3
//1900 1916 1890
//1049 1103 1065

//3*8+2 1*8-4 ++
//1893 1889 1878
//1059 1072 1070

//3*8+2 1*8-5
//1859 1909 1901
//1061 1091 1045

//3*8+2 1*8-6 +++
//1876 1862 1865  heap 1888 1831   с 2432 2349 2341   rust 2990   go 2472  cs 29814 29558
//1066 1056 1056  heap 1031 1016   с 1233 1259 1223   rust 1919   go 830   cs 38709 38409

//3*8+2 1*8-7
//1864 1928 1911 1888
//1079 1082 1073 1111

//3*8+3 1*8-6
//2152
//1179

//3*8+1 1*8-6
//1921 1942
//1285 1288

//3*8-23 1*8-7 hash overflow

//3*8-22 1*8-6
//8072 8338
//8726 8952

//3*8-21 1*8-5
//6381 6473
//6659 6688

//3*8-20 1*8-4
//4842 4756
//4852 4725

//3*8-19 1*8-3
//4022
//3695

//3*8-18 1*8-2
//3543
//3102

//3*8-17 1*8-1
//3830
//3060




//16'000'000 heap

//3*8+2 1*8-6 +++
//4234
//2700

//3*8+7 1*8 - 4 +++
//rust oom
//


//32'000'000 heap


//c++ umap 28686 28591
//c++ umap 12873 12456

//rust HashMap 16700
//rust HashMap 9936

//go-map 10925
//go-map 5155

//csd 9861
//csd 15098

//3*8+2 1*8-6
//9169
//6967

//22717 23009
//22665 23026

//3*8+1 1*8 - 6
// с 13557    rust 9353
// с 8156     rust 11247

//3*8+2 1*8 - 6
//22844    с 10470   rust 9044
//22823    с 6103    rust 9615

//3*8+3 1*8 - 7
// с 10173 +++   rust 9429
// с 5845        rust 8911

//3*8+3 1*8 - 6
//19207    с 10204    rust 9461   go 12316
//18580    с 5875     rust 8662   go 4920

//3*8+3 1*8 - 5
// с 10155   rust 9507 +++
// с 5862    rust 8496

//3*8+3 1*8 - 4
// с 10181    rust 10269
// с 5825     rust 8460

//3*8+3 1*8 - 3
//rust 11149
//rust 8435

//3*8+3 1*8 - 0
// с 10296
// с 5841

//3*8+4 1*8 - 3
//rust 12460   go 10988
//rust 8008    go 4198

//3*8+4 1*8 - 4
//rust 11895   go 10940 +++  
//rust 8008    go 4020

//3*8+4 1*8 - 5
//rust 11597   go 10956
//rust 8011    go 4049

//3*8+4 1*8 - 6
//18378    с 10909   rust 12074   go 11021    10999  cs 155978   
//17201    с 6640    rust 8205    go 4064+++  4098   cs 200232   

//3*8+4 1*8 - 7
//rust 11438   go 10920
//rust 8218    go 4487

//3*8+5 1*8 - 6
//14940    с 11487   go 10359
//13481    с 6969    go 4252

//3*8+6 1*8 - 6
//13441    с 12371   go 11398
//9476     с 7156    go 4096

//3*8+7 1*8 - 6
//13887    с 13787 13706   go 13486
//8002     с 7173  7146    go 3960

//3*8+7 1*8 - 5
//13435    с 13774
//7748     с 7176

//3*8+7 1*8 - 4 +++
//13634    с 13714 13728 13707    rust oom   go oom  cs oom
//7650     с 7141  7148  7138

//3*8+7 1*8 - 3
//13629    с 13717
//7660     с 7140

//3*8+7 1*8 - 2 oom

//3*8+7 1*8 oom


//64'000'000 heap

//3*8+2 1*8-6 +++ hash overflow
//55798

//3*8+2 1*8 oom

//3*8+3 1*8 oom

//3*8+4 1*8 oom

//3*8+4 1*8 - 4
//go oom
//go 

//3*8+5 1*8 oom
