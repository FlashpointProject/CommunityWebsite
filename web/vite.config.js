import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import * as path from 'path';
import { dirname } from 'node:path'
import { fileURLToPath } from 'node:url'

const _dirname = typeof __dirname !== 'undefined'
  ? __dirname
  : dirname(fileURLToPath(import.meta.url));

// https://vitejs.dev/config/
export default defineConfig({
  base: '/',
  build: {
    minify: false,
  },
  plugins: [react()],
  css: {
    preprocessorOptions: {
      less: {
        math: 'always',
        relativeUrls: true,
        javascriptEnabled: true,
      },
    },
  },
  resolve: {
    // semantic-ui theming requires the import path of theme.config to be rewritten to our local theme.config file
    alias: {
      'semantic-ui/site': path.resolve(_dirname, './src/semantic-ui/site'),
      '../../theme.config': path.resolve(
        _dirname,
        './src/semantic-ui/theme.config'
      ),
    },
  }
});
