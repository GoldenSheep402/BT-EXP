package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/GoldenSheep402/BT-EXP/utils/tracker"
	"os"

	"github.com/anacrolix/torrent/bencode"
)

type FileInfo struct {
	Length uint64   `bencode:"length"`
	Path   []string `bencode:"path"`
}

type BencodeInfo struct {
	Files       *[]FileInfo `bencode:"files,omitempty"`
	Name        string      `bencode:"name"`
	Length      *uint64     `bencode:"length,omitempty"`
	Md5sum      *string     `bencode:"md5sum,omitempty"`
	Pieces      string      `bencode:"pieces"`
	PieceLength uint64      `bencode:"piece length"`
	Private     *int        `bencode:"private,omitempty"`
	Source      *string     `bencode:"source,omitempty"`
}

type BencodeTorrent struct {
	Announce  string      `bencode:"announce"`
	CreatedBy *string     `bencode:"created by,omitempty"`
	CreatedAt *int        `bencode:"creation date,omitempty"`
	Info      BencodeInfo `bencode:"info"`
}

func readAndPrintTorrent(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := bencode.NewDecoder(file)
	bencodeTorrent := &BencodeTorrent{}
	err = decoder.Decode(bencodeTorrent)
	if err != nil {
		return err
	}

	fmt.Printf("Announce URL: %s\n", bencodeTorrent.Announce)
	if bencodeTorrent.CreatedBy != nil {
		fmt.Printf("Created By: %s\n", *bencodeTorrent.CreatedBy)
	}
	if bencodeTorrent.CreatedAt != nil {
		fmt.Printf("Creation Date: %d\n", *bencodeTorrent.CreatedAt)
	}
	fmt.Printf("Name: %s\n", bencodeTorrent.Info.Name)
	if bencodeTorrent.Info.Length != nil {
		fmt.Printf("Length: %d\n", *bencodeTorrent.Info.Length)
	}
	fmt.Printf("Piece Length: %d\n", bencodeTorrent.Info.PieceLength)
	if bencodeTorrent.Info.Private != nil {
		fmt.Printf("Private: %d\n", *bencodeTorrent.Info.Private)
	}

	marshaledInfo, err := bencode.Marshal(bencodeTorrent.Info)
	if err != nil {
		return err
	}
	hash := sha1.Sum(marshaledInfo)
	fmt.Printf("Info Hash (SHA1): %x\n", hash)

	return nil
}

func main() {
	// filePath := "/Users/GS/Downloads/test.torrent"
	//
	// err := readAndPrintTorrent(filePath)
	// if err != nil {
	// 	fmt.Printf("读取种子文件出错: %v\n", err)
	// }
	tracker.Init()
}
