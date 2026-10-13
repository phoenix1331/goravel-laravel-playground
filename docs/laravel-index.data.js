// Builds the Laravel snippet index at docs build time, so it never goes stale.
// For every file in docs/snippets/laravel it finds the chapter that shows
// it and the Go file shown next to it in the same code group.
import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const docs = path.dirname(fileURLToPath(import.meta.url))
const repo = 'https://github.com/phoenix1331/goravel-laravel-playground/blob/main'

export default {
  watch: ['./snippets/laravel/**/*', './chapters/*.md', './*.md'],

  load() {
    const rows = new Map()

    // Chapters first, then the other pages (such as Reading the comparisons)
    const pages = [
      ...fs.readdirSync(path.join(docs, 'chapters')).sort().map((f) => 'chapters/' + f),
      ...fs.readdirSync(docs).filter((f) => f.endsWith('.md')).sort(),
    ]

    for (const page of pages) {
      const source = fs.readFileSync(path.join(docs, page), 'utf8')
      const title =
        (source.match(/^title: "(.+)"$/m) || [])[1] || (source.match(/^# (.+)$/m) || [])[1] || page
      const link = '/' + page.replace(/\.md$/, '')

      // Each code group pairs Laravel snippets with Go imports
      const groups = source.split('::: code-group').slice(1).map((part) => part.split(':::')[0])
      const goInChapter = [...source.matchAll(/<<< @\/\.\.\/([^#{\s]+\.go)/g)].map((m) => m[1])

      for (const match of source.matchAll(/<<< @\/snippets\/laravel\/(ch\d\d\/[^\s\[{]+)/g)) {
        const snippet = match[1]
        if (rows.has(snippet)) continue

        const group = groups.find((g) => g.includes(snippet))
        const goInGroup = group ? [...group.matchAll(/<<< @\/\.\.\/([^#{\s]+\.go)/g)].map((m) => m[1]) : []
        const goFiles = [...new Set(goInGroup.length ? goInGroup : goInChapter)]

        rows.set(snippet, {
          snippet,
          snippetUrl: `${repo}/docs/snippets/laravel/${snippet}`,
          chapter: title,
          chapterLink: link,
          go: goFiles.map((f) => ({ file: f, url: `${repo}/${f}` })),
        })
      }
    }

    return [...rows.values()].sort((a, b) => a.snippet.localeCompare(b.snippet))
  },
}
