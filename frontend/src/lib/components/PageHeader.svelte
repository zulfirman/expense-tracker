<script>
  export let title = '';
  export let subtitle = '';
  export let pageCode = '';
  export let actions = null; // Slot for action buttons
</script>

<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
  <div class="flex items-start gap-2">
    <div>
      {#if pageCode}
        <span class="text-xs font-mono text-base-content/30 mr-2">{pageCode}</span>
      {/if}
      <h1 class="text-2xl font-bold">{title}</h1>
      {#if subtitle}
        <p class="text-sm text-base-content/70 mt-1">{subtitle}</p>
      {/if}
    </div>
  </div>
  {#if actions}
    <div class="flex gap-2 flex-shrink-0">
      <slot name="actions" />
    </div>
  {/if}
</div>



