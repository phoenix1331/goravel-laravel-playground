import { defineConfig } from 'vitepress'

const repo = 'https://github.com/phoenix1331/goravel-laravel-playground'

export default defineConfig({
  title: 'Goravel playground for Laravel devs',
  description: 'A hands-on Goravel playground for Laravel developers',
  lang: 'en-GB',
  base: '/goravel-laravel-playground/',
  cleanUrls: true,
  lastUpdated: true,
  themeConfig: {
    nav: [
      { text: 'Go for PHP devs', link: '/go-for-php-devs' },
      { text: 'Chapters', link: '/chapters/01-installation' },
      { text: 'Cheat sheet', link: '/cheat-sheet' }
    ],
    sidebar: [
      {
        text: 'Start here',
        items: [
          { text: 'Go for PHP devs', link: '/go-for-php-devs' },
          { text: 'Reading the comparisons', link: '/reading-the-comparisons' },
          { text: 'Cheat sheet', link: '/cheat-sheet' }
        ]
      },
      {
        text: 'Core web',
        items: [
          { text: '01. Installation and structure', link: '/chapters/01-installation' },
          { text: '02. Configuration', link: '/chapters/02-configuration' },
          { text: '03. Routing', link: '/chapters/03-routing' },
          { text: '04. Controllers', link: '/chapters/04-controllers' },
          { text: '05. Service container and providers', link: '/chapters/05-service-container' },
          { text: '06. Facades', link: '/chapters/06-facades' },
          { text: '07. Form requests and validation', link: '/chapters/07-validation' },
          { text: '08. Middleware', link: '/chapters/08-middleware' }
        ]
      },
      {
        text: 'Data layer',
        items: [
          { text: '09. ORM models', link: '/chapters/09-models' },
          { text: '10. Relationships', link: '/chapters/10-relationships' },
          { text: '11. Migrations', link: '/chapters/11-migrations' },
          { text: '12. Factories and seeders', link: '/chapters/12-factories-seeders' },
          { text: '13. Query builder and raw SQL', link: '/chapters/13-query-builder' }
        ]
      },
      {
        text: 'Background and console',
        items: [
          { text: '14. Artisan commands', link: '/chapters/14-console-commands' },
          { text: '15. Events and listeners', link: '/chapters/15-events' },
          { text: '16. Jobs and queues', link: '/chapters/16-queues' },
          { text: '17. Task scheduling', link: '/chapters/17-scheduling' },
          { text: '18. Cache', link: '/chapters/18-cache' }
        ]
      },
      {
        text: 'Security',
        items: [
          { text: '19. Authentication', link: '/chapters/19-authentication' },
          { text: '20. Gates and policies', link: '/chapters/20-authorisation' }
        ]
      },
      {
        text: 'Integration and testing',
        items: [
          { text: '21. HTTP client', link: '/chapters/21-http-client' },
          { text: '22. Logging', link: '/chapters/22-logging' },
          { text: '23. Testing', link: '/chapters/23-testing' },
          { text: '24. Error handling', link: '/chapters/24-error-handling' }
        ]
      }
    ],
    socialLinks: [{ icon: 'github', link: repo }],
    editLink: {
      pattern: `${repo}/edit/main/docs/:path`,
      text: 'Edit this page on GitHub'
    },
    search: { provider: 'local' },
    outline: 'deep'
  }
})
