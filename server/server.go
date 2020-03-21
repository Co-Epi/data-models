package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/wolkdb/cloudstore/log"
)

const (
	sslBaseDir     = "/etc/pki/tls/certs/wildcard"
	sslKeyFileName = "www.wolk.com.key"
	caFileName     = "www.wolk.com.bundle"

	defaultPort  uint16 = 8080
	defaultProject  = "us-west1-wlk"
	defaultInstance = "co-epi"
)

// Server manages HTTP connections
type Server struct {
	backend  *Backend
	Handler  http.Handler
	HTTPPort uint16
}

// NewServer returns an HTTP Server to handle simple-api-process-flow https://github.com/Co-Epi/data-models/blob/master/simple-api-process-flow.md
func NewServer(httpPort uint16, project, instance string) (s *Server, err error) {
	s = &Server{
		HTTPPort: httpPort,
	}
	backend, err := NewBackend(project, instance)
	if err != nil {
		return s, err
	}
	s.backend = backend

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.getConnection)
	s.Handler = mux
	go s.Start()
	return s, nil
}

func (s *Server) exposureAndSymptomsHandler(w http.ResponseWriter, r *http.Request) {
	// Read Post Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()

	// Parse body as ExposureAndSymptoms
	var payload ExposureAndSymptoms
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle ExposureAndSymptoms payload
	err = s.backend.processExposureAndSymptoms(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte("OK"))
}

func (s *Server) exposureCheckHandler(w http.ResponseWriter, r *http.Request) {
	var payload ExposureCheck

	// Read Post Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()

	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	symptoms, err := s.backend.processExposureCheck(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	symptomsJSON, err := json.Marshal(symptoms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(symptomsJSON)
}

func (s *Server) getConnection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Info("[CoEpiServer:getConnection] processing req for url path", "r.URL.Path", r.URL.Path, "method", r.Method)
	if strings.Compare(r.URL.Path, "/exposurecheck") == 0 {
		s.exposureCheckHandler(w, r)
	} else if strings.Compare(r.URL.Path, "/exposureandsymptoms") == 0 {
		s.exposureAndSymptomsHandler(w, r)
	}
}

// Start kicks off the HTTP Server
func (s *Server) Start() (err error) {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.HTTPPort),
		Handler:      s.Handler,
		ReadTimeout:  600 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	SSLKeyFile := path.Join(sslBaseDir, sslKeyFileName)
	CAFile := path.Join(sslBaseDir, caFileName)

	// Note: bringing the intermediate certs with CAFile into a cert pool and the tls.Config is *necessary*
	certpool, err := x509.SystemCertPool() // https://stackoverflow.com/questions/26719970/issues-with-tls-connection-in-golang -- instead of x509.NewCertPool()
	if err != nil {
		return fmt.Errorf("SystemCertPool: %v", err)
	}

	pem, err := ioutil.ReadFile(CAFile)
	if err != nil {
		return fmt.Errorf("Failed to read client certificate authority: %v", err)
	}
	if !certpool.AppendCertsFromPEM(pem) {
		return fmt.Errorf("Can't parse client certificate authority")
	}

	config := tls.Config{
		ClientCAs:  certpool,
		ClientAuth: tls.NoClientCert, // tls.RequireAndVerifyClientCert,
	}
	config.BuildNameToCertificate()

	srv.TLSConfig = &config

	err = srv.ListenAndServeTLS(CAFile, SSLKeyFile)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	
	port := defaultPort
	project := defaultProject
	instance := defaultInstance

	_, err := NewServer(port, project, instance)
	if err != nil {
		panic(err)
	}
	fmt.Printf("CoEpi Go Server - Listening on port %d...\n", port)
	select {}
}
