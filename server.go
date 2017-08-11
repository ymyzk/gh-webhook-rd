package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Server struct {
	Config *Config
	Logger *log.Logger
	Client *Client
	hooks  map[string]HookConfig
	re     *regexp.Regexp
}

type Payload struct {
	Ref string
}

func NewServer(config *Config, logger *log.Logger) (*Server, error) {
	discardLogger := log.New(ioutil.Discard, "", log.LstdFlags)
	if logger == nil {
		logger = discardLogger
	}

	client, err := NewClient(config.Rundeck.URL, config.Rundeck.AuthToken, logger)
	if err != nil {
		return nil, err
	}

	hooks := map[string]HookConfig{}
	for _, hook := range config.Hooks {
		hooks[hook.URL] = hook
	}

	server := Server{
		Config: config,
		Logger: logger,
		Client: client,
		hooks:  hooks,
		re:     regexp.MustCompile("^/webhook/([a-zA-Z0-9_-]+)$"),
	}
	return &server, nil
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	m := s.re.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	// Validate Webhook URL
	hook, ok := s.hooks[m[1]]
	if !ok {
		http.NotFound(w, r)
		return
	}

	s.Logger.Printf("Hook: %+v\n", hook)

	// Validate Content-type
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be 'application/json'", http.StatusUnsupportedMediaType)
		return
	}

	// Validate X-GitHub-Event
	if r.Header.Get("X-GitHub-Event") != "push" {
		http.Error(w, "Unsupported event", http.StatusAccepted)
		return
	}

	// Parse payload
	var payload Payload
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&payload)
	if err != nil {
		http.Error(w, "Failed to parse body", http.StatusBadRequest)
		return
	}
	s.Logger.Printf("Payload: %+v\n", payload)

	// Validate branch
	if hook.Branch == payload.Ref {
		http.Error(w, "Unsupported ref", http.StatusAccepted)
		return
	}

	// Run a job
	s.Logger.Println("Deploy requested")
	ctx := r.Context()
	ok, err = s.Client.RunJob(ctx, hook.JobID)
	if err != nil {
		http.Error(w, "Failed to run a job", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "OK")
}

func (s *Server) Start() {
	sc := s.Config.Server
	hp := fmt.Sprintf("%s:%d", sc.Host, sc.Port)
	http.HandleFunc("/webhook/", s.handler)
	http.ListenAndServe(hp, nil)
}
