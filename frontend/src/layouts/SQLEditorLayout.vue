<template>
  <div class="relative h-screen overflow-hidden flex flex-col">
    <BannersWrapper />
    <nav class="bg-white border-b border-block-border">
      <div class="max-w-full mx-auto px-4">
        <EditorHeader />
      </div>
    </nav>
    <!-- Suspense is experimental, be aware of the potential change -->
    <Suspense>
      <template #default>
        <ProvideSQLEditorContext>
          <router-view />
        </ProvideSQLEditorContext>
      </template>
      <template #fallback>
        <div class="flex flex-row justify-between p-4 space-x-2">
          <span class="items-center flex">Loading...</span>
          <button
            class="items-center flex justify-center btn-normal"
            @click.prevent="ping"
          >
            Ping
          </button>
        </div>
      </template>
    </Suspense>
  </div>
</template>

<script lang="ts" setup>
import ProvideSQLEditorContext from "@/components/ProvideSQLEditorContext.vue";
import { ServerInfo } from "@/types";
import { pushNotification, useActuatorStore } from "@/store";
import EditorHeader from "@/views/sql-editor/EditorHeader.vue";
import BannersWrapper from "@/components/BannersWrapper.vue";

const actuatorStore = useActuatorStore();

const ping = () => {
  actuatorStore.fetchServerInfo().then((info: ServerInfo) => {
    pushNotification({
      module: "bytebase",
      style: "SUCCESS",
      title: JSON.stringify(info),
    });
  });
};
</script>
