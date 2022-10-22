package commission

import (
	"etherenum-api/etherenum-service/api/pkg/hex"
	"fmt"
	"strconv"
)

type Calculator struct {
	converter *hex.Converter
}

func NewCommissionCalculator(converter *hex.Converter) *Calculator {
	return &Calculator{converter: converter}
}

func (cc *Calculator) GetCommission(gas, price string) string {
	return fmt.Sprintf("0x"+strconv.FormatInt(cc.converter.HexaNumberToInteger(gas) * cc.converter.HexaNumberToInteger(price), 16))
}
