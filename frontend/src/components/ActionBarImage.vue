<template>
  <v-row class="fill-height pa-0 ma-0">
    <v-col class="d-flex flex-column justify-space-between ma-0 pa-0">
      <v-container class="pt-4 px-5 pb-0">
        <v-card-title class="text-h6 pa-0">Target</v-card-title>
        <v-select
          :items="items"
          v-model="config.target"
          @update:model-value="store.setConfig()"
          density="compact"
          label="Target"
          class="mt-2"
          :disabled="isProcessing"
          variant="solo-inverted"
          flat
        >
        </v-select>
        <v-slider
          v-model="slider"
          @end="store.setConfig()"
          max="100"
          min="1"
          density="compact"
          class="mb-0 mt-2"
          :disabled="
            (config.target == target.WEBP && config.webpOpt.lossless) ||
            isProcessing
          "
          thumb-label
          step="1"
        >
          <template v-slot:prepend> Quality </template>
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
            :disabled="isProcessing"
          ></v-switch>
        </div>
      </v-container>
      <v-divider></v-divider>
      <v-container class="pa-5">
        <v-card-title class="text-h6 pa-0 pb-2">Resizing</v-card-title>
        <v-select
          :items="config.sizes"
          item-title="name"
          item-value="name"
          v-model="config.activeSize"
          clearable
          @update:model-value="store.setConfig()"
          variant="solo-inverted"
          flat
          label="Size"
          density="compact"
          class="mt-2"
          :disabled="isProcessing"
        ></v-select>
        <v-btn
            width="100%"
            @click="store.sizeModal = true"
            variant="text"
            border
        >
          edit sizes
        </v-btn>
      </v-container>
      <v-divider></v-divider>
      <v-container class="pt-5 px-5 pb-0">
        <v-card-title class="text-h6 pa-0">Output directory</v-card-title>
        <v-tooltip bottom>
          <template v-slot:activator="{ props }">
            <v-card
              class="pa-4 mt-2 text-truncate"
              v-bind="props"
              @click="store.setOutDirImages()"
              :disabled="isProcessing"
            >
              {{ config.outDir }}
            </v-card>
          </template>
          <span>Select output directory</span>
        </v-tooltip>
        <v-btn
          class="mt-5"
          prepend-icon="mdi-folder-open"
          variant="text"
          border
          width="100%"
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
          :disabled="isProcessing"
          v-model="config.preserveCreationTime"
        ></v-switch>
      </v-container>
      <v-spacer class="pa-0 ma-0"></v-spacer>
      <v-container class="d-flex justify-center align-center pa-5 pt-0">
        <v-btn
          color="primary"
          block
          class="mb-0"
          @click="imageStore.process()"
          v-if="!isProcessing"
        >
          Process
          {{ imageStore.files.length ? `${imageStore.files.length} files` : "" }}</v-btn
        >
        <v-btn
          color="error"
          block
          class="mb-0"
          @click="imageStore.stopProcess()"
          v-else
        >
          Cancel
        </v-btn>
      </v-container>
    </v-col>
    <sizes-manager></sizes-manager>
  </v-row>
</template>

<script setup lang="ts">
import { target } from "@/store/types";
import { useMainStore, useImageStore } from "@/store";
import { storeToRefs } from "pinia";
import { computed, onBeforeMount } from "vue";
import SizesManager from "@/components/SizesManager.vue";

const store = useMainStore();
const imageStore = useImageStore();
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
