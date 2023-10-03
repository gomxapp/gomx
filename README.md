# go-htmx

### Setup
First, install [Air](https://github.com/cosmtrek/air)
```bash
go install github.com/csmtrek/air@latest
```
Run the app with Air
```bash
cd go-htmx
air
```
Start the Tailwind processor
```bash
cd go-htmx
./tailwindcss -i ./client/assets/input.css -o ./client/assets/output.css --watch
```