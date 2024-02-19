# go-htmx

### Setup
First, install [Air](https://github.com/cosmtrek/air)
```bash
go install github.com/cosmtrek/air@latest
```
Run the app with Air
```bash
cd go-htmx
```
The first time you run with Air, use:
```bash
air -c .air.toml
```
Afterward, you can use:
```bash
cd go-htmx
air
```
Make sure you have Tailwind v3.4.1 installed via npm

Start the Tailwind processor
```bash
cd go-htmx
npx tailwindcss -i ./app/assets/input.css -o ./app/assets/output.css --watch
```