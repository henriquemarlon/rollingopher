import { defineConfig } from 'vocs'

export default defineConfig({
  basePath: '/rollingopher',
  baseUrl: 'https://henriquemarlon.github.io',
  description: 'Rollingopher High Level Framework Documentation',
  title: 'Rollingopher Docs',
  sidebar: [
    {
      text: 'Getting Started',
      link: '/getting-started',
    },
    {
      text: 'Example',
      link: '/example',
    },
  ],
})
