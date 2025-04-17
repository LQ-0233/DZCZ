// uno.config.ts
import presetAttributify from '@unocss/preset-attributify'
import presetUno from '@unocss/preset-uno'
import transformerDirectives from '@unocss/transformer-directives'
import { defineConfig } from 'unocss'


export default defineConfig({
  // ...UnoCSS options
  rules: [
    [/^bg-color-#(.*)$/, ([, color]) => ({ 'background-color': `#${color}` })],
  ],
  presets: [
    presetUno(),
    presetAttributify({
      /* preset options */
    })
  ],

  shortcuts: { 
  },
  transformers: [
    transformerDirectives(),
  ],
})