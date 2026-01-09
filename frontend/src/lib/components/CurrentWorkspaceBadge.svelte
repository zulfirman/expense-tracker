<script>
  import { workspace } from '$lib/stores/workspace';

  $: wsState = $workspace;
  $: current = wsState.list.find((w) => w.id === wsState.currentId);
</script>

{#if wsState.list.length > 0}
  <div class="mt-2 mb-3 text-base text-base-content flex items-center gap-3">
    <span class="font-bold uppercase tracking-wide text-xs">Workspace:</span>
    {#if current}
      <span class="badge badge-primary badge-md px-4 py-3 rounded-full font-semibold">
        {current.name}
      </span>
    {:else}
      <span class="badge badge-outline badge-lg px-4 py-3 rounded-full font-semibold">
        Not selected
      </span>
    {/if}
  </div>
{:else}
  <div class="mt-2 mb-3 text-base text-warning font-semibold">
    ⚠️ No workspace yet. Create one in your profile.
  </div>
{/if}





