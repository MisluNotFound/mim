# mim
go实现的im系统 

## issue
1. 连接池
2. redis设计（消息设计 在线用户设计 群用户表设计 消息队列设计）
3. messaging层
4. 重写消息收发逻辑
# 流程
用户登录，建立长连接 
获取未读消息 拉取
服务器记录长连接 向logic层查询用户加入的群聊 
    在redis中查找是否有该群
        有则更新用户在群内的状态
        无则将该群的user记录下来
用户下线更新两个表中的状态

当用户发送消息时，交由messaging层处理
    messaging层根据消息类型进行不同操作：
        ack消息：表示客户端接收到消息 服务端记录消息已送达 持久化该条消息 记录已读状态 如果在一定时间内没有收到ack消息则记录状态为未读并持久化
        单聊：生成全局唯一序列号，查询用户是否在线：
            在线：发送到用户的channel中，由ws manager发送给用户
            离线：直接交由logic层持久化 记录未读状态   
        群聊：生成全局唯一序列号，查询所有用户并判断是否在线：
            与单聊一致



以下是一个基于Redis的简单聊天应用示例，以演示Redis在聊天应用中的功能和优势：

创建一个实时聊天室：
使用Redis的发布/订阅功能创建一个频道。
每当用户在聊天室中发送消息时，将该消息发布到频道中。
所有订阅了该频道的用户将实时收到该消息，实现实时聊天效果。
实现在线用户管理：
使用Redis的集合数据类型，将每个在线用户的唯一标识添加到集合中。
当用户上线时，将其添加到集合中；当用户下线时，将其从集合中移除。
通过查询集合中的成员数量，可以获取在线用户的数量。
创建一个聊天消息历史记录：
使用Redis的列表数据类型，将每条聊天消息添加到列表的尾部。
当需要获取聊天历史记录时，可以从列表的头部按照时间顺序遍历。
通过控制列表的长度，可以限制历史记录的大小，保持存储在Redis中的记录数量不过度增长。