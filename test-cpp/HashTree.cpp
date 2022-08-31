#include "HashTree.h"

template<class T>
void HashTree<T>::Add(byte* const& hash, T const& object) {
	byte* const hashCopy = new byte[HASH_SIZE];
	copy(hash, hash + HASH_SIZE, hashCopy);
	T* const pObjectCopy = new T(object);
	mRootNode.SetObject(hashCopy, pObjectCopy);
}

template<class T>
T& HashTree<T>::operator[] (byte* const& hash) {
	return mRootNode.GetObject(hash);
}
