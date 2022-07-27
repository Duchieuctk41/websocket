package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *hubImpl) serveHTTP() {
	tag := "[hubImpl.serveHTTP]"
	log := h.log.WithField("tag", tag)

	mux := http.NewServeMux()
	s := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", h.config.Port),
		Handler: mux,
	}

	mux.HandleFunc("/ws", h.wsHandler)

	h.httpSever = s
	logrus.Infof("start listening at %q", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Errorf("failed to start http server: %v", err)
	}
}
