<template>
  <v-row class="fill-height pa-0 ma-0">
    <v-col class="d-flex flex-column justify-space-between ma-0 pa-0">
      <v-container class="pt-5 px-5 pb-0">
        <v-card-title class="text-h6 pa-0">Target</v-card-title>
        <v-slider
          v-model="slider"
          @end="store.setConfig()"
          max="16"
          min="1"
          density="compact"
          class="mb-0 mt-2"
          color="primary"
          :disabled="config.target == target.WEBP && config.webpOpt.lossless"
          thumb-label
          step="0.1"
        >
          <template v-slot:prepend> Bitrate </template>
          <!-- <template v-slot:append> <v-text-field
                variant="solo-filled"
                density="compact"
                v-model="slider"
                hide-details
                single-line
                disabled=""
                solo
                class="mt-2"
                style="width: 50px"
                ></v-text-field>
            <v-card
              rounded="md"
              color="#2a2a2a"
              class="py-2 text-center"
              width="50"
            >
              {{ slider + "%" }}
            </v-card>
          </template> -->
        </v-slider>
        <div v-if="config.target == target.WEBP">
          <v-switch
            color="primary"
            density="compact"
            inset
            label="Lossless"
            flat
            @update:model-value="store.setConfig()"
            v-model="config.webpOpt.lossless"
          ></v-switch>
        </div>
        <div v-if="config.target == target.JPG">
          <v-switch
            color="primary"
            density="compact"
            inset
            label="Preserve metadata"
            flat
            @update:model-value="store.setConfig()"
            v-model="config.jpegOpt.preserveMetadata"
          ></v-switch>
        </div>
      </v-container>
      <v-divider></v-divider>
      <v-container class="pt-5 px-5 pb-0">
        <v-card-title class="text-h6 pa-0">Output directory</v-card-title>
        <v-tooltip bottom>
          <template v-slot:activator="{ props }">
            <v-card
              class="pa-4 mt-2"
              color="#2a2a2a"
              v-bind="props"
              @click="store.setOutDirImages()"
            >
              {{ config.outDir }}
            </v-card>
          </template>
          <span>Select output directory</span>
        </v-tooltip>
        <v-btn
          class="mt-5"
          prepend-icon="mdi-folder-open"
          width="100%"
          color="grey-darken-3"
          @click="store.openOutDirImages()"
        >
          Open output Directory
        </v-btn>
        <v-switch
          color="primary"
          class="mt-2"
          inset
          label="Preserve file creation time"
          flat
          @update:model-value="store.setConfig()"
          v-model="config.preserveCreationTime"
        ></v-switch>
      </v-container>
      <v-spacer></v-spacer>
      <v-container class="d-flex justify-center align-center pa-5">
        <v-btn
          color="primary"
          block
          class="mb-0"
          @click="videoStore.process()"
          v-if="!isProcessing"
        >
          Process Video
          {{ videoStore.files.length ? `${videoStore.files.length} files` : "" }}</v-btn
        >
        <v-btn
          color="error"
          block
          class="mb-0"
          @click="videoStore.stopProcess()"
          v-else
        >
          Cancel
        </v-btn>
      </v-container>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { target } from "@/store/types";
import { useMainStore, useVideoStore } from "@/store";
import { storeToRefs } from "pinia";
import { computed, ref, onBeforeMount } from "vue";

const store = useMainStore();
const videoStore = useVideoStore();
const { config, isProcessing } = storeToRefs(store);

onBeforeMount(async () => {
  await store.getConfig();
});

const items = Object.values(target);

const slider = computed({
  get: () => {
    switch (config.value.target) {
      case target.JPG:
        return config.value.jpegOpt.quality;
      case target.PNG:
        return config.value.pngOpt.quality;
      case target.WEBP:
        return config.value.webpOpt.quality;
    }
  },
  set: (value) => {
    if (!store.configLoaded) return;
    switch (config.value.target) {
      case target.JPG:
        config.value.jpegOpt.quality = Math.floor(value);
        break;
      case target.PNG:
        config.value.pngOpt.quality = Math.floor(value);
        break;
      case target.WEBP:
        config.value.webpOpt.quality = Math.floor(value);
        break;
    }
  },
});
</script>

<style scoped>
.clickable {
  cursor: pointer;
}
</style>
