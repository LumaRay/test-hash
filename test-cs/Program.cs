using System;
using System.Collections.Generic;

namespace test_cs
{
    class Program
    {
		const int MAX_HASHES = 32 * 1024 * 1024;
		static byte[][] hashes;

		static void GenerateRandomHashes()
		{
			hashes = new byte[MAX_HASHES][];

			Random rnd = new Random();
			for (int idxHash = 0; idxHash < MAX_HASHES; idxHash++)
			{
				hashes[idxHash] = new byte[HashNode<int>.HASH_SIZE];
				rnd.NextBytes(hashes[idxHash]);
				//BitConverter.GetBytes()
				//for (int idxByte = 0; idxByte < HashNode<int>.HASH_SIZE; idxByte += 2)
				//	*(uint16_t*)&hashes[idxHash][idxByte] = 0x0101 | (uint16_t)rand();
			}
		}

		static void MakeHashTree(HashTree<int> hashTree)
		{
			for (int idxHash = 0; idxHash < MAX_HASHES; idxHash++)
			{
				byte[] hash = hashes[idxHash];
				int obj = (int)hash[0];
				hashTree.Add(ref hash, ref obj);
			}
		}

		static void TestHashTree(HashTree<int> hashTree)
		{
			byte[] hash;
			for (int idxHash = 0; idxHash < MAX_HASHES; idxHash++)
			{
				hash = hashes[idxHash];
				//ref int b = ref hashTree.Get(ref hash);//[hash];//.Get(ref hash);//[hash];
				//byte hashByte = hash[0];
				if (hashTree.Get(ref hash) != (int)hashes[idxHash][0]) 
					Console.WriteLine("TestHashTree!");
			}
		}

		static void MakeDict(ref Dictionary<byte[], int> hashMap)
		{
			for (int idxHash = 0; idxHash < MAX_HASHES; idxHash++)
			{
				byte[] hash = hashes[idxHash];
				hashMap[hash] = hash[0];
			}
		}

		static void TestDict(ref Dictionary<byte[], int> hashMap)
		{
			long l = 0;
			for (int idxHash = 0; idxHash < MAX_HASHES; idxHash++)
			{
				byte[] hash = hashes[idxHash];
				l += hash[0];
				if (hashMap[hash] != hash[0]) 
					Console.WriteLine("TestDict!");
			}
		}

		static void TestDictWrongGet(ref Dictionary<byte[], int> hashMap)
		{
			Random rnd = new Random();
			for (int idxHash = 0; idxHash < MAX_HASHES; idxHash++)
			{
				var hash = new byte[HashNode<int>.HASH_SIZE];
				rnd.NextBytes(hash);
				if (hashMap.ContainsKey(hash)) 
					Console.WriteLine("TestStdUMapWrongGet!");
			}
		}

		static void Main(string[] args)
        {
			Console.Write("Generating hashes...");

			GenerateRandomHashes();

			Console.WriteLine("Done!");

			HashTree<int> hashTree = new HashTree<int>();

			var watch = new System.Diagnostics.Stopwatch();

			//watch.Start();
			//MakeHashTree(hashTree);
			//watch.Stop();
			//Console.WriteLine($"MakeHashTree: {watch.ElapsedMilliseconds}");

			//watch.Start();
			//TestHashTree(hashTree);
			//watch.Stop();
			//Console.WriteLine($"TestHashTree: {watch.ElapsedMilliseconds}");

			var hashMap = new Dictionary<byte[], int>();

			watch.Start();
			MakeDict(ref hashMap);
			watch.Stop();
			Console.WriteLine($"MakeDict: {watch.ElapsedMilliseconds}");

			watch.Start();
			TestDict(ref hashMap);
			watch.Stop();
			Console.WriteLine($"TestDict: {watch.ElapsedMilliseconds}");

			TestDictWrongGet(ref hashMap);
		}
    }
}
