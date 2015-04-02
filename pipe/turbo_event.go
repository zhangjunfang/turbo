package pipe

import (
	"github.com/blackbeans/turbo/client"
	"github.com/blackbeans/turbo/packet"
)

type IEvent interface {
}

//返回的事件
type IBackwardEvent interface {
	IEvent
}

//向前的事件
type IForwardEvent interface {
	IEvent
}

type PacketEvent struct {
	IForwardEvent
	Packet       *packet.Packet //本次的数据包
	RemoteClient *client.RemotingClient
}

func NewPacketEvent(remoteClient *client.RemotingClient, packet *packet.Packet) *PacketEvent {
	return &PacketEvent{Packet: packet, RemoteClient: remoteClient}
}

type HeartbeatEvent struct {
	IForwardEvent
	RemoteClient *client.RemotingClient
	Opaque       int32
	Version      int64
}

//心跳事件
func NewHeartbeatEvent(remoteClient *client.RemotingClient, opaque int32, version int64) *HeartbeatEvent {
	return &HeartbeatEvent{
		Version:      version,
		Opaque:       opaque,
		RemoteClient: remoteClient}

}

//远程操作事件
type RemotingEvent struct {
	Event      IForwardEvent
	futures    chan map[string]chan interface{} //所有的回调的future
	TargetHost []string                         //发送的特定hostport
	GroupIds   []string                         //本次发送的分组
	Packet     *packet.Packet                   //tlv的packet数据
}

func NewRemotingEvent(packet *packet.Packet, targetHost []string, groupIds ...string) *RemotingEvent {
	revent := &RemotingEvent{
		TargetHost: targetHost,
		GroupIds:   groupIds,
		Packet:     packet,
		futures:    make(chan map[string]chan interface{}, 1)}
	return revent
}

func (self *RemotingEvent) AttachEvent(Event IForwardEvent) {
	self.Event = Event
}

//等待响应
func (self *RemotingEvent) Wait() map[string]chan interface{} {
	return <-self.futures
}

//网络回调事件
type RemoteFutureEvent struct {
	IBackwardEvent
	*RemotingEvent
	Futures map[string] /*groupid*/ chan interface{} //网络结果回调
}

//网络回调事件
func NewRemoteFutureEvent(remoteEvent *RemotingEvent, futures map[string]chan interface{}) *RemoteFutureEvent {
	fe := &RemoteFutureEvent{Futures: futures}
	fe.RemotingEvent = remoteEvent
	return fe
}

//到头的事件
type SunkEvent struct {
	IForwardEvent
}
