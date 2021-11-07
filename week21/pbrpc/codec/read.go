package codec

import (
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"net"

	"gitee.com/infraboard/go-course/day21/pbrpc/codec/pb"
	"github.com/golang/snappy"
	"google.golang.org/protobuf/proto"
)

func ReadRequestHeader(r io.Reader, header *pb.RequestHeader) (err error) {
	// recv header (more)
	pbHeader, err := recvFrame(r, int(pb.Const_MAX_REQUEST_HEADER_LEN))
	if err != nil {
		return err
	}

	// Marshal Header
	err = proto.Unmarshal(pbHeader, header)
	if err != nil {
		return err
	}

	return nil
}

func recvFrame(r io.Reader, maxSize int) (data []byte, err error) {
	size, err := readUvarint(r)
	if err != nil {
		return nil, err
	}
	if maxSize > 0 {
		if int(size) > maxSize {
			return nil, fmt.Errorf("protorpc: varint overflows maxSize(%d)", maxSize)
		}
	}
	if size != 0 {
		data = make([]byte, size)
		if err = read(r, data); err != nil {
			return nil, err
		}
	}
	return data, nil
}

// ReadUvarint reads an encoded unsigned integer from r and returns it as a uint64.
func readUvarint(r io.Reader) (uint64, error) {
	var x uint64
	var s uint
	for i := 0; ; i++ {
		var b byte
		b, err := readByte(r)
		if err != nil {
			return 0, err
		}
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return x, errors.New("protorpc: varint overflows a 64-bit integer")
			}
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
}

func readByte(r io.Reader) (c byte, err error) {
	data := make([]byte, 1)
	if err = read(r, data); err != nil {
		return 0, err
	}
	c = data[0]
	return
}

func read(r io.Reader, data []byte) error {
	for index := 0; index < len(data); {
		n, err := r.Read(data[index:])
		if err != nil {
			if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
				return err
			}
		}
		index += n
	}
	return nil
}

func ReadRequestBody(r io.Reader, header *pb.RequestHeader, request proto.Message) error {
	maxBodyLen := maxUint32(header.RawRequestLen, header.SnappyCompressedRequestLen)

	// recv body (end)
	compressedPbRequest, err := recvFrame(r, int(maxBodyLen))
	if err != nil {
		return err
	}

	// checksum
	if header.Checksum != 0 {
		if crc32.ChecksumIEEE(compressedPbRequest) != header.Checksum {
			return fmt.Errorf("protorpc.readRequestBody: unexpected checksum")
		}
	}

	var pbRequest []byte
	if header.SnappyCompressedRequestLen != 0 {
		// decode the compressed data
		pbRequest, err = snappy.Decode(nil, compressedPbRequest)
		if err != nil {
			return err
		}
		// check wire header: rawMsgLen
		if uint32(len(pbRequest)) != header.RawRequestLen {
			return fmt.Errorf("protorpc.readRequestBody: Unexcpeted header.RawRequestLen")
		}
	} else {
		pbRequest = compressedPbRequest
	}

	// Unmarshal to proto message
	if request != nil {
		err = proto.Unmarshal(pbRequest, request)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadResponseHeader(r io.Reader, header *pb.ResponseHeader) error {
	// recv header (more)
	pbHeader, err := recvFrame(r, 0)
	if err != nil {
		return err
	}

	// Marshal Header
	err = proto.Unmarshal(pbHeader, header)
	if err != nil {
		return err
	}

	return nil
}
