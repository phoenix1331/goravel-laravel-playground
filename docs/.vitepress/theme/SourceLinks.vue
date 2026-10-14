<script setup>
import { computed } from 'vue'
import { useData } from 'vitepress'

const { frontmatter, theme } = useData()
const files = computed(() => frontmatter.value.sourceFiles || [])
</script>

<template>
  <div v-if="files.length" class="source-links">
    <span class="label">Code in this repo:</span>
    <a v-for="file in files" :key="file" :href="`${theme.repo}/blob/main/${file}`" target="_blank" rel="noopener">
      <code>{{ file }}</code>
    </a>
  </div>
</template>

<style scoped>
.source-links {
  display: flex;
  flex-wrap: wrap;
  gap: 6px 10px;
  align-items: center;
  margin-bottom: 24px;
  padding: 10px 14px;
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  background: var(--vp-c-bg-soft);
  font-size: 13px;
}
.label {
  color: var(--vp-c-text-2);
  font-weight: 600;
}
</style>
