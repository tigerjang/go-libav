package swscale

/*
#include <libswscale/swscale.h>

#cgo pkg-config: libswscale
*/
import "C"
import (
	"github.com/tigerjang/go-libav/avutil"
	"unsafe"
)

type Flags C.int

const (
	SWS_FAST_BILINEAR Flags = C.SWS_FAST_BILINEAR
	SWS_BILINEAR Flags = C.SWS_BILINEAR
	SWS_BICUBIC Flags = C.SWS_BICUBIC
	SWS_X Flags = C.SWS_X
	SWS_POINT Flags = C.SWS_POINT
	SWS_AREA Flags = C.SWS_AREA
	SWS_BICUBLIN Flags = C.SWS_BICUBLIN
	SWS_GAUSS Flags = C.SWS_GAUSS
	SWS_SINC Flags = C.SWS_SINC
	SWS_LANCZOS Flags = C.SWS_LANCZOS
	SWS_SPLINE Flags = C.SWS_SPLINE
	SWS_PARAM_DEFAULT Flags = C.SWS_PARAM_DEFAULT
	SWS_PRINT_INFO Flags = C.SWS_PRINT_INFO
	SWS_FULL_CHR_H_INT Flags = C.SWS_FULL_CHR_H_INT
	SWS_FULL_CHR_H_INP Flags = C.SWS_FULL_CHR_H_INP
	SWS_DIRECT_BGR Flags = C.SWS_DIRECT_BGR
	SWS_ACCURATE_RND Flags = C.SWS_ACCURATE_RND
	SWS_BITEXACT Flags = C.SWS_BITEXACT
	SWS_ERROR_DIFFUSION Flags = C.SWS_ERROR_DIFFUSION
)

type SwsContext struct {
	CSwsContext *C.struct_SwsContext // struct_SwsContext
}

func GetContext(
		srcW, srcH int, srcFormat avutil.PixelFormat,
		dstW, dstH int, dstFormat avutil.PixelFormat,
		flags Flags, srcFilter, dstFilter *SwsFilter, param []float64) (*SwsContext, error) {
	var C_sf *C.SwsFilter = nil
	var C_df *C.SwsFilter = nil
	if srcFilter != nil {
		C_sf = srcFilter.CSwsFilter
	}
	if dstFilter != nil {
		C_df = dstFilter.CSwsFilter
	}

	var C_Param *C.double = nil
	if param != nil {
		C_Param = (*C.double)(unsafe.Pointer(&param[0]))
	}

	ret := C.sws_getContext(
		C.int(srcW), C.int(srcH), C.enum_AVPixelFormat(srcFormat),
		C.int(dstW), C.int(dstH), C.enum_AVPixelFormat(dstFormat),
		C.int(flags), C_sf, C_df, C_Param)

	if ret == nil {
		return nil, avutil.NewErrorFromCode(0)
	}
	return &SwsContext{ret}, nil
}

type SwsFilter struct {
	CSwsFilter *C.SwsFilter
}

//const uint8_t *const srcSlice[], const int srcStride[], int srcSliceY, int srcSliceH, uint8_t *const dst[], const int dstStride[])
func (ctx *SwsContext) SwsScale(
		srcSlice []*uint8, srcStride []int32,
		srcSliceY, srcSliceH int,
		dst []*uint8, dstStride []int32) int {
	ret := C.sws_scale(
		ctx.CSwsContext,
		(**C.uint8_t)(unsafe.Pointer(&srcSlice[0])),
		(*C.int)(unsafe.Pointer(&srcStride[0])),
		C.int(srcSliceY),
		C.int(srcSliceH),
		(**C.uint8_t)(unsafe.Pointer(&dst[0])),
		(*C.int)(unsafe.Pointer(&dstStride[0])))
	return int(ret)
}

func (ctx *SwsContext) Free() {
	C.sws_freeContext(ctx.CSwsContext)
}

