<script setup>
import { data as rows } from './laravel-index.data.js'
</script>

# Laravel snippet index

Every Laravel file in [`docs/snippets/laravel`](https://github.com/phoenix1331/goravel-laravel-playground/tree/main/docs/snippets/laravel), the chapter that uses it, and the Go code it is compared with. The list is built from the chapter pages each time the site is built.

<table>
  <thead>
    <tr><th>Laravel snippet</th><th>Go counterpart</th><th>Chapter</th></tr>
  </thead>
  <tbody>
    <tr v-for="row in rows" :key="row.snippet">
      <td><a :href="row.snippetUrl"><code>{{ row.snippet }}</code></a></td>
      <td>
        <div v-for="go in row.go" :key="go.file"><a :href="go.url"><code>{{ go.file }}</code></a></div>
        <span v-if="!row.go.length">Shown on its own</span>
      </td>
      <td><a :href="row.chapterLink">{{ row.chapter }}</a></td>
    </tr>
  </tbody>
</table>
