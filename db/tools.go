package db

import "os"

func Init(dir string) {
	// todo
	os.MkdirAll(dir, 0755)
}

func Get(key string, value interface{}) {
	// todo
}

func Set(key string, value interface{}) {
	// todo
}
