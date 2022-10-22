package hex

import (
	"etherenum-api/etherenum-service/api/pkg/logger"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

type Converter struct {
	Logger logger.Logger
}

func NewConverter(logger logger.Logger) *Converter {
	return &Converter{Logger: logger.Named("hexConverter")}
}
func (c *Converter) HexaNumberToInteger(hexaString string) int64 {
	logger := c.Logger.
		Named("HexaNumberToInteger").
		With("hex", hexaString)

	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)

	output, err := strconv.ParseInt(numberStr, 16, 64)
	if err != nil {
		logger.Error("failed to get block: body is empty. ", "err", err)
		return 0
	}

	return output
}

func (c *Converter) BigFloatConverter(hex string) float64 {
	logger := c.Logger.
		Named("BigFloatConverter").
		With("hex", hex)

	numberStr := strings.Replace(hex, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)

	n, err := strconv.ParseUint(numberStr, 16, 64)
	if err != nil {
		logger.Error("failed to parse uInt ", "err", err)
		return 0
	}

	n2 := uint64(n)
	f := fmt.Sprintf("%v",*(*float64)(unsafe.Pointer(&n2)))
	if len(f) >= 4 {
		summary, err := strconv.ParseFloat(f[:4], 64)
		if err != nil {
			logger.Error("failed to parse uInt ", "err", err)
			return 0
		}
		return summary
	}

	summary, err := strconv.ParseFloat(f, 64)
	if err != nil {
		logger.Error("failed to parse uInt ", "err", err)
		return 0
	}
	return summary
}
