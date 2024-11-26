package file

import (
	"context"
	"encoding/gob"
	"errors"
	"log"
	"net"
	"os"
	"slices"

	"github.com/google/uuid"
	"github.com/ophum/tinyipam/pkg/ipam"
)

func init() {
	gob.Register(&data{})
}

type data struct {
	CIDR    *net.IPNet
	Records []*ipam.IP
}

type FileService struct {
	path string
}

var _ ipam.Interface = (*FileService)(nil)

func New(path string) (*FileService, error) {
	return &FileService{
		path: path,
	}, nil
}

func (s *FileService) createIP(name string, addr uint32) (*ipam.IP, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	svcIP := ipam.IP{
		ID:   id.String(),
		Name: name,
		Addr: addr,
	}

	d, err := s.load()
	if err != nil {
		return nil, err
	}

	d.Records = append(d.Records, &svcIP)

	if err := s.save(d); err != nil {
		return nil, err
	}
	return &svcIP, nil
}

func (s *FileService) load() (*data, error) {
	f, err := os.Open(s.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var data data
	if err := gob.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *FileService) save(data *data) error {
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	slices.SortStableFunc(data.Records, func(a, b *ipam.IP) int {
		return int(a.Addr) - int(b.Addr)
	})
	if err := gob.NewEncoder(f).Encode(&data); err != nil {
		return err
	}
	return nil
}
func (s *FileService) Init(ctx context.Context, cidr *net.IPNet, isReserveNetworkAddress, isReserveBroadcastAddress bool) error {
	data := data{
		CIDR: cidr,
	}

	if err := s.save(&data); err != nil {
		return err
	}

	baseAddr := ipam.IPtoUint32(cidr.IP)
	if isReserveNetworkAddress {
		if _, err := s.createIP("Network", baseAddr); err != nil {
			return err
		}
	}
	if isReserveBroadcastAddress {
		reversedMask := data.CIDR.Mask
		for i := range reversedMask {
			reversedMask[i] ^= 0xff
		}

		num := ipam.IPtoUint32(net.IP(reversedMask))
		if _, err := s.createIP("Broadcast", baseAddr+num); err != nil {
			return err
		}
	}
	return nil
}

func (s *FileService) List(ctx context.Context) ([]*ipam.IP, error) {
	d, err := s.load()
	if err != nil {
		return nil, err
	}

	var addrs []*ipam.IP
	for _, r := range d.Records {
		addrs = append(addrs, &ipam.IP{
			ID:   r.ID,
			Name: r.Name,
			Addr: r.Addr,
		})
	}
	return addrs, nil
}

func (s *FileService) Get(ctx context.Context, id string) (*ipam.IP, error) {
	d, err := s.load()
	if err != nil {
		return nil, err
	}

	i := slices.IndexFunc(d.Records, func(r *ipam.IP) bool {
		return r.ID == id
	})
	if i == -1 {
		return nil, errors.New("not found")
	}
	return &ipam.IP{
		ID:   d.Records[i].ID,
		Name: d.Records[i].Name,
		Addr: d.Records[i].Addr,
	}, nil
}

func (s *FileService) AcquireIP(ctx context.Context, opts ...ipam.Option) (*ipam.IP, error) {
	opt := ipam.ApplyOption(opts...)

	d, err := s.load()
	if err != nil {
		return nil, err
	}

	baseAddr := ipam.IPtoUint32(d.CIDR.IP)

	num := baseAddr
	if len(d.Records) == 0 {
	} else if len(d.Records) == 1 {
		if d.Records[0].Addr == baseAddr {
			num = baseAddr + 1
		}
	} else {
		found := false
		for i := 0; i < len(d.Records)-1; i++ {
			if d.Records[i+1].Addr-d.Records[i].Addr == 1 {
				continue
			}

			num = d.Records[i].Addr + 1
			found = true
			break
		}

		if !found {
			log.Println("not found")
			num = d.Records[len(d.Records)-1].Addr + 1
		}
	}

	addr, err := s.createIP(opt.Name, num)
	if err != nil {
		return nil, err
	}

	return addr, nil
}

func (s *FileService) Update(ctx context.Context, ip *ipam.IP) error {
	d, err := s.load()
	if err != nil {
		return err
	}

	i := slices.IndexFunc(d.Records, func(v *ipam.IP) bool {
		return v.ID == ip.ID
	})
	if i == -1 {
		return errors.New("not found")
	}

	addr := ipam.IPtoUint32(ip.IP())
	if slices.ContainsFunc(d.Records, func(v *ipam.IP) bool {
		return v.ID != ip.ID && v.Addr == addr
	}) {
		return errors.New("ip duplicated")
	}

	d.Records[i].Name = ip.Name
	d.Records[i].Addr = ip.Addr

	return s.save(d)
}

func (s *FileService) Delete(ctx context.Context, id string) error {
	d, err := s.load()
	if err != nil {
		return err
	}

	d.Records = slices.DeleteFunc(d.Records, func(v *ipam.IP) bool {
		log.Println(v.ID, id)
		return v.ID == id
	})
	if err := s.save(d); err != nil {
		return err
	}
	return nil
}
