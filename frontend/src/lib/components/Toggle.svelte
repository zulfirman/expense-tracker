<script>
  import { createEventDispatcher } from 'svelte';
  
  export let checked = false;
  export let disabled = false;
  export let size = 'md'; // sm, md, lg
  export let color = 'primary'; // primary, success, error, warning, info
  
  const dispatch = createEventDispatcher();
  
  function handleChange(e) {
    if (!disabled) {
      checked = e.target.checked;
      dispatch('change', checked);
    }
  }
  
  $: sizeClasses = {
    sm: 'w-10 h-5',
    md: 'w-12 h-6',
    lg: 'w-16 h-8'
  };
  
  $: thumbSizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-5 h-5',
    lg: 'w-7 h-7'
  };
  
  $: translateClasses = {
    sm: checked ? 'translate-x-5' : 'translate-x-0.5',
    md: checked ? 'translate-x-6' : 'translate-x-0.5',
    lg: checked ? 'translate-x-8' : 'translate-x-0.5'
  };
  
  $: bgColorClasses = {
    primary: checked ? 'bg-primary' : 'bg-base-300',
    success: checked ? 'bg-success' : 'bg-base-300',
    error: checked ? 'bg-error' : 'bg-base-300',
    warning: checked ? 'bg-warning' : 'bg-base-300',
    info: checked ? 'bg-info' : 'bg-base-300'
  };
</script>

<label class="relative inline-flex items-center {disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'}">
  <input
    type="checkbox"
    {checked}
    {disabled}
    on:change={handleChange}
    class="sr-only peer"
  />
  <div
    class="relative {sizeClasses[size]} rounded-full transition-all duration-300 ease-in-out {bgColorClasses[color]} {disabled ? 'opacity-50' : ''} peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-offset-2 peer-focus:ring-primary/50"
  >
    <div
      class="absolute top-0.5 left-0.5 {thumbSizeClasses[size]} bg-white rounded-full shadow-lg transition-all duration-300 ease-in-out {translateClasses[size]}"
    ></div>
  </div>
</label>
