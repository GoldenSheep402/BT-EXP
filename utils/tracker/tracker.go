package tracker

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/GoldenSheep402/BT-EXP/utils/rds"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zeebo/bencode"
)

type trackerReq struct {
	InfoHash      string // 必需
	PeerID        string // 必需
	Port          int    // 必需
	Uploaded      int    // 必需
	Downloaded    int    // 必需
	Left          int    // 必需
	Event         string // 可选
	IP            string // 可选
	NumWant       int    // 可选
	Key           string // 可选
	Compact       int    // 可选
	NoPeerID      int    // 可选
	TrackerID     string // 可选
	Corrupt       int    // 非标准参数
	SupportCrypto int    // 非标准参数
	Redundant     int    // 非标准参数
}

type trackerResp struct {
	FailureReason  string      `bencode:"failure reason,omitempty"`  // 可选
	WarningMessage string      `bencode:"warning message,omitempty"` // 可选
	Interval       int         `bencode:"interval"`                  // 必需
	MinInterval    int         `bencode:"min interval,omitempty"`    // 可选
	TrackerID      string      `bencode:"tracker id,omitempty"`      // 可选
	Complete       int         `bencode:"complete"`                  // 必需
	Incomplete     int         `bencode:"incomplete"`                // 必需
	Peers          interface{} `bencode:"peers"`                     // 必需
}

type Peer struct {
	PeerID   string `bencode:"peer id"` // 对等节点的唯一标识符
	IP       string `bencode:"ip"`      // 对等节点的 IP 地址
	Port     int    `bencode:"port"`    // 对等节点的端口号
	LastSeen time.Time
}

var RDB *redis.Client
var ctx = context.Background()

func Init() {
	var err error
	RDB, err = rds.CreateClient()
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.GET("/tracker", handleTrackerRequest)
	r.Run(":18312")
}

func AddPeerToRedis(infoHash string, peer *Peer) error {
	key := "torrent:" + infoHash
	field := peer.PeerID
	value := fmt.Sprintf("%s:%d:%d", peer.IP, peer.Port, time.Now().Unix())
	err := RDB.HSet(ctx, key, field, value).Err()
	if err != nil {
		return err
	}

	RDB.Expire(ctx, key, time.Hour)
	return nil
}

func GetPeersFromRedis(infoHash string, numWant int) ([]*Peer, error) {
	key := "torrent:" + infoHash
	peersData, err := RDB.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	peers := make([]*Peer, 0, len(peersData))
	now := time.Now().Unix()
	for peerID, data := range peersData {
		parts := strings.Split(data, ":")
		if len(parts) != 3 {
			continue
		}
		ip := parts[0]
		port, _ := strconv.Atoi(parts[1])
		lastSeen, _ := strconv.ParseInt(parts[2], 10, 64)
		if now-lastSeen > 3600 {
			RDB.HDel(ctx, key, peerID)
			continue
		}
		peers = append(peers, &Peer{
			PeerID:   peerID,
			IP:       ip,
			Port:     port,
			LastSeen: time.Unix(lastSeen, 0),
		})
		if len(peers) >= numWant {
			break
		}
	}
	return peers, nil
}

func handleTrackerRequest(c *gin.Context) {
	req := trackerReq{}

	infoHash := c.Query("info_hash")
	req.PeerID = c.Query("peer_id")
	portStr := c.Query("port")
	uploadedStr := c.Query("uploaded")
	downloadedStr := c.Query("downloaded")
	leftStr := c.Query("left")
	req.Event = c.Query("event")
	req.IP = c.Query("ip")
	numWantStr := c.Query("numwant")
	req.Key = c.Query("key")
	compactStr := c.Query("compact")
	noPeerIDStr := c.Query("no_peer_id")
	req.TrackerID = c.Query("trackerid")
	corruptStr := c.Query("corrupt")
	supportCryptoStr := c.Query("supportcrypto")
	redundantStr := c.Query("redundant")

	var err error
	decodedHash, err := url.QueryUnescape(infoHash)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid info_hash")
		return
	}

	infoHashBytes := []byte(decodedHash)
	hexString := fmt.Sprintf("%x", infoHashBytes)

	req.InfoHash = hexString

	fmt.Printf("解码后的 info_hash: %x\n", decodedHash)

	if req.Port, err = strconv.Atoi(portStr); err != nil {
		c.String(http.StatusBadRequest, "Invalid port")
		return
	}
	if req.Uploaded, err = strconv.Atoi(uploadedStr); err != nil {
		c.String(http.StatusBadRequest, "Invalid uploaded")
		return
	}
	if req.Downloaded, err = strconv.Atoi(downloadedStr); err != nil {
		c.String(http.StatusBadRequest, "Invalid downloaded")
		return
	}
	if req.Left, err = strconv.Atoi(leftStr); err != nil {
		c.String(http.StatusBadRequest, "Invalid left")
		return
	}
	if numWantStr != "" {
		req.NumWant, _ = strconv.Atoi(numWantStr)
	} else {
		req.NumWant = 50 // 默认值
	}
	if compactStr != "" {
		req.Compact, _ = strconv.Atoi(compactStr)
	} else {
		req.Compact = 1 // 默认值，返回紧凑形式
	}
	if noPeerIDStr != "" {
		req.NoPeerID, _ = strconv.Atoi(noPeerIDStr)
	}
	if corruptStr != "" {
		req.Corrupt, _ = strconv.Atoi(corruptStr)
	}
	if supportCryptoStr != "" {
		req.SupportCrypto, _ = strconv.Atoi(supportCryptoStr)
	}
	if redundantStr != "" {
		req.Redundant, _ = strconv.Atoi(redundantStr)
	}

	if req.IP == "" {
		req.IP = c.ClientIP()
	}

	switch req.Event {
	case "started", "":
		peer := &Peer{
			PeerID: req.PeerID,
			IP:     req.IP,
			Port:   req.Port,
		}
		err = AddPeerToRedis(req.InfoHash, peer)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to add peer")
			return
		}
	case "stopped":
		key := "torrent:" + req.InfoHash
		RDB.HDel(ctx, key, req.PeerID)
		c.String(http.StatusOK, "")
		return
	case "completed":
	default:
	}

	peers, err := GetPeersFromRedis(req.InfoHash, req.NumWant)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to get peers")
		return
	}

	resp := trackerResp{
		Interval:   1800, // 示例值，客户端应等待的秒数
		Complete:   0,    // Seeder 数量
		Incomplete: 0,    // Leecher 数量
	}

	if req.Left == 0 {
		resp.Complete++
	} else {
		resp.Incomplete++
	}

	if req.Compact == 1 {
		var peersData []byte
		for _, peer := range peers {
			ip := net.ParseIP(peer.IP).To4()
			if ip == nil {
				continue
			}
			portBytes := make([]byte, 2)
			binary.BigEndian.PutUint16(portBytes, uint16(peer.Port))
			peersData = append(peersData, ip...)
			peersData = append(peersData, portBytes...)
		}
		resp.Peers = peersData
	} else {
		var peersList []Peer
		for _, peer := range peers {
			if req.NoPeerID == 1 {
				peer.PeerID = ""
			}
			peersList = append(peersList, *peer)
		}
		resp.Peers = peersList
	}

	encodedResp, err := bencode.EncodeBytes(resp)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to encode response")
		return
	}
	c.Data(http.StatusOK, "text/plain", encodedResp)
}
