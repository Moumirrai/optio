// Utilities
import { createPinia } from 'pinia'
import { useMainStore } from './MainStore'
import { useVideoStore } from './video/VideoStore'
import { useImageStore } from './image/ImageStore'
import { editorMode } from './types'

export default createPinia()

export { useMainStore, useImageStore, useVideoStore, editorMode }
