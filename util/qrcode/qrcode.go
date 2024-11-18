package qrcode

import (
	"bytes"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
)

type Option struct {
	width  int                     // 宽度
	height int                     // 高度
	level  qr.ErrorCorrectionLevel // 容错等级
}

// Options 选项
type Options func(*Option)

// WithSize 设置二维码尺寸
func WithSize(width, height int) Options {
	return func(o *Option) {
		o.width = width
		o.height = height
	}
}

// WithErrorCorrectionLevel 设置二维码容错等级
func WithErrorCorrectionLevel(level qr.ErrorCorrectionLevel) Options {
	return func(o *Option) {
		o.level = level
	}
}

// GenerateQrcode 生成二维码
func GenerateQrcode(data string, opts ...Options) ([]byte, error) {
	opt := applyOptions(opts...)

	code, err := generateQrCode(data, opt.level)
	if err != nil {
		return nil, err
	}

	code, err = scaleQrCode(code, opt.width, opt.height)
	if err != nil {
		return nil, err
	}

	return encodeToPng(code)
}

func applyOptions(opts ...Options) Option {
	opt := Option{
		width:  256,
		height: 256,
		level:  qr.L,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// 二维码容错等级设置
func generateQrCode(data string, level qr.ErrorCorrectionLevel) (barcode.Barcode, error) {
	return qr.Encode(data, level, qr.Auto)
}

// 二维码大小设置
func scaleQrCode(code barcode.Barcode, width, height int) (barcode.Barcode, error) {
	return barcode.Scale(code, width, height)
}

// 将二维码改成png格式
func encodeToPng(code barcode.Barcode) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, code)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
