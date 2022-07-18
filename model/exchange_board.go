package model

import "fmt"

//交换机
type ExchangeBoard struct {
	MacPort map[string]int//mac与端口映射表
	PortList PortSenderM//所有的端口 端口后可能是电脑，也可能是交换机
	ExchangeBoardPortM ExchangeBoardPortM//交换机与别的交换机连接的端口，比如A的端口上是B  就是A.ExchangeBoardPortM= [B]1
}
func NewExchangeBoard()*ExchangeBoard{
	return &ExchangeBoard{
		MacPort:  map[string]int{},
		PortList: PortSenderM{},
		ExchangeBoardPortM:ExchangeBoardPortM{},
	}
}
func (e *ExchangeBoard)NewMessage(ip,msg string)*Message{//交换机本身没有自己的初始消息
	return nil
}
func (exchangeBoard *ExchangeBoard)SetSender(portSenderM PortSenderM){
	for k,v:=range portSenderM{
		if toExchangeBoard,ok:=v.(*ExchangeBoard);ok{
			exchangeBoard.ExchangeBoardPortM[toExchangeBoard]=k
			exchangeBoard.PortList[k]=v
		}else{
			exchangeBoard.PortList[k]=v
		}
	}
}
func (exchangeBoard *ExchangeBoard)GetSender(port int)Sender{
	return exchangeBoard.PortList[port]
}
func (e *ExchangeBoard)SendMessage(message *Message){
	e.MacPort[message.Head.FromMac]=message.Head.FromPort//更新mac与端口映射表，每次都更新，防止机器换端口
	if message.Head.ToMac==""{//目标mac为空，发送广播
		for k,v:=range e.PortList{
			if k==message.Head.FromPort{//发送方端口不对其发送消息
				continue
			}
			message.Head.IsArpReq=true
			e.commonSendMessage(v,message)
		}
	}else{//有mac
		if port,ok:=e.MacPort[message.Head.ToMac];ok{
			e.commonSendMessage(e.PortList[port],message)
		}else{//广播，更新mac地址表，不存在了
			fmt.Println("?")
		}
	}
}
//给指定端口机器发送消息
func (e *ExchangeBoard)commonSendMessage(v interface{},message *Message){
	switch v.(type) {
	case *Computer:
		v.(*Computer).MsgCh<-message
	case *ExchangeBoard:
		message.Head.FromPort=v.(*ExchangeBoard).ExchangeBoardPortM[e]//交换机的的fromport需要替换
		v.(*ExchangeBoard).SendMessage(message)
	}
}
