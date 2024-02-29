<template>
  <v-container class="main-wrapper fill-height pa-0" :fluid="true">
    <v-navigation-drawer
      location="right"
      temporary
      permanent
      width="350"
      color="background"
    >
      <action-bar-image v-if="mode == 'image'"></action-bar-image>
      <action-bar-video v-else></action-bar-video>
    </v-navigation-drawer>
    <v-toolbar color="background" class="border-bottom">
      <v-toolbar-title class="text-h5 font-weight-bold">Optio</v-toolbar-title>
      <v-card class="pa-1 mr-10" rounded="lg" :disabled="!canSwitchMode">
        <v-btn
          variant="flat"
          :color="mode == editorMode.IMAGE ? 'background' : 'surface'"
          rounded="lg"
          @click="switchTo(editorMode.IMAGE)"
        >image
        </v-btn
        >
        <v-btn
          variant="flat"
          :color="mode == editorMode.VIDEO ? 'background' : 'surface'"
          rounded="lg"
          @click="switchTo(editorMode.VIDEO)"
        >video
        </v-btn
        >
      </v-card>
      <v-btn icon="mdi-cog" rounded="lg" class="pa-1" height="44" width="44" color="grey-lighten-1" @click="store.configModal = true">
      </v-btn>
    </v-toolbar>

    <config-component></config-component>

    <v-container :class="$style.file_list" class="px-0">
      <keep-alive>
        <editor-component-image v-if="mode == editorMode.IMAGE"></editor-component-image>
        <editor-component-video v-else></editor-component-video>
      </keep-alive>

    </v-container>
  </v-container>
</template>

<script setup lang="ts">
import EditorComponentImage from "@/components/EditorComponentImage.vue";
import EditorComponentVideo from "@/components/EditorComponentVideo.vue";
import ActionBarImage from "@/components/ActionBarImage.vue";
import ActionBarVideo from "@/components/ActionBarVideo.vue";
import {editorMode, useMainStore} from "@/store";
import {storeToRefs} from "pinia";
import {ref} from "vue";
import ConfigComponent from "@/components/ConfigComponent.vue";

const store = useMainStore();
//generate list of 1000 numbers
const numbers = Array.from(Array(200).keys());

const {canSwitchMode, mode} = storeToRefs(store);

function switchTo(newMode: editorMode) {
  mode.value = newMode;
}
</script>

<style lang="scss" module>
.file_list {
  overflow-y: auto;
  height: calc(100vh - 64px);
}
</style>
