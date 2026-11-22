/// <reference types="vite/client" />
const apiUrl = import.meta.env.VITE_API_URL;

console.log(apiUrl);

interface ImportMetaEnv {
  readonly VITE_API_URL: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
