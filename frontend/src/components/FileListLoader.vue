<template>

  <v-container :fluid="true" class="mx-1 pr- pa-3 scrollable-part" style="width: 100%">
    <v-skeleton-loader>pepe</v-skeleton-loader>
  </v-container>
</template>

<script setup lang="ts">
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
