package db

import "context"

type DB struct {
	ctx     context.Context
	String  *String
	HashMap *HashMap
}

func NewDB(ctx context.Context) *DB {
	return &DB{
		ctx:     ctx,
		String:  NewStrings(),
		HashMap: NewHashMap(),
	}
}
