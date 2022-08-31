using System;
using System.Collections.Generic;
using System.Text;

namespace test_cs
{
    class HashTree<T>
    {
        HashNode<T> mRootNode;

        public HashTree()
        {
            mRootNode = new HashNode<T>();
        }
        //public ref T this[byte[] hash]
        //{
        //    get { return ref Get(ref hash); }
        //    //set { Add(ref hash, ref value); }
        //}

        public ref T Get(ref byte[] hash)
        {
            return ref mRootNode.GetObject(ref hash);
        }

        public void Add(ref byte[] hash, ref T obj)
        {
            byte[] hashCopy = new byte[HashNode<T>.HASH_SIZE];
            Array.Copy(hash, hashCopy, HashNode<T>.HASH_SIZE);
            T pObjectCopy = obj;
            mRootNode.SetObject(ref hashCopy, ref pObjectCopy);
        }
    }
}
