package nfqueue

import (
	"net"
	"unsafe"
)

import "C"

/*
Cast argument to Queue* before calling the real callback

Notes:
  - export cannot be done in the same file (nfqueue.go) else it
    fails to build (multiple definitions of C functions)
    See https://github.com/golang/go/issues/3497
    See https://github.com/golang/go/wiki/cgo
  - this cast is caused by the fact that cgo does not support
    exporting structs
    See https://github.com/golang/go/wiki/cgo

This function must _nerver_ be called directly.
*/
/*
BUG(GoCallbackWrapper): The return value from the Go callback is used as a
verdict. This works, and avoids packets without verdict to be queued, but
prevents using out-of-order replies.
*/
//export GoCallbackWrapper
func GoCallbackWrapper(ptr_q, ptr_nfad, ptr_packet_hw *unsafe.Pointer) int {
	q := (*Queue)(unsafe.Pointer(ptr_q))
	packet_hw := C.GoBytes(unsafe.Pointer(ptr_packet_hw), 8)
	payload := build_payload(q.c_qh, ptr_nfad)
	return q.cb(q, payload, net.HardwareAddr(packet_hw))
}
