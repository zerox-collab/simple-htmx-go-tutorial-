package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Data for our 'Click to Edit' template (now Exercise 6)
type Contact struct {
	Name  string
	Email string
}

// CORS middleware to allow cross-origin requests
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Call the next handler
		next(w, r)
	}
}

func main() {
	// ----------------------------------------------------------------------------------
	// HANDLER FOR THE MAIN PAGE
	// ----------------------------------------------------------------------------------
	http.HandleFunc("/", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html")) // Assumes index.html is in the root
		tmpl.Execute(w, nil)
	}))

	// ----------------------------------------------------------------------------------
	// HANDLERS FOR HTMX EXERCISES (Re-ordered from simplest to hardest)
	// ----------------------------------------------------------------------------------

	// Exercise 1: Click to Change Text (NEW)
	http.HandleFunc("/exercise1", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// This replaces the button that was clicked
		fmt.Fprint(w, `<button id="ex1-target" class="btn btn-success" hx-post="/exercise1" hx-swap="outerHTML">Clicked! ✅</button>`)
	}))
	http.HandleFunc("/exercise1/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<button id="ex1-target" class="btn btn-primary" hx-post="/exercise1" hx-swap="outerHTML">Click Me</button>`)
	}))

	// Exercise 2: Simple Click to Load
	http.HandleFunc("/exercise2", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, HTMX! This content was loaded from the server. 🎉")
	}))
	http.HandleFunc("/exercise2/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "") // Reset to an empty div
	}))

	// Exercise 3: Polling for Updates (NEW)
	http.HandleFunc("/exercise3", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// Return the current server time
		fmt.Fprintf(w, "Server time is: <strong>%s</strong>", time.Now().Format("03:04:05 PM"))
	}))
	http.HandleFunc("/exercise3/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Loading server time...")
	}))

	// Exercise 4: Echo User Input
	http.HandleFunc("/exercise4", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userInput := r.URL.Query().Get("user-input")
		fmt.Fprintf(w, "You typed: <strong>%s</strong>", userInput)
	}))
	http.HandleFunc("/exercise4/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "")
	}))

	// Exercise 5: Form Submission
	http.HandleFunc("/exercise5/submit", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second) // Simulate server work
		name := r.PostFormValue("name")
		log.Println("Received form submission:", name)
		fmt.Fprintf(w, `<div class="alert alert-success" id="ex5-response">Thank you, %s! Your message has been received.</div>`, name)
	}))
	http.HandleFunc("/exercise5/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("form-reset").Parse(`
			<div id="ex5-response">
				<form hx-post="/exercise5/submit" hx-target="#ex5-response" hx-swap="outerHTML" hx-indicator="#ex5-indicator">
					<div class="mb-3">
						<label for="name" class="form-label">Name</label>
						<input type="text" id="name" name="name" class="form-control" required>
					</div>
					<button type="submit" class="btn btn-success">
						Submit <span class="spinner-border spinner-border-sm htmx-indicator" id="ex5-indicator"></span>
					</button>
				</form>
			</div>
		`))
		tmpl.Execute(w, nil)
	}))

	// Exercise 6: Click to Edit
	http.HandleFunc("/exercise6/contact/1", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		data := Contact{Name: "Jane Doe", Email: "jane.doe@example.com"}

		if r.Method == http.MethodPut {
			data.Name = r.PostFormValue("name")
			data.Email = r.PostFormValue("email")
			// Return display view
			tmpl, _ := template.New("contact-view").Parse(contactViewTmpl)
			tmpl.Execute(w, data)
			return
		}
		// Return edit form
		tmpl, _ := template.New("contact-edit").Parse(contactEditTmpl)
		tmpl.Execute(w, data)
	}))
	http.HandleFunc("/exercise6/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		data := Contact{Name: "Jane Doe", Email: "jane.doe@example.com"}
		tmpl, _ := template.New("contact-view").Parse(contactViewTmpl)
		tmpl.Execute(w, data)
	}))

	// ----------------------------------------------------------------------------------
	// CODE DISPLAY ENDPOINTS
	// ----------------------------------------------------------------------------------
	addCodeEndpoints()

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

// Templates for Exercise 6 to keep the main function clean
var contactViewTmpl = `
<div id="contact-1" class="p-2 border rounded" hx-target="this" hx-swap="outerHTML">
	<p class="mb-1"><strong>Name:</strong> {{.Name}}</p>
	<p class="mb-2"><strong>Email:</strong> {{.Email}}</p>
	<button class="btn btn-primary btn-sm" hx-get="/exercise6/contact/1">Click To Edit</button>
</div>`

var contactEditTmpl = `
<div id="contact-1" hx-target="this" hx-swap="outerHTML">
	<form class="p-2 border rounded" hx-put="/exercise6/contact/1">
		<div class="mb-2">
			<label class="form-label small">Name</label>
			<input type="text" name="name" class="form-control form-control-sm" value="{{.Name}}">
		</div>
		<div class="mb-3">
			<label class="form-label small">Email</label>
			<input type="email" name="email" class="form-control form-control-sm" value="{{.Email}}">
		</div>
		<button type="submit" class="btn btn-success btn-sm">Save</button>
		<button type="button" class="btn btn-secondary btn-sm" hx-get="/exercise6/reset" hx-target="#contact-1" hx-swap="outerHTML">Cancel</button>
	</form>
</div>`

// Helper function to register all /code/* endpoints
func addCodeEndpoints() {
	http.HandleFunc("/code/exercise1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Exercise 1: Click to Change Text</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://unpkg.com/htmx.org@1.9.12"></script>
</head>
<body>
    <div class="container mt-5">
        <h1>Exercise 1: Click to Change Text</h1>
        <p>Click the button below to see it change!</p>
        
        <!-- The button that changes when clicked -->
        <button id="main-button" class="btn btn-primary"
                hx-post="/exercise1"
                hx-swap="outerHTML">
            Click Me
        </button>
        
        <!-- Reset button to restore original state -->
        <div class="mt-3">
            <button class="btn btn-secondary"
                    hx-get="/exercise1/reset"
                    hx-target="#main-button"
                    hx-swap="outerHTML">
                Reset
            </button>
        </div>
    </div>
</body>
</html>`)
	})

	http.HandleFunc("/code/exercise2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Exercise 2: Click to Load Content</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://unpkg.com/htmx.org@1.9.12"></script>
</head>
<body>
    <div class="container mt-5">
        <h1>Exercise 2: Click to Load Content</h1>
        <p>Click the button to load content from the server into the target div.</p>
        
        <!-- Button that loads content into another element -->
        <button class="btn btn-primary"
                hx-get="/exercise2"
                hx-target="#ex2-target">
            Load Content
        </button>
        
        <!-- Target div where content will be loaded -->
        <div id="ex2-target" class="mt-3 p-3 bg-light rounded border" style="min-height: 50px;">
            <!-- Content will appear here -->
        </div>
        
        <!-- Reset button -->
        <div class="mt-3">
            <button class="btn btn-secondary"
                    hx-get="/exercise2/reset"
                    hx-target="#ex2-target">
                Reset
            </button>
        </div>
    </div>
</body>
</html>`)
	})

	http.HandleFunc("/code/exercise3", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Exercise 3: Polling for Updates</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://unpkg.com/htmx.org@1.9.12"></script>
</head>
<body>
    <div class="container mt-5">
        <h1>Exercise 3: Polling for Updates</h1>
        <p>This div automatically updates every 2 seconds with the current server time.</p>
        
        <!-- Div that polls the server every 2 seconds -->
        <div class="alert alert-info"
             hx-get="/exercise3"
             hx-trigger="load, every 2s">
            Loading server time...
        </div>
        
        <!-- Reset button -->
        <div class="mt-3">
            <button class="btn btn-secondary"
                    hx-get="/exercise3/reset"
                    hx-target=".alert">
                Reset
            </button>
        </div>
    </div>
</body>
</html>`)
	})

	http.HandleFunc("/code/exercise4", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Exercise 4: Send User Input</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://unpkg.com/htmx.org@1.9.12"></script>
</head>
<body>
    <div class="container mt-5">
        <h1>Exercise 4: Send User Input</h1>
        <p>Type in the input field below. The server will echo your input with a 500ms delay after you stop typing.</p>
        
        <!-- Input that sends data to server on keyup with debouncing -->
        <div class="mb-3">
            <label for="user-input" class="form-label">Type something:</label>
            <input type="text" 
                   id="user-input"
                   class="form-control"
                   name="user-input"
                   hx-get="/exercise4"
                   hx-trigger="keyup changed delay:500ms"
                   hx-target="#ex4-output"
                   placeholder="Type here...">
        </div>
        
        <!-- Output area where server response appears -->
        <div class="mt-2">
            Server response: <strong id="ex4-output" class="text-primary"></strong>
        </div>
        
        <!-- Reset button -->
        <div class="mt-3">
            <button class="btn btn-secondary"
                    hx-get="/exercise4/reset"
                    hx-target="#ex4-output"
                    onclick="document.getElementById('user-input').value = ''">
                Reset
            </button>
        </div>
    </div>
</body>
</html>`)
	})

	http.HandleFunc("/code/exercise5", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Exercise 5: Form Submission & Loading Indicators</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://unpkg.com/htmx.org@1.9.12"></script>
    <style>
        /* HTMX indicator styles */
        .htmx-indicator { display: none; }
        .htmx-request .htmx-indicator { display: inline-block; }
    </style>
</head>
<body>
    <div class="container mt-5">
        <h1>Exercise 5: Form Submission & Loading Indicators</h1>
        <p>Submit the form below. Notice the loading spinner that appears during submission.</p>
        
        <!-- Form container that gets replaced with response -->
        <div id="ex5-response">
            <form hx-post="/exercise5/submit"
                  hx-target="#ex5-response"
                  hx-swap="outerHTML"
                  hx-indicator="#ex5-indicator">
                
                <div class="mb-3">
                    <label for="name" class="form-label">Name</label>
                    <input type="text" 
                           id="name" 
                           name="name" 
                           class="form-control" 
                           required>
                </div>
                
                <button type="submit" class="btn btn-success">
                    Submit 
                    <!-- Loading spinner (hidden by default, shown during request) -->
                    <span class="spinner-border spinner-border-sm htmx-indicator" 
                          id="ex5-indicator"></span>
                </button>
            </form>
        </div>
        
        <!-- Reset button -->
        <div class="mt-3">
            <button class="btn btn-secondary"
                    hx-get="/exercise5/reset"
                    hx-target="#ex5-response"
                    hx-swap="outerHTML">
                Reset
            </button>
        </div>
    </div>
</body>
</html>`)
	})

	http.HandleFunc("/code/exercise6", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Exercise 6: Click To Edit</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://unpkg.com/htmx.org@1.9.12"></script>
</head>
<body>
    <div class="container mt-5">
        <h1>Exercise 6: Click To Edit</h1>
        <p>Click "Click To Edit" to switch to edit mode. The server controls the UI state.</p>
        
        <!-- Contact display/edit area - this entire div gets swapped -->
        <div id="contact-1" class="p-3 border rounded" hx-target="this" hx-swap="outerHTML">
            <!-- Display mode (initial state) -->
            <p class="mb-1"><strong>Name:</strong> Jane Doe</p>
            <p class="mb-2"><strong>Email:</strong> jane.doe@example.com</p>
            <button class="btn btn-primary btn-sm" 
                    hx-get="/exercise6/contact/1">
                Click To Edit
            </button>
        </div>
        
        <!-- Reset button -->
        <div class="mt-3">
            <button class="btn btn-secondary"
                    hx-get="/exercise6/reset"
                    hx-target="#contact-1"
                    hx-swap="outerHTML">
                Reset
            </button>
        </div>
        
        <!-- 
        Note: When you click "Click To Edit", the server will return an edit form like this:
        
        <div id="contact-1" hx-target="this" hx-swap="outerHTML">
            <form class="p-3 border rounded" hx-put="/exercise6/contact/1">
                <div class="mb-2">
                    <label class="form-label small">Name</label>
                    <input type="text" name="name" class="form-control form-control-sm" value="Jane Doe">
                </div>
                <div class="mb-3">
                    <label class="form-label small">Email</label>
                    <input type="email" name="email" class="form-control form-control-sm" value="jane.doe@example.com">
                </div>
                <button type="submit" class="btn btn-success btn-sm">Save</button>
                <button type="button" class="btn btn-secondary btn-sm" hx-get="/exercise6/reset" hx-target="#contact-1" hx-swap="outerHTML">Cancel</button>
            </form>
        </div>
        -->
    </div>
</body>
</html>`)
	})
}