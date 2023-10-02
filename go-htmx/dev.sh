#!/bin/bash
echo Starting in development mode
air & ./tailwindcss -i ./client/assets/input.css -o ./client/assets/output.css --watch
