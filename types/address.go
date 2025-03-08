package types

import "encoding/hex"

type Address [20]uint8

func AddressFromBytes(b []byte) Address {
	if len(b) != 20 {
		panic("address length must be 20 bytes")
	}

	var res [20]uint8

	for i := 0; i < 20; i++ {
		res[i] = b[i]
	}

	return res
}

func (a Address) String() string {
	return hex.EncodeToString(a[:])
}
