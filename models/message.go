/*
 * @Descripttion:
 * @Author:
 * @Date: 2023-05-08 08:06:37
 * @LastEditTime: 2023-07-08 11:10:02
 */
package models

import (
	"chatapp/utils"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FormId   int64  //发送者
	TargetId int64  //接收者
	Type     int    //发送类型 群聊 私聊 广播
	Media    int    //消息类型 文字 表情包 图片 音频
	Content  string //消息内容
	Pic      string //消息内容
	Url      string
	Desc     string
	Amount   int //其他数据统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn          *websocket.Conn
	Addr          string        //客户端地址
	FirstTime     uint64        //首次连接时间
	HeartbeatTime uint64        //心跳时间
	LoginTime     uint64        //登录时间
	DataQueue     chan []byte   //消息
	GroupSets     set.Interface //好友/群
}

//映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

//读写锁
var rwLocker sync.RWMutex

// 发送者ID 接收者ID 消息类型 发送内容 发送类型
func Chat(writer http.ResponseWriter, request *http.Request) {
	//校验token
	//token := query.Get("token")
	query := request.URL.Query()
	Id := query.Get("userId")

	//	targetId := query.Get("target")
	//	context := query.Get("context")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	isvalid := true //checkToken
	conn, err := (&websocket.Upgrader{
		//token 校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalid
		},
	}).Upgrade(writer, request, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取Conn
	node := &Node{
		Conn: conn,

		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// 用户关系

	// userid 绑定 node 并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	// 完成发送逻辑
	go sendProc(node)
	// 完成接受逻辑
	go recvProc(node)
	sendMsg(userId, []byte("欢迎进入聊天系统"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("<<<", data)
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udprecvProc()
}

// 完成udp数据发送携程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}

	for {
		select {
		case data := <-udpsendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成udp数据接收携程
func udprecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私信
		sendMsg(msg.TargetId, data)
	case 2: //群发
		//sendGroupMsg()
	case 3: //广播
	//	sendAllMsg()
	case 4:

	}
}

func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	jsonMsg := Message{}
	json.Unmarshal(msg, &jsonMsg)
	ctx := context.Background()
	targetIdStr := strconv.Itoa(int(userId))
	userIdStr := strconv.Itoa(int(jsonMsg.ID))

	r, err := utils.Red.Get(ctx, "online_"+userIdStr).Result()
	if err != nil {
		fmt.Println(err)
	}
	if r != "" {
		if ok {
			fmt.Println("sendMsg >>> userId :", userId, "msg :", string(msg))
		}
	}
	var key string
	if userId > jsonMsg.FormId {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}
	res, err := utils.Red.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
	}
	score := float64(cap(res)) + 1
	ress, e := utils.Red.ZAdd(ctx, key, &redis.Z{score, msg}).Result()
	if e != nil {
		fmt.Println(e)
	}

	if ok {
		fmt.Println(ress)
		node.DataQueue <- msg
	}
}

func JoinGroup(userId uint, comId string) (int, string) {
	contact := Contact{}
	contact.OwnerId = userId
	//contact.TargetId = comId
	contact.Type = 2

	community := Community{}
	utils.DB.Where("id = ? or name = ? ", comId, comId).Find(&community)

	if community.Name == "" {
		return -1, "没有找到群"
	}
	utils.DB.Where("owner_id = ? and target_id = ? and type = 2", userId, comId).Find(&contact)
	if !contact.CreatedAt.IsZero() {
		return -1, "已加过此群"
	} else {
		contact.TargetId = community.ID
		utils.DB.Create(&contact)
		return 0, "加群成功"
	}
}

func (msg Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(msg)
}

//获取缓存里的消息

func ReadRedisMsg(userIdA int64, userIdB int64, start, end int64, isRev bool) []string {
	rwLocker.RLock()
	//_, _ := clientMap[userIdA]
	rwLocker.RUnlock()

	ctx := context.Background()
	userIdStr := strconv.Itoa(int(userIdA))
	targetIdStr := strconv.Itoa(int(userIdB))
	var key string
	if userIdA > userIdB {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}
	var rels []string
	var err error
	if isRev {
		rels, err = utils.Red.ZRange(ctx, key, 0, 10).Result()

	} else {
		rels, err = utils.Red.ZRevRange(ctx, key, 0, 10).Result()

	}
	//rels, err := utils.Red.ZRange(ctx, key, 0, 10).Result()
	//rels, err := utils.Red.ZRange(ctx, key, 0, 10).Result()
	if err != nil {
		fmt.Println(err)
	}
	// 发送推送消息
	return rels
}
