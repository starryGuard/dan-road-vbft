package vbft

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"

	"dan-road-vbft/common"
	"dan-road-vbft/consensus/vbft/config"
	"dan-road-vbft/core/types"
	goSdk "github.com/ontio/ontology-go-sdk"
)

func TestBlock_getProposer(t *testing.T) {
	blk, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed: %v", err)
	}
	type fields struct {
		Block *types.Block
		Info  *vconfig.VbftBlockInfo
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name:   "test",
			fields: fields{Block: blk.Block, Info: blk.Info},
			want:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blk := &Block{
				Block: tt.fields.Block,
				Info:  tt.fields.Info,
			}
			if got := blk.getProposer(); got != tt.want {
				t.Errorf("Block.getProposer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlock_getBlockNum(t *testing.T) {
	blk, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed: %v", err)
	}
	type fields struct {
		Block *types.Block
		Info  *vconfig.VbftBlockInfo
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name:   "test",
			fields: fields{Block: blk.Block, Info: blk.Info},
			want:   uint32(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blk := &Block{
				Block: tt.fields.Block,
				Info:  tt.fields.Info,
			}
			if got := blk.getBlockNum(); got != tt.want {
				t.Errorf("Block.getBlockNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlock_getPrevBlockHash(t *testing.T) {
	blk, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed: %v", err)
	}
	type fields struct {
		Block *types.Block
		Info  *vconfig.VbftBlockInfo
	}
	tests := []struct {
		name   string
		fields fields
		want   common.Uint256
	}{
		{
			name:   "test",
			fields: fields{Block: blk.Block, Info: blk.Info},
			want:   common.Uint256{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blk := &Block{
				Block: tt.fields.Block,
				Info:  tt.fields.Info,
			}
			if got := blk.getPrevBlockHash(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Block.getPrevBlockHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlock_getLastConfigBlockNum(t *testing.T) {
	blk, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed: %v", err)
	}

	type fields struct {
		Block *types.Block
		Info  *vconfig.VbftBlockInfo
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name:   "test",
			fields: fields{Block: blk.Block, Info: blk.Info},
			want:   uint32(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blk := &Block{
				Block: tt.fields.Block,
				Info:  tt.fields.Info,
			}
			if got := blk.getLastConfigBlockNum(); got != tt.want {
				t.Errorf("Block.getLastConfigBlockNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlock_getNewChainConfig(t *testing.T) {
	blk, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed: %v", err)
	}
	type fields struct {
		Block *types.Block
		Info  *vconfig.VbftBlockInfo
	}
	tests := []struct {
		name   string
		fields fields
		want   *vconfig.ChainConfig
	}{
		{
			name:   "test",
			fields: fields{Block: blk.Block, Info: blk.Info},
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blk := &Block{
				Block: tt.fields.Block,
				Info:  tt.fields.Info,
			}
			if got := blk.getNewChainConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Block.getNewChainConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerialize(t *testing.T) {
	blk, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed: %v", err)
	}
	_, err = blk.Serialize()
	if err != nil {
		t.Errorf("Block Serialize failed :%v", err)
		return
	}
	t.Log("Block Serialize succ")
}

func TestInitVbftBlock(t *testing.T) {
	blk, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed: %v", err)
	}
	_, err = initVbftBlock(blk.Block)
	if err != nil {
		t.Errorf("initVbftBlock failed: %v", err)
		return
	}
	t.Log("TestInitVbftBlock succ")
}

func TestSignTX(t *testing.T) {
	sdk := goSdk.NewOntologySdk()
	//sdk.NewRpcClient().SetAddress("http://172.17.0.2:20336")
	wallet, err := sdk.OpenWallet("/Users/lixiaohan/Library/go/src/dan-road-vbft/wallet.dat")
	if err != nil {
		fmt.Println(err)
	}
	from, err := wallet.GetAccountByAddress("ANRryVJESVNodWvtVcrkLqMimzqGe7jBSZ", []byte("1"))
	to, err := wallet.GetAccountByAddress("AGna5UaixJTZcUkijCnM8n8hifmjCySjTc", []byte("2"))

	var f *os.File
	var err1 error
	for i := 0; i < 1000; i++ {
		tx, err := sdk.Native.Ont.NewTransferTransaction(500, 30000, from.Address, to.Address, 1)
		err = sdk.SignToTransaction(tx, from)
		txbf := new(bytes.Buffer)

		txt, err := tx.IntoImmutable()
		txt.Serialize(txbf)
		if err != nil {
			fmt.Println(err)
		}
		hexCode := common.ToHexString(txbf.Bytes())
		fmt.Println(hexCode)

		f, err1 = os.OpenFile("/Users/lixiaohan/Desktop/cvbft.txt", os.O_APPEND, 0666)
		w := bufio.NewWriter(f)
		fmt.Fprintln(w, hexCode)
		if err1 != nil {
			fmt.Println(err1)
		}

	}
	f.Close()
}
