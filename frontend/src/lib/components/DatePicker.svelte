<script>
  import { onMount, onDestroy } from 'svelte';
  import flatpickr from 'flatpickr';
  import 'flatpickr/dist/flatpickr.min.css';

  export let value = '';
  export let placeholder = 'Select date';
  export let id = '';
  export let disabled = false;
  export let minDate = null;
  export let maxDate = null;
  export let dateFormat = 'Y-m-d';
  export let mode = 'single'; // 'single' or 'range'

  let inputElement;
  let flatpickrInstance;

  $: if (value && flatpickrInstance) {
    flatpickrInstance.setDate(value, false);
  }

  onMount(() => {
    if (inputElement) {
      const options = {
        dateFormat: dateFormat,
        defaultDate: value || undefined,
        minDate: minDate || undefined,
        maxDate: maxDate || undefined,
        mode: mode,
        disableMobile: true, // Use flatpickr on mobile too
        onChange: (selectedDates, dateStr) => {
          value = dateStr;
          // Dispatch custom event
          const event = new CustomEvent('dateChange', { detail: { date: dateStr, dates: selectedDates } });
          inputElement.dispatchEvent(event);
        }
      };

      // Special handling for month picker
      if (dateFormat === 'Y-m') {
        options.enableTime = false;
        options.dateFormat = 'Y-m';
        // Create a custom format function for month display
        options.onReady = (selectedDates, dateStr, instance) => {
          // Set up month navigation
          instance.config.mode = 'single';
        };
      }

      flatpickrInstance = flatpickr(inputElement, options);
    }
  });

  onDestroy(() => {
    if (flatpickrInstance) {
      flatpickrInstance.destroy();
    }
  });

  function handleInput(e) {
    // Allow manual input
    value = e.target.value;
  }
</script>

<input
  bind:this={inputElement}
  {id}
  type="text"
  {placeholder}
  {disabled}
  value={value}
  on:input={handleInput}
  class="input input-bordered w-full border-2"
/>

<style>
  /* Flatpickr theme customization to match DaisyUI */
  :global(.flatpickr-calendar) {
    background: hsl(var(--b2));
    border: 1px solid hsl(var(--b3));
    border-radius: 1rem;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
    font-family: inherit;
  }

  :global(.flatpickr-months) {
    background: hsl(var(--p));
    border-radius: 1rem 1rem 0 0;
    padding: 0.75rem;
  }

  :global(.flatpickr-month) {
    color: hsl(var(--pc));
  }

  :global(.flatpickr-current-month) {
    color: hsl(var(--pc));
    font-weight: 600;
  }

  :global(.flatpickr-prev-month),
  :global(.flatpickr-next-month) {
    color: hsl(var(--pc));
    fill: hsl(var(--pc));
  }

  :global(.flatpickr-prev-month:hover),
  :global(.flatpickr-next-month:hover) {
    background: hsl(var(--p) / 0.2);
    border-radius: 0.5rem;
  }

  :global(.flatpickr-weekdays) {
    background: hsl(var(--b2));
    border-bottom: 1px solid hsl(var(--b3));
  }

  :global(.flatpickr-weekday) {
    color: hsl(var(--bc) / 0.7);
    font-weight: 600;
    font-size: 0.875rem;
  }

  :global(.flatpickr-days) {
    background: hsl(var(--b2));
  }

  :global(.flatpickr-day) {
    color: hsl(var(--bc));
    border: 1px solid transparent;
  }

  :global(.flatpickr-day:hover) {
    background: hsl(var(--b3));
    border-color: hsl(var(--p));
  }

  :global(.flatpickr-day.selected) {
    background: hsl(var(--p));
    color: hsl(var(--pc));
    border-color: hsl(var(--p));
  }

  :global(.flatpickr-day.today) {
    border-color: hsl(var(--p));
    font-weight: 600;
  }

  :global(.flatpickr-day.flatpickr-disabled) {
    color: hsl(var(--bc) / 0.4);
    opacity: 0.4;
    cursor: not-allowed;
  }

  :global(.flatpickr-day.prevMonthDay),
  :global(.flatpickr-day.nextMonthDay) {
    color: hsl(var(--bc) / 0.5);
    opacity: 0.5;
  }
</style>
