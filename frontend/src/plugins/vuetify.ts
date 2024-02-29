/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Composables
import { createVuetify } from 'vuetify'
import { md3 } from 'vuetify/blueprints'

const backdrop = "#09090b"

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  //blueprint: md3,
  theme: {
    defaultTheme: "dark",
    themes: {
      dark: {
        colors: {
          primary: '#ffcc00',
          //primary: '#ffffff',
          background: backdrop,
          //surface: backdrop,
          "on-surface": "#ffffff",
          "on-background": "#ffffff",
          secondary: '#5CBBF6',
        },
      },
    },
  },
})
