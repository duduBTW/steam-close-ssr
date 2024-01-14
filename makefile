run:
	npx tailwindcss -i ./static/input.css -o ./static/output.css
	templ generate
	go run main.go
