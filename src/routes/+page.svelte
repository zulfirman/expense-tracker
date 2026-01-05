<script>
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { auth } from '$lib/stores/auth';

  onMount(() => {
    auth.init();
    if ($auth.isAuthenticated) {
      goto('/expenses');
    } else {
      goto('/login');
    }
  });
</script>










