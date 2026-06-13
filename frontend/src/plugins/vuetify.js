import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createVuetify } from 'vuetify'

// Brand theme: trustworthy blue + energetic orange (handyman feel).
export default createVuetify({
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        colors: {
          primary: '#1565C0',   // blue – trust
          secondary: '#FB8C00', // orange – action
          accent: '#00897B',
          error: '#E53935',
          success: '#2E7D32',
          surface: '#FFFFFF',
          background: '#F5F7FA',
        },
      },
    },
  },
  defaults: {
    VBtn: { rounded: 'lg', style: 'text-transform: none; font-weight: 600;' },
    VCard: { rounded: 'xl' },
  },
})
