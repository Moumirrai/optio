<template>
  <v-container :fluid="true" class="fill-height pa-0">
    <v-card class="mb-4 py-2 px-3 mx-4" width="100%" rounded="lg">
      <v-row>
        <v-col cols="3"><strong>Name</strong></v-col>
        <v-col cols="2"><strong>Size</strong></v-col>
        <v-col cols="2"><strong>Date Created</strong></v-col>
        <v-col><strong>New Size</strong></v-col>
        <v-col><strong>Ratio</strong></v-col>
        <v-col><strong>Complete</strong></v-col>
      </v-row>
    </v-card>
    <!--    <v-lazy
            :min-height="200"
            :options="{'threshold':0.5}"
            transition="fade-transition"
            class="fill-height"
            style="width: 100%"
        >-->
    <v-container :fluid="true" class="mx-1 pr-4 pa-0 scrollable-part" style="width: 100%">
      <v-card
        v-for="file in files"
        width="100%"
        :key="file.name"
        class="pa-3 mb-2 mr-2 ml-3 animate_border_color"
        color="background"
        border
        rounded="lg"
      >
        <v-row class="data_row">

          <v-col cols="3">
            <v-tooltip :text="file.name" open-delay="500">
              <template v-slot:activator="{ props }">
                <div v-bind="props" class="text-truncate">{{ file.name }}</div>
              </template>
            </v-tooltip>
          </v-col>
          <v-col cols="2">{{ formatSize(file.size) }}</v-col>
          <v-col>{{ formatDate(file.dateCreated).slice(0, -3) }}</v-col>
          <v-col>{{ file.converted ? formatSize(file.convertedSize) : "" }}</v-col>
          <v-col>{{ file.converted ? file.ratio + "%" : "" }}</v-col>
          <v-col class="text-center">
            <v-icon v-if="file.converted" color="primary" size="20" class="mr-2"
            >mdi-check
            </v-icon
            >
            <v-icon v-else color="white" size="20" class="mr-2">mdi-minus</v-icon>
          </v-col>
        </v-row>
      </v-card>
      <v-skeleton-loader v-for="i in 10" width="100%" color="background" height="48" class="rounded-lg mb-2 ml-3 border" v-if="loading"
                         type="list-item"></v-skeleton-loader>
    </v-container>
  </v-container>
</template>

<script setup lang="ts">
import {useImageStore} from "@/store";
import {storeToRefs} from "pinia";

import {formatDate, formatSize} from "@/utils/format";
import FileListLoader from "@/components/FileListLoader.vue";

const imageStore = useImageStore();

const {files, loading} = storeToRefs(imageStore);
</script>

<style scoped lang="scss">
.scrollable-part {
  overflow-y: scroll;
  height: calc(100% - 50px);
}

.data_row {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.animate_border_color {
  animation: animate_border_color 0.5s ease;
}

.surface_border {
  border-color: #212121;
}
</style>
