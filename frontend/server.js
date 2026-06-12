import express from 'express';
import { createServer } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

const app = express();
const port = 5173;

const vite = await createServer({
	plugins: [svelte()],
	server: { middlewareMode: true }
});

app.use(vite.middlewares);
app.use(express.static('src'));

app.listen(port, () => {
	console.log(`Frontend running at http://localhost:${port}`);
});
