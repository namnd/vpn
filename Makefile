.PHONY: css
css:
	tailwindcss -i ./css/styles.css -m -o ./assets/styles.css --watch
