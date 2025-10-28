# Thaiyyal — Workflow Builder MVP (Frontend)

This repository contains a very small Next.js + React + TypeScript frontend MVP for a visual workflow builder using React Flow. The app provides a left canvas to compose a minimal workflow (two numeric inputs, an operation node, a visualization node) and a right pane that shows the generated JSON payload.

Quick start

1. Install dependencies

```bash
npm install
```

2. Run the dev server

```bash
npm run dev
```

Open the URL printed by Next (usually http://localhost:3000).

Notes

- The UI is intentionally minimal. Click nodes in the canvas to edit their values (prompts). Use the "Show payload" button to toggle the generated JSON on the right.
- This front-end generates a JSON payload of the form:

```json
{
	"nodes": [{"id": "1", "data": {...}}, ...],
	"edges": [{"id": "e1-3", "source": "1", "target": "3"}, ...]
}
```

Next steps (suggested)

- Add Tailwind/PostCSS build wiring if you want to use Tailwind utilities fully.
- Add a simple API route to accept the generated payload and evaluate or persist workflows.
