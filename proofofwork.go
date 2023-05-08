package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

//难度，多少位
const targetBits = 4

// ProofOfWork represents a proof-of-work
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork builds and returns a ProofOfWork
/**
这个函数用于创建一个新的工作量证明实例 ProofOfWork。
它接受一个指向 Block 实例的指针 b 作为参数，并返回一个指向新创建对象的指针 ProofOfWork。
函数体内，首先创建了一个大整数实例 target，它的二进制表示为 00...01，长度为256个比特。
实际上，此处创建的 target 是难度阈值，用来衡量区块哈希的最大能力。
然后，使用 Lsh() 方法（左移）从 target 中去掉 targetBits 个比特，其余位将被设置为零。
targetBits 是一个整数值，标识当前的目标难度值。
最后，创建一个新的 ProofOfWork 实例，并使用指向 Block 的指针 b 和 target 作为其属性值。
*/
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// Run performs a proof-of-work
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	/**
	循环的每一步都会执行以下操作：

	调用 pow.prepareData() 函数准备区块数据 data；
	使用 sha256.Sum256() 方法计算整个数据的哈希值，并将结果返回到 hash 数组的固定长度中；
	将 hash 转化成整数值 hashInt；
	比较 hashInt 和 pow.target，如果 hashInt 小于 pow.target，则找到了一个合法的哈希值，否则继续增加随机值 nonce。
	*/
	for nonce < maxNonce {
		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validate validates block's PoW
//验证Pow是否正确，就用大整数来比较
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
