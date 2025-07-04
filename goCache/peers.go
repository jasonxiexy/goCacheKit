// 使用一致性哈希选择节点        是                                    是
//     |-----> 是否是远程节点 -----> HTTP 客户端访问远程节点 --> 成功？-----> 服务端返回返回值
//                     |  否                                    ↓  否
//                     |----------------------------> 回退到本地节点处理。

package gocache

// peerPicker is the interface that must be implemented to locate
// the per that owns a specific key.
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
