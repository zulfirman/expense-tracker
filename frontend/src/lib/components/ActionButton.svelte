<script>
  export let variant = 'primary'; // primary, soft, error, success
  export let size = 'sm'; // xs, sm, md, lg
  export let loading = false;
  export let disabled = false;
  export let onClick = null;
  
  let buttonClass = `btn btn-${variant} btn-${size}`;
</script>

<button
  class={buttonClass}
  on:click={onClick}
  {disabled}
>
  {#if loading}
    <span class="loading loading-spinner loading-xs"></span>
  {/if}
  <slot />
</button>


