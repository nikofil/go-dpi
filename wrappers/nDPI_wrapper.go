package wrappers

// #cgo CFLAGS: -I/home/av/workspace/nDPI/src/include
// #cgo LDFLAGS: /home/av/workspace/nDPI/lib/libndpi.a -lpcap -lm -pthread
// #include "nDPI_wrapper_impl.h"
import "C"
import (
	"github.com/mushorg/go-dpi"
	"github.com/pkg/errors"
	"unsafe"
)

var ndpiCodeToProtocol = map[uint32]godpi.Protocol{
	7:   godpi.Http,
	5:   godpi.Dns,
	92:  godpi.Ssh,
	127: godpi.Rpc,
	3:   godpi.Smtp,
	88:  godpi.Rdp,
	16:  godpi.Smb,
	81:  godpi.Icmp,
	1:   godpi.Ftp,
	91:  godpi.Ssl,
	64:  godpi.Ssl,
	10:  godpi.Netbios,
}

type NDPIWrapper struct{}

func (_ NDPIWrapper) InitializeWrapper() error {
	if C.ndpi_initialize() != 0 {
		return errors.New("nDPI global structure initialization failed")
	}
	return nil
}

func (_ NDPIWrapper) DestroyWrapper() error {
	C.ndpi_destroy()
	return nil
}

func (_ NDPIWrapper) ClassifyFlow(flow *godpi.Flow) (godpi.Protocol, error) {
	for _, ppacket := range flow.Packets {
		var pktHeader C.struct_pcap_pkthdr
		packet := *ppacket
		pktHeader.ts.tv_sec = C.__time_t(packet.Metadata().Timestamp.Second())
		pktHeader.ts.tv_usec = 0
		pktHeader.caplen = C.bpf_u_int32(packet.Metadata().CaptureLength)
		pktHeader.len = C.bpf_u_int32(packet.Metadata().Length)
		pktDataSlice := packet.Data()
		pktDataPtr := unsafe.Pointer(&pktDataSlice[0])
		ndpiProto := C.pcap_packet_callback(&pktHeader, (*C.u_char)(pktDataPtr))
		if proto, found := ndpiCodeToProtocol[uint32(ndpiProto)]; found {
			return proto, nil
		} else if ndpiProto < 0 {
			switch ndpiProto {
			case -10:
				return godpi.Unknown, errors.New("nDPI wrapper does not support IPv6")
			case -11:
				return godpi.Unknown, errors.New("Received fragmented packet")
			case -12:
				return godpi.Unknown, errors.New("Error creating nDPI flow")
			default:
				return godpi.Unknown, errors.New("nDPI unknown error")
			}
		}
	}
	return godpi.Unknown, nil
}
