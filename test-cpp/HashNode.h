#pragma once

#define HASH_SIZE	32
#define CHILD_NODES	256

#include <memory>

using namespace std;

typedef unsigned char byte;

template<class T>
class HashNode;

template<class T>
class HashNode {
	T* mpObject;
	byte* mpHash;
	//unique_ptr<HashNode<T>> mChildNodes[CHILD_NODES];
	//HashNode<T>* mChildNodes[CHILD_NODES];
	HashNode<T>** mChildNodes;
	int mAddressBits;
public:
	HashNode();

	HashNode(byte* const& pHash, T* const& pObject);

	HashNode(byte*&& rpHash, T&& rpObject) noexcept;

	~HashNode();

	void SetObject(byte* const& pHash, T* const& pObject, const int iByteOffset = 0);

	T& GetObject(byte* const& pHash, const int iByteOffset = 0);
};

#include "HashNode.cpp"