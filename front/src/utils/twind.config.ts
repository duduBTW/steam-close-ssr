import { defineConfig } from "@twind/core";
import presetAutoprefix from "@twind/preset-autoprefix";
import presetTailwind from "@twind/preset-tailwind";

export default defineConfig({
  presets: [presetAutoprefix(), presetTailwind()],
  /* @twind/with-web-components will use
   * hashed class names in production by default
   * If you don't want this, uncomment the next line
   */
  // hash: false,
});
