package hex

import (
	"etherenum-api/etherenum-service/api/pkg/logger"
	"strconv"
	"strings"
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
		With("hexaString", hexaString)

	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)

	output, err := strconv.ParseInt(numberStr, 16, 64)
	if err != nil {
		logger.Error("failed to get block: body is empty. ", "err", err)
		return 0
	}

	return output
}
