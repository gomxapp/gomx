# GOMX

### Setup
First, install [Air](https://github.com/cosmtrek/air)
```bash
go install github.com/cosmtrek/air@latest
```
Run the app with Air

The first time you run with Air, use:
```bash
air -c .air.toml
```
Afterward, you can use:
```bash
air
```
Make sure you have Tailwind v3.4.1 installed via npm

Start the Tailwind processor
```bash
npx tailwindcss -i ./app/assets/input.css -o ./app/assets/output.css --watch
```