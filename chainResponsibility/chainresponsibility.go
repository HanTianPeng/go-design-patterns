/*
Author: Conk
Function: 责任链模式
CreateTime: 2021/5/25 9:59 上午
UpdateTime: 2021/5/25 9:59 上午
责任链模式:
	将请求的发送和接收解耦，让多个接收对象都有机会处理这个请求。将这些接收对象串成一条链，并沿着这条链传递这个请求，直到链上的某个接收对象能够处理它为止。
优点:
	请求的发送者和接收者解耦
缺点:
	不能保证请求一定被接收
应用场景:
	OA审批,中间件
*/

package main

import "fmt"

type Manager interface {
	HaveRight(day int) bool
	HandleFeeRequest(name string, day int) bool
}

// RequestChain: 责任链抽象处理类
type RequestChain struct {
	Manager
	nextRequestChain *RequestChain
}

// SetRequestChain: 挂节点
func (rs *RequestChain) SetRequestChain(r *RequestChain) {
	rs.nextRequestChain = r
}

// HandleFeeRequest: 处理请求
func (rs *RequestChain) HandleFeeRequest(name string, day int) bool {
	// 判断谁有权限,属于谁的责任就归谁处理
	if rs.Manager.HaveRight(day) {
		return rs.Manager.HandleFeeRequest(name, day)
	}
	if rs.nextRequestChain != nil {
		return rs.nextRequestChain.HandleFeeRequest(name, day)
	}
	return false
}

// HaveRight: 判断是否有权限
func (rs *RequestChain) HaveRight(day int) bool {
	return true
}

// DirectManager: 直接领导类
type DirectManager struct {}

func (dm *DirectManager) HaveRight(day int) bool {
	return day < 3
}

// HandleFeeRequest: 直接领导权限校验规则: 请假天数是基本校验规则, 也可以根据针对性审批
func (dm *DirectManager) HandleFeeRequest(name string, day int) bool {
	if name == "conk" {
		fmt.Printf("DirectManager allows %s to  take a  %d-day \n", name, day)
		return true
	}
	fmt.Printf("DirectManager rejects %s to have a %d-day \n", name, day)
	return false
}

func NewDirectManagerChain() *RequestChain {
	return &RequestChain{
		Manager: &DirectManager{},
	}
}


// DeptManager: 部门经理类
type DeptManager struct {}

func (dem *DeptManager) HaveRight(day int) bool {
	return day < 7
}

func (dem *DeptManager) HandleFeeRequest(name string, day int) bool {
	if name == "pht" {
		fmt.Printf("DeptManager allows %s to  take a %d-day \n", name, day)
		return true
	}
	fmt.Printf("DeptManager rejects %s to have a %d-day \n", name, day)
	return false
}

func NewDeptManagerChain() *RequestChain {
	return &RequestChain{
		Manager: &DeptManager{},
	}
}


// GeneralManager: 总经理类
type GeneralManager struct {}

func (gm *GeneralManager) HaveRight(day int) bool {
	return day >= 7
}

func (gm *GeneralManager) HandleFeeRequest(name string, day int) bool {
	if name == "penghantian" {
		fmt.Printf("GeneralManager allows %s to  take a %d-day \n", name, day)
		return true
	}
	fmt.Printf("GeneralManager rejects %s to have a %d-day \n", name, day)
	return false
}

func NewGeneralManagerChain() *RequestChain {
	return &RequestChain{
		Manager: &GeneralManager{},
	}
}

// ChainRepquestFactory: 实例化处理器链
func ChainRepquestFactory() Manager {
	cr1 := NewDirectManagerChain()
	cr2 := NewDeptManagerChain()
	cr3 := NewGeneralManagerChain()

	cr1.SetRequestChain(cr2)
	cr2.SetRequestChain(cr3)

	return cr1
}

func main() {
	cr := ChainRepquestFactory()
	cr.HandleFeeRequest("conk", 2)
	cr.HandleFeeRequest("pht", 5)
	cr.HandleFeeRequest("penghantian", 10)
	cr.HandleFeeRequest("zzm", 4)
}

