/*
 * @Descripttion:定时器
 * @Author: DW
 * @Date: 2023-07-08 10:51:24
 * @LastEditTime: 2023-07-08 10:56:00
 */
package utils

import "time"

type TimeFunc func(interface{}) bool

func Timer(delay, tick time.Duration, fun TimeFunc, param interface{}) {
	go func() {
		if fun == nil {
			return
		}
		tm := time.NewTimer(delay)
		for {
			select {
			case <-tm.C:
				if fun(param) == false {
					return
				}
				tm.Reset(tick)
			}
		}
	}()
}
