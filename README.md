# mim
go实现的im系统 

## issue
0. 
1. 存消息已读未读 时间戳
2. 推拉隔离
3. 优化查询等等
4. 群聊消息

# 流程
用户登录，建立长连接 
放入在线表 logic层获取用户群组在redis中维护在线用户以及用户的群聊和会话 logic开一个长连接业务
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


拉模式：
    用户上线之后主动拉取会话列表，并拉取一定量的消息（所有未读+多少条已读）/（一共多少条？）
    短链接logic先获取用户id，找session列表，从每个session列表拉若干消息id，查询消息信息拼接返回



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

表设计：
会话列表:
    sortedset prefix-senderid: senderid+targetid msgid/seq 长期存在
消息列表：
    sortedset prefix-senderid+targetid: msgid/seq msgid/seq
离线消息列表：
    set prefix-userid: senderid
    list prefix-senderid: msgid/seq
消息记录：
    senderid， targetid， content， isread


缓存设计：
    消息存入缓存，异步写入数据库
    对于文本消息 直接存储内容，对于二进制消息，文件系统存储内容，缓存存储url
    写策略：先写redis

mq层设计：
    connect与logic层的交互使用mq
    connect-->logic:
        根据logic层消费者数量确定queue      logic初始化queue
        connect将消息放入mq 保证消息只被消费一次        
        logic处理消息放入mq         
    logic-->connect:
        server对应exchange
        bucket对应queue
        bucket作为消费者
        bucket获取消息，在自己的map中找到user发送