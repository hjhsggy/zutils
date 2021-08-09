package random

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

var sli = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
	"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

// GenRand 生成指定长度的安全随机数
func GenRand(len int) (ret string) {

	var buf bytes.Buffer
	for i := 0; i < len; i++ {
		idx, _ := rand.Int(rand.Reader, big.NewInt(36))
		buf.WriteString(sli[int(idx.Int64())])
	}
	ret = buf.String()
	fmt.Println("", zap.Any("GenRand ret : ", ret))
	return
}

// GenUnique 生成唯一ID
func GenUnique() (id string) {

	_uuid, _ := uuid.NewV4()
	_id := fmt.Sprintf("%s", _uuid)

	sum := md5.Sum([]byte(_id))
	id = strings.ToUpper(fmt.Sprintf("%x", sum))
	fmt.Println("", zap.Any("GenUnique id : ", id))
	return
}
