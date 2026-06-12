import express from 'express';
import { createServer } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const app = express();
const PORT = 5173;

let vite;

async function start() {
	vite = await createServer({
		plugins: [svelte()],
		server: { middlewareMode: true }
	});

	app.use(vite.middlewares);

	// Serve client-side assets
	app.use(express.static(path.join(__dirname, 'src')));

	// Fallback to index.html for SPA routing
	app.use((req, res) => {
		const appHtml = path.join(__dirname, 'src/app.html');
		if (fs.existsSync(appHtml)) {
			res.setHeader('Content-Type', 'text/html');
			res.send(fs.readFileSync(appHtml, 'utf-8'));
		} else {
			res.status(404).send('Not found');
		}
	});

	app.listen(PORT, () => {
		console.log(`\n✓ Frontend dev server ready at http://localhost:${PORT}\n`);
	});
}

start().catch((err) => {
	console.error(err);
	process.exit(1);
});
