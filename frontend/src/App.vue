<template>
  <v-app>
    <v-main class="main-content fill-height">
      <EditorLayout />
    </v-main>
    <!--
    <v-main>
      <HelloWorld />
    </v-main>
    -->
  </v-app>
</template>

<script setup lang="ts">
import EditorLayout from "@/pages/EditorLayout.vue";
import { onMounted, watch } from "vue";
import { WailsEventManager } from "@/services/WailsEventManager";
import { HotkeyManager } from "@/services/HotkeyManager";
import "@/assets/styles/root.scss";
import { useMainStore, useImageStore, useVideoStore } from "@/store";

const store = useMainStore();
const imageStore = useImageStore();
const videoStore = useVideoStore();

onMounted(() => {
  new WailsEventManager(store, imageStore, videoStore);
  new HotkeyManager(store, imageStore, videoStore);
});

document.addEventListener("contextmenu", (event) => {
  //event.preventDefault();
});
</script>

<style>
html {
  overflow-y: auto;
}
</style>
