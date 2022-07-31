package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/spuxy/barber/handlers/debug"
	"go.uber.org/zap"
)

func dedfaultDebugMux() *gin.Engine {
	r := gin.Default()
	r.GET("/debug/pprof", gin.WrapF(pprof.Index))
	r.GET("/debug/pprof/cmdline", gin.WrapF(pprof.Index))
	r.GET("/debug/pprof/profile", gin.WrapF(pprof.Profile))
	r.GET("/debug/pprof/symbol", gin.WrapF(pprof.Symbol))
	r.GET("/debug/pprof/trace", gin.WrapF(pprof.Trace))
	r.GET("/debug/vars", gin.WrapH(expvar.Handler()))
	return r
}

func DebugMux(logger *zap.SugaredLogger, ver, commit, date string) http.Handler {
	mux := dedfaultDebugMux()

	dh := debug.Handler{
		Log:          logger,
		BuildVersion: ver,
		BuildCommit:  commit,
		BuildDate:    date,
	}

	mux.GET("/debug/readiness", gin.WrapF(dh.Readiness))
	mux.GET("/debug/liveness", gin.WrapF(dh.Liveness))

	return mux
}
