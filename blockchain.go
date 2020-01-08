package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	bk "myblockchain/block"
	"myblockchain/consensus"
	"time"

	"github.com/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

// Blockchain keeps a sequence of Blocks
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// BlockchainIterator is used to iterate over blockchain blocks
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// AddBlock saves provided data as a block in the blockchain
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash) //改成pbft
	if newBlock == nil {
		return
	}

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})
}

// Iterator ...
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

// Next returns next block starting from the tip
func (i *BlockchainIterator) Next() *bk.Block {
	var block *bk.Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = bk.DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// NewBlock creates and returns Block
func NewBlock(data string, prevBlockHash []byte) *bk.Block {

	block := &bk.Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0} //区块结构
	// 添加区块流程 1 构建区块信息 2，广播 3，收到恢复后计算hash上链
	//发起请求
	sendreq()

	//这里需要阻塞一个tcp连接

	if handlereq() {
		block.Hash = gethash(block)
		block.Nonce = 0
		return block
	} else {
		fmt.Println("未得到回复，请求失败")
		return nil
	}

}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock() *bk.Block {
	data := "Genesis Block"
	var prevBlockHash []byte

	block := &bk.Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}

	block.Hash = gethash(block)
	//hash := consensus.gethash(block)
	block.Nonce = 0

	return block
	//return NewBlock("Genesis Block", []byte{})
}

func gethash(block *bk.Block) []byte {

	data := bytes.Join(
		[][]byte{
			block.PrevBlockHash,
			block.Data,
			consensus.IntToHex(block.Timestamp),
			consensus.IntToHex(int64(0)),
		},
		[]byte{},
	)
	hash := sha256.Sum256(data)
	return hash[:]
}
