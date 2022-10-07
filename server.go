package conn

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	Server  *echo.Echo
	HTTP2   bool
	LogFile string
	Debug   bool
}

var loggerMiddleware = middleware.LoggerConfig{
	CustomTimeFormat: "02.01.2006 15:04:05",
	Format:           "${remote_ip} ${time_custom} '${method} ${uri}' ${status} : ${bytes_in} >> ${bytes_out}\n",
}

var debugMiddleware = middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	log.Printf("\n%s\n", reqBody)
})

// Start server instance
func InitServer() *Server {
	return &Server{
		Server:  echo.New(),
		HTTP2:   true,
		LogFile: "server.log",
		Debug:   false,
	}
}

// Run Echo server
func (s Server) Start() {

	var srv http.Server

	if s.HTTP2 {
		srv = setupHTTP2(s.Server)
	} else {
		srv = http.Server{
			Addr:    ":8080",
			Handler: s.Server,
		}
	}

	if s.LogFile != "" {
		file := setupLogger(s.LogFile)
		loggerMiddleware.Output = file
		s.Server.Use(middleware.LoggerWithConfig(loggerMiddleware))
		if s.Debug {
			s.Server.Use(debugMiddleware)
		}
	}

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// Config and set logs file
func setupLogger(filepath string) (file *os.File) {

	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}

	log.SetOutput(file)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	return
}

// Enable HTTP2 to Echo server
func setupHTTP2(e *echo.Echo) http.Server {

	h2s := &http2.Server{
		MaxConcurrentStreams: 250,
		MaxReadFrameSize:     1048576,
		IdleTimeout:          10 * time.Second,
	}

	return http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(e, h2s),
	}

}
