<template>
  <div class="main-wrapper">
    <v-navigation-drawer location="right" temporary permanent width="350">
      <action-bar-image v-if="mode == 'image'"></action-bar-image>
      <action-bar-video v-else></action-bar-video>
    </v-navigation-drawer>
    <v-app-bar app flat class="app-bar">
      <v-app-bar-title>Optio</v-app-bar-title>
      <div class="d-flex align-center flex-row pa-4">
        <v-btn
          :active="mode == editorMode.IMAGE"
          :disabled="!canSwitchMode"
          density="default"
          border
          variant="text"
          width="40"
          class="mx-3"
          @click="switchTo(editorMode.IMAGE)"
          >Image</v-btn
        >
        <v-btn
          :active="mode == editorMode.VIDEO"
          :disabled="!canSwitchMode"
          density="default"
          border
          variant="text"
          width="40"
          @click="switchTo(editorMode.VIDEO)"
          >Video</v-btn
        >
      </div>
    </v-app-bar>

    <div class="editor-wrapper">
      <v-main v-resize="onresize" class="pt-0 px-0">
        <div class="pt-0 px-0" style="border-radius: 20px">
          <v-sheet
            class="scrollable pt-0 px-0"
            color="background"
            :style="`height: ${getHeight()}`"
          >
            <editor-component-image
              v-if="mode == 'image'"
            ></editor-component-image>
            <editor-component-video v-else></editor-component-video>
          </v-sheet>
        </div>
      </v-main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, nextTick, ref } from "vue";
import EditorComponentImage from "@/components/EditorComponentImage.vue";
import EditorComponentVideo from "@/components/EditorComponentVideo.vue";
import ActionBarImage from "@/components/ActionBarImage.vue";
import ActionBarVideo from "@/components/ActionBarVideo.vue";
import { useMainStore, editorMode } from "@/store";
import { storeToRefs } from "pinia";

const store = useMainStore();

const { canSwitchMode, mode } = storeToRefs(store);

function switchTo(mode: editorMode) {
  store.mode = mode;
  nextTick(onresize);
}

onMounted(() => {
  nextTick(onresize);
});

function onresize() {
  //64px is v-app-bar height in your case
  const scrollable = document.querySelector<HTMLElement>(".scrollable");
  if (scrollable) {
    scrollable.style.height = window.innerHeight - 64 + "px";
  }
}

function getHeight() {
  return window.innerHeight - 64 + "px";
}
</script>

<style scoped lang="scss">
.main-wrapper {
  background-color: #212121;
}

.editor-wrapper {
  //border-top-right-radius: 20px;
  overflow: hidden; /* to ensure the border radius applies to all child elements */
}
.app-container {
  height: 100vh;
}

.app-bar {
  width: 100%;
}

.sidebar {
  overflow-y: auto;
}

.editorWindowBackground {
  background-color: #212121;
}

.editorWindow {
  border-top-right-radius: 20px;
  overflow-y: auto;
}

.customScroll {
  overflow-y: auto;
  overflow-x: hidden;
  //height: 100%;
}

.scrollableCustom {
  overflow-y: scroll;
}
</style>
