#pragma once

#include "HashNode.h"

template<class T>
class HashTree {
	HashNode<T> mRootNode;
public:
	void Add(byte* const& hash, T const& object);

	T& operator[] (byte* const& hash);
};

#include "HashTree.cpp"