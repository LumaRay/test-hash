use rand::prelude::*;
//use std::ops::{Index, IndexMut};
//use math::round;
use std::time::{SystemTime};//, UNIX_EPOCH};
use std::collections::HashMap;

const HASH_SIZE: usize = 32;

//const MAX_HASHES: usize = 8 * 1024;// * 1024;
//const MAX_HASHES: usize = 1 * 1024 * 1024;

// const MAX_HASHES: usize = 8 * 1024 * 1024;
// const ROOT_NODE_ADDRESS_BITS: u32 = 3*8 + 2;
// const CHILD_NODE_ADDRESS_BITS: u32 = 1*8 - 6;

const MAX_HASHES: usize = 32 * 1024 * 1024;
const ROOT_NODE_ADDRESS_BITS: u32 = 3*8 + 3;
const CHILD_NODE_ADDRESS_BITS: u32 = 1*8 - 5;

struct HashNode<T: Clone> {
	object: Option<T>,
	hash: Option<[u8; HASH_SIZE]>,//Option<*mut [u8; HASH_SIZE]>,
	child_nodes: Option<Vec<Option<HashNode<T>>>>,//Option<*mut [HashNode<T>]>,
	address_bits: u32
}

struct HashTree<T: Clone> {
	root_node: HashNode<T>
}

impl<T: Clone> HashTree<T> {
    fn create() -> HashTree<T> {
        HashTree::<T> { 
            root_node: HashNode::<T> {
                child_nodes: None,
                object: None,
                hash: None,
                address_bits: ROOT_NODE_ADDRESS_BITS
            }
        }
    }
    fn add(&mut self, hash: [u8; HASH_SIZE], object: T) {
        let opt_hash: Option<[u8; HASH_SIZE]> = Some(hash.clone());
        let opt_object: Option<T> = Some(object.clone());
        self.root_node.set_object(opt_hash, opt_object, 0);
    }
    fn get(&mut self, hash: [u8; HASH_SIZE]) -> T {
        self.root_node.get_object(hash, 0).unwrap()
    }
}

impl<T: Clone> HashNode<T> {
    fn create(hash: Option<[u8; HASH_SIZE]>, object: Option<T>) -> HashNode::<T> { 
        HashNode::<T> {
            child_nodes: None,
            object,
            hash,
            address_bits: CHILD_NODE_ADDRESS_BITS
        }
    }
    fn set_object(&mut self, hash: Option<[u8; HASH_SIZE]>, object: Option<T>, byte_offset: usize) {
        if self.child_nodes.is_none() {
            //this->mChildNodes = calloc((uint64_t)1 << this->mAddressBits, sizeof(struct HashNode));
            //self.child_nodes = Some(vec![None; 1 << self.address_bits]);
            self.child_nodes = Some(Vec::with_capacity(1 << self.address_bits));
            self.child_nodes.as_mut().unwrap().resize_with(1 << self.address_bits, || None::<HashNode<T>>);
        }
        if self.hash.is_some() && self.object.is_some() {
            let existing_hash = self.hash;
            let existing_hash_address:u32 = 
                !(u32::MAX << self.address_bits) 
                & u32::from_le_bytes(
                    existing_hash.unwrap()[byte_offset..byte_offset+4 as usize]
                    .try_into().expect("slice with incorrect length"));
            self.child_nodes.as_mut().unwrap()[existing_hash_address as usize] = 
                Some(HashNode::<T>::create(self.hash.take(), self.object.take()));
            // self.hash = None;
            // self.object = None;
        }
        let hash_address:u32 = !(u32::MAX << self.address_bits) 
            & u32::from_le_bytes(hash.unwrap()[byte_offset..byte_offset+4 as usize]
            .try_into().expect("slice with incorrect length"));
        let child_node = self.child_nodes.as_mut().unwrap()[hash_address as usize].as_ref();
        if child_node.is_none() || byte_offset == HASH_SIZE - 1 {
            self.child_nodes.as_mut().unwrap()[hash_address as usize] = Some(HashNode::create(hash, object));
        }
        else {
            self.child_nodes.as_mut().unwrap()[hash_address as usize].as_mut().unwrap().set_object(hash, object, 
                byte_offset + math::round::ceil((self.address_bits as f64) / 8.0, 0) as usize);
        }
    }
    fn get_object(&mut self, hash: [u8; HASH_SIZE], byte_offset: usize) -> Option<T> {
        if self.hash.is_some() && self.object.is_some() 
        && (self.hash.unwrap()[byte_offset..] == hash[byte_offset..]) {
            return self.object.clone();
        }
        let hash_address: u32 = !(u32::MAX << self.address_bits) 
            & u32::from_le_bytes(hash[byte_offset..byte_offset+4 as usize]
            .try_into().expect("slice with incorrect length"));
        let child_node = self.child_nodes.as_mut().unwrap()[hash_address as usize].as_ref();
        if child_node.is_none() 
        || byte_offset + math::round::ceil(self.address_bits as f64 / 8.0, 0) as usize >= HASH_SIZE - 1 {
            panic!("Hash overflow!");
            //return None;
        }
        return self.child_nodes.as_mut().unwrap()[hash_address as usize].as_mut().unwrap().get_object(hash, 
            byte_offset + math::round::ceil(self.address_bits as f64 / 8.0, 0) as usize);
    }
}

fn generate_random_hashes(hashes: &mut Vec<[u8; HASH_SIZE]>) {
    for idx_hash in 0..MAX_HASHES {
        //thread_rng().try_fill(&mut hashes[..])?;
        //rand::thread_rng().gen_range(0..256).fill(&mut hashes[idxHash]);
        for idx_byte in 0..HASH_SIZE {
            hashes[idx_hash][idx_byte] = random();
        }
        //	*(uint16_t*)&hashes[idxHash][idxByte] = (uint16_t)rand();
    }
}

fn make_hash_tree(hashes: &mut Vec<[u8; HASH_SIZE]>, hash_tree: &mut HashTree<i32>) {
	for idx_hash in 0..MAX_HASHES {
		let hash = hashes[idx_hash];
		let inp = hash[0];
		hash_tree.add(hash, inp as i32);
	}
}

fn test_hash_tree(hashes: &mut Vec<[u8; HASH_SIZE]>, hash_tree: &mut HashTree<i32>) {
	for idx_hash in 0..MAX_HASHES {
		let hash = hashes[idx_hash];
		//let b: i32 = hash_tree.get(hash);
		//let hash_byte: u8 = hash[0];
		if hash_tree.get(hash) != hash[0] as i32 {
            println!("TestHashTree!\n");
        }
	}
}

fn make_crate_hash_map(hashes: &mut Vec<[u8; HASH_SIZE]>, map: &mut HashMap<Vec<u8>, i32>) {
    for idx_hash in 0..MAX_HASHES {
        map.insert(hashes[idx_hash].to_vec(), hashes[idx_hash][0] as i32);
        //println!("make_crate_hash_map - insert {:?}", idx_hash);
    }
}

fn test_crate_hash_map(hashes: &mut Vec<[u8; HASH_SIZE]>, map: &mut HashMap<Vec<u8>, i32>) {
    for idx_hash in 0..MAX_HASHES {
        let hash = hashes[idx_hash];
        //println!("test_crate_hash_map - test {:?}", idx_hash);
		if map.get(&hash.to_vec()).unwrap() != &(hash[0] as i32) {
            println!("TestHashTree!\n");
        }
    }
}

fn main() {
    //[#Fill]
    let mut hashes = vec![[0u8; HASH_SIZE]; MAX_HASHES];

    print!("Generating hashes...");
    generate_random_hashes(&mut hashes);
    println!("Done {0}", hashes[4][0]);

    let mut hash_tree = HashTree::<i32>::create();

    let start = SystemTime::now();
	make_hash_tree(&mut hashes, &mut hash_tree);
    println!("MakeHashTree {:?}", SystemTime::now().duration_since(start).unwrap());
	
    let start = SystemTime::now();
	test_hash_tree(&mut hashes, &mut hash_tree);
    println!("TestHashTree {:?}", SystemTime::now().duration_since(start).unwrap());

    let mut crate_hash_map = HashMap::<Vec<u8>, i32>::new();
    
    let start = SystemTime::now();
    make_crate_hash_map(&mut hashes, &mut crate_hash_map);
    println!("make_crate_hash_map {:?}", SystemTime::now().duration_since(start).unwrap());

    let start = SystemTime::now();
    test_crate_hash_map(&mut hashes, &mut crate_hash_map);
    println!("test_crate_hash_map {:?}", SystemTime::now().duration_since(start).unwrap());
}
