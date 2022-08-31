#include "HashNode.h"
#include <cmath>

template<class T>
HashNode<T>::HashNode() 
	: mChildNodes{ nullptr }
	, mpObject(nullptr)
	, mpHash(nullptr)
	, mAddressBits(3*8 + 7) 
{
	//mChildNodes = new HashNode<T>*[CHILD_NODES];
	//cout << m_childNodes[10];
}

template<class T>
HashNode<T>::HashNode(byte* const& pHash, T* const& pObject) 
	: mChildNodes{ nullptr }
	, mpObject(pObject)
	, mpHash(pHash)
	, mAddressBits(1*8 - 5) 
{
	//mChildNodes = new (HashNode<T>*)[CHILD_NODES];
	//cout << m_childNodes[10];
}

//template<class T>
//HashNode<T>::HashNode(byte*&& rpHash, T&& rpObject) noexcept {
//	mpHash = move(rpHash);
//	mpObject = move(rpObject);
//}

template<class T>
HashNode<T>::~HashNode() {
	if (mpObject) {
		delete mpObject;
	}

	if (mpHash) {
		delete [] mpHash;
	}

	if (mChildNodes) {
		for (int idxChildNode = 0; idxChildNode < 1 << mAddressBits; idxChildNode++)
			if (mChildNodes[idxChildNode])
				delete mChildNodes[idxChildNode];
		delete[] mChildNodes;
	}
}

template<class T>
void HashNode<T>::SetObject(byte* const& pHash, T* const& pObject, const int iByteOffset) {
	if (!mChildNodes)
		mChildNodes = new HashNode<T>*[(uint64_t)1 << mAddressBits]();
	if (mpHash && mpObject) {
		//T* pExistingObject = mpObject.get();
		byte* pExistingHash = mpHash;
		//int iExistingHashAddress = 0;
		//memcpy(&iExistingHashAddress, &pExistingHash[iByteOffset], mAddressBytes);
		int iExistingHashAddress = ~((uint32_t)-1 << mAddressBits) & *((uint32_t*)&pExistingHash[iByteOffset]);
		//unique_ptr<HashNode<T>>* ppExistingChildNode = &mChildNodes[ucExistingHashByte];
		//auto pExistingNewChildNode = unique_ptr<HashNode<T>>(new HashNode<T>(pExistingHash, pExistingObject));
		//auto pExistingNewChildNode = unique_ptr<HashNode<T>>(new HashNode<T>(move(mpHash), move(mpObject)));
		//*ppExistingChildNode = move(pExistingNewChildNode);
		//mChildNodes[ucExistingHashByte] = unique_ptr<HashNode<T>>(new HashNode<T>(move(mpHash), move(mpObject)));
		mChildNodes[iExistingHashAddress] = new HashNode<T>(mpHash, mpObject);
		//mpHash.release();
		//mpObject.release();
		mpHash = nullptr;
		mpObject = nullptr;
	}
	//int iHashAddress = 0;
	//memcpy(&iHashAddress, &pHash[iByteOffset], mAddressBytes);
	int iHashAddress = ~((uint32_t)-1 << mAddressBits) & *((uint32_t*)&pHash[iByteOffset]);
	HashNode<T>** ppChildNode = &mChildNodes[iHashAddress];
	if (*ppChildNode == nullptr || iByteOffset == HASH_SIZE - 1) {
		//auto pNewChildNode = unique_ptr<HashNode<T>>(new HashNode<T>(pHash, pObject));
		//*ppChildNode = move(pNewChildNode);
		*ppChildNode = new HashNode<T>(pHash, pObject);
	}
	else {
		(*ppChildNode)->SetObject(pHash, pObject, iByteOffset + (int)ceil((float)mAddressBits / 8));
	}
}

template<class T>
T& HashNode<T>::GetObject(byte* const& pHash, const int iByteOffset) {
	if (mpHash && mpObject) if (memcmp(mpHash + iByteOffset, pHash + iByteOffset, (size_t)HASH_SIZE - iByteOffset) == 0) {
		return *mpObject;
	}
	//int iHashAddress = 0;
	//memcpy(&iHashAddress, &pHash[iByteOffset], mAddressBytes);
	int iHashAddress = ~((uint32_t)-1 << mAddressBits) & *((uint32_t*)&pHash[iByteOffset]);
	HashNode<T>** ppChildNode = &mChildNodes[iHashAddress];
	if (*ppChildNode == nullptr || iByteOffset + (int)ceil((float)mAddressBits / 8) >= HASH_SIZE - 1) {
		throw -1;
	}
	return (*ppChildNode)->GetObject(pHash, iByteOffset + (int)ceil((float)mAddressBits / 8));
}
