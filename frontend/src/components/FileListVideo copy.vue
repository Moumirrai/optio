<template>
  <div>
    <v-card class="mb-3 pa-3">
      <v-row>
        <v-col cols="3"><strong>Name</strong></v-col>
        <v-col cols="2"><strong>Size</strong></v-col>
        <v-col cols="3"><strong>Date Modified</strong></v-col>
        <v-col cols="1,5"><strong>Progress</strong></v-col>
        <v-col cols="1,5"><strong>Converted</strong></v-col>
      </v-row>
    </v-card>
    <div class="scrollable-part">
      <v-card v-for="file in files" :key="file.name" class="m-3 pa-3 my-2">
        <v-row class="slim-row">
          <v-col cols="3">{{ file.name }}</v-col>
          <v-col cols="2">{{ formatSize(file.size) }}</v-col>
          <v-col cols="3">{{ formatDate(file.dateCreated) }}</v-col>
          <v-col cols="1,5">{{}}</v-col>
          <v-col cols="1,5">
            <v-icon v-if="file.converted" color="green" size="20" class="mr-2"
              >mdi-check</v-icon
            >
            <v-icon v-else color="red" size="20" class="mr-2">mdi-close</v-icon>
          </v-col>
        </v-row>
        <v-expand-transition>
          <v-progress-linear
            v-if="file.progress && file.progress > 0"
            color="primary"
            height="5"
            v-model="file.progress"
            absolute
            location="bottom"
          ></v-progress-linear>
        </v-expand-transition>
        <!-- <v-row>
        <v-col cols="12">
          <div class="full-width-row">Your content here</div>
        </v-col>
      </v-row> -->
      </v-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useVideoStore } from "@/store";
import { storeToRefs } from "pinia";
import { ref } from "vue";

import { formatDate, formatSize } from "@/utils/format";

const videoStore = useVideoStore();

const progress = ref(20);

const { files } = storeToRefs(videoStore);
</script>

<style scoped lang="scss">
.slim-row {
  //make text only single line
}

.scrollable-part {
  overflow-y: scroll;
  height: 100%;
}
</style>
