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

func main() {
	// ----------------------------------------------------------------------------------
	// HANDLER FOR THE MAIN PAGE
	// ----------------------------------------------------------------------------------
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html")) // Assumes index.html is in the root
		tmpl.Execute(w, nil)
	})

	// ----------------------------------------------------------------------------------
	// HANDLERS FOR HTMX EXERCISES (Re-ordered from simplest to hardest)
	// ----------------------------------------------------------------------------------

	// Exercise 1: Click to Change Text (NEW)
	http.HandleFunc("/exercise1", func(w http.ResponseWriter, r *http.Request) {
		// This replaces the button that was clicked
		fmt.Fprint(w, `<button class="btn btn-success" hx-post="/exercise1" hx-swap="outerHTML">Clicked! ✅</button>`)
	})
	http.HandleFunc("/exercise1/reset", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<button class="btn btn-primary" hx-post="/exercise1" hx-swap="outerHTML">Click Me</button>`)
	})

	// Exercise 2: Simple Click to Load
	http.HandleFunc("/exercise2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, HTMX! This content was loaded from the server. 🎉")
	})
	http.HandleFunc("/exercise2/reset", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "") // Reset to an empty div
	})

	// Exercise 3: Polling for Updates (NEW)
	http.HandleFunc("/exercise3", func(w http.ResponseWriter, r *http.Request) {
		// Return the current server time
		fmt.Fprintf(w, "Server time is: <strong>%s</strong>", time.Now().Format("03:04:05 PM"))
	})
	http.HandleFunc("/exercise3/reset", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Loading server time...")
	})

	// Exercise 4: Echo User Input
	http.HandleFunc("/exercise4", func(w http.ResponseWriter, r *http.Request) {
		userInput := r.URL.Query().Get("user-input")
		fmt.Fprintf(w, "You typed: <strong>%s</strong>", userInput)
	})
	http.HandleFunc("/exercise4/reset", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "")
	})

	// Exercise 5: Form Submission
	http.HandleFunc("/exercise5/submit", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second) // Simulate server work
		name := r.PostFormValue("name")
		log.Println("Received form submission:", name)
		fmt.Fprintf(w, `<div class="alert alert-success" id="ex5-response">Thank you, %s! Your message has been received.</div>`, name)
	})
	http.HandleFunc("/exercise5/reset", func(w http.ResponseWriter, r *http.Request) {
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
	})

	// Exercise 6: Click to Edit
	http.HandleFunc("/exercise6/contact/1", func(w http.ResponseWriter, r *http.Request) {
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
	})
	http.HandleFunc("/exercise6/reset", func(w http.ResponseWriter, r *http.Request) {
		data := Contact{Name: "Jane Doe", Email: "jane.doe@example.com"}
		tmpl, _ := template.New("contact-view").Parse(contactViewTmpl)
		tmpl.Execute(w, data)
	})

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
		<button class="btn btn-secondary btn-sm" hx-get="/exercise6/contact/1">Cancel</button>
	</form>
</div>`

// Helper function to register all /code/* endpoints
func addCodeEndpoints() {
	http.HandleFunc("/code/exercise1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<button class="btn btn-primary"
        hx-post="/exercise1"
        hx-swap="outerHTML">
    Click Me
</button>
`)
	})

	http.HandleFunc("/code/exercise2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<button class="btn btn-primary"
        hx-get="/exercise2"
        hx-target="#ex2-target">
    Load Content
</button>
<div id="ex2-target" class="mt-3 p-3 bg-light rounded border"></div>
`)
	})

	http.HandleFunc("/code/exercise3", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<div class="alert alert-info"
     hx-get="/exercise3"
     hx-trigger="every 2s">
    Loading server time...
</div>
`)
	})

	http.HandleFunc("/code/exercise4", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<input type="text" class="form-control"
       name="user-input"
       hx-get="/exercise4"
       hx-trigger="keyup changed delay:500ms"
       hx-target="#ex4-output"
       placeholder="Type here...">
<div class="mt-2">Server response: <strong id="ex4-output"></strong></div>
`)
	})

	http.HandleFunc("/code/exercise5", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<form hx-post="/exercise5/submit"
      hx-target="#ex5-response"
      hx-swap="outerHTML"
      hx-indicator="#ex5-indicator">
    ...
</form>
`)
	})

	http.HandleFunc("/code/exercise6", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<div hx-target="this" hx-swap="outerHTML">
  ...
  <button hx-get="/exercise6/contact/1">Edit</button>
</div>

<form hx-put="/exercise6/contact/1">
  ...
  <button type="submit">Save</button>
  <button hx-get="/exercise6/contact/1">Cancel</button>
</form>
`)
	})
}