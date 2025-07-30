package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os" // <-- Import the "os" package
	"time"
)

// Data for our 'Click to Edit' template
type Contact struct {
	Name  string
	Email string
}

// CORS middleware to allow cross-origin requests
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, HX-Request, HX-Trigger, HX-Target, HX-Current-URL, HX-Boosted, HX-Trigger-Name, HX-Prompt")
		w.Header().Set("Access-Control-Expose-Headers", "HX-Location, HX-Push-Url, HX-Redirect, HX-Refresh, HX-Replace-Url, HX-Reswap, HX-Retarget, HX-Reselect, HX-Trigger, HX-Trigger-After-Settle, HX-Trigger-After-Swap")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	// --- Dynamic Endpoint Configuration ---
	// Check for APP_ENV and create a helper function to format endpoints.
	isProduction := os.Getenv("APP_ENV") == "production"
	baseURL := "https://simple-htmx-go-tutorial-production.up.railway.app"

	endpoint := func(path string) string {
		if isProduction {
			return baseURL + path
		}
		return path
	}

	// ----------------------------------------------------------------------------------
	// HANDLER FOR THE MAIN PAGE
	// ----------------------------------------------------------------------------------
	http.HandleFunc("/", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, nil)
	}))

	// ----------------------------------------------------------------------------------
	// HANDLERS FOR HTMX EXERCISES
	// ----------------------------------------------------------------------------------

	// Exercise 1: Click to Change Text
	http.HandleFunc("/exercise1", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<button id="ex1-target" class="btn btn-success" hx-post="%s" hx-swap="outerHTML">Clicked! ✅</button>`, endpoint("/exercise1"))
	}))
	http.HandleFunc("/exercise1/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<button id="ex1-target" class="btn btn-primary" hx-post="%s" hx-swap="outerHTML">Click Me</button>`, endpoint("/exercise1"))
	}))

	// Exercise 2: Simple Click to Load
	http.HandleFunc("/exercise2", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, HTMX! This content was loaded from the server. 🎉")
	}))
	http.HandleFunc("/exercise2/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "")
	}))

	// Exercise 3: Polling for Updates
	http.HandleFunc("/exercise3", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
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
		time.Sleep(1 * time.Second)
		name := r.PostFormValue("name")
		log.Println("Received form submission:", name)
		fmt.Fprintf(w, `<div class="alert alert-success" id="ex5-response">Thank you, %s! Your message has been received.</div>`, name)
	}))
	http.HandleFunc("/exercise5/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// Pass the dynamic URL into the template
		tmpl := template.Must(template.New("form-reset").Parse(`
            <div id="ex5-response">
                <form hx-post="{{.SubmitURL}}" hx-target="#ex5-response" hx-swap="outerHTML" hx-indicator="#ex5-indicator">
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
		tmpl.Execute(w, map[string]string{
			"SubmitURL": endpoint("/exercise5/submit"),
		})
	}))

	// Exercise 6: Click to Edit
	http.HandleFunc("/exercise6/contact/1", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// Pass dynamic URLs into the contact templates
		data := map[string]interface{}{
			"Name":      "Jane Doe",
			"Email":     "jane.doe@example.com",
			"ActionURL": endpoint("/exercise6/contact/1"),
			"ResetURL":  endpoint("/exercise6/reset"),
		}

		if r.Method == http.MethodPut {
			data["Name"] = r.PostFormValue("name")
			data["Email"] = r.PostFormValue("email")
			tmpl, _ := template.New("contact-view").Parse(contactViewTmpl)
			tmpl.Execute(w, data)
			return
		}

		tmpl, _ := template.New("contact-edit").Parse(contactEditTmpl)
		tmpl.Execute(w, data)
	}))
	http.HandleFunc("/exercise6/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Name":      "Jane Doe",
			"Email":     "jane.doe@example.com",
			"ActionURL": endpoint("/exercise6/contact/1"),
		}
		tmpl, _ := template.New("contact-view").Parse(contactViewTmpl)
		tmpl.Execute(w, data)
	}))

	addCodeEndpoints()

	// ----------------------------------------------------------------------------------
	// SERVER STARTUP
	// ----------------------------------------------------------------------------------
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default for local development
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

// Templates for Exercise 6 now use template variables for URLs
var contactViewTmpl = `
<div id="contact-1" class="p-2 border rounded" hx-target="this" hx-swap="outerHTML">
    <p class="mb-1"><strong>Name:</strong> {{.Name}}</p>
    <p class="mb-2"><strong>Email:</strong> {{.Email}}</p>
    <button class="btn btn-primary btn-sm" hx-get="{{.ActionURL}}">Click To Edit</button>
</div>`

var contactEditTmpl = `
<div id="contact-1" hx-target="this" hx-swap="outerHTML">
    <form class="p-2 border rounded" hx-put="{{.ActionURL}}">
        <div class="mb-2">
            <label class="form-label small">Name</label>
            <input type="text" name="name" class="form-control form-control-sm" value="{{.Name}}">
        </div>
        <div class="mb-3">
            <label class="form-label small">Email</label>
            <input type="email" name="email" class="form-control form-control-sm" value="{{.Email}}">
        </div>
        <button type="submit" class="btn btn-success btn-sm">Save</button>
        <button type="button" class="btn btn-secondary btn-sm" hx-get="{{.ResetURL}}" hx-target="#contact-1" hx-swap="outerHTML">Cancel</button>
    </form>
</div>`


// Helper function to register all /code/* endpoints
func addCodeEndpoints() {
	baseURL := "https://simple-htmx-go-tutorial-production.up.railway.app"

	http.HandleFunc("/code/exercise1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<!DOCTYPE html>
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
        
        <button id="ex1-target" class="btn btn-primary"
                hx-post="%s/exercise1"
                hx-swap="outerHTML">
            Click Me
        </button>
        
        <div class="mt-3">
            <button class="btn btn-secondary"
                    hx-get="%s/exercise1/reset"
                    hx-target="#ex1-target"
                    hx-swap="outerHTML">
                Reset
            </button>
        </div>
    </div>
</body>
</html>`, baseURL, baseURL)
	})

	http.HandleFunc("/code/exercise2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<!DOCTYPE html>
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
        
        <button class="btn btn-primary"
                hx-get="%s/exercise2"
                hx-target="#ex2-target">
            Load Content
        </button>
        
        <div id="ex2-target" class="mt-3 p-3 bg-light rounded border" style="min-height: 50px;">
            </div>
        
        <div class="mt-3">
            <button class="btn btn-secondary"
                    hx-get="%s/exercise2/reset"
                    hx-target="#ex2-target">
                Reset
            </button>
        </div>
    </div>
</body>
</html>`, baseURL, baseURL)
	})

	http.HandleFunc("/code/exercise3", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<!DOCTYPE html>
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
        
        <div class="alert alert-info"
             hx-get="%s/exercise3"
             hx-trigger="load, every 2s">
            Loading server time...
        </div>
        
        <div class="mt-3">
            <button class="btn btn-secondary"
                    hx-get="%s/exercise3/reset"
                    hx-target=".alert">
                Reset
            </button>
        </div>
    </div>
</body>
</html>`, baseURL, baseURL)
	})

	http.HandleFunc("/code/exercise4", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<!DOCTYPE html>
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
        
        <div class="mb-3">
            <label for="user-input" class="form-label">Type something:</label>
            <input type="text" 
                   id="user-input"
                   class="form-control"
                   name="user-input"
                   hx-get="%s/exercise4"
                   hx-trigger="keyup changed delay:500ms"
                   hx-target="#ex4-output"
                   placeholder="Type here...">
        </div>
        
        <div class="mt-2">
            Server response: <strong id="ex4-output" class="text-primary"></strong>
        </div>
        
        <div class="mt-3">
            <button class="btn btn-secondary"
                    hx-get="%s/exercise4/reset"
                    hx-target="#ex4-output"
                    onclick="document.getElementById('user-input').value = ''">
                Reset
            </button>
        </div>
    </div>
</body>
</html>`, baseURL, baseURL)
	})

	http.HandleFunc("/code/exercise5", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<!DOCTYPE html>
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
        
        <div id="ex5-response">
            <form hx-post="%s/exercise5/submit"
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
                    <span class="spinner-border spinner-border-sm htmx-indicator" 
                          id="ex5-indicator"></span>
                </button>
            </form>
        </div>
        
        <div class="mt-3">
            <button class="btn btn-secondary"
                    hx-get="%s/exercise5/reset"
                    hx-target="#ex5-response"
                    hx-swap="outerHTML">
                Reset
            </button>
        </div>
    </div>
</body>
</html>`, baseURL, baseURL)
	})

	http.HandleFunc("/code/exercise6", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<!DOCTYPE html>
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
        
        <div id="contact-1" class="p-3 border rounded" hx-target="this" hx-swap="outerHTML">
            <p class="mb-1"><strong>Name:</strong> Jane Doe</p>
            <p class="mb-2"><strong>Email:</strong> jane.doe@example.com</p>
            <button class="btn btn-primary btn-sm" 
                    hx-get="%s/exercise6/contact/1">
                Click To Edit
            </button>
        </div>
        
        <div class="mt-3">
            <button class="btn btn-secondary"
                    hx-get="%s/exercise6/reset"
                    hx-target="#contact-1"
                    hx-swap="outerHTML">
                Reset
            </button>
        </div>
        
        </div>
</body>
</html>`, baseURL, baseURL)
	})

	http.HandleFunc("/code/exercise1/go", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, `// Exercise 1: Click to Change Text
http.HandleFunc("/exercise1", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<button id=\"ex1-target\" class=\"btn btn-success\" hx-post=\"%s\" hx-swap=\"outerHTML\">Clicked! ✅</button>", endpoint("/exercise1"))
}))
http.HandleFunc("/exercise1/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<button id=\"ex1-target\" class=\"btn btn-primary\" hx-post=\"%s\" hx-swap=\"outerHTML\">Click Me</button>", endpoint("/exercise1"))
}))`)
	})

	http.HandleFunc("/code/exercise2/go", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, `// Exercise 2: Simple Click to Load
http.HandleFunc("/exercise2", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, HTMX! This content was loaded from the server. 🎉")
}))
http.HandleFunc("/exercise2/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "")
}))`)
	})

	http.HandleFunc("/code/exercise3/go", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, `// Exercise 3: Polling for Updates
http.HandleFunc("/exercise3", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Server time is: <strong>%s</strong>", time.Now().Format("03:04:05 PM"))
}))
http.HandleFunc("/exercise3/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Loading server time...")
}))`)
	})

	http.HandleFunc("/code/exercise4/go", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, `// Exercise 4: Echo User Input
http.HandleFunc("/exercise4", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
    userInput := r.URL.Query().Get("user-input")
    fmt.Fprintf(w, "You typed: <strong>%s</strong>", userInput)
}))
http.HandleFunc("/exercise4/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "")
}))`)
	})

	http.HandleFunc("/code/exercise5/go", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, `// Exercise 5: Form Submission
http.HandleFunc("/exercise5/submit", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
    time.Sleep(1 * time.Second)
    name := r.PostFormValue("name")
    log.Println("Received form submission:", name)
    fmt.Fprintf(w, "<div class=\"alert alert-success\" id=\"ex5-response\">Thank you, %s! Your message has been received.</div>", name)
}))
http.HandleFunc("/exercise5/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
    // Pass the dynamic URL into the template
    tmpl := template.Must(template.New("form-reset").Parse("\n        <div id=\"ex5-response\">\n            <form hx-post=\"{{.SubmitURL}}\" hx-target=\"#ex5-response\" hx-swap=\"outerHTML\" hx-indicator=\"#ex5-indicator\">\n                <div class=\"mb-3\">\n                    <label for=\"name\" class=\"form-label\">Name</label>\n                    <input type=\"text\" id=\"name\" name=\"name\" class=\"form-control\" required>\n                </div>\n                <button type=\"submit\" class=\"btn btn-success\">\n                    Submit <span class=\"spinner-border spinner-border-sm htmx-indicator\" id=\"ex5-indicator\"></span>\n                </button>\n            </form>\n        </div>\n    "))
    tmpl.Execute(w, map[string]string{
        "SubmitURL": endpoint("/exercise5/submit"),
    })
}))`)
	})

	http.HandleFunc("/code/exercise6/go", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, `// Exercise 6: Click to Edit
http.HandleFunc("/exercise6/contact/1", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
    // Pass dynamic URLs into the contact templates
    data := map[string]interface{}{
        "Name":      "Jane Doe",
        "Email":     "jane.doe@example.com",
        "ActionURL": endpoint("/exercise6/contact/1"),
        "ResetURL":  endpoint("/exercise6/reset"),
    }

    if r.Method == http.MethodPut {
        data["Name"] = r.PostFormValue("name")
        data["Email"] = r.PostFormValue("email")
        tmpl, _ := template.New("contact-view").Parse(contactViewTmpl)
        tmpl.Execute(w, data)
        return
    }

    tmpl, _ := template.New("contact-edit").Parse(contactEditTmpl)
    tmpl.Execute(w, data)
}))
http.HandleFunc("/exercise6/reset", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{
        "Name":      "Jane Doe",
        "Email":     "jane.doe@example.com",
        "ActionURL": endpoint("/exercise6/contact/1"),
    }
    tmpl, _ := template.New("contact-view").Parse(contactViewTmpl)
    tmpl.Execute(w, data)
}))`)
	})
}