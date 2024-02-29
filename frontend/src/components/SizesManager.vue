<template>
  <v-dialog width="500" v-model="sizeModal" scrim="#000" opacity>

    <v-card color="background" border rounded="lg" class="animate_size">
      <v-card-title>
        <v-row class="px-4 py-5 text-h5">
          Sizes
          <v-spacer></v-spacer>
          <v-btn icon="mdi-close" variant="plain" rounded="lg" @click="sizeModal = false">
          </v-btn>
        </v-row>
      </v-card-title>
      <v-window v-model="step">
        <v-window-item :value="1">
          <v-container class="px-5 pt-0">
            <v-sheet max-height="350" style="overflow-y: auto" rounded="lg" color="transparent">
              <v-card v-for="(size, index) in config.sizes " :key="size.name" border color="background" rounded="lg"
                      class="pa-2 pl-5 mb-2">
                <v-row>
                  <v-col cols="3" class="text-truncate">
                    {{ size.name }}
                  </v-col>
                  <v-col cols="4" class="text-truncate">
                    {{ size.width }}x{{ size.height }}
                  </v-col>
                  <v-col cols="2">
                    {{ resizeStrategy[size.strategy] }}
                  </v-col>
                  <v-col>
                    <v-icon size="20" class="mr-2" @click="editSize(index)">mdi-pencil</v-icon>
                  </v-col>
                  <v-col>
                    <v-icon color="red" size="20" class="mr-2" @click="store.removeSize(index)">mdi-delete</v-icon>
                  </v-col>
                </v-row>
              </v-card>
            </v-sheet>
          </v-container>
        </v-window-item>
        <v-window-item :value="2">
          <v-container class="px-5 pt-2">
            <v-form
              ref="form"
            >
              <v-text-field
                v-model="currentSize.name"
                label="Name"
                variant="outlined"
                rounded="lg"
                :rules="nameRule"
                required
              ></v-text-field>
              <v-row class="px-3 py-5 align-center">
                <v-text-field
                  v-model="currentSize.width"
                  label="Width"
                  variant="outlined"
                  rounded="lg"
                  :rules="sizeRule"
                  hide-details
                  required
                ></v-text-field>
                <v-icon class="mx-5">
                  mdi-close
                </v-icon>
                <v-text-field
                  v-model="currentSize.height"
                  required
                  variant="outlined"
                  rounded="lg"
                  hide-details
                  label="Height"
                  :rules="sizeRule"
                ></v-text-field>
              </v-row>
              <v-select :items="resizeStrategySelect" class="pt-5" item-title="text" item-value="value" variant="outlined"
                        rounded="lg" label="Resize strategy" v-model="currentSize.strategy" hide-details
              ></v-select>
            </v-form>
          </v-container>
        </v-window-item>
      </v-window>

      <v-card-actions class="px-4 pb-4">
        <v-btn
          v-if="step > 1"
          variant="text"
          @click="back()"
        >
          Back
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn
          v-if="step < 2"
          color="primary"
          variant="flat"
          @click="createNew()"
        >
          Create new
        </v-btn>
        <v-btn
          v-else
          color="primary"
          variant="flat"
          @click="validateAndSave"
        >
          Save
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import {useMainStore} from "@/store";
import {resizeStrategy, Size} from "@/store/types";
import {storeToRefs} from "pinia";
import {reactive, ref} from "vue";

const form = ref();
const step = ref(1);

const store = useMainStore();
const {sizeModal, config} = storeToRefs(store);


const illegalChars = ['\\', '/', ':', '*', '?', '"', '<', '>', '|'];

const illegalCharsRegex = new RegExp(`[${illegalChars.map(char => '\\' + char).join('')}]`);

const currentSize: Size = reactive({
  strategy: resizeStrategy.Fill,
  width: 1920,
  height: 1080,
  name: '',
})

const isEditing = ref(false);
let backupSize: Size;

const resizeStrategySelect = [
  {
    text: "Fill",
    value: resizeStrategy.Fill,
  },
  {
    text: 'Fit',
    value: resizeStrategy.Fit,
  },
  {
    text: 'Smart',
    value: resizeStrategy.Smart,
  }
]

const sizeRule = [
  (value: number) => {
    if (isNaN(value) || !Number.isInteger(Number(value))) return 'Only whole numbers are allowed.'
    if (value < 1) return 'Must be greater than 0.'
    if (value > 10000) return 'Must be less than 10000.'
    return true;
  },
]

const nameRule = [
  (value: string) => {
    if (value.length < 1) return 'You must enter a name.'
    if (config.value.sizes.find(size => size.name === value)) return 'Name must be unique.'
    if (value.length > 100) return 'Name must be less than 100 characters long.'
    if (illegalCharsRegex.test(value)) return 'Name cannot contain any of the following characters: \\ / : * ? " < > |'
    return true;
  },
]

function back() {
  if (isEditing.value == true) {
    store.addSize(backupSize);
    isEditing.value = false;
  }
  step.value = 1;
}

function createNew() {
  currentSize.name = '';
  currentSize.width = 1920;
  currentSize.height = 1080;
  currentSize.strategy = resizeStrategy.Fill;
  isEditing.value = false;
  step.value = 2;
}

async function validateAndSave() {
  const result = await form.value.validate();
  if (result.valid) {
    console.log('valid')
    isEditing.value = false;
    store.addSize({...currentSize}) // create a new object
    step.value = 1;
  } else {
    console.log('invalid')
  }
}

function editSize(index: number) {
  const size = config.value.sizes[index];
  currentSize.name = size.name;
  currentSize.width = size.width;
  currentSize.height = size.height;
  currentSize.strategy = size.strategy;
  backupSize = {...size}; // create a new object
  isEditing.value = true;
  store.removeSize(index);
  step.value = 2;
}

</script>

<style scoped lang="scss">

</style>
