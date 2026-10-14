// Extends the default theme with a "Code in this repo" box at the top of
// every page that imports Go source. The file list is worked out at build
// time by transformPageData in config.mts.
import { h } from 'vue'
import DefaultTheme from 'vitepress/theme'
import SourceLinks from './SourceLinks.vue'

export default {
  extends: DefaultTheme,
  Layout() {
    return h(DefaultTheme.Layout, null, {
      'doc-before': () => h(SourceLinks),
    })
  },
}
