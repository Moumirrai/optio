<!--<template>
<v-container class="editorContainer">
    <v-responsive class="fill-height">
      <v-card outlined class="pb-1 fixed-card">
        &lt;!&ndash; <v-progress-linear
          :indeterminate="true"
          location="bottom center"
          color="amber"
          height="25"
        ></v-progress-linear> &ndash;&gt;

        <v-row>
          <v-col cols="6">
            <v-card-text>
              <div class="text-subtitle-1">
                Selected: {{ imageStore.files.length }}
              </div>
              <div class="text-subtitle-1">
                Total Size:
                {{
                  imageStore.totalSize
                    ? formatSize(imageStore.totalSize)
                    : "0 KB"
                }}
              </div>
            </v-card-text>
            <v-card-actions>
              <v-btn
                color="primary"
                @click="imageStore.addFiles()"
                :disabled="isProcessing"
              >
                <v-icon left>mdi-plus</v-icon>
                Add Images
              </v-btn>
              <v-btn
                color="error"
                @click="imageStore.clear()"
                :disabled="isProcessing"
              >
                <v-icon left>mdi-delete</v-icon>
                Clear
              </v-btn>
            </v-card-actions>
          </v-col>
          <v-col cols="6">
            &lt;!&ndash; Reserved for statistics after conversion &ndash;&gt;
          </v-col>
        </v-row>
        <v-expand-transition>
          <v-progress-linear
            v-if="progress > 0"
            color="primary"
            height="10"
            v-model="progress"
            absolute
            location="bottom"
          ></v-progress-linear>
        </v-expand-transition>
      </v-card>
      <v-card class="mt-5 scrollable-list">
        <file-list-image></file-list-image>
      </v-card>
    </v-responsive>
  </v-container>
</template>-->
<template>
  <v-container :fluid="true" class="fill-height flex-container">
    <v-toolbar
      height="200"
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
        <v-card-title>Images</v-card-title>
        <v-card-actions>
          <v-btn
            color="primary"
            @click="imageStore.addFiles()"
            :disabled="ongoingProcess"
          >
            <v-icon left>mdi-plus</v-icon>
            Add Images
          </v-btn>
          <v-btn
            color="error"
            @click="imageStore.clear()"
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
        <file-list-image></file-list-image>
      </div>
    </v-container>
  </v-container>
</template>

<script setup lang="ts">
//import usestore
import { useMainStore, useImageStore } from "@/store";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import FileListImage from "./FileListImage.vue";

import { formatSize } from "@/utils/format";
import { storeToRefs } from "pinia";
import { ref, onMounted, onBeforeUnmount } from "vue";
import FileListVideo from "@/components/FileListVideo.vue";

const store = useMainStore();
const imageStore = useImageStore();

const { isProcessing, progress, ongoingProcess } = storeToRefs(store);

EventsOn("debug:test", (files) => {
  console.log(files);
});

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
