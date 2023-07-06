/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_WS_ENDPOINT: string
  readonly VITE_WS_PROTOCOL: string
  // more env variables...
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
