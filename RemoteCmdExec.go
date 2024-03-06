package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Web      WebConfig       `yaml:"web"`
	Commands []CommandConfig `yaml:"commands"`
	TLS      TLSConfig       `yaml:"tls"`
}

type WebConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	UseTLS  bool   `yaml:"useTLS"`
}

type CommandConfig struct {
	Name    string   `yaml:"name"`
	Command string   `yaml:"command"`
	Params  []string `yaml:"params"`
}

type TLSConfig struct {
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

func loadConfig(filePath string) (Config, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func handleIndex(w http.ResponseWriter, r *http.Request, config Config) {
	tmpl, err := template.New("index").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Command Executor</title>
		</head>
		<body>
			<h1>Available Commands</h1>
			<ul>
				{{range .Commands}}
					<li><a href="/run?name={{.Name}}">{{.Name}}</a></li>
				{{end}}
			</ul>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, config)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func handleRun(w http.ResponseWriter, r *http.Request, config Config) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var selectedCommand CommandConfig
	for _, cmd := range config.Commands {
		if cmd.Name == name {
			selectedCommand = cmd
			break
		}
	}

	if selectedCommand.Name == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	cmd := exec.Command(selectedCommand.Command, selectedCommand.Params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, "Command Execution Failed", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("result").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Command Result</title>
		</head>
		<body>
			<h1>Command: {{.Name}}</h1>
			<p>Command Output:</p>
			<pre>{{.Output}}</pre>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Name   string
		Output string
	}{
		Name:   selectedCommand.Name,
		Output: string(output),
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func main() {
	configPath := flag.String("c", "config.yaml", "Path to the configuration file")
	flag.Parse()

	config, err := loadConfig(*configPath)

	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleIndex(w, r, config)
	})

	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		handleRun(w, r, config)
	})

	listenAddr := fmt.Sprintf("%s:%d", config.Web.Address, config.Web.Port)

	if config.Web.UseTLS {
		log.Printf("Server is listening on https://%s\n", listenAddr)
		log.Fatal(http.ListenAndServeTLS(listenAddr, config.TLS.CertFile, config.TLS.KeyFile, nil))
	} else {
		log.Printf("Server is listening on http://%s\n", listenAddr)
		log.Fatal(http.ListenAndServe(listenAddr, nil))
	}
}
