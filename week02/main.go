package main

import (
	"fmt"
	"github.com/pkg/errors"
)

/*
1. 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，
抛给上层。为什么，应该怎么做请写出代码？
*/

/*
答：在dao层应该wrap这个error，因为我们想要在biz层知道详细的错误信息及堆栈信息
*/

// dao

var NoRowsError = errors.New("No errors")

func Get() error {
	sql := ""
	return errors.Wrap(NoRowsError, fmt.Sprintf("not found,sql:%s", sql))
}

// biz

func main() {
	err := Get()
	if errors.Is(err, NoRowsError) {
		println("not found", err.Error())
	}
}
