package domain

type RespPair[T any] struct {
	Resp T
	Err  error
}
