package geecache

import pb "geeCache/geecache/geecachepb"

type PeerPicker interface {
	PickPeer(key string) (peerGetter PeerGetter, ok bool)
}

type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
