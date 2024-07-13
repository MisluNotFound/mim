# mim
go实现的im系统后端
[前端地址](https://github.com/MisluNotFound/mim-front/tree/master)
## 实现的功能
- 支持好友管理
- 支持群聊管理
- 支持单聊
- 支持群聊
- 支持媒体文件
## 技术栈
- rpc框架：rpcx
- web框架：gin
- mq：rabbitMQ
- 长连接：websocket
- 缓存：redis
- 网关：nginx
## 系统架构
根据连接类型和处理的功能将系统分成三层：
- api：负责接收短链接业务
- connect：负责长连接维护
- logic：负责处理业务，分为短链接业务层和长连接业务层，与db和redis交互
### 各层之间的交互
api层通过rpcx与logic层通信，connect层与logic层通过mq通信，connect与api层无交互。
## 系统详细设计
### mq设计
#### connect->logic
connect层生产者的数量由server数量决定，往同一个exchange投放消息，logic层设定多个消费者采用work方式消费消息。
#### logic->connect
logic层设定多个生产者，根据消息接收方的serverId， bucketIdx确定exchange和routingKey，此时每个bucket是一个消费者。
### connect层设计
用户发起长连接请求，nginx将请求发到目标服务器，用户与connect层建立长连接之后，由服务器维护用户的在线列表。为了减少锁的竞争，服务器分为若干个桶，每个用户根据id和散列值放入对应的桶中。每个桶根据用户的活跃时间维护一个小根堆用于清理长时间待机用户和解决当桶满时有新用户的情况。connect层监听client的消息和发送client收到的消息，当client发送消息之后，将消息投放到mq中，由logic处理，同时connect层消费logic层处理过的消息，根据消息内容转发给对应的client。
### logic层设计
logic层通过rpcx处理api层的短链接业务，通过mq层处理connect层的消息并投递回mq。
### 消息设计
发消息采用写扩散的方式，对于在线的用户，发往该用户的消息直接通过ws长连接发往用户。
收消息采用读扩散的方式，无论单聊还是群聊，每条消息只写一份，根据不同的条件进行读取。
#### 消息的不丢
客户端本地存储消息缓存和已接受的消息seq
当推送一条消息时，客户端告诉服务端已接受，服务端更新接收方的ack信息如果失败或超时则更新接收方的err信息
客户端接收消息失败，向服务端发送请求获取seq > lastAck的消息然后去重
维护一个用户(receiver)ack表，记录lastAck的消息。通过lastAck可以获取离线消息
1. 添加lastErr记录失败的消息，保证lastErr>=lastAck
2. 当 lastAckErr == null && msgId1 > lastAck 时，更新 lastAck 为 msgId1
3. 当 lastAckErr != null :
   msgId1 < lastAckErr 则更新 lastAck 为 msgId1
   msgId1 >= lastAckErr 则不做处理  
#### 离线消息设计
维护一个用户的会话列表，记录会话的状态：已读/未读，再维护一个会话中离线消息数，有离线消息时自增。当用户上线时，获取所有未读数和会话中最后一条消息，减少对db的访问，当用户查看会话时再进行拉取，更新状态和清零离线消息数。
