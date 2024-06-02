<template>
  <v-row class="fill-height pa-0 ma-0">
    <v-col class="d-flex flex-column justify-space-between ma-0 pa-0">
      <v-container class="pt-4 px-5 pb-6">
        <v-card-title class="text-h6 pa-0">Target</v-card-title>

        <v-select
          :items="codecs"
          v-model="config.videoOpt.codec"
          @update:model-value="store.setConfig()"
          density="compact"
          label="Codec"
          class="mt-2"
          :disabled="isProcessing"
          variant="solo-inverted"
          flat
        >
        </v-select>
        <v-switch
          color="primary"
          density="compact"
          inset
          label="Percentage bitrate mode"
          flat
          @update:model-value="store.setConfig()"
          v-model="config.videoOpt.percentageMode"
        ></v-switch>
        <v-slider
          v-if="config.videoOpt.percentageMode"
          v-model="config.videoOpt.percentage"
          @end="store.setConfig()"
          max="100"
          min="1"
          density="compact"
          class="mb-0"
          thumb-label
          step="1"
          hide-details
        >
          <template v-slot:prepend> Bitrate </template>
          <template v-slot:thumb-label="{ modelValue }">
            {{ modelValue }}%
          </template>
        </v-slider>
        <v-row v-else>
          <v-col cols="6">
            <v-text-field
              v-model="bitrate"
              @input="store.setConfig()"
              label="Bitrate"
              type="number"
              variant="solo"
              density="compact"
              hide-details
              dense
              outlined
              :disabled="config.videoOpt.percentageMode || isProcessing"
            ></v-text-field>
          </v-col>
          <v-col cols="6" class="text-center my-auto">
            {{formatSize(config.videoOpt.bitrate * 1000)}}/s
          </v-col>
        </v-row>
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
import {editorMode, target, codec} from "@/store/types";
import { useMainStore, useVideoStore } from "@/store";
import { storeToRefs } from "pinia";
import { computed, ref, onBeforeMount } from "vue";
import {formatSize} from "@/utils/format";

const store = useMainStore();
const videoStore = useVideoStore();
const { config, isProcessing } = storeToRefs(store);

onBeforeMount(async () => {
  await store.getConfig();
});

const codecs = Object.values(codec);

function setPercentageMode(value: boolean) {
  config.value.videoOpt.percentageMode = value;
  store.setConfig();
}

const bitrate = computed({
  get: () => {
    return config.value.videoOpt.bitrate.toString();
  },
  set: (value: string) => {
    const intValue = parseInt(value, 10);
    if (!isNaN(intValue) && intValue > 0) {
      config.value.videoOpt.bitrate = intValue;
      store.setConfig();
    }
  },
});
</script>

<style scoped>
.clickable {
  cursor: pointer;
}
</style>
