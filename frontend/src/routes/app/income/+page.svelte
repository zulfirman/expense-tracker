<script>
  import { page } from '$app/stores';
  import InputIncome from '$components/InputIncome.svelte';
  import PageHeader from '$lib/components/PageHeader.svelte';
  import PageContainer from '$lib/components/PageContainer.svelte';
  import CurrentWorkspaceBadge from '$lib/components/CurrentWorkspaceBadge.svelte';
  import { getPageCode } from '$lib/utils/pageCodes';

  $: pageCode = getPageCode($page.url.pathname);
</script>

<PageContainer pageCode={pageCode}>
  <PageHeader
    title="Add Income"
    subtitle="Record your income and track your earnings."
    pageCode={pageCode}
  />
  <CurrentWorkspaceBadge />
  <div class="card bg-base-100 shadow-xl border border-base-300">
    <InputIncome />
  </div>
</PageContainer>