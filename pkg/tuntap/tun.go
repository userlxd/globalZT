package tuntap

import (
	"errors"
	"fmt"
	"globalZT/tools/log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

const (
	IFF_NO_PI = 0x10
	IFF_TUN   = 0x01
	IFF_TAP   = 0x02
	TUNSETIFF = 0x400454CA

	IPv6_HEADER_LENGTH = 40
)

type Tun interface {
	Read(ch chan []byte) error
	Write(ch chan []byte) error
	Close() error
}

type tuntap struct {
	mtu  int
	name string
	fd   *os.File
}

func (tun *tuntap) Write(ch chan []byte) error {
	for {
		select {
		case data := <-ch:
			if _, err := tun.fd.Write(data); err != nil {
				return err
			}
		}
	}
}

func (tun *tuntap) Read(ch chan []byte) error {
	buf := make([]byte, tun.mtu)
	for {
		n, err := tun.fd.Read(buf)
		if err != nil {
			return err
		}
		// check length.
		totalLen := 0
		switch buf[0] & 0xf0 {
		case 0x40:
			totalLen = 256*int(buf[2]) + int(buf[3])
		case 0x60:
			totalLen = 256*int(buf[4]) + int(buf[5]) + IPv6_HEADER_LENGTH
		}
		if totalLen != n {
			return fmt.Errorf("read n(%v)!=total(%v)", n, totalLen)
		}
		send := make([]byte, totalLen)
		copy(send, buf)
		ch <- send
	}
}

func (tun *tuntap) Close() error {
	return tun.fd.Close()
}

func Open(addr net.IP, network net.IP, mask net.IP) (Tun, error) {
	var tun = new(tuntap)
	var dev = "/dev/net/tun"
	var err error

	tun.fd, err = os.OpenFile(dev, os.O_RDWR, 0)
	if err != nil {
		log.Log.Errorw("[Load Tun File Error]", "msg", err, "obj", dev)
		return tun, err
	}

	ifr := make([]byte, 18)
	ifr[17] = IFF_NO_PI
	ifr[16] = IFF_TUN
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(tun.fd.Fd()), uintptr(TUNSETIFF),
		uintptr(unsafe.Pointer(&ifr[0])))
	if errno != 0 {
		log.Log.Errorw("[Syscall Tun Error]", "msg", errno, "obj", dev)
		return tun, errors.New("syscall tun error")
	}

	tun.name = string(ifr)
	tun.name = tun.name[:strings.Index(tun.name, "\000")]
	log.Log.Infof("[TUN/TAP dev %s opend]", tun.name)
	tun.mtu = 1500

	var cmd = fmt.Sprintf("ifconfig %s %s netmask %s mtu %d", tun.name, addr.String(), mask.String(), tun.mtu)
	if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
		log.Log.Errorw("[Exec Create Tun Cmd Error]", "msg", err, "obj", cmd)
		return tun, err
	}

	cmd = "route add -net 1.1.1.0 netmask 255.255.255.0 gw 4.4.4.4"
	if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
		log.Log.Errorw("[Exec Create Tun Cmd Error]", "msg", err, "obj", cmd)
		return tun, err
	}
	return tun, nil
}
