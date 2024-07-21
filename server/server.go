package server

import "log/slog"

// Server
// @Description: 消息队列服务器，接收客户端发送的消息
type Server struct {
	*Config

	topics map[string]Storage // 存储所有主题

	consumers   []Consumer    // 订阅了特定主题的消费者
	producers   []Producer    // 向特定主题发送消息的生产者
	produceChan chan Message  // 接收生产者发送的消息的通道
	quitChan    chan struct{} // 通知服务器退出
}

func NewServer(conf *Config) (*Server, error) {
	produceChan := make(chan Message)
	return &Server{
		Config:      conf,
		topics:      make(map[string]Storage),
		quitChan:    make(chan struct{}),
		produceChan: produceChan,
		producers: []Producer{
			NewHttpProducer(conf.ListenAddr, produceChan),
		},
	}, nil
}

func (s *Server) Start() {
	for _, producer := range s.producers { //启动所有生产者
		go producer.Start()
	}
	s.loop()
}

// loop
// @Description: 监听produceChan里的消息
// @receiver s
func (s *Server) loop() {
	for {
		select {
		case <-s.quitChan:
			return
		case msg := <-s.produceChan:
			offset, err := s.publish(msg)
			if err != nil {
				slog.Error("publish message failed", "error", err)
			} else {
				slog.Info("publish message success", "offset", offset)
			}
		}
	}
}

// publish
// @Description: 将消息（Message）发布到指定的主题（Topic）
// @receiver s
// @param msg
// @return int
// @return error
func (s *Server) publish(msg Message) (int, error) {
	store := s.getStoreForTopic(msg.Topic)
	return store.Push(msg.Data)
}

func (s *Server) getStoreForTopic(topic string) Storage {
	if _, ok := s.topics[topic]; !ok { //如果topic不存在，则创建一个新的存储并将其添加到topics映射中。
		s.topics[topic] = s.StoreProducerFunc()
		slog.Info("create new storage for topic", "topic", topic)
	}
	return s.topics[topic] //获取与给定topic关联的存储（Storage）
}
