package avformat

/*
#include <libavutil/avutil.h>
#include <libavformat/avformat.h>
 */
import "C"
import (
	"unsafe"
)

/*
不能使用 C.GoBytes(unsafe.Pointer(buf), buf_size) 将C数组转换为Go Slice, 因为C.GoBytes是拷贝
要用: (*[1 << 30]byte)(unsafe.Pointer(buf))[:buf_size:buf_size]

为了避免 "cgo argument has Go pointer to Go pointer" 错误:
	https://github.com/golang/go/wiki/cgo#function-variables
关闭cgo运行时检查, 设置环境变量: export GODEBUG=cgocheck=0
造成这个错误是因为向C函数传递的参数customIOContextOpaque结构中可能存在go指针(如闭包函数)
官方解决方案太复杂
*/

//export avio_ctx_rcb_wrapper
func avio_ctx_rcb_wrapper(opaque unsafe.Pointer, buf *C.uint8_t, buf_size C.int) C.int {
	customCtx := (*customIOContextOpaque)(opaque)
	goBufSlice := (*[1 << 30]byte)(unsafe.Pointer(buf))[:buf_size:buf_size]
	return C.int(customCtx.readCallback(goBufSlice, int(buf_size)))
}

//export avio_ctx_wcb_wrapper
func avio_ctx_wcb_wrapper(opaque unsafe.Pointer, buf *C.uint8_t, buf_size C.int) C.int {
	customCtx := (*customIOContextOpaque)(opaque)
	goBufSlice := (*[1 << 30]byte)(unsafe.Pointer(buf))[:buf_size:buf_size]
	return C.int(customCtx.writeCallback(goBufSlice, int(buf_size)))
}

//export avio_ctx_scb_wrapper
func avio_ctx_scb_wrapper(opaque unsafe.Pointer, offset C.int64_t, whence C.int) C.int64_t {
	customCtx := (*customIOContextOpaque)(opaque)
	return C.int64_t(customCtx.seekCallback(int64(offset), int(whence)))
}
