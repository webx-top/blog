package session

import (
	"strconv"

	"github.com/webx-top/echo"
	"github.com/webx-top/webx/lib/middleware/session"
)

func Middleware(engine string, setting interface{}) echo.MiddlewareFunc {
	store := Store(engine, setting)
	return session.Sessions("XSESSION", store)
}

func Store(engine string, setting interface{}) (store session.Store) {
	switch engine {
	case `file`:
		s := setting.(map[string]string)
		path, _ := s["path"]
		key, _ := s["key"]
		store = session.NewFilesystemStore(path, []byte(key))
	case `redis`:
		s := setting.(map[string]string)
		sizeStr, _ := s["size"]
		network, _ := s["network"]
		address, _ := s["address"]
		password, _ := s["password"]
		key, _ := s["key"]
		size, _ := strconv.Atoi(sizeStr)
		if size < 1 {
			size = 10
		}
		var err error
		store, err = session.NewRedisStore(size, network, address, password, []byte(key))
		if err != nil {
			panic(err)
		}
	case `cookie`:
		fallthrough
	default:
		s := setting.(string)
		store = session.NewCookieStore([]byte(s))
	}
	return
}
