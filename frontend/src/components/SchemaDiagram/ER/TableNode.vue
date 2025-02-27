<template>
  <div
    class="absolute overflow-hidden rounded-md shadow-lg border-b border-gray-200 bg-white w-[16rem] divide-y z-[10]"
    bb-node-type="table"
    :bb-node-id="dummy ? `dummy-${idOfTable(table)}` : idOfTable(table)"
    :bb-status="tableStatus(table)"
    :style="{
      left: `${rect.x}px`,
      top: `${rect.y}px`,
    }"
  >
    <h3
      class="group font-medium leading-6 text-white px-2 py-2 rounded-t-md gap-x-1 relative text-center whitespace-pre-wrap break-words"
      :style="{
        'background-color': tableColor,
      }"
    >
      <template v-if="schema.name !== ''">
        <span :class="[isTableDropped && 'line-through']">
          {{ schema.name }}
        </span>
        <span>.</span>
      </template>
      <span :class="[isTableDropped && 'line-through']">{{ table.name }}</span>
      <span v-if="tableStatusText" class="ml-1 text-sm">
        ({{ tableStatusText }})
      </span>

      <button
        v-if="editable"
        class="invisible group-hover:visible absolute top-[50%] -mt-2 right-1 text-main bg-white/70 hover:bg-gray-200 p-0.5 rounded"
        @click="events.emit('edit-table', { schema, table })"
      >
        <heroicons-outline:pencil class="w-4 h-4" />
      </button>
    </h3>

    <table class="w-full text-sm table-fixed">
      <tr
        v-for="(column, i) in table.columns"
        :key="i"
        :bb-column-name="dummy ? `dummy-${column.name}` : column.name"
        :bb-status="columnStatus(column)"
        :class="columnClasses(column)"
      >
        <td class="w-5 py-1.5 relative">
          <heroicons-outline:key
            v-if="isPrimaryKey(table, column)"
            class="w-3 h-3 mx-auto text-amber-500"
          />
          <tabler:diamonds
            v-else-if="isIndex(table, column)"
            class="w-3 h-3 mx-auto text-gray-500"
          />
        </td>
        <td class="w-auto text-xs py-1.5">
          <div
            class="whitespace-pre-wrap break-words pr-1.5"
            :class="editable && 'hover:!text-accent'"
            @click="handleClickColumn(column, 'name')"
          >
            {{ column.name }}
          </div>
        </td>
        <td class="w-[6rem] text-xs text-gray-400 py-1.5 text-right">
          <div
            class="truncate pr-1.5"
            :class="editable && 'hover:!text-accent'"
            @click="handleClickColumn(column, 'type')"
          >
            {{ column.type }}
          </div>
        </td>
      </tr>
    </table>
  </div>
</template>

<script lang="ts" setup>
import { computed } from "vue";

import { hashCode } from "@/bbkit/BBUtil";
import {
  ColumnMetadata,
  SchemaMetadata,
  TableMetadata,
} from "@/types/proto/store/database";
import {
  useSchemaDiagramContext,
  isPrimaryKey,
  isIndex,
  useGeometry,
} from "../common";
import { VueClass } from "@/utils";

const props = withDefaults(
  defineProps<{
    schema: SchemaMetadata;
    table: TableMetadata;
  }>(),
  {}
);
const {
  dummy,
  editable,
  idOfTable,
  rectOfTable,
  schemaStatus,
  tableStatus,
  columnStatus,
  panning,
  events,
} = useSchemaDiagramContext();

const COLOR_LIST = [
  "#64748B",
  "#EF4444",
  "#F97316",
  "#EAB308",
  "#84CC16",
  "#22C55E",
  "#10B981",
  "#06B6D4",
  "#0EA5E9",
  "#3B82F6",
  "#6366F1",
  "#8B5CF6",
  "#A855F7",
  "#D946EF",
  "#EC4899",
  "#F43F5E",
];

const tableColor = computed(() => {
  const index = (hashCode(props.table.name) & 0xfffffff) % COLOR_LIST.length;
  return COLOR_LIST[index];
});

const isTableDropped = computed(() => {
  return (
    schemaStatus(props.schema) === "dropped" ||
    tableStatus(props.table) === "dropped"
  );
});

const isTableCreated = computed(() => {
  return (
    schemaStatus(props.schema) === "created" ||
    tableStatus(props.table) === "created"
  );
});

const isTableChanged = computed(() => {
  return tableStatus(props.table) === "changed";
});

const columnClasses = (column: ColumnMetadata): VueClass => {
  const classes: string[] = [];
  if (editable.value) {
    classes.push("cursor-pointer");
  }

  const status = columnStatus(column);
  if (status === "changed") {
    classes.push("text-yellow-700", "bg-yellow-50");
  } else if (status === "created") {
    classes.push("text-green-700", "bg-green-50");
  } else if (status === "dropped") {
    classes.push("text-red-700", "bg-red-50", "line-through");
  }

  return classes;
};

const handleClickColumn = (column: ColumnMetadata, target: "name" | "type") => {
  if (!editable.value) return;
  if (panning.value) return;
  events.emit("edit-column", {
    schema: props.schema,
    table: props.table,
    column,
    target,
  });
};

const tableStatusText = computed(() => {
  if (isTableCreated.value) return "Created";
  if (isTableDropped.value) return "Dropped";
  if (isTableChanged.value) return "Changed";
  return "";
});

const rect = computed(() => rectOfTable(props.table));

useGeometry(rect);
</script>
