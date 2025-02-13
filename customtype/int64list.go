package customtype

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Int64List []int64

// Scan 扫描
// 该方法前缀必须使用*号，其他方法不可使用*，否则会造成数据读写异常
func (l *Int64List) Scan(value interface{}) (err error) {
	*l = Int64List{}
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
			*l = append(*l, int64(num))
		}
	}
	return nil
}

func (l Int64List) Value() (driver.Value, error) {
	str := l.ToString()
	return str, nil
}

func (l Int64List) ToString() string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(l)), ","), "[]")
}

// Exist 判断是否存在
func (l Int64List) Exist(member int64) bool {
	set := make(map[int64]struct{})
	for _, item := range l {
		set[item] = struct{}{}
	}
	if _, exist := set[member]; exist {
		return true
	}
	return false
}

// ToMap 转换成map对象
func (l Int64List) ToMap() map[int64]struct{} {
	set := make(map[int64]struct{})
	for _, item := range l {
		set[item] = struct{}{}
	}
	return set
}
