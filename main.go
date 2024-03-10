package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

var templates = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Shared Memory Web App</title>
</head>
<body>
    <h1>Shared Memory Web App</h1>
    <p>Shared Variable: {{.SharedVariable}}</p>
    <form action="/update" method="post">
        <label for="inputText">Update Shared Variable:</label>
        <input type="text" id="inputText" name="inputText" required>
        <button type="submit">Update</button>
    </form>
</body>
</html>
`))

type PageVariables struct {
	SharedVariable string
}

var mutex sync.Mutex

func runCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func readSharedVariable() string {
	data, err := os.ReadFile("shared_variable.txt")
	if err != nil {
		return ""
	}
	return string(data)
}

func writeSharedVariable(value string) error {
	mutex.Lock()
	defer mutex.Unlock()

	return os.WriteFile("shared_variable.txt", []byte(value), 0644)
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		SharedVariable: readSharedVariable(),
	}

	err := templates.Execute(w, pageVariables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		newValue := r.FormValue("inputText")

		err := writeSharedVariable(newValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func startHTTPServer(namespace string, port int, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHTTP)
	mux.HandleFunc("/update", handleUpdate)

	err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), mux)
	if err != nil {
		fmt.Printf("Error creating server in namespace %s: %v\n", namespace, err)
		return
	}
}

func main() {
	runCommand("ip netns delete s1")
	runCommand("ip netns delete s2")

	runCommand("ip netns add s1")
	runCommand("ip netns add s2")

	var wg sync.WaitGroup

	serverPortS1 := 8083 // Port for s1
	serverPortS2 := 8084 // Port for s2

	wg.Add(2)
	go startHTTPServer("s1", serverPortS1, &wg)
	go startHTTPServer("s2", serverPortS2, &wg)

	wg.Wait()
}
