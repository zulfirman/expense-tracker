<script>
  import { createEventDispatcher } from 'svelte';

  export let value = '';
  export let placeholder = 'Select date';
  export let id = '';
  export let disabled = false;
  export let minDate = null;
  export let maxDate = null;
  export let label = 'Select date';

  const dispatch = createEventDispatcher();

  let inputEl;

  function handleInput(e) {
    value = e.target.value;
    dispatch('dateChange', { date: value });
  }

  function openNativePicker() {
    if (inputEl?.showPicker) {
      inputEl.showPicker();
    } else {
      inputEl?.focus();
    }
  }
</script>

<fieldset class="fieldset w-full">
  <legend class="fieldset-legend">{label}</legend>
  <input
    bind:this={inputEl}
    {id}
    type="date"
    {placeholder}
    {disabled}
    bind:value
    min={minDate}
    max={maxDate}
    on:input={handleInput}
    on:focus={openNativePicker}
    on:click={openNativePicker}
    class="input input-bordered w-full border-2"
  />
</fieldset>
