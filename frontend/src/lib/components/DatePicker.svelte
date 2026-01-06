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
  class="date-picker-input"
/>

<style>
  :global(.date-picker-input) {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    font-size: 1rem;
    font-family: inherit;
    background: var(--background);
    color: var(--text-primary);
    cursor: pointer;
  }

  :global(.date-picker-input:focus) {
    outline: none;
    border-color: var(--primary-color);
  }

  :global(.date-picker-input:disabled) {
    opacity: 0.6;
    cursor: not-allowed;
  }

  /* Flatpickr theme customization */
  :global(.flatpickr-calendar) {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 0.5rem;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    font-family: inherit;
  }

  :global(.flatpickr-months) {
    background: var(--primary-color);
    border-radius: 0.5rem 0.5rem 0 0;
    padding: 0.5rem;
  }

  :global(.flatpickr-month) {
    color: white;
  }

  :global(.flatpickr-current-month) {
    color: white;
    font-weight: 600;
  }

  :global(.flatpickr-prev-month),
  :global(.flatpickr-next-month) {
    color: white;
    fill: white;
  }

  :global(.flatpickr-prev-month:hover),
  :global(.flatpickr-next-month:hover) {
    background: rgba(255, 255, 255, 0.2);
    border-radius: 0.25rem;
  }

  :global(.flatpickr-weekdays) {
    background: var(--background);
    border-bottom: 1px solid var(--border);
  }

  :global(.flatpickr-weekday) {
    color: var(--text-secondary);
    font-weight: 600;
    font-size: 0.875rem;
  }

  :global(.flatpickr-days) {
    background: var(--surface);
  }

  :global(.flatpickr-day) {
    color: var(--text-primary);
    border: 1px solid transparent;
  }

  :global(.flatpickr-day:hover) {
    background: var(--background);
    border-color: var(--primary-color);
  }

  :global(.flatpickr-day.selected) {
    background: var(--primary-color);
    color: white;
    border-color: var(--primary-color);
  }

  :global(.flatpickr-day.today) {
    border-color: var(--primary-color);
    font-weight: 600;
  }

  :global(.flatpickr-day.flatpickr-disabled) {
    color: var(--text-secondary);
    opacity: 0.4;
    cursor: not-allowed;
  }

  :global(.flatpickr-day.prevMonthDay),
  :global(.flatpickr-day.nextMonthDay) {
    color: var(--text-secondary);
    opacity: 0.5;
  }
</style>

