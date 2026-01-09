<script>
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { auth } from '$lib/stores/auth';

  onMount(() => {
    auth.init();
    if ($auth.isAuthenticated) {
      goto('/app/expenses');
    } else {
      goto('/app/login');
    }
  });
</script>
