import { defineConfig } from 'astro/config';

import remarkToc from 'remark-toc';

// https://astro.build/config
import tailwind from "@astrojs/tailwind";

// https://astro.build/config
export default defineConfig({
  site: 'https://scadagobr.com',
  integrations: [tailwind()],
  markdown: {
    remarkPlugins: [remarkToc],
    extendDefaultPlugins: true,
  },
});