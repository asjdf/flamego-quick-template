package router

import (
	"github.com/flamego/flamego"
	"github.com/wujunyi792/flamego-quick-template/internal/app/file/handle"
)

func AliGroup(e *flamego.Flame) {
	e.Group("/ali", func() {
		e.Get("/token", handle.HandelGetAliUploadToken)
		e.Post("/upload", handle.HandelAliUpLoad)
	})
}
