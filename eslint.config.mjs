import js from "@eslint/js";
import globals from "globals";
import { defineConfig } from "eslint/config";

export default defineConfig([
    {
        files: ["**/*.{js,mjs,cjs}"],
        rules: {
            semi: ["warn", "always"],
            indent: ["error", 4],
            "no-trailing-spaces": ["error"]
        },
        plugins: { js },
        extends: ["js/recommended"],
        languageOptions: { globals: globals.browser }
    },
]);
