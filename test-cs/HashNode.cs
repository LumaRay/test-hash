using System;
using System.Collections.Generic;
using System.Text;

namespace test_cs
{
	class HashNode<T>
	{
		public const int HASH_SIZE = 32;
		const int CHILD_NODES = 256;

		HashNode<T>[] mChildNodes;
		T mpObject;
		byte[] mHash;
		int mAddressBits;
		public HashNode()
		{
			mChildNodes = null;
			//mpObject = null;
			mHash = null;
			mAddressBits = 3 * 8 + 4;
		}

		HashNode(byte[] hash, T obj)
		{
			mChildNodes = null;
			mpObject = obj;
			mHash = hash;
			mAddressBits = 1 * 8 - 6;
		}

		public void SetObject(ref byte[] hash, ref T obj, int iByteOffset = 0) 
		{
			if (mChildNodes == null)
				mChildNodes = new HashNode<T>[(UInt64)1 << mAddressBits];
			if (mHash != null && mpObject != null) {
				byte[] existingHash = mHash;
				UInt32 iExistingHashAddress = ~(UInt32.MaxValue << mAddressBits) & (UInt32)existingHash[iByteOffset];
				mChildNodes[iExistingHashAddress] = new HashNode<T>(mHash, mpObject);
				mHash = null;
				//mpObject = null;
			}
			UInt32 iHashAddress = ~(UInt32.MaxValue << mAddressBits) & (UInt32)hash[iByteOffset];
			HashNode<T> childNode = mChildNodes[iHashAddress];
			if (childNode == null || iByteOffset == HASH_SIZE - 1) 
			{
				mChildNodes[iHashAddress] = new HashNode<T>(hash, obj);
			}
			else
			{
				childNode.SetObject(ref hash, ref obj, iByteOffset + (int)Math.Ceiling((float)mAddressBits / 8));
			}
		}

		public ref T GetObject(ref byte[] hash, int iByteOffset = 0) 
		{
			if (mHash != null && mpObject != null)
			{
				bool eq = true;
				for (int i = iByteOffset; i < HASH_SIZE; i++)
                {
					if (mHash[i] != hash[i])
                    {
						eq = false;
						break;
                    }
                }
				if (eq)
					return ref mpObject;
			}
			UInt32 iHashAddress = ~(UInt32.MaxValue << mAddressBits) & (UInt32)hash[iByteOffset];
			HashNode<T> childNode = mChildNodes[iHashAddress];
			if (childNode == null || iByteOffset + (int) Math.Ceiling((float) mAddressBits / 8) >= HASH_SIZE - 1) {
				throw new Exception("Errror");
			}
			return ref childNode.GetObject(ref hash, iByteOffset + (int)Math.Ceiling((float)mAddressBits / 8));
		}
	}
}
