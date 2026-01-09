<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { auth } from '$lib/stores/auth';
  import { workspace } from '$lib/stores/workspace';
  import { get } from 'svelte/store';
  import InputExpenses from '$components/InputExpenses.svelte';
  import PageHeader from '$lib/components/PageHeader.svelte';
  import PageContainer from '$lib/components/PageContainer.svelte';
  import CurrentWorkspaceBadge from '$lib/components/CurrentWorkspaceBadge.svelte';
  import WorkspaceModal from '$lib/components/WorkspaceModal.svelte';
  import { getPageCode } from '$lib/utils/pageCodes';

  let showWorkspaceModal = false;

  $: pageCode = getPageCode($page.url.pathname);

  onMount(async () => {
    // Check if user needs to create first workspace
    if ($auth.isAuthenticated && $auth.user) {
      if (!$auth.user.firstSigninCompleted) {
        await workspace.init();
        const wsState = get(workspace);
        if (wsState.list.length === 0) {
          setTimeout(() => {
            showWorkspaceModal = true;
          }, 300);
        }
      }
    }
  });
</script>

<PageContainer pageCode={pageCode}>
  <PageHeader
    title="Add Expense"
    subtitle="Record your expenses and track your spending."
    pageCode={pageCode}
  />
  <CurrentWorkspaceBadge />
  <div class="card bg-base-100 shadow-xl border border-base-300">
    <InputExpenses/>
  </div>
</PageContainer>

<WorkspaceModal
  bind:open={showWorkspaceModal}
  required={true}
  on:saved={() => {
    showWorkspaceModal = false;
  }}
/>
