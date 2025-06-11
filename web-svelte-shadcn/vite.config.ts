import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import path from 'path';

export default defineConfig({
    plugins: [tailwindcss(), svelte()],
    base: '/osc/',
    resolve: {
        alias: {
            $lib: path.resolve("./src/lib"),
        },
    },
});
