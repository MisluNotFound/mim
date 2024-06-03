# mim
go实现的im系统
## 实现的功能
- 支持好友管理
- 支持群聊管理
- 支持单聊
- 支持群聊
- 支持消息漫游
## 系统架构

首先，消息缓存存储在客户端
当推送一条消息时，客户端告诉服务端已接受，服务端更新接收方的ack信息如果失败或超时则更新接收方的err信息
客户端接收消息失败，向服务端发送请求获取seq > lastAck的消息然后去重

消息类型添加 type，url，timer，is_group
维护一个用户(receiver)ack表 记录lastAck的消息 通过lastAck可以获取离线消息 -> redis中
    添加lastErr记录失败的消息，保证lastErr>=lastAck
    当 lastAckErr === null && msgId1 > lastAck 时，更新 lastAck 为 msgId1
    当 lastAckErr !== null :
        msgId1 < lastAckErr 则更新 lastAck 为 msgId1
        msgId1 >= lastAckErr 则不做处理  

获取未读数api
    lastRead记录在缓存中，ZSet user: session score count, session: lastRead 当用户点开会话时更新，lastRead = seq count = 0
    先获取所有的lastRead，再去数据库拉取

获取未读消息api
    本地缓存没有时，发送请求，根据请求的session获取lastRead，去拉取消息
