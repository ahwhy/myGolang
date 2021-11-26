package codec

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"net"

	"github.com/golang/snappy"
	"google.golang.org/protobuf/proto"

	"github.com/ahwhy/myGolang/rpc/pbrpc/codec/pb"
)

func maxUint32(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

var (
	UseSnappy            = true
	UseCrc32ChecksumIEEE = true
)

func WriteResponse(w io.Writer, id uint64, serr string, response proto.Message) (err error) {
	// clean response if error
	if serr != "" {
		response = nil
	}

	// marshal response
	pbResponse := []byte{}
	if response != nil {
		pbResponse, err = proto.Marshal(response)
		if err != nil {
			return err
		}
	}

	// compress serialized proto data
	compressedPbResponse := snappy.Encode(nil, pbResponse)

	// generate header
	header := &pb.ResponseHeader{
		Id:                          id,
		Error:                       serr,
		RawResponseLen:              uint32(len(pbResponse)),
		SnappyCompressedResponseLen: uint32(len(compressedPbResponse)),
		Checksum:                    crc32.ChecksumIEEE(compressedPbResponse),
	}

	if !UseSnappy {
		header.SnappyCompressedResponseLen = 0
		compressedPbResponse = pbResponse
	}
	if !UseCrc32ChecksumIEEE {
		header.Checksum = 0
	}

	// check header size
	pbHeader, err := proto.Marshal(header)
	if err != err {
		return
	}

	// send header (more)
	if err = sendFrame(w, pbHeader); err != nil {
		return
	}

	// send body (end)
	if err = sendFrame(w, compressedPbResponse); err != nil {
		return
	}

	return nil
}

func sendFrame(w io.Writer, data []byte) (err error) {
	// Allocate enough space for the biggest uvarint
	var size [binary.MaxVarintLen64]byte

	if data == nil || len(data) == 0 {
		n := binary.PutUvarint(size[:], uint64(0))
		if err = write(w, size[:n], false); err != nil {
			return
		}
		return
	}

	// Write the size and data
	n := binary.PutUvarint(size[:], uint64(len(data)))
	if err = write(w, size[:n], false); err != nil {
		return
	}
	if err = write(w, data, false); err != nil {
		return
	}
	return
}

func write(w io.Writer, data []byte, onePacket bool) error {
	if onePacket {
		if _, err := w.Write(data); err != nil {
			return err
		}
		return nil
	}
	for index := 0; index < len(data); {
		n, err := w.Write(data[index:])
		if err != nil {
			if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
				return err
			}
		}
		index += n
	}
	return nil
}

func WriteRequest(w io.Writer, id uint64, method string, request proto.Message) error {
	// marshal request
	pbRequest := []byte{}
	if request != nil {
		var err error
		pbRequest, err = proto.Marshal(request)
		if err != nil {
			return err
		}
	}

	// compress serialized proto data
	compressedPbRequest := snappy.Encode(nil, pbRequest)

	// generate header
	header := &pb.RequestHeader{
		Id:                         id,
		Method:                     method,
		RawRequestLen:              uint32(len(pbRequest)),
		SnappyCompressedRequestLen: uint32(len(compressedPbRequest)),
		Checksum:                   crc32.ChecksumIEEE(compressedPbRequest),
	}

	if !UseSnappy {
		header.SnappyCompressedRequestLen = 0
		compressedPbRequest = pbRequest
	}
	if !UseCrc32ChecksumIEEE {
		header.Checksum = 0
	}

	// check header size
	pbHeader, err := proto.Marshal(header)
	if err != err {
		return err
	}
	if len(pbHeader) > int(pb.Const_MAX_REQUEST_HEADER_LEN) {
		return fmt.Errorf("protorpc.writeRequest: header larger than max_header_len: %d.", len(pbHeader))
	}

	// send header (more)
	if err := sendFrame(w, pbHeader); err != nil {
		return err
	}

	// send body (end)
	if err := sendFrame(w, compressedPbRequest); err != nil {
		return err
	}

	return nil
}

func ReadResponseBody(r io.Reader, header *pb.ResponseHeader, response proto.Message) error {
	maxBodyLen := int(maxUint32(header.RawResponseLen, header.SnappyCompressedResponseLen))

	// recv body (end)
	compressedPbResponse, err := recvFrame(r, maxBodyLen)
	if err != nil {
		return err
	}

	// checksum
	if header.Checksum != 0 {
		if crc32.ChecksumIEEE(compressedPbResponse) != header.Checksum {
			return fmt.Errorf("protorpc.readResponseBody: unexpected checksum")
		}
	}

	var pbResponse []byte
	if header.SnappyCompressedResponseLen != 0 {
		// decode the compressed data
		pbResponse, err = snappy.Decode(nil, compressedPbResponse)
		if err != nil {
			return err
		}
		// check wire header: rawMsgLen
		if uint32(len(pbResponse)) != header.RawResponseLen {
			return fmt.Errorf("protorpc.readResponseBody: Unexcpeted header.RawResponseLen")
		}
	} else {
		pbResponse = compressedPbResponse
	}

	// Unmarshal to proto message
	if response != nil {
		err = proto.Unmarshal(pbResponse, response)
		if err != nil {
			return err
		}
	}

	return nil
}
