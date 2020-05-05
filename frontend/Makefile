run:
	@rm -rf dist && npm run dev

deploy:
	@rm -rf dist
	@npm run build
	@npm run deploy:netlify
	@rm -rf dist

.PHONY: deploy run