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
        <v-card-title>Pepe</v-card-title>
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

import { formatSize } from "@/utils/format";
import { storeToRefs } from "pinia";

const store = useMainStore();
const videoStore = useVideoStore();

const { ongoingProcess } = storeToRefs(store);

function test(e: Event) {
  console.log(e);
}

const lastScrollTop = ref(0);
const smallToolbar = ref(false);

/* const handleScroll = (e: Event) => {
  const target = e.target as HTMLElement;
  const scrollTop = target.scrollTop;

  if (Math.abs(scrollTop - lastScrollTop.value) > 100) {
    if (scrollTop > lastScrollTop.value) {
      smallToolbar.value = true;
    } else {
      smallToolbar.value = false;
    }
    lastScrollTop.value = scrollTop <= 0 ? 0 : scrollTop; // For Mobile or negative scrolling
  }
}; */

//
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
