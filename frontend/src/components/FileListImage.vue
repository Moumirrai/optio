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
    <v-container class="mx-1 mb-0 pa-0 scrollable-part" style="width: 100%">
      <v-row v-for="(item, index) in files" class="ml-0 mr-5 mt-2">
        <v-sheet width="100%" min-height="64" color="transparent">
          <v-lazy
            :options="{threshold: 0, rootMargin: '100%'}"
            transition="none"
            class="fill-height">
            <v-card
              width="100%"
              class="pa-3 mb-2 mr-2 ml-3 animate_border_color"
              color="background"
              border
              rounded="lg"
            >
              <v-row class="data_row">
                <v-col cols="3">
                  <v-tooltip :text="item.name" open-delay="500">
                    <template v-slot:activator="{ props }">
                      <div v-bind="props" class="text-truncate">{{ item.name }}</div>
                    </template>
                  </v-tooltip>
                </v-col>
                <v-col cols="2">{{ formatSize(item.size) }}</v-col>
                <v-col>{{ formatDate(item.dateCreated).slice(0, -3) }}</v-col>
                <v-col>{{ item.converted ? formatSize(item.convertedSize) : "" }}</v-col>
                <v-col>{{ item.converted ? item.ratio + "%" : "" }}</v-col>
                <v-col class="text-center">
                  <v-icon v-if="item.converted" color="primary" size="20" class="mr-2"
                  >mdi-check
                  </v-icon
                  >
                  <v-icon v-else color="white" size="20" class="mr-2">mdi-minus</v-icon>
                </v-col>
              </v-row>
            </v-card>
          </v-lazy>
        </v-sheet>
      </v-row>
    </v-container>
  </v-container>
</template>

<script setup lang="ts">
import {useImageStore} from "@/store";
import {storeToRefs} from "pinia";
import {formatDate, formatSize} from "@/utils/format";

const imageStore = useImageStore();

const {files} = storeToRefs(imageStore);
</script>

<style scoped lang="scss">
.scrollable-part {
  overflow-y: scroll;
  height: calc(100% - 56px);
}

.data_row {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.animate_border_color {
  animation: animate_border_color 0.5s ease;
}

</style>
