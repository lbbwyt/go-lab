package endian

import (
	"encoding/binary"
	"fmt"
)

func BigEndianAndLittleEndianByLibrary() {
	var value uint32 = 10
	by := make([]byte, 4)
	binary.BigEndian.PutUint32(by, value)
	fmt.Println("转换成大端后 ", by)
	fmt.Println("使用大端字节序输出结果：", binary.BigEndian.Uint32(by))
	little := binary.LittleEndian.Uint32(by)
	fmt.Println("大端字节序使用小端输出结果：", little)
}
