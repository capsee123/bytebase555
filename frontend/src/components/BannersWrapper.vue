<template>
  <template v-if="shouldShowDemoBanner">
    <BannerDemo />
  </template>
  <template v-if="shouldShowDebugBanner">
    <BannerDebug />
  </template>
  <template v-if="shouldShowSubscriptionBanner">
    <BannerSubscription />
  </template>
  <template v-if="shouldShowReadonlyBanner">
    <div
      class="px-3 py-1 w-full text-lg font-medium bg-yellow-500 text-white flex justify-center items-center"
    >
      {{ $t("banner.readonly") }}
    </div>
  </template>
  <template v-if="shouldShowExternalUrlBanner">
    <BannerExternalUrl />
  </template>
</template>

<script lang="ts" setup>
import { storeToRefs } from "pinia";
import { computed } from "vue";
import {
  useActuatorStore,
  useCurrentUser,
  useDebugStore,
  useSubscriptionStore,
} from "@/store/modules";
import { hasWorkspacePermission, isDev } from "@/utils";
import BannerDemo from "@/views/BannerDemo.vue";
import BannerDebug from "@/views/BannerDebug.vue";
import BannerExternalUrl from "@/views/BannerExternalUrl.vue";
import BannerSubscription from "@/views/BannerSubscription.vue";

const actuatorStore = useActuatorStore();
const currentUser = useCurrentUser();
const debugStore = useDebugStore();
const subscriptionStore = useSubscriptionStore();

const { isDemo, isReadonly, needConfigureExternalUrl } =
  storeToRefs(actuatorStore);
const { isDebug } = storeToRefs(debugStore);
const { isExpired, isTrialing } = storeToRefs(subscriptionStore);

const shouldShowDemoBanner = computed(() => {
  // Only show demo banner if it's the default demo (as opposed to the feature demo).
  return actuatorStore.serverInfo?.demoName == "default";
});

// For now, debug mode is a global setting and will affect all users.
// So we only allow DBA and Owner to toggle it and thus show a banner
// reminding them to turn off
const shouldShowDebugBanner = computed(() => {
  return (
    isDebug.value &&
    hasWorkspacePermission(
      "bb.permission.workspace.debug",
      currentUser.value.role
    )
  );
});

const shouldShowSubscriptionBanner = computed(() => {
  return isExpired.value || isTrialing.value;
});

const shouldShowReadonlyBanner = computed(() => {
  return !isDemo.value && isReadonly.value;
});

const shouldShowExternalUrlBanner = computed(() => {
  return !isDev && needConfigureExternalUrl.value;
});
</script>
