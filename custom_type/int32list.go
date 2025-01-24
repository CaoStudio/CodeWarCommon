package customtype

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	caolog "github.com/CaoStudio/cao-log"
	"strconv"
	"strings"
)

type Int32List []int32

// Scan 扫描
// 该方法前缀必须使用*号，其他方法不可使用*，否则会造成数据读写异常
//
//goland:noinspection GoMixedReceiverTypes
func (l *Int32List) Scan(value interface{}) (err error) {
	*l = Int32List{}
	tx := &sql.NullString{}
	if err = tx.Scan(value); err != nil {
		return err
	}
	if tx.Valid {
		if tx.String == "" {
			return
		}
		strs := strings.Split(tx.String, ",")
		for _, str := range strs {
			var num int
			num, err = strconv.Atoi(str)
			if err != nil {
				return err
			}
			*l = append(*l, int32(num))
		}
	}
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (l Int32List) Value() (driver.Value, error) {
	str := l.ToString()
	caolog.Info("列表数据", l, str)
	return str, nil
}

//goland:noinspection GoMixedReceiverTypes
func (l Int32List) ToString() string {
	//sp := fmt.Sprint(*l)
	//strs := strings.Fields(sp)
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(l)), ","), "[]")
}

// Exist 判断是否存在
//
//goland:noinspection GoMixedReceiverTypes
func (l Int32List) Exist(member int32) bool {
	set := make(map[int32]struct{})
	for _, item := range l {
		set[item] = struct{}{}
	}
	if _, exist := set[member]; exist {
		return true
	}
	return false
}

// ToMap 转换成map对象
//
//goland:noinspection GoMixedReceiverTypes
func (l Int32List) ToMap() map[int32]struct{} {
	set := make(map[int32]struct{})
	for _, item := range l {
		set[item] = struct{}{}
	}
	return set
}
