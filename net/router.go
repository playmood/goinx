package net

import "goinx/iface"

// 实现router时，先嵌入BaseRouter基类，根据需要对这个基类方法重写即可
type BaseRouter struct{}

func (b *BaseRouter) PreHandle(request iface.IRequest) {}

func (b *BaseRouter) Handle(request iface.IRequest) {}

func (b *BaseRouter) PostHandle(request iface.IRequest) {}
