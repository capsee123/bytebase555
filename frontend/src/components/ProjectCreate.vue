<template>
  <form class="w-96 py-2 space-y-6 divide-y divide-block-border">
    <div class="grid gap-y-6 gap-x-4 grid-cols-1">
      <div class="col-span-1">
        <label for="name" class="text-base leading-6 font-medium text-control">
          {{ $t("project.create-modal.project-name") }}
          <span class="text-red-600">*</span>
        </label>
        <BBTextField
          class="mt-4 w-full"
          :required="true"
          :placeholder="'Project name'"
          :value="state.project.name"
          @input="
            state.project.name = ($event.target as HTMLInputElement).value
          "
        />
      </div>
      <div class="col-span-1">
        <label for="name" class="text-base leading-6 font-medium text-control">
          {{ $t("project.create-modal.key") }}
          <span class="text-red-600">*</span>
          <span class="ml-1 text-sm font-normal">
            {{ $t("project.create-modal.key-hint") }}
          </span>
        </label>
        <BBTextField
          class="mt-4 w-full uppercase"
          :required="true"
          :value="state.project.key"
          @input="state.project.key = ($event.target as HTMLInputElement).value"
        />
      </div>
      <div class="col-span-1">
        <div for="name" class="text-base leading-6 font-medium text-control">
          {{ $t("common.mode") }}
          <span class="text-red-600">*</span>
        </div>
        <div class="mt-2 textlabel">
          <div class="radio-set-row">
            <div class="radio">
              <input
                v-model="state.project.tenantMode"
                tabindex="-1"
                type="radio"
                class="btn"
                value="DISABLED"
              />
              <label class="label">{{ $t("project.mode.standard") }}</label>
            </div>
            <div class="radio">
              <input
                v-model="state.project.tenantMode"
                tabindex="-1"
                type="radio"
                class="btn"
                value="TENANT"
              />
              <label class="label">{{ $t("project.mode.tenant") }}</label>
              <FeatureBadge
                feature="bb.feature.multi-tenancy"
                class="text-accent"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- Create button group -->
    <div class="pt-4 flex justify-end">
      <button
        type="button"
        class="btn-normal py-2 px-4"
        @click.prevent="cancel"
      >
        {{ $t("common.cancel") }}
      </button>
      <button
        class="btn-primary ml-3 inline-flex justify-center py-2 px-4"
        :disabled="!allowCreate"
        @click.prevent="create"
      >
        {{ $t("common.create") }}
      </button>
    </div>
  </form>

  <div
    v-if="state.isCreating"
    class="absolute inset-0 bg-white/50 flex justify-center items-center"
  >
    <BBSpin />
  </div>

  <FeatureModal
    v-if="state.showFeatureModal"
    feature="bb.feature.multi-tenancy"
    @cancel="state.showFeatureModal = false"
  />
</template>

<script lang="ts">
import { computed, reactive, defineComponent } from "vue";
import { useRouter } from "vue-router";
import { isEmpty } from "lodash-es";
import { useI18n } from "vue-i18n";
import { useEventListener } from "@vueuse/core";
import { projectSlug, randomString } from "@/utils";
import { Project, ProjectCreate } from "@/types";
import {
  hasFeature,
  pushNotification,
  useUIStateStore,
  useProjectStore,
} from "@/store";

interface LocalState {
  project: ProjectCreate;
  showFeatureModal: boolean;
  isCreating: boolean;
}

export default defineComponent({
  name: "ProjectCreate",
  emits: ["dismiss"],
  setup(props, { emit }) {
    const router = useRouter();
    const { t } = useI18n();
    const projectStore = useProjectStore();

    const state = reactive<LocalState>({
      project: {
        name: "New Project",
        key: randomString(3).toUpperCase(),
        tenantMode: "DISABLED",
        dbNameTemplate: "",
      } as Project,
      showFeatureModal: false,
      isCreating: false,
    });

    useEventListener("keydown", (e) => {
      if (e.code == "Escape") {
        emit("dismiss");
      }
    });

    const allowCreate = computed(() => {
      if (isEmpty(state.project.name)) return false;

      return true;
    });

    const create = async () => {
      if (
        state.project.tenantMode == "TENANT" &&
        !hasFeature("bb.feature.multi-tenancy")
      ) {
        state.showFeatureModal = true;
        return;
      }

      try {
        state.isCreating = true;

        const createdProject = await projectStore.createProject(state.project);
        useUIStateStore().saveIntroStateByKey({
          key: "project.visit",
          newState: true,
        });

        pushNotification({
          module: "bytebase",
          style: "SUCCESS",
          title: t("project.create-modal.success-prompt", {
            name: createdProject.name,
          }),
        });

        const url = {
          path: `/project/${projectSlug(createdProject)}`,
          hash: "",
        };
        router.push(url);
        emit("dismiss");
      } finally {
        state.isCreating = false;
      }
    };

    const cancel = () => {
      emit("dismiss");
    };

    return {
      state,
      allowCreate,
      cancel,
      create,
    };
  },
});
</script>
