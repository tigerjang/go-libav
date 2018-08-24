package avcodec

/*
#include <libavcodec/avcodec.h>
#include <libavutil/hwcontext.h>

static enum AVPixelFormat CBWrapper_get_format(AVCodecContext *ctx,
                                        const enum AVPixelFormat *pix_fmts)
{
	long (*goCB)(void *, const enum AVPixelFormat *, unsigned int);
	void ** opaque = ((void **)(ctx -> opaque));
	//printf("opaque[0]: %d", opaque[0]);
	//printf("opaque[1]: %d", opaque[1]);
	goCB = (long (*)(void *, const enum AVPixelFormat *, unsigned int))(opaque[1]);
    const enum AVPixelFormat *p;
    unsigned int size = 0;
    for (p = pix_fmts; *p != -1; p++) {
        size ++;
    }
    return (enum AVPixelFormat)(goCB(opaque[0], pix_fmts, size));
}

static void set_get_format_callback(AVCodecContext *ctx) {
	ctx->get_format = CBWrapper_get_format;
}

*/
import "C"
import (
	"unsafe"
	"syscall"
	"github.com/tigerjang/go-libav/avutil"
	//"log"
)

// ************** HWDeviceContent **************
type HWDeviceType C.enum_AVHWDeviceType

const (
	//HWDeviceTypeVDPAU        HWDeviceType = C.AV_HWDEVICE_TYPE_VDPAU
	//HWDeviceTypeCUDA         HWDeviceType = C.AV_HWDEVICE_TYPE_CUDA
	//HWDeviceTypeVAAPI        HWDeviceType = C.AV_HWDEVICE_TYPE_VAAPI
	//HWDeviceTypeDXVA2        HWDeviceType = C.AV_HWDEVICE_TYPE_DXVA2
	//HWDeviceTypeQSV          HWDeviceType = C.AV_HWDEVICE_TYPE_QSV
	//HWDeviceTypeVIDEOTOOLBOX HWDeviceType = C.AV_HWDEVICE_TYPE_VIDEOTOOLBOX
	HWDeviceTypeNONE         HWDeviceType = C.AV_HWDEVICE_TYPE_NONE
	//HWDeviceTypeD3D11VA      HWDeviceType = C.AV_HWDEVICE_TYPE_D3D11VA
	//HWDeviceTypeDRM          HWDeviceType = C.AV_HWDEVICE_TYPE_DRM
)

type HWFrameTransferDirection C.enum_AVHWFrameTransferDirection

const (
	// Transfer the data from the queried hw frame
	HWFrameTransferDirectionFrom HWFrameTransferDirection = C.AV_HWFRAME_TRANSFER_DIRECTION_FROM
	// Transfer the data to the queried hw frame.
	HWFrameTransferDirectionTo HWFrameTransferDirection = C.AV_HWFRAME_TRANSFER_DIRECTION_TO
)

func (ctx *Context) GetFormatCallback(callback func(codecCtx *Context, availPxlFmts []string) string) {
	//ctx.CAVCodecContext.get_format = C.CBWrapper_get_format
	C.set_get_format_callback((*C.AVCodecContext)(unsafe.Pointer(ctx.CAVCodecContext)))
	ctx.opaque[1] = syscall.NewCallback(
		func (goCtxPtr uintptr, pxlFmts *avutil.PixelFormat, pixFmtSize C.uint) int64 {
			pf_arr := (*[1 << 30]avutil.PixelFormat)(unsafe.Pointer(pxlFmts))[:pixFmtSize:pixFmtSize]
			pf_names := make([]string, pixFmtSize, pixFmtSize)
			for idx, pf := range pf_arr {
				pf_names[idx] = pf.Name()
			}
			r_pf_name := callback((*Context)(unsafe.Pointer(goCtxPtr)), pf_names)
			r_pf, fmt_exist := avutil.FindPixelFormatByName(r_pf_name)
			ctx.hwPxlFmt = r_pf
			if !fmt_exist {
				return int64(avutil.PixelFormatNone)
			}
			return int64(r_pf)
		})
}

func (ctx *Context) GetHwCtxPixelFormat() avutil.PixelFormat {
	return ctx.hwPxlFmt
}


func HWDeviceFindTypeByName(name string) HWDeviceType {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	hwType := (HWDeviceType)(C.av_hwdevice_find_type_by_name(cName))
	return hwType
}


type HWDeviceContext struct {
	CAVBufferRef *C.AVBufferRef
}

func NewHWDeviceContext(deviceType HWDeviceType, device string, options *avutil.Dictionary, flags int) (*HWDeviceContext, error) {
	var hwDeviceCtx *C.AVBufferRef

	var cDevice *C.char
	if device == "" {
		cDevice = nil
	} else {
		cDevice = C.CString(device)
		defer C.free(unsafe.Pointer(cDevice))
	}

	var cOptions *C.AVDictionary
	if options != nil {
		cOptions = *(**C.AVDictionary)(options.Pointer())
	}

	if code := int(C.av_hwdevice_ctx_create(&hwDeviceCtx, (C.enum_AVHWDeviceType)(deviceType), cDevice, cOptions, C.int(flags))); code < 0 {
		return nil, ErrAllocationError
	}
	return &HWDeviceContext{hwDeviceCtx}, nil
}

func (ctx *Context) SetHWDeviceContext(hwCtx *HWDeviceContext) {
	ctx.CAVCodecContext.hw_device_ctx = C.av_buffer_ref(hwCtx.CAVBufferRef)
}

func (ctx *HWDeviceContext) Free() {
	C.av_free(unsafe.Pointer(ctx.CAVBufferRef))
}

func HWFrameTransferData(dst, src *avutil.Frame, flags int) error {
	if code := C.av_hwframe_transfer_data((*C.AVFrame)(unsafe.Pointer(dst.CAVFrame)), (*C.AVFrame)(unsafe.Pointer(src.CAVFrame)), C.int(flags)); code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return nil
}

func HWFrameTransferGetFormats(hwframe_ctx_ref *avutil.BufferRef,
	dir HWFrameTransferDirection, flags int) ([]avutil.PixelFormat, error) {
	var cRet *avutil.PixelFormat

	if code := C.av_hwframe_transfer_get_formats(
		(*C.AVBufferRef)(unsafe.Pointer(hwframe_ctx_ref)),
		C.enum_AVHWFrameTransferDirection(dir),
		(**C.enum_AVPixelFormat)(unsafe.Pointer(&cRet)),
		C.int(flags)); code < 0 {
		return nil, avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	eleSize := unsafe.Sizeof(*cRet)
	arrLen := 0
	for addr := uintptr(unsafe.Pointer(cRet));
			*((*avutil.PixelFormat)(unsafe.Pointer(addr))) != avutil.PixelFormatNone;
			addr += eleSize {
		arrLen += 1
	}
	return (*[1 << 30]avutil.PixelFormat)(unsafe.Pointer(cRet))[:arrLen:arrLen], nil
}

// TODO: Free !!!!!!!!!!!
// ************** HWDeviceContent **************
