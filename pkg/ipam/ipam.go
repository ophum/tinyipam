package ipam

import (
	"context"
	"net"
)

type IP struct {
	ID   string
	Name string
	Addr uint32
}

func (ip *IP) IP() net.IP {
	r := net.IP{0, 0, 0, 0}
	r[3] = byte(ip.Addr & 0xff)
	r[2] = byte(ip.Addr & 0xff00 >> 8)
	r[1] = byte(ip.Addr & 0xff0000 >> 16)
	r[0] = byte(ip.Addr & 0xff000000 >> 24)
	return r
}

func IPtoUint32(ip net.IP) uint32 {
	var intip uint32 = 0
	for i, v := range ip {
		intip += uint32(v) << ((3 - i) * 8)
	}

	return intip
}

type Interface interface {
	Init(ctx context.Context, cidr *net.IPNet, isReserveNetworkAddress, isReserveBroaddcastAddress bool) error
	List(ctx context.Context) ([]*IP, error)
	Get(ctx context.Context, id string) (*IP, error)
	AcquireIP(ctx context.Context, opts ...Option) (*IP, error)
	Update(ctx context.Context, ip *IP) error
	Delete(ctx context.Context, id string) error
}

type Option interface {
	apply(o *Options)
}

type Options struct {
	Name string
	IP   string
}

func ApplyOption(opts ...Option) *Options {
	var opt Options
	for _, o := range opts {
		o.apply(&opt)
	}
	return &opt
}

type ip struct {
	v string
}

func StaticIP(ipaddr string) Option {
	return &ip{
		v: ipaddr,
	}
}

func (ip *ip) apply(o *Options) {
	o.IP = ip.v
}

type withName struct {
	name string
}

func Name(name string) Option {
	return &withName{
		name: name,
	}
}
func (o *withName) apply(opts *Options) {
	opts.Name = o.name
}
