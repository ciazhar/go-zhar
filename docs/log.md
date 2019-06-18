
util
```
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func LogRequest(handler http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		logger.Info(" " +r.RemoteAddr +" "+ r.Method + " "+ r.URL.String() +" "+ strconv.Itoa(lrw.statusCode) )
		handler.ServeHTTP(w, r)
	})
}
```

use in server
```
if err := http.ListenAndServe(":"+port, util.LogRequest(handler,logger) ); err != nil {
    logger.Fatal(err.Error())
}
```
