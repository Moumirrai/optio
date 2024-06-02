<template>
  <v-container :fluid="true" class="fill-height flex-container">
    <v-toolbar
      :height="smallToolbar ? 100 : 200"
      color="background"
      :floating="true"
      class="position-sticky"
    >
      <v-card
        class="fill-height mx-4"
        width="100%"
        border
        rounded="lg"
        color="background"
      >
        <v-card-title>Video</v-card-title>
        <v-card-actions>
          <v-btn
            color="primary"
            @click="videoStore.addFiles()"
            :disabled="ongoingProcess"
          >
            <v-icon left>mdi-plus</v-icon>
            Add Videos
          </v-btn>
          <v-btn
            color="error"
            @click="videoStore.clear()"
            :disabled="ongoingProcess"
          >
            <v-icon left>mdi-delete</v-icon>
            Clear
          </v-btn>
        </v-card-actions>
        <v-card-subtitle v-if="videoStore.current.eta > 0">
          ETA: {{formatTime(videoStore.current.eta)}}
        </v-card-subtitle>
      </v-card>
    </v-toolbar>
    <v-container
      :fluid="true"
      class="flex-item pa-0 pt-4"
      style="overflow: hidden"
    >
      <div style="height: 100%; width: 100%; overflow: hidden">
        <file-list-video></file-list-video>
      </div>
    </v-container>
  </v-container>
</template>

<script setup lang="ts">
import { useMainStore, useVideoStore } from "@/store";
import FileListVideo from "./FileListVideo.vue";
import { ref } from "vue";

import {formatDate, formatSize, formatTime} from "@/utils/format";
import { storeToRefs } from "pinia";

const store = useMainStore();
const videoStore = useVideoStore();

const { ongoingProcess } = storeToRefs(store);

const smallToolbar = ref(false);

</script>

<style scoped lang="scss">
.flex-container {
  display: flex;
  flex-direction: column;
  flex-wrap: nowrap; /* Prevent items from wrapping to the next line */
}

.flex-item {
  flex-grow: 1;
  overflow-y: auto; /* Add a vertical scrollbar when the content is too long */
}
.scrollable {
  overflow-y: auto;
}
</style>
